package files

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/IceBear-CreditEase-LLM/aigc-admin/src/api"
	"github.com/IceBear-CreditEase-LLM/aigc-admin/src/repository"
	"github.com/IceBear-CreditEase-LLM/aigc-admin/src/repository/types"
	"github.com/IceBear-CreditEase-LLM/aigc-admin/src/util"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/google/uuid"
	"io"
	"mime/multipart"
	"strings"
	"time"
)

type Service interface {
	CreateFile(ctx context.Context, request FileRequest) (file File, err error)
	ListFiles(ctx context.Context, request ListFileRequest) (files FileList, err error)
	GetFile(ctx context.Context, fileId string) (file File, err error)
	DeleteFile(ctx context.Context, fileId string) (err error)
	UploadToS3(ctx context.Context, file multipart.File, fileType string, isPublicBucket bool) (s3Url string, err error)
}

type service struct {
	logger  log.Logger
	traceId string
	store   repository.Repository
	apiSvc  api.Service
	Config
}

func (s *service) UploadToS3(ctx context.Context, file multipart.File, fileType string, isPublicBucket bool) (s3Url string, err error) {
	// 如果 isPublicBucket 为true,则使用公有桶，否则使用私有桶
	bucketName := s.S3.BucketName //默认私有桶
	if isPublicBucket == true {
		bucketName = s.S3.BucketNamePublic
	}
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId), "method", "UploadToS3")
	paths := strings.Split(time.Now().Format(time.RFC3339), "-")
	id, _ := util.GenShortId(24)
	fileName := fmt.Sprintf("%s.%s", id, fileType)
	targetPath := fmt.Sprintf("%s/%s/%s/%s", fileType, paths[0], paths[1], fileName)
	err = s.apiSvc.S3Client(ctx).Upload(ctx, bucketName, targetPath, file, "")
	if err != nil {
		return "", err
	}

	shareUrl, err := s.apiSvc.S3Client(ctx).ShareGen(ctx, bucketName, targetPath, 60*24*31*12)
	if err != nil {
		_ = level.Error(logger).Log("pan", "ShareGen", "err", err.Error())
		return
	}

	return shareUrl, nil
}

func (s *service) CreateFile(ctx context.Context, request FileRequest) (file File, err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId), "method", "CreateFile")

	defer request.file.Close()
	// 计算文件md5
	hash := md5.New()
	if _, err = io.Copy(hash, request.file); err != nil {
		_ = level.Error(logger).Log("io.Copy", "err", err.Error())
		return file, err
	}
	md5Str := hex.EncodeToString(hash.Sum(nil))
	if _, err = request.file.Seek(0, io.SeekStart); err != nil {
		_ = level.Error(logger).Log("request.file.Seek", "err", err.Error())
		return file, err
	}

	// 根据md5查询文件是否已经存在
	if res, err := s.store.Files().FindFileByMd5(ctx, md5Str); err == nil && res.ID > 0 {
		file = File{
			FileId:    res.FileID,
			FileName:  res.Name,
			Size:      res.Size,
			FileType:  res.Type,
			S3Url:     res.S3Url,
			Purpose:   res.Purpose,
			TenantId:  res.TenantID,
			CreatedAt: res.CreatedAt,
			LineCount: res.LineCount,
		}
		return file, nil
	}

	// 文件上传云盘
	//url, err := s.UploadToS3(ctx, request.file, request.FileType, false)
	//if err != nil {
	//	return file, err
	//}
	// 文件上传到本地

	// 保存文件信息到数据库
	data := &types.Files{
		FileID:     uuid.New().String(),
		Name:       request.header.Filename,
		Size:       request.header.Size,
		Type:       request.FileType,
		Md5:        md5Str,
		S3Url:      "",
		Purpose:    request.Purpose,
		TenantID:   request.TenantId,
		LineCount:  request.LineCount,
		TokenCount: request.TokenCount,
	}
	err = s.store.Files().CreateFile(ctx, data)
	if err != nil {
		_ = level.Error(logger).Log("store.Chat", "CreateFile", "err", err.Error())
		return file, err
	}
	file = File{
		FileId:    data.FileID,
		FileName:  data.Name,
		Size:      data.Size,
		FileType:  data.Type,
		S3Url:     data.S3Url,
		Purpose:   data.Purpose,
		TenantId:  data.TenantID,
		CreatedAt: data.CreatedAt,
		LineCount: data.LineCount,
	}
	return file, nil
}

func (s *service) ListFiles(ctx context.Context, request ListFileRequest) (files FileList, err error) {
	files.Files = make([]File, 0)
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId), "method", "ListFiles")
	res, total, err := s.store.Files().ListFiles(ctx, request.TenantId, request.Purpose, request.FileName, request.FileType, request.Page, request.PageSize)
	if err != nil {
		_ = level.Error(logger).Log("store.Chat", "ListFiles", "err", err.Error())
		return files, err
	}
	for _, v := range res {
		files.Files = append(files.Files, File{
			FileId:    v.FileID,
			FileName:  v.Name,
			Size:      v.Size,
			FileType:  v.Type,
			S3Url:     v.S3Url,
			Purpose:   v.Purpose,
			TenantId:  v.TenantID,
			CreatedAt: v.CreatedAt,
		})
	}
	files.Total = total
	return files, nil
}

func (s *service) GetFile(ctx context.Context, fileId string) (file File, err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId), "method", "GetFile")
	res, err := s.store.Files().FindFileByFileId(ctx, fileId)
	if err != nil {
		_ = level.Error(logger).Log("store.Chat", "FindFileByFileId", "err", err.Error(), "fileId", fileId)
		return file, err
	}
	file = File{
		FileId:    res.FileID,
		FileName:  res.Name,
		Size:      res.Size,
		FileType:  res.Type,
		S3Url:     res.S3Url,
		Purpose:   res.Purpose,
		TenantId:  res.TenantID,
		CreatedAt: res.CreatedAt,
	}
	return file, nil
}

func (s *service) DeleteFile(ctx context.Context, fileId string) (err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId), "method", "DeleteFile")
	err = s.store.Files().DeleteFile(ctx, fileId)
	if err != nil {
		_ = level.Error(logger).Log("store.Chat", "DeleteFile", "err", err.Error(), "fileId", fileId)
		return err
	}
	return nil
}

type Config struct {
	S3 struct {
		AccessKey        string
		SecretKey        string
		BucketName       string //默认私有桶配置
		BucketNamePublic string //公有桶配置
		ProjectName      string
	}
}

func NewService(logger log.Logger, traceId string, store repository.Repository, apiSvc api.Service, cfg Config) Service {
	return &service{
		logger:  logger,
		traceId: traceId,
		store:   store,
		Config:  cfg,
		apiSvc:  apiSvc,
	}
}

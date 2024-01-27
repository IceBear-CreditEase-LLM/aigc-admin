package files

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"mime/multipart"
)

type tracing struct {
	next   Service
	tracer opentracing.Tracer
}

func (t tracing) UploadToS3(ctx context.Context, file multipart.File, fileType string, isPublicBucket bool) (s3Url string, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, t.tracer, "", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "pkg.files",
	})
	defer func() {
		span.LogKV("file", file, "fileType", fileType, "isPublicBucket", isPublicBucket, "err", err)
		span.SetTag(string(ext.Error), err != nil)
		span.Finish()
	}()
	return t.next.UploadToS3(ctx, file, fileType, isPublicBucket)
}

func (t tracing) CreateFile(ctx context.Context, request FileRequest) (file File, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, t.tracer, "CreateFile", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "pkg.files",
	})
	defer func() {
		span.LogKV("channelId", request.TenantId, "purpose", request.Purpose, "err", err)
		span.SetTag("err", err != nil)
		span.Finish()
	}()
	return t.next.CreateFile(ctx, request)
}

func (t tracing) ListFiles(ctx context.Context, request ListFileRequest) (files FileList, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, t.tracer, "ListFiles", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "pkg.files",
	})
	defer func() {
		span.LogKV("tenantId", request.TenantId, "purpose", request.Purpose, "fileName", request.FileName, "fileType", request.FileType, "page", request.Page, "pageSize", request.PageSize, "err", err)
		span.SetTag("err", err != nil)
		span.Finish()
	}()
	return t.next.ListFiles(ctx, request)
}

func (t tracing) GetFile(ctx context.Context, fileId string) (file File, err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, t.tracer, "GetFile", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "pkg.files",
	})
	defer func() {
		span.LogKV("fileId", fileId, "err", err)
		span.SetTag("err", err != nil)
		span.Finish()
	}()
	return t.next.GetFile(ctx, fileId)
}

func (t tracing) DeleteFile(ctx context.Context, fileId string) (err error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, t.tracer, "DeleteFile", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "pkg.files",
	})
	defer func() {
		span.LogKV("fileId", fileId, "err", err)
		span.SetTag("err", err != nil)
		span.Finish()
	}()
	return t.next.DeleteFile(ctx, fileId)
}

func NewTracing(otTracer opentracing.Tracer) Middleware {
	return func(next Service) Service {
		return &tracing{
			next:   next,
			tracer: otTracer,
		}
	}
}

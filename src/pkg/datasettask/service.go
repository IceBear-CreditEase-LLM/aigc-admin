package datasettask

import (
	"context"
	"fmt"
	"github.com/IceBear-CreditEase-LLM/aigc-admin/src/middleware"
	"github.com/IceBear-CreditEase-LLM/aigc-admin/src/repository"
	"github.com/IceBear-CreditEase-LLM/aigc-admin/src/repository/types"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"strconv"
	"strings"
	"time"
)

// Service is the interface that provides datasettask methods.
type Service interface {
	// CreateTask creates a new task.
	CreateTask(ctx context.Context, tenantId uint, req taskCreateRequest) (err error)
	// ListTasks returns all tasks.
	ListTasks(ctx context.Context, tenantId uint, name string, page, pageSize int) (res []taskDetail, total int64, err error)
	// DeleteTask deletes a task.
	DeleteTask(ctx context.Context, tenantId uint, uuid string) (err error)
	// GetTaskSegmentNext 获取一条待标注任务样本
	GetTaskSegmentNext(ctx context.Context, tenantId uint, taskId string) (res taskSegmentDetail, err error)
	// AnnotationTaskSegment 标注一条任务样本
	AnnotationTaskSegment(ctx context.Context, tenantId uint, taskId, taskSegmentId string, req taskSegmentAnnotationRequest) (err error)
	// AbandonTaskSegment 放弃一条标注任务样本
	AbandonTaskSegment(ctx context.Context, tenantId uint, taskId string) (err error)
	// SyncCheckTaskDatasetSimilar 同步检查标注任务的数据集相似
	SyncCheckTaskDatasetSimilar(ctx context.Context, tenantId uint, taskId string) (err error)
	// SplitAnnotationDataSegment 将标注数据拆分成训练集和测试集
	SplitAnnotationDataSegment(ctx context.Context, tenantId uint, taskId string, req taskSplitAnnotationDataRequest) (err error)
	// ExportAnnotationData 导出标注任务数据
	ExportAnnotationData(ctx context.Context, tenantId uint, taskId string, formatType string) (err error)
	// DeleteAnnotationTask 删除标注任务
	DeleteAnnotationTask(ctx context.Context, tenantId uint, taskId string) (err error)
	// CleanAnnotationTask 清理标注任务
	CleanAnnotationTask(ctx context.Context, tenantId uint, taskId string) (err error)
}

type service struct {
	traceId    string
	logger     log.Logger
	repository repository.Repository
}

func (s *service) CleanAnnotationTask(ctx context.Context, tenantId uint, taskId string) (err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId))
	task, err := s.repository.DatasetTask().GetTask(ctx, tenantId, taskId)
	if err != nil {
		err = errors.Wrap(err, "get task failed")
		_ = level.Warn(logger).Log("msg", "get task failed", "err", err)
		return
	}

	if task.Status != types.DatasetAnnotationStatusProcessing {
		err = errors.New("the annotation task is not processing, cannot be cleaned")
		_ = level.Warn(logger).Log("msg", "the annotation task is not processing, cannot be cleaned", "err", err)
		return
	}
	now := time.Now()
	task.Status = types.DatasetAnnotationStatusCleaned
	task.CompletedAt = &now
	if err = s.repository.DatasetTask().UpdateTask(ctx, task); err != nil {
		err = errors.Wrap(err, "update task failed")
		_ = level.Error(logger).Log("msg", "update task failed", "err", err)
		return
	}

	return
}

func (s *service) ListTasks(ctx context.Context, tenantId uint, name string, page, pageSize int) (res []taskDetail, total int64, err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId))
	tasks, total, err := s.repository.DatasetTask().ListTasks(ctx, tenantId, name, page, pageSize)
	if err != nil {
		err = errors.Wrap(err, "list tasks failed")
		_ = level.Warn(logger).Log("msg", "list tasks failed", "err", err)
		return
	}
	for _, task := range tasks {
		sequence := strings.Split(task.DataSequence, "-")
		start, _ := strconv.Atoi(sequence[0])
		end, _ := strconv.Atoi(sequence[1])
		res = append(res, taskDetail{
			UUID:           task.UUID,
			Name:           task.Name,
			Remark:         task.Remark,
			AnnotationType: task.AnnotationType,
			Principal:      task.Principal,
			Status:         string(task.Status),
			Total:          task.Total,
			Completed:      task.Completed,
			DataSequence:   []int{start, end},
			CreatedAt:      task.CreatedAt,
			CompletedAt:    task.CompletedAt,
			Abandoned:      task.Abandoned,
			TrainTotal:     task.TrainTotal,
			TestTotal:      task.TestTotal,
			TestReport:     task.TestReport,
		})
	}
	return
}

func (s *service) DeleteTask(ctx context.Context, tenantId uint, uuid string) (err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId))
	task, err := s.repository.DatasetTask().GetTask(ctx, tenantId, uuid)
	if err != nil {
		err = errors.Wrap(err, "get task failed")
		_ = level.Warn(logger).Log("msg", "get task failed", "err", err)
		return
	}
	if task.Status == types.DatasetAnnotationStatusProcessing {
		err = errors.New("the annotation task is processing, cannot be deleted")
		_ = level.Warn(logger).Log("msg", "the annotation task is processing, cannot be deleted", "err", err)
		return
	}
	if err = s.repository.DatasetTask().DeleteTask(ctx, tenantId, uuid); err != nil {
		err = errors.Wrap(err, "delete task failed")
		_ = level.Error(logger).Log("msg", "delete task failed", "err", err)
		return
	}
	return
}

func (s *service) GetTaskSegmentNext(ctx context.Context, tenantId uint, taskId string) (res taskSegmentDetail, err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId))
	task, err := s.repository.DatasetTask().GetTask(ctx, tenantId, taskId)
	if err != nil {
		err = errors.Wrap(err, "get task failed")
		_ = level.Warn(logger).Log("msg", "get task failed", "err", err)
		return
	}
	if task.Status != types.DatasetAnnotationStatusProcessing && task.Status != types.DatasetAnnotationStatusPending {
		err = errors.New("the annotation task is not processing, cannot be annotated")
		_ = level.Warn(logger).Log("msg", "the annotation task is not processing, cannot be annotated", "err", err)
		return
	}
	segment, err := s.repository.DatasetTask().GetTaskOneSegment(ctx, task.ID, types.DatasetAnnotationStatusPending)
	if err != nil {
		err = errors.Wrap(err, "get task segment next failed")
		_ = level.Warn(logger).Log("msg", "get task segment next failed", "err", err)
		return
	}
	res = taskSegmentDetail{
		UUID:           segment.UUID,
		AnnotationType: segment.AnnotationType,
		SegmentContent: segment.SegmentContent,
		Status:         string(segment.Status),
		CreatedAt:      segment.CreatedAt,
		Document:       segment.Document,
		Instruction:    segment.Instruction,
		Input:          segment.Input,
		Question:       segment.Question,
		Intent:         segment.Intent,
		Output:         segment.Output,
		CreatorEmail:   segment.CreatorEmail,
	}
	return
}

func (s *service) AnnotationTaskSegment(ctx context.Context, tenantId uint, taskId, taskSegmentId string, req taskSegmentAnnotationRequest) (err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId))
	email, _ := ctx.Value(middleware.ContextKeyUserEmail).(string)
	task, err := s.repository.DatasetTask().GetTask(ctx, tenantId, taskId)
	if err != nil {
		err = errors.Wrap(err, "get task failed")
		_ = level.Warn(logger).Log("msg", "get task failed", "err", err)
		return
	}
	if task.Status != types.DatasetAnnotationStatusProcessing && task.Status != types.DatasetAnnotationStatusPending {
		err = errors.New("the annotation task is not processing, cannot be annotated")
		_ = level.Warn(logger).Log("msg", "the annotation task is not processing, cannot be annotated", "err", err)
		return
	}
	segment, err := s.repository.DatasetTask().GetTaskSegmentByUUID(ctx, task.ID, taskSegmentId)
	if err != nil {
		err = errors.Wrap(err, "get task segment by uuid failed")
		_ = level.Warn(logger).Log("msg", "get task segment by uuid failed", "err", err)
		return
	}
	segment.Document = req.Document
	segment.Instruction = req.Instruction
	segment.Input = req.Input
	segment.Question = req.Question
	segment.Intent = req.Intent
	segment.Output = req.Output
	segment.CreatorEmail = email
	segment.SegmentType = types.DatasetAnnotationSegmentTypeTrain
	segment.Status = types.DatasetAnnotationStatusCompleted
	if err = s.repository.DatasetTask().UpdateTaskSegment(ctx, segment); err != nil {
		err = errors.Wrap(err, "update task segment failed")
		_ = level.Error(logger).Log("msg", "update task segment failed", "err", err)
		return
	}

	task.Completed = task.Completed + 1
	task.TrainTotal = task.TrainTotal + 1

	if task.Status != types.DatasetAnnotationStatusProcessing {
		task.Status = types.DatasetAnnotationStatusProcessing
	}
	if task.TrainTotal+task.Abandoned >= task.Total {
		task.Status = types.DatasetAnnotationStatusCompleted
		now := time.Now()
		task.CompletedAt = &now
	}
	if err = s.repository.DatasetTask().UpdateTask(ctx, task); err != nil {
		err = errors.Wrap(err, "update task failed")
		_ = level.Error(logger).Log("msg", "update task failed", "err", err)
		return err
	}

	return
}

func (s *service) AbandonTaskSegment(ctx context.Context, tenantId uint, taskId string) (err error) {
	//TODO implement me
	panic("implement me")
}

func (s *service) SyncCheckTaskDatasetSimilar(ctx context.Context, tenantId uint, taskId string) (err error) {
	//TODO implement me
	panic("implement me")
}

func (s *service) SplitAnnotationDataSegment(ctx context.Context, tenantId uint, taskId string, req taskSplitAnnotationDataRequest) (err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId))
	annotationTask, err := s.repository.DatasetTask().GetTask(ctx, tenantId, taskId)
	if err != nil {
		err = errors.Wrap(err, "get task failed")
		_ = level.Warn(logger).Log("msg", "get task failed", "err", err)
		return
	}
	if annotationTask.Status != types.DatasetAnnotationStatusCompleted {
		err = errors.New("the annotation task is not completed, cannot be split")
		_ = level.Warn(logger).Log("msg", "the annotation task is not completed, cannot be split", "err", err)
		return
	}
	if req.TestPercent <= 0 || req.TestPercent >= 1 {
		err = errors.New("the test percent must be between 0 and 1")
		_ = level.Warn(logger).Log("msg", "the test percent must be between 0 and 1", "err", err)
		return
	}

	taskSegments, err := s.repository.DatasetTask().GetTaskSegmentByRand(ctx, annotationTask.ID, req.TestPercent,
		types.DatasetAnnotationStatusCompleted, types.DatasetAnnotationSegmentTypeTrain)
	if err != nil {
		err = errors.Wrap(err, "get dataset document segment by rand failed")
		_ = level.Warn(logger).Log("msg", "get dataset document segment by rand failed", "err", err)
		return
	}

	segmentIds := make([]uint, 0)
	for _, segment := range taskSegments {
		segmentIds = append(segmentIds, segment.ID)
	}

	if err = s.repository.DatasetTask().UpdateTaskSegmentType(ctx, segmentIds, types.DatasetAnnotationSegmentTypeTest); err != nil {
		err = errors.Wrap(err, "update task segment type failed")
		_ = level.Error(logger).Log("msg", "update task segment type failed", "err", err)
		return
	}
	annotationTask.TestTotal = len(segmentIds)
	annotationTask.TrainTotal = annotationTask.Total - annotationTask.TestTotal
	if err = s.repository.DatasetTask().UpdateTask(ctx, annotationTask); err != nil {
		err = errors.Wrap(err, "update task failed")
		_ = level.Error(logger).Log("msg", "update task failed", "err", err)
		return
	}

	return
}

func (s *service) ExportAnnotationData(ctx context.Context, tenantId uint, taskId string, formatType string) (err error) {
	//TODO implement me
	panic("implement me")
}

func (s *service) DeleteAnnotationTask(ctx context.Context, tenantId uint, taskId string) (err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId))
	annotationTask, err := s.repository.DatasetTask().GetTask(ctx, tenantId, taskId)
	if err != nil {
		err = errors.Wrap(err, "get task failed")
		_ = level.Warn(logger).Log("msg", "get task failed", "err", err)
		return
	}
	if annotationTask.Status != types.DatasetAnnotationStatusCompleted && annotationTask.Status != types.DatasetAnnotationStatusCleaned {
		err = errors.New("the annotation task is not completed or cleaned, cannot be deleted")
		_ = level.Warn(logger).Log("msg", "the annotation task is not completed or cleaned, cannot be deleted", "err", err)
		return
	}
	return
}

func (s *service) CreateTask(ctx context.Context, tenantId uint, req taskCreateRequest) (err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId))
	//email, _ := ctx.Value(middleware.ContextKeyUserEmail).(string)

	datasetDocument, err := s.repository.DatasetTask().GetDatasetDocumentByUUID(ctx, tenantId, req.DatasetId)
	if err != nil {
		err = errors.Wrap(err, "get dataset document by uuid failed")
		_ = level.Warn(logger).Log("msg", "get dataset document by uuid failed", "err", err)
		return
	}

	if datasetDocument == nil {
		err = errors.New("dataset document not found")
		_ = level.Warn(logger).Log("msg", "dataset document not found", "err", err)
		return
	}

	total := req.DataSequence[1] - req.DataSequence[0]
	if total <= 0 {
		err = errors.Wrap(err, "total less than or equal to 0")
		_ = level.Warn(logger).Log("msg", "total less than or equal to 0", "err", err)
		return
	}

	datasetTask := types.DatasetAnnotationTask{
		DatasetDocumentId: datasetDocument.ID,
		UUID:              "task-" + uuid.New().String(),
		Name:              req.Name,
		Remark:            req.Remark,
		AnnotationType:    req.AnnotationType,
		TenantID:          tenantId,
		Principal:         req.Principal,
		Status:            types.DatasetAnnotationStatusPending,
		DataSequence:      fmt.Sprintf("%d-%d", req.DataSequence[0], req.DataSequence[1]),
	}

	if err = s.repository.DatasetTask().CreateTask(ctx, &datasetTask); err != nil {
		err = errors.Wrap(err, "create task failed")
		_ = level.Warn(logger).Log("msg", "create task failed", "err", err)
		return
	}

	documentSegments, err := s.repository.DatasetTask().GetDatasetDocumentSegmentByRange(ctx, datasetDocument.ID,
		req.DataSequence[0], req.DataSequence[1])
	if err != nil {
		err = errors.Wrap(err, "get dataset document segment by range failed")
		_ = level.Warn(logger).Log("msg", "get dataset document segment by range failed", "err", err)
		return
	}
	datasetTask.Total = len(documentSegments)

	var taskSegments []types.DatasetAnnotationTaskSegment
	for _, segment := range documentSegments {
		taskSegments = append(taskSegments, types.DatasetAnnotationTaskSegment{
			DataAnnotationID: datasetTask.ID,
			UUID:             "das-" + uuid.New().String(),
			AnnotationType:   datasetTask.AnnotationType,
			SegmentContent:   segment.SegmentContent,
			Status:           types.DatasetAnnotationStatusPending,
			SegmentID:        segment.ID,
		})
	}

	if err = s.repository.DatasetTask().AddTaskSegments(ctx, taskSegments); err != nil {
		err = errors.Wrap(err, "add task segments failed")
		_ = level.Error(logger).Log("msg", "add task segments failed", "err", err)
		return
	}

	return
}

func New(traceId string, logger log.Logger, repository repository.Repository) Service {
	logger = log.With(logger, "service", "datasettask")
	return &service{
		traceId:    traceId,
		logger:     logger,
		repository: repository,
	}
}

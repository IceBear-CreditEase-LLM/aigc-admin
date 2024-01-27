package model

import (
	"context"
	"github.com/IceBear-CreditEase-LLM/aigc-admin/src/helpers/page"
	"github.com/IceBear-CreditEase-LLM/aigc-admin/src/repository/types"
	"gorm.io/gorm"
	"time"
)

type Middleware func(Service) Service

type ListModelRequest struct {
	Page         int    `json:"page"`
	PageSize     int    `json:"pageSize"`
	ModelName    string `json:"modelName"`
	Enabled      *bool  `json:"enabled"`
	IsPrivate    *bool  `json:"isPrivate"`
	IsFineTuning *bool  `json:"isFineTuning"`
	ProviderName string `json:"providerName"`
}

type ListEvalRequest struct {
	Page        int    `json:"page"`
	PageSize    int    `json:"pageSize"`
	ModelName   string `json:"modelName"`
	MetricName  string `json:"metricName"`
	Status      string `json:"status"`
	DatasetType string `json:"datasetType"`
}

type Service interface {
	// ListModels 模型分页列表
	ListModels(ctx context.Context, request ListModelRequest) (res []types.Models, total int64, err error)
	// CreateModel 创建模型
	CreateModel(ctx context.Context, data *types.Models) (err error)
	// GetModel 获取模型
	GetModel(ctx context.Context, id uint, preload ...string) (res types.Models, err error)
	// UpdateModel 更新模型
	UpdateModel(ctx context.Context, request UpdateModelRequest) (err error)
	// DeleteModel 删除模型
	DeleteModel(ctx context.Context, id uint) (err error)
	// FindModelsByTenantId 根据租户id查询模型
	FindModelsByTenantId(ctx context.Context, tenantId uint) (res []types.Models, err error)
	// GetModelByModelName 根据模型名称查询模型
	GetModelByModelName(ctx context.Context, modelName string) (res types.Models, err error)
	// CreateEval 创建评估任务
	CreateEval(ctx context.Context, data *types.LLMEvalResults) (err error)
	// ListEval 评估任务分页列表
	ListEval(ctx context.Context, request ListEvalRequest) (res []types.LLMEvalResults, total int64, err error)
	// UpdateEval 更新评估任务
	UpdateEval(ctx context.Context, data *types.LLMEvalResults) (err error)
	// GetEval 获取评估任务
	GetEval(ctx context.Context, id uint) (res types.LLMEvalResults, err error)
	// DeleteEval 删除评估任务
	DeleteEval(ctx context.Context, id uint) (err error)
}

type service struct {
	db *gorm.DB
}

func (s *service) DeleteEval(ctx context.Context, id uint) (err error) {
	err = s.db.WithContext(ctx).Where("id = ?", id).Delete(&types.LLMEvalResults{}).Error
	return
}

func (s *service) UpdateEval(ctx context.Context, data *types.LLMEvalResults) (err error) {
	err = s.db.WithContext(ctx).Save(data).Error
	return
}

func (s *service) GetEval(ctx context.Context, id uint) (res types.LLMEvalResults, err error) {
	err = s.db.WithContext(ctx).Where("id = ?", id).First(&res).Error
	return
}

func (s *service) CreateEval(ctx context.Context, data *types.LLMEvalResults) (err error) {
	err = s.db.WithContext(ctx).Create(data).Error
	return
}

func (s *service) ListEval(ctx context.Context, request ListEvalRequest) (res []types.LLMEvalResults, total int64, err error) {
	query := s.db.WithContext(ctx).Model(&types.LLMEvalResults{})
	if request.ModelName != "" {
		query = query.Where("model_name = ?", request.ModelName)
	}
	if request.MetricName != "" {
		query = query.Where("metric_name = ?", request.MetricName)
	}
	if request.Status != "" {
		query = query.Where("status = ?", request.Status)
	}

	if request.DatasetType != "" {
		query = query.Where("dataset_type = ?", request.DatasetType)
	}
	limit, offset := page.Limit(request.Page, request.PageSize)
	err = query.Count(&total).Order("id DESC").Limit(limit).Offset(offset).Find(&res).Error
	return
}

func (s *service) GetModelByModelName(ctx context.Context, modelName string) (res types.Models, err error) {
	err = s.db.WithContext(ctx).Where("model_name = ?", modelName).First(&res).Error
	return
}

func (s *service) ListModels(ctx context.Context, request ListModelRequest) (res []types.Models, total int64, err error) {
	query := s.db.WithContext(ctx).Model(&types.Models{})
	if request.ModelName != "" {
		query = query.Where("model_name LIKE ?", "%"+request.ModelName+"%")
	}
	if request.Enabled != nil {
		query = query.Where("enabled = ?", *request.Enabled)
	}
	if request.IsPrivate != nil {
		query = query.Where("is_private = ?", *request.IsPrivate)
	}
	if request.IsFineTuning != nil {
		query = query.Where("is_fine_tuning = ?", *request.IsFineTuning)
	}
	if request.ProviderName != "" {
		query = query.Where("provider_name = ?", request.ProviderName)
	}
	limit, offset := page.Limit(request.Page, request.PageSize)
	err = query.Count(&total).Order("updated_at DESC").Limit(limit).Offset(offset).Preload("Tenants").Preload("ModelDeploy").Find(&res).Error
	return
}

func (s *service) CreateModel(ctx context.Context, data *types.Models) (err error) {
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err = tx.Model(&types.Models{}).Create(data).Error; err != nil {
			return err
		}
		if len(data.TenantId) > 0 {
			models := make([]types.TenantModelAssociations, 0)
			for _, v := range data.TenantId {
				models = append(models, types.TenantModelAssociations{
					ModelID:  data.ID,
					TenantID: v,
				})
			}
			if err = tx.Model(&types.TenantModelAssociations{}).Create(&models).Error; err != nil {
				return err
			}
		}
		return nil
	})
	return
}

func (s *service) GetModel(ctx context.Context, id uint, preload ...string) (res types.Models, err error) {
	db := s.db.WithContext(ctx)
	for _, v := range preload {
		db = db.Preload(v)
	}
	err = db.Where("id = ?", id).First(&res).Error
	return
}

type UpdateModelRequest struct {
	Id        uint
	TenantId  *[]uint
	MaxTokens *int
	Enabled   *bool
	Remark    *string
}

func (s *service) UpdateModel(ctx context.Context, request UpdateModelRequest) (err error) {
	data, err := s.GetModel(ctx, request.Id)
	if err != nil {
		return
	}
	if request.MaxTokens != nil {
		data.MaxTokens = *request.MaxTokens
	}
	if request.Enabled != nil {
		data.Enabled = *request.Enabled
	}
	if request.Remark != nil {
		data.Remark = *request.Remark
	}
	data.UpdatedAt = time.Now()
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.WithContext(ctx).Save(data).Error; err != nil {
			return err
		}
		if request.TenantId != nil {
			// 删除原有的关联关系
			if err := tx.WithContext(ctx).Where("model_id = ?", request.Id).Delete(&types.TenantModelAssociations{}).Error; err != nil {
				return err
			}
			// 创建新的关联关系
			if len(*request.TenantId) > 0 {
				models := make([]types.TenantModelAssociations, 0)
				for _, v := range *request.TenantId {
					models = append(models, types.TenantModelAssociations{
						ModelID:  request.Id,
						TenantID: v,
					})
				}
				if err := tx.WithContext(ctx).Model(&types.TenantModelAssociations{}).Create(&models).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
	return
}

func (s *service) DeleteModel(ctx context.Context, id uint) (err error) {
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err = tx.WithContext(ctx).Where("id = ?", id).Delete(&types.Models{}).Error; err != nil {
			return err
		}
		if err = tx.WithContext(ctx).Where("model_id = ?", id).Delete(&types.TenantModelAssociations{}).Error; err != nil {
			return err
		}
		if err = tx.WithContext(ctx).Where("model_id = ?", id).Delete(&types.ChannelModelAssociations{}).Error; err != nil {
			return err
		}
		return nil
	})
	return
}

func (s *service) FindModelsByTenantId(ctx context.Context, tenantId uint) (res []types.Models, err error) {
	var tenant types.Tenants
	err = s.db.WithContext(ctx).Where("id = ?", tenantId).Preload("Models").First(&tenant).Error
	res = tenant.Models
	return
}

func New(db *gorm.DB) Service {
	return &service{
		db: db,
	}
}

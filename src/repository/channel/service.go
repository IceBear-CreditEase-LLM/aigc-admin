package channel

import (
	"context"
	"github.com/IceBear-CreditEase-LLM/aigc-admin/src/helpers/page"
	"github.com/IceBear-CreditEase-LLM/aigc-admin/src/repository/types"
	"gorm.io/gorm"
)

type Middleware func(service Service) Service
type Service interface {
	// ListChannels 渠道分页列表
	ListChannels(ctx context.Context, request ListChannelRequest) (res []types.ChatChannels, total int64, err error)
	// CreateChannel 创建渠道
	CreateChannel(ctx context.Context, data *types.ChatChannels) (err error)
	// GetChannel 获取渠道
	GetChannel(ctx context.Context, id uint) (res types.ChatChannels, err error)
	// UpdateChannel 更新渠道
	UpdateChannel(ctx context.Context, data *types.ChatChannels) (err error)
	// DeleteChannel 删除渠道
	DeleteChannel(ctx context.Context, id uint) (err error)
}

type service struct {
	db *gorm.DB
}
type ListChannelRequest struct {
	Page        int     `json:"page"`
	PageSize    int     `json:"pageSize"`
	Name        *string `json:"name"`
	Email       *string `json:"email"`
	ProjectName *string `json:"projectName"`
	ServiceName *string `json:"serviceName"`
	TenantId    uint    `json:"tenantId"`
}

func (s *service) ListChannels(ctx context.Context, request ListChannelRequest) (res []types.ChatChannels, total int64, err error) {
	query := s.db.WithContext(ctx).Model(&types.ChatChannels{}).Where("tenant_id = ?", request.TenantId)
	if request.Name != nil {
		query = query.Where("name LIKE ?", "%"+*request.Name+"%")
	}
	if request.Email != nil {
		query = query.Where("email LIKE ?", "%"+*request.Email+"%")
	}
	if request.ProjectName != nil {
		query = query.Where("project_name LIKE ?", "%"+*request.ProjectName+"%")
	}
	if request.ServiceName != nil {
		query = query.Where("service_name LIKE ?", "%"+*request.ServiceName+"%")
	}
	limit, offset := page.Limit(request.Page, request.PageSize)
	err = query.Count(&total).Order("id DESC").Limit(limit).Offset(offset).Preload("ChannelModels").Find(&res).Error
	return
}
func (s *service) CreateChannel(ctx context.Context, data *types.ChatChannels) (err error) {
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err = tx.WithContext(ctx).Model(&types.ChatChannels{}).Create(data).Error; err != nil {
			return err
		}
		if len(data.ModelId) > 0 {
			channelModels := make([]types.ChannelModelAssociations, 0)
			for _, v := range data.ModelId {
				channelModels = append(channelModels, types.ChannelModelAssociations{
					ChannelID: data.ID,
					ModelID:   v,
				})
			}
			if err = tx.WithContext(ctx).Model(&types.ChannelModelAssociations{}).Create(channelModels).Error; err != nil {
				return err
			}
		}
		return err
	})
	return
}
func (s *service) GetChannel(ctx context.Context, id uint) (res types.ChatChannels, err error) {
	err = s.db.WithContext(ctx).Where("id = ?", id).First(&res).Error
	return
}
func (s *service) UpdateChannel(ctx context.Context, data *types.ChatChannels) (err error) {
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err = tx.Save(data).Error
		if err != nil {
			return err
		}
		if len(data.ModelId) > 0 {
			if err = tx.Where("channel_id = ?", data.ID).Delete(&types.ChannelModelAssociations{}).Error; err != nil {
				return err
			}
			channelModels := make([]types.ChannelModelAssociations, 0)
			for _, v := range data.ModelId {
				channelModels = append(channelModels, types.ChannelModelAssociations{
					ChannelID: data.ID,
					ModelID:   v,
				})
			}
			if err = tx.WithContext(ctx).Model(&types.ChannelModelAssociations{}).Create(channelModels).Error; err != nil {
				return err
			}
		}
		return err
	})
	return
}

func (s *service) DeleteChannel(ctx context.Context, id uint) (err error) {
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err = tx.WithContext(ctx).Where("channel_id = ?", id).Delete(&types.ChannelModelAssociations{}).Error
		if err != nil {
			return err
		}
		err = tx.WithContext(ctx).Where("id = ?", id).Delete(&types.ChatChannels{}).Error
		return err
	})
	return
}
func New(db *gorm.DB) Service {
	return &service{db: db}
}

package types

import (
	"gorm.io/gorm"
	"time"
)

// ChatChannels 渠道表
type ChatChannels struct {
	gorm.Model
	Name          string   `gorm:"column:name;size:64;not null;unique;index;comment:名称" json:"name"`
	Alias         string   `gorm:"column:alias;size:64;null;default:'New Channel';comment:渠道名称" json:"alias"`
	Remark        string   `gorm:"column:remark;size:128;null;comment:备注" json:"remark"`
	Quota         int      `gorm:"column:quota;null;default:10;comment:配额" json:"quota"`
	Models        string   `gorm:"column:models;size:255;null;comment:支持模型" json:"models"`
	OnlyOpenAI    bool     `gorm:"column:only_openai;null;default:false;comment:仅使用openai" json:"only_openai"`
	ApiKey        string   `gorm:"column:api_key;index;unique;size:128;comment:ApiKey" json:"api_key"`
	Email         string   `gorm:"column:email;size:128;null;comment:邮箱" json:"email"`
	LastOperator  string   `gorm:"column:last_operator;size:100;null;comment:最后操作人" json:"last_operator"`
	TenantId      uint     `gorm:"column:tenant_id;type:bigint(20) unsigned;NOT NULL"` // 租户ID
	ChannelModels []Models `gorm:"many2many:channel_model_associations;foreignKey:id;joinForeignKey:channel_id;References:id;joinReferences:model_id"`
	ModelId       []uint   `gorm:"-" json:"modelId"`
	Tenant        Tenants  `gorm:"foreignKey:TenantId;references:ID"`
}

func (*ChatChannels) TableName() string {
	return "chat_channels"
}

// ChannelModelAssociations 渠道和模型中间表
type ChannelModelAssociations struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	ChannelID uint `gorm:"column:channel_id;type:bigint(20) unsigned;NOT NULL"` // 渠道表主键ID channels.id
	ModelID   uint `gorm:"column:model_id;type:bigint(20) unsigned"`            // 模型表主键ID models.id
}

func (m *ChannelModelAssociations) TableName() string {
	return "channel_model_associations"
}

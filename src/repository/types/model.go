package types

import (
	"gorm.io/gorm"
	"time"
)

// Models 模型表
type Models struct {
	gorm.Model
	ProviderName       ModelProvider      `gorm:"column:provider_name;type:varchar(50);default:localai;NOT NULL"`      // 模型供应商 openai、localai
	ModelType          ModelType          `gorm:"column:model_type;type:varchar(30);default:text-generation;NOT NULL"` // 模型类型 text-generation、embeddings、whisper
	ModelName          string             `gorm:"column:model_name;type:varchar(50);NOT NULL"`                         // 模型名称
	MaxTokens          int                `gorm:"column:max_tokens;type:int(11);default:2048;NOT NULL"`                // 最长上下文
	IsPrivate          bool               `gorm:"column:is_private;type:tinyint(1);default:0;NOT NULL"`                // 是否是私有模型
	IsFineTuning       bool               `gorm:"column:is_fine_tuning;type:tinyint(1);default:0;NOT NULL"`            // 是否是微调模型
	Enabled            bool               `gorm:"column:enabled;type:tinyint(1);default:0"`                            // 是否启用
	Remark             string             `gorm:"column:remark;size:255;null;comment:备注"`
	ModelDeploy        ModelDeploy        `gorm:"foreignKey:ModelID;references:ID"`
	Tenants            []Tenants          `gorm:"many2many:tenant_model_associations;foreignKey:ID;references:ID;joinForeignKey:ModelID;joinReferences:TenantID"`
	TenantId           []uint             `gorm:"-"`
	FineTuningTrainJob FineTuningTrainJob `gorm:"foreignKey:FineTunedModel;references:ModelName"`
	Parameters         float64            `gorm:"column:parameters;type:decimal(7,2);default:0;NOT NULL"` // 模型参数量
	LastOperator       string             `gorm:"column:last_operator;size:100;null;comment:最后操作人"`
}

func (m *Models) TableName() string {
	return "models"
}

func (m *Models) CanDelete() bool {
	return m.IsPrivate && (m.ModelDeploy.Status == "" || m.ModelDeploy.Status == ModelDeployStatusFailed.String())
}

func (m *Models) CanDeploy() bool {
	return m.IsPrivate && (m.ModelDeploy.Status == "" || m.ModelDeploy.Status == ModelDeployStatusFailed.String())
}

func (m *Models) CanUndeploy() bool {
	return m.IsPrivate && (m.ModelDeploy.Status == ModelDeployStatusPending.String() ||
		m.ModelDeploy.Status == ModelDeployStatusRunning.String() ||
		m.ModelDeploy.Status == ModelDeployStatusSuccess.String())
}

// ModelDeploy 模型部署
type ModelDeploy struct {
	gorm.Model
	ModelID     uint   `gorm:"column:model_id;type:bigint(20) unsigned;NOT NULL"` // 模型表主键 models.id
	ModelPath   string `gorm:"column:model_path;type:varchar(255);NOT NULL"`      // 模型部署路径
	Status      string `gorm:"column:status;type:varchar(32)"`                    // 部署状态
	PaasJobName string `gorm:"column:paas_job_name;type:varchar(255)"`            // paas job name
}

func (m *ModelDeploy) TableName() string {
	return "model_deploy"
}

type TenantModelAssociations struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	TenantID  uint `gorm:"column:tenant_id;type:bigint(20) unsigned;NOT NULL"` // 租户ID
	ModelID   uint `gorm:"column:model_id;type:bigint(20) unsigned;NOT NULL"`  // 模型ID
}

func (m *TenantModelAssociations) TableName() string {
	return "tenant_model_associations"
}

type ModelProvider string
type ModelType string
type ModelDeployStatus string

const (
	ModelProviderOpenAI  ModelProvider = "OpenAI"
	ModelProviderLocalAI ModelProvider = "LocalAI"

	ModelTypeTextGeneration ModelType = "text-generation"
	ModelTypeEmbeddings     ModelType = "embeddings"
	ModelTypeWhisper        ModelType = "whisper"

	ModelDeployStatusPending ModelDeployStatus = "pending"
	ModelDeployStatusRunning ModelDeployStatus = "running"
	ModelDeployStatusSuccess ModelDeployStatus = "success"
	ModelDeployStatusFailed  ModelDeployStatus = "failed"
)

func (m ModelProvider) String() string {
	return string(m)
}

func (m ModelType) String() string {
	return string(m)
}

func (m ModelDeployStatus) String() string {
	return string(m)
}

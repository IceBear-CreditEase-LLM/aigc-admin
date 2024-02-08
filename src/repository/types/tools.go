package types

import (
	"gorm.io/gorm"
)

type ToolType string

const (
	// ToolTypeCodeInterpreter 代码解释器
	ToolTypeCodeInterpreter ToolType = "code_interpreter"
	// ToolTypeRetrieval 检索
	ToolTypeRetrieval ToolType = "retrieval"
	// ToolTypeTypeFunction 功能
	ToolTypeTypeFunction ToolType = "function"
)

// Tools 工具
type Tools struct {
	gorm.Model
	// 工具唯一ID
	UUID string `gorm:"column:uuid;size:64;not null;unique;index;comment:UUID"`
	// Name 名称
	Name string `gorm:"column:name;size:64;not null;uniqueIndex:idx_name_tenant;comment:名称"`
	// TenantId 租户ID
	TenantId uint `gorm:"column:tenant_id;size:64;not null;uniqueIndex:idx_name_tenant;comment:租户ID"`
	// Description 描述
	Description string `gorm:"column:description;size:5000;null;comment:描述"`
	// ToolType 工具类型
	ToolType ToolType `gorm:"column:type;size:64;not null;index;comment:工具类型"`
	// Metadata 工具元数据
	Metadata string `gorm:"column:metadata;type:varchar(10240);not null;comment:工具元数据"`
	// Operator 操作人
	Operator string `gorm:"column:operator;size:100;null;comment:操作人"`
	// Remark 备注
	Remark     string       `gorm:"column:remark;size:255;null;comment:备注"`
	Assistants []Assistants `gorm:"many2many:assistant_tool_associations;foreignKey:ID;references:ID;joinReferences:AssistantId;joinForeignKey:ToolId"`
}

// TableName sets the insert table name for this struct type
func (*Tools) TableName() string {
	return "tools"
}

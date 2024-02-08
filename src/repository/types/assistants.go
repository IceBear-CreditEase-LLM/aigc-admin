package types

import "gorm.io/gorm"

// Assistants 助手
type Assistants struct {
	gorm.Model
	// 唯一ID
	UUID string `gorm:"column:uuid;size:64;not null;unique;index;comment:UUID"`
	// TenantId 租户ID
	TenantId uint `gorm:"column:tenant_id;size:64;not null;index;comment:租户ID"`
	// Name 名称
	Name string `gorm:"column:name;size:64;not null;unique;index;comment:名称"`
	// Avatar 头像
	Avatar string `gorm:"column:avatar;size:255;null;comment:头像"`
	//// Alias 别名
	//Alias string `gorm:"column:name;size:64;not null;index;comment:别名"`
	// Remark 描述
	Remark string `gorm:"column:remark;size:500;null;comment:描述"`
	// ModelName 模型名称
	ModelName string `gorm:"column:model_name;size:128;not null;index;comment:模型名称"`
	// Description 助手描述
	Description string `gorm:"column:description;type:varchar(1000);not null;comment:助手描述"`
	// Instructions 助手使用说明
	Instructions string `gorm:"column:instructions;type:varchar(4096);not null;comment:助手使用说明"`
	// Metadata 助手元数据
	Metadata string `gorm:"column:metadata;type:text;not null;comment:助手元数据"`
	// Operator 操作人
	Operator string `gorm:"column:operator;size:64;not null;comment:操作人"`

	//Tools []AssistantToolAssociations `gorm:"foreignKey:AssistantId;references:ID"`
	Tools []Tools `gorm:"many2many:assistant_tool_associations;foreignKey:ID;references:ID;joinReferences:AssistantId;joinForeignKey:ToolId"`
	//AssistantFiles []Files                     `gorm:"many2many:assistant_file_associations;foreignKey:AssistantId;joinForeignKey:AssistantId;References:ID;JoinReferences:ID"`
}

// AssistantMessages 助手消息
type AssistantMessages struct {
	gorm.Model
	// AssistantId 助手ID
	AssistantId uint `gorm:"column:assistant_id;index;not null;comment:助手ID"`
	// Request 请求
	Request string `gorm:"column:request;type:varchar(40960);not null;comment:请求"`
	// Response 响应
	Response string `gorm:"column:response;type:varchar(10240);not null;comment:响应"`
	// Messages 消息
	Messages string `gorm:"column:messages;type:longtext;not null;comment:消息"`

	Tools []Tools `gorm:"many2many:assistant_tool_associations;foreignKey:ID;references:ID;joinReferences:AssistantId;joinForeignKey:ToolId"`
}

// AssistantToolAssociations 助手工具
type AssistantToolAssociations struct {
	//gorm.Model
	// AssistantId 助手ID
	AssistantId uint `gorm:"column:assistant_id;index;not null;comment:助手ID"`
	// ToolId 工具ID
	ToolId uint `gorm:"column:tool_id;index;not null;comment:工具ID"`
}

// AssistantFileAssociations 助手文件关联
type AssistantFileAssociations struct {
	AssistantId uint `gorm:"column:assistant_id;not null;index;comment:助手ID"`
	FileId      uint `gorm:"column:file_id;not null;index;comment:文件ID"`
}

// AssistantToolFunctions 助手工具功能
type AssistantToolFunctions struct {
	gorm.Model
	// Description 功能描述
	Description string `gorm:"column:description;type:varchar(1000);not null;comment:功能描述"`
	// Name 功能名称
	Name string `gorm:"column:name;index;type:varchar(255);not null;comment:功能名称"`
	// Parameters 功能参数
	Parameters string `gorm:"column:parameters;type:varchar(10240);not null;comment:功能参数"`
}

// TableName sets the insert table name for this struct type
func (*Assistants) TableName() string {
	return "assistants"
}

// TableName sets the insert table name for this struct type
func (*AssistantFileAssociations) TableName() string {
	return "assistant_file_associations"
}

// TableName sets the insert table name for this struct type
func (*AssistantToolAssociations) TableName() string {
	return "assistant_tool_associations"
}

// TableName sets the insert table name for this struct type
func (*AssistantToolFunctions) TableName() string {
	return "assistant_tool_functions"
}

// TableName sets the insert table name for this struct type
func (*AssistantMessages) TableName() string {
	return "assistant_messages"
}

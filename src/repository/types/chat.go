/**
 * @Time : 2023/5/12 15:09
 * @Author : solacowa@gmail.com
 * @File : chat
 * @Software: GoLand
 */

package types

import (
	"gorm.io/gorm"
	"strings"
	"time"
)

type ChatModel string
type ChatPromptType string
type ChatSendType string
type FileExt string
type AudioType string

func (c ChatPromptType) String() string {
	return string(c)
}
func (c ChatModel) String() string {
	return string(c)
}

func (c ChatSendType) String() string {
	return string(c)
}

func (c FileExt) String() string {
	return string(c)
}

func (c AudioType) String() string {
	return string(c)
}

const (
	ChatModelDefault           ChatModel = "default"
	ChatModelVicuna13b         ChatModel = "vicuna-13b"
	ChatModelVicuna7b          ChatModel = "vicuna-7b"
	ChatModelAlpaca7b          ChatModel = "alpaca-7b"
	ChatModelAlpaca13b         ChatModel = "alpaca-13b"
	ChatModelAlpaca33b         ChatModel = "alpaca-33b"
	ChatModelLLaMA7b           ChatModel = "llama-7b"
	ChatModelLLaMA13b          ChatModel = "llama-13b"
	ChatModelLLaMA30b          ChatModel = "llama-30b"
	ChatModelLLaMA65b          ChatModel = "llama-65b"
	ChatModelGPT3Dot5Turbo     ChatModel = "gpt-3.5-turbo"
	ChatModelGPT3Dot5Turbo0613 ChatModel = "gpt-3.5-turbo-0613"
	ChatModelGPT4              ChatModel = "gpt-4"
	ChatModelGPT40314          ChatModel = "gpt-4-0314"
	ChatModelGPT40613          ChatModel = "gpt-4-0613"

	ChatPromptTypeImagine       ChatPromptType = "imagine"
	ChatPromptTypeText          ChatPromptType = "text"
	ChatPromptTypeQA            ChatPromptType = "qa"
	ChatSendTypeText            ChatSendType   = "text"
	ChatSendTypeAudio           ChatSendType   = "audio"
	ChatSendTypeImagine         ChatSendType   = "imagine"
	ChatSendTypeFile            ChatSendType   = "file"
	FileExtMp3                  FileExt        = "mp3"
	FileExtPng                  FileExt        = "png"
	FileExtPdf                  FileExt        = "pdf"
	AudioTypeTextToAudio        AudioType      = "text_to_audio"
	AudioTypeAudioTranslation   AudioType      = "audio_translation"
	AudioTypeAudioTranscription AudioType      = "audio_transcription"
)

// Chat 对话消息表
type Chat struct {
	gorm.Model
	ChatModel      ChatModel      `gorm:"column:chat_model;size:64;null;index;default:'default';comment:聊天模型" json:"chat_model"`
	Email          string         `gorm:"column:email;size:128;not null;index;comment:邮箱" json:"email"`
	Prompt         string         `gorm:"column:prompt;type:text;size:8192;not null;comment:提示" json:"prompt"`
	Response       string         `gorm:"column:response;type:longtext;null;comment:内容" json:"response"`
	BeginTime      time.Time      `gorm:"column:begin_time;null;comment:开始时间" json:"begin_time"`
	EndTime        time.Time      `gorm:"column:end_time;null;comment:结束时间" json:"end_time"`
	Temperature    float64        `gorm:"column:temperature;null;comment:温度" json:"temperature"`
	TopP           float64        `gorm:"column:top_p;null;comment:top_p" json:"top_p"`
	Status         int            `gorm:"column:status;null;comment:状态" json:"status"`
	RoleId         uint           `gorm:"column:role_id;null;index;comment:角色id" json:"role_id"`
	TimeCost       string         `gorm:"column:time_cost;size:16;null;comment:耗时" json:"time_cost"`
	Error          bool           `gorm:"column:error;null;default:false;comment:错误" json:"error"`
	ChatId         string         `gorm:"column:chat_id;size:64;null;index;comment:聊天id" json:"chat_id"`
	MaxLength      int            `gorm:"column:max_length;null;comment:最大长度" json:"max_length"`
	ConversationId uint           `gorm:"column:conversation_id;null;index;comment:对话ID" json:"conversation_id"`
	PromptType     ChatPromptType `gorm:"column:prompt_type;size:32;null;default:'text';comment:提示类型" json:"prompt_type"`
	PanUrl         string         `gorm:"column:pan_url;size:50;null;comment:文件云盘地址" json:"pan_url"`
	SendType       ChatSendType   `gorm:"column:send_type;size:10;null;default:'text';comment:发送类型" json:"send_type"`
	Ext            string         `gorm:"column:ext;size:128;null;comment:文件后缀" json:"ext"`

	Conversation ChatConversation `gorm:"foreignKey:ConversationId;references:ID" json:"-"`
}

// ChatConversation 对话表
type ChatConversation struct {
	gorm.Model
	Uuid        string    `gorm:"column:uuid;size:40;not null;unique;comment:UUID" json:"uuid"`
	Alias       string    `gorm:"column:alias;size:64;null;default:'New Chat';comment:别名" json:"alias"`
	Email       string    `gorm:"column:email;size:128;not null;index;comment:邮箱" json:"email"`
	ChatModel   ChatModel `gorm:"column:chat_model;size:64;null;index;default:'default';comment:聊天模型" json:"chat_model"`
	ChannelId   uint      `gorm:"column:channel_id;null;index;comment:渠道ID" json:"channel_id"`
	Temperature float64   `gorm:"column:temperature;null;comment:温度" json:"temperature"`
	TopP        float64   `gorm:"column:top_p;null;comment:采样性" json:"top_p"`
	MaxTokens   int       `gorm:"column:max_tokens;default:2048;null;comment:最大支持Tokens" json:"max_tokens"`
	SysPrompt   string    `gorm:"column:sys_prompt;size:255;null;comment:系统提示" json:"sys_prompt"`

	Channel ChatChannels `gorm:"foreignKey:channel_id;references:id" json:"-"`
}

// ChatSystemPrompt 系统提示表
type ChatSystemPrompt struct {
	gorm.Model
	ChatModel  ChatModel      `gorm:"column:chat_model;size:64;notnull;index;default:'default';comment:聊天模型" json:"chat_model"`
	Content    string         `gorm:"column:content;type:text;size:8192;not null;comment:内容" json:"content"`
	PromptType ChatPromptType `gorm:"column:prompt_type;size:32;notnull;index;default:'text';comment:提示类型" json:"prompt_type"`
}

// ChatPromptTypes 提示类型表
type ChatPromptTypes struct {
	gorm.Model
	Name   string `gorm:"column:name;size:32;not null;index;comment:名称" json:"name"`
	Alias  string `gorm:"column:alias;size:64;null;default:'New Chat';comment:别名" json:"alias"`
	Remark string `gorm:"column:remark;size:128;null;comment:备注" json:"remark"`
}

// ChatPrompts 提示表
type ChatPrompts struct {
	gorm.Model
	Title      string         `gorm:"column:title;size:64;not null;index;comment:标题" json:"title"`
	Content    string         `gorm:"column:content;type:text;size:8192;not null;comment:内容" json:"content"`
	PromptType ChatPromptType `gorm:"column:prompt_type;size:32;notnull;index;default:'text';comment:提示类型" json:"prompt_type"`
}

// ChatChannelModels 渠道模型表
type ChatChannelModels struct {
	gorm.Model
	ChannelId    uint         `gorm:"column:channel_id;null;index;comment:渠道ID" json:"channel_id"`
	ModelName    ChatModel    `gorm:"column:model;size:64;notnull;index;default:'default';comment:聊天模型" json:"model"`
	MaxTokens    int          `gorm:"column:max_tokens;default:2048;null;comment:最大支持Tokens" json:"max_tokens"`
	IsPrivate    bool         `gorm:"column:is_private;null;default:false;comment:是否为本地私有模型" json:"is_private"`
	ChatChannels ChatChannels `gorm:"foreignKey:id;references:channel_id" json:"-"`
}

// ChatMessages 消息表
type ChatMessages struct {
	gorm.Model
	ModelName      ChatModel `gorm:"column:model;size:64;notnull;index;default:'default';comment:聊天模型" json:"model"`
	ChannelId      uint      `gorm:"column:channel_id;null;index;comment:渠道ID" json:"channel_id"`
	Response       string    `gorm:"column:response;type:longtext;size:65536;null;comment:回复" json:"response"`
	Prompt         string    `gorm:"column:prompt;type:text;size:32768;not null;comment:问题" json:"prompt"`
	PromptTokens   int       `gorm:"column:prompt_tokens;default:0;null;comment:问题Tokens" json:"prompt_tokens"`
	ResponseTokens int       `gorm:"column:response_tokens;default:0;null;comment:回复Tokens" json:"response_tokens"`
	Finished       bool      `gorm:"column:finished;default:false;null;comment:是否完成" json:"finished"`
	TimeCost       string    `gorm:"column:time_cost;size:32;null;comment:耗时" json:"time_cost"`
	Temperature    float64   `orm:"column:temperature;default:0.9;null;comment:温度" json:"temperature"`
	TopP           float64   `orm:"column:top_p;default:0.9;null;comment:核心采样" json:"top_p"`
	N              int       `orm:"column:n;default:1;null;comment:聊天完成选项" json:"n"`
	User           string    `orm:"column:user;size:64;null;comment:用户" json:"user"`
	MessageId      string    `orm:"column:message_id;size:128;null;comment:消息ID" json:"message_id"`
	Object         string    `orm:"column:object;size:128;null;comment:对象" json:"object"`
	Created        int64     `gorm:"column:created;null;comment:创建时间" json:"created"`
	// 想想后面函数要怎么搞
}

// 正向标签 masterpiece, best quality, top quality, ultra highres, 8k hdr, 8k wallpaper, RAW, huge file size, intricate details, sharp focus, natural lighting, realistic, professional, delicate, amazing, CG, finely detailed, beautiful detailed, colourful
// 反向标签 paintings, sketches, lowres, normal quality, worst quality, low quality, cropped, dot, mole, ugly, grayscale, monochrome, duplicate, morbid, mutilated, missing fingers, extra fingers, too many fingers, fused fingers, mutated hands, bad hands, poorly drawn hands, poorly drawn face, poorly drawn eyebrows, bad anatomy, cloned face, long neck, extra legs, extra arms, missing arms missing legs, malformed limbs, deformed, simple background, bad proportions, disfigured, skin spots, skin blemishes, age spot, bad feet, error, text, extra digit, fewer digits, jpeg artifacts, signature, username, blurry, watermark, mask, logo

type ChatRole struct {
	gorm.Model
	Name  string `gorm:"column:name;size:32;not null;index;comment:角色名称" json:"name"`
	Alias string `gorm:"column:alias;size:32;null;default:'New Chat';comment:角色别名" json:"alias"`
	Email string `gorm:"column:email;size:128;not null;index;comment:邮箱" json:"email"`
}

type ChatAllowUser struct {
	gorm.Model
	Email string `gorm:"column:email;size:128;not null;index;comment:邮箱" json:"email"`
}

func (c *Chat) TableName() string {
	return "chat"
}

func (*ChatRole) TableName() string {
	return "chat_role"
}

func (*ChatConversation) TableName() string {
	return "chat_conversation"
}

func (*ChatAllowUser) TableName() string {
	return "chat_allow_user"
}

func (*ChatSystemPrompt) TableName() string {
	return "chat_system_prompt"
}

func (*ChatPromptTypes) TableName() string {
	return "chat_prompt_types"
}

func (*ChatPrompts) TableName() string {
	return "chat_prompts"
}

func (*ChatChannelModels) TableName() string {
	return "chat_channel_models"
}

// RemoveGPT4Models 移除含有GPT4的模型
func RemoveGPT4Models(models []ChatChannelModels) []ChatChannelModels {
	if len(models) == 0 {
		return models
	}
	var newModels []ChatChannelModels
	for _, model := range models {
		if !strings.Contains(strings.ToLower(model.ModelName.String()), "gpt-4") {
			newModels = append(newModels, model)
		}
	}
	return newModels
}

// ChatAudio 音频文件
type ChatAudio struct {
	gorm.Model
	ChannelId     uint      `json:"channel_id" gorm:"channel_id"`         // 渠道ID
	FileName      string    `json:"file_name" gorm:"file_name"`           // 上传文件名
	TargetPath    string    `json:"target_path" gorm:"target_path"`       // 网盘路径
	PanUrl        string    `json:"pan_url" gorm:"pan_url"`               // 生成网盘地址
	AudioText     string    `json:"audio_text" gorm:"audio_text"`         // 音频对应的文本
	AudioDuration float64   `json:"audio_duration"`                       // 音频时长
	AudioType     AudioType `json:"audio_type"`                           // 音频处理类型  text_to_audio audio_translation audio_transcription
	TranslateText string    `json:"translate_text" gorm:"translate_text"` // 音频文本翻译
}

// TableName 表名称
func (*ChatAudio) TableName() string {
	return "chat_audio"
}

package types

import (
	"time"
)

// Project 项目模型
type Project struct {
	ID          string            `json:"id" jsonschema:"项目唯一标识"`
	Name        string            `json:"name" jsonschema:"项目名称"`
	Description string            `json:"description" jsonschema:"项目描述"`
	UserIds     []string          `json:"user_ids" jsonschema:"可以访问此项目的用户ID列表"`
	CreatedAt   *time.Time        `json:"created_at" jsonschema:"项目创建时间"`
	UpdatedAt   *time.Time        `json:"updated_at" jsonschema:"项目更新时间"`
	Metadata    map[string]string `json:"metadata" jsonschema:"项目元数据"`
}

// DocumentType 文档类型枚举
type DocumentType string

const (
	DocumentTypeAPI  DocumentType = "api"  // 接口文档
	DocumentTypeText DocumentType = "text" // 文本文档
)

// Document 文档模型
type Document struct {
	ID          string                 `json:"id,omitempty" jsonschema:"文档唯一标识"`
	ProjectID   string                 `json:"project_id,omitempty" jsonschema:"所属项目ID"`
	CreatorID   string                 `json:"creator_id,omitempty" jsonschema:"创建者ID"`
	Type        DocumentType           `json:"type,omitempty" jsonschema:"文档类型：api或text"`
	Name        string                 `json:"name,omitempty" jsonschema:"文档名称"`
	Description string                 `json:"description,omitempty" jsonschema:"文档描述"`
	APIContent  APIDocumentContent     `json:"api_content,omitempty" jsonschema:"接口文档内容"`
	TextContent TextDocumentContent    `json:"text_content,omitempty" jsonschema:"文本文档内容"`
	Tags        []string               `json:"tags,omitempty" jsonschema:"文档标签"`
	CreatedAt   *time.Time             `json:"created_at,omitempty" jsonschema:"文档创建时间"`
	UpdatedAt   *time.Time             `json:"updated_at,omitempty" jsonschema:"文档更新时间"`
	Metadata    map[string]interface{} `json:"metadata,omitempty" jsonschema:"文档元数据"`
}

// APIDocumentContent API文档内容
type APIDocumentContent struct {
	Method          string                 `json:"method" jsonschema:"HTTP方法"`
	Path            string                 `json:"path" jsonschema:"接口路径"`
	Body            map[string]interface{} `json:"body,omitempty" jsonschema:"存在body参数时传入"`
	Query           map[string]interface{} `json:"query,omitempty" jsonschema:"存在query参数时传入"`
	PathParams      map[string]interface{} `json:"path_params,omitempty" jsonschema:"存在path参数时传入"`
	ResponseBizCode string                 `json:"response_biz_code,omitempty" jsonschema:"response 业务码，存在时传入，通常不需要传递该字段，除非用户要求"`
	Header          map[string]interface{} `json:"header,omitempty" jsonschema:"API请求头，可选。键值对格式，通常不需要传递该字段，除非用户要求"`
}

// TextDocumentContent 文本文档内容
type TextDocumentContent struct {
	ContentType string            `json:"content_type" jsonschema:"内容类型，如prompt、note等"`
	Content     interface{}       `json:"content" jsonschema:"文本内容"`
	Variables   map[string]string `json:"variables,omitempty" jsonschema:"模板变量，用于提示词模板"`
}

// UserContext 用户上下文
type UserContext struct {
	UserID    string    `json:"user_id" jsonschema:"用户ID"`
	ProjectID string    `json:"project_id" jsonschema:"当前项目ID"`
	Project   *Project  `json:"project,omitempty" jsonschema:"当前项目信息"`
	ExpiresAt time.Time `json:"expires_at" jsonschema:"上下文过期时间"`
}

// SearchRequest 搜索请求
type SearchRequest struct {
	Query  string        `json:"query" jsonschema:"搜索关键词"`
	Type   *DocumentType `json:"type,omitempty" jsonschema:"文档类型过滤"`
	Tags   []string      `json:"tags,omitempty" jsonschema:"标签过滤"`
	Limit  int           `json:"limit,omitempty" jsonschema:"返回数量限制，默认20"`
	Offset int           `json:"offset,omitempty" jsonschema:"偏移量，默认0"`
}

// SearchResult 搜索结果
type SearchResult struct {
	Documents []Document `json:"documents" jsonschema:"匹配的文档列表"`
	Total     int64      `json:"total" jsonschema:"总数量"`
	Offset    int        `json:"offset" jsonschema:"当前偏移量"`
	Limit     int        `json:"limit" jsonschema:"当前限制"`
}

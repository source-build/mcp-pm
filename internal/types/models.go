package types

import (
	"time"
)

// Project 项目模型
type Project struct {
	ID          string            `json:"id" jsonschema:"项目ID，必填。项目的唯一标识符，用于区分不同项目。建议使用UUID或类似的唯一ID格式。"`
	CreatorId   string            `json:"creator_id,omitempty" jsonschema:"创建者ID，必填。项目创建者的用户ID，用于标识项目的所有者。通常为系统分配的唯一用户标识。"`
	Name        string            `json:"name,omitempty" jsonschema:"项目名称，必填。项目的显示名称，用于界面展示和识别。建议使用有意义的名称，如：'用户管理系统'、'API文档项目'等。长度2-100字符。"`
	Description string            `json:"description,omitempty" jsonschema:"项目描述，可选。项目的详细说明，包括项目用途、功能介绍、技术栈等。支持Markdown格式，建议包含项目背景和使用场景。"`
	UserIds     []string          `json:"user_ids,omitempty" jsonschema:"用户ID列表，必填。可以访问此项目的用户ID数组。包括创建者和被授权的用户。格式：['user-123', 'user-456']。用于权限控制。"`
	CreatedAt   *time.Time        `json:"created_at,omitempty" jsonschema:"创建时间，可选。项目的创建时间戳，格式为RFC3339。自动生成，通常不需要手动设置。"`
	UpdatedAt   *time.Time        `json:"updated_at,omitempty" jsonschema:"更新时间，可选。项目的最后更新时间戳，格式为RFC3339。自动更新，用于跟踪项目修改历史。"`
	Metadata    map[string]string `json:"metadata,omitempty" jsonschema:"项目元数据，可选。项目的附加信息键值对，如：{'env': 'production', 'version': '1.0.0'}。用于存储自定义属性。"`
}

// DocumentType 文档类型枚举
type DocumentType string

const (
	DocumentTypeAPI  DocumentType = "api"  // 接口文档
	DocumentTypeText DocumentType = "text" // 文本文档
)

// Document 文档模型
type Document struct {
	ID          string                 `json:"id,omitempty" jsonschema:"文档ID，可选。文档的唯一标识符，系统自动生成UUID格式。创建时可不传，系统会自动分配。"`
	ProjectID   string                 `json:"project_id,omitempty" jsonschema:"项目ID，必填。文档所属的项目ID，用于关联项目。必须为有效的项目标识符。"`
	CreatorID   string                 `json:"creator_id,omitempty" jsonschema:"创建者ID，必填。文档创建者的用户ID，用于权限控制和作者标识。必须为有效的用户标识符。"`
	Type        DocumentType           `json:"type,omitempty" jsonschema:"文档类型，必填。文档的类型标识：'api'表示接口文档，'text'表示文本文档。必须为有效的枚举值。"`
	Name        string                 `json:"name,omitempty" jsonschema:"文档名称，必填。文档的显示名称，用于搜索和识别。建议使用有意义的名称，如：'用户登录接口'、'系统配置说明'等。长度2-100字符。"`
	Description string                 `json:"description,omitempty" jsonschema:"文档描述，可选。文档的详细说明，包括用途、使用场景、注意事项等。支持Markdown格式，建议包含具体的使用示例。"`
	APIContent  APIDocumentContent     `json:"api_content,omitempty" jsonschema:"API文档内容，可选。当文档类型为'api'时必需，包含HTTP方法、路径、请求响应结构等接口详细信息。"`
	TextContent TextDocumentContent    `json:"text_content,omitempty" jsonschema:"文本文档内容，可选。当文档类型为'text'时必需，包含文本内容、内容类型、模板变量等信息。"`
	Tags        []string               `json:"tags,omitempty" jsonschema:"文档标签，可选。用于分类和搜索的标签数组，如：['auth', 'user', 'v1']。支持多个标签，便于文档管理和筛选。"`
	CreatedAt   *time.Time             `json:"created_at,omitempty" jsonschema:"创建时间，可选。文档的创建时间戳，格式为RFC3339。自动生成，通常不需要手动设置。"`
	UpdatedAt   *time.Time             `json:"updated_at,omitempty" jsonschema:"更新时间，可选。文档的最后更新时间戳，格式为RFC3339。自动更新，用于跟踪文档修改历史。"`
	Metadata    map[string]interface{} `json:"metadata,omitempty" jsonschema:"文档元数据，可选。文档的附加信息键值对，如：{'version': '1.0', 'status': 'active'}。用于存储自定义属性。"`
}

// APIDocumentContent API文档内容
type APIDocumentContent struct {
	Method          string                 `json:"method,omitempty" jsonschema:"HTTP方法，必填。API的请求方法。必须是标准HTTP方法：GET(查询)、POST(创建)、PUT(更新)、DELETE(删除)、PATCH(部分更新)、HEAD、OPTIONS、TRACE。请根据API的实际功能选择。"`
	Path            string                 `json:"path,omitempty" jsonschema:"API路径，必填。API的访问路径，如：'/api/users/{id}'、'/api/orders'。支持路径参数，使用'{paramName}'格式。必须以'/'开头，避免使用特殊字符。"`
	Body            map[string]interface{} `json:"body,omitempty" jsonschema:"请求体，可选。POST/PUT请求的JSON数据，如：{'name': '张三', 'age': 25}。GET/DELETE请求通常不需要。请提供完整的请求体结构和字段说明。"`
	Query           map[string]interface{} `json:"query,omitempty" jsonschema:"查询参数，可选。URL查询参数，如：{'page': 1, 'size': 10}。用于GET请求的参数传递，格式为键值对。"`
	PathParams      map[string]interface{} `json:"path_params,omitempty" jsonschema:"路径参数，可选。URL路径中的参数，如：{'id': 123}。对应路径中的'{paramName}'占位符，用于动态路径构建。"`
	ResponseBizCode string                 `json:"response_biz_code,omitempty" jsonschema:"业务状态码，可选。自定义的业务响应码，如：'SUCCESS'、'USER_NOT_FOUND'。用于业务逻辑判断，通常与HTTP状态码配合使用。"`
	Header          map[string]interface{} `json:"header,omitempty" jsonschema:"请求头，可选。API请求所需的HTTP头部信息，如：{'Authorization': 'Bearer token', 'Content-Type': 'application/json'}。通常不需要传递，除非有特殊认证要求。"`
}

// TextDocumentContent 文本文档内容
type TextDocumentContent struct {
	ContentType string            `json:"content_type" jsonschema:"内容类型，必填。文本文档的具体类型，支持的值：'readme'(README文档)、'prompt'(提示词模板)、'config'(配置文件)、'note'(普通笔记)、'spec'(技术规范)。必须为有效枚举值。"`
	Content     interface{}       `json:"content" jsonschema:"文档内容，必填。文本文档的实际内容，可以是字符串、JSON对象或任意可转换为字符串的内容。根据content_type自动格式化存储。"`
	Variables   map[string]string `json:"variables,omitempty" jsonschema:"模板变量，可选。用于提示词模板的变量键值对，如：{'username': '用户名', 'api_key': 'API密钥'}。仅在content_type为'prompt'时有效。"`
}

// UserContext 用户上下文
type UserContext struct {
	UserID    string    `json:"user_id" jsonschema:"用户ID，必填。当前用户的唯一标识符，用于身份验证和权限控制。必须为有效的用户ID格式。"`
	ProjectID string    `json:"project_id" jsonschema:"项目ID，必填。用户当前操作的项目标识符，用于确定操作上下文。必须为有效的项目ID。"`
	Project   *Project  `json:"project,omitempty" jsonschema:"项目信息，可选。当前项目的详细信息对象，包含项目名称、描述等。通常由系统自动填充，无需手动设置。"`
	ExpiresAt time.Time `json:"expires_at" jsonschema:"过期时间，必填。用户上下文的过期时间戳，格式为RFC3339。用于控制会话有效期，超时后需要重新认证。"`
}

// SearchRequest 搜索请求
type SearchRequest struct {
	Query  string        `json:"query" jsonschema:"搜索关键词，必填。用于搜索文档的关键词，可以是文档名称、路径、描述中的任意词汇。如：'user'、'login'、'/api/users'。支持模糊搜索。"`
	Type   *DocumentType `json:"type,omitempty" jsonschema:"文档类型过滤，可选。按文档类型过滤结果，'api'表示接口文档，'text'表示文本文档。不传则返回所有类型的文档。"`
	Tags   []string      `json:"tags,omitempty" jsonschema:"标签过滤，可选。按标签过滤文档，如：['user', 'auth', 'v1']。支持多个标签组合筛选。不传则忽略标签过滤。"`
	Limit  int           `json:"limit,omitempty" jsonschema:"返回数量限制，可选。控制返回结果数量，默认20，最大100。用于分页浏览。"`
	Offset int           `json:"offset,omitempty" jsonschema:"偏移量，可选。分页查询的起始位置，默认0。用于跳过前面的结果，获取后续数据。"`
}

// SearchResult 搜索结果
type SearchResult struct {
	Documents []Document `json:"documents" jsonschema:"文档列表，必填。匹配搜索条件的文档对象数组，每个文档包含完整的文档信息和内容。"`
	Total     int64      `json:"total" jsonschema:"总数量，必填。匹配条件的文档总数，用于分页计算和结果统计。"`
	Offset    int        `json:"offset" jsonschema:"当前偏移量，必填。当前结果集的起始位置，与请求中的offset参数对应。"`
	Limit     int        `json:"limit" jsonschema:"当前限制，必填。当前返回结果的数量限制，与请求中的limit参数对应。"`
}

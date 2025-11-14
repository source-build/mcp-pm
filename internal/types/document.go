package types

import "time"

// APIDocument 表示接口文档的数据结构
type APIDocument struct {
	ID          string            `json:"id" jsonschema:"Interface document ID"`                    // 接口唯一标识
	Name        string            `json:"name" jsonschema:"Interface document name"`                // 接口名称
	Method      string            `json:"method" jsonschema:"HTTP method"`                          // HTTP方法
	Path        string            `json:"path" jsonschema:"Interface document path"`                // 接口路径
	Description string            `json:"description" jsonschema:"Interface document description"`  // 接口描述
	Request     interface{}       `json:"request" jsonschema:"Request structure"`                   // 请求参数结构
	Response    interface{}       `json:"response" jsonschema:"Response structure"`                 // 响应结构
	Tags        []string          `json:"tags" jsonschema:"Interface document tags"`                // 标签
	CreatedAt   time.Time         `json:"created_at" jsonschema:"Interface document creation time"` // 创建时间
	UpdatedAt   time.Time         `json:"updated_at" jsonschema:"Interface document update time"`   // 更新时间
	Metadata    map[string]string `json:"metadata" jsonschema:"Interface document metadata"`        // 额外元数据
}

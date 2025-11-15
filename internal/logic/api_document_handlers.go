package logic

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/source-build/mcp-pm/internal/config"
	"github.com/source-build/mcp-pm/internal/es"
	"github.com/source-build/mcp-pm/internal/types"
	"github.com/source-build/mcp-pm/internal/utils"
)

/*
API文档管理处理器

注意：这些工具专门用于管理API接口文档，包括：
- RESTful API的端点文档
- GraphQL schema文档
- RPC服务文档
- 微服务接口文档

适用场景：
- API文档自动生成和维护
- 接口规范管理
- 团队协作时的API文档中心
- 与API测试工具集成

数据结构说明：
- DocumentType: "api" 标识这是API文档
- content.API: 包含method, path, request/response等API特定字段
- 标签：用于分类不同类型或版本的API（如REST, GraphQL, v1, v2等）
*/

// CreateAPIDocument 创建API文档
func CreateAPIDocument(_ context.Context, req *mcp.CallToolRequest, input struct {
	Name            string                 `json:"name" jsonschema:"文档名称，必填。可使用中文或英文，比如 userLogin、用户登录接口、用户微信登录接口等"`
	Description     string                 `json:"description" jsonschema:"文档详细描述，必填。清晰说明API的功能、业务场景和使用方法"`
	Method          string                 `json:"method" jsonschema:"HTTP方法，必填。标准HTTP方法之一：GET、POST、PUT、DELETE、PATCH、HEAD、OPTIONS、TRACE等"`
	Path            string                 `json:"path" jsonschema:"API相对路径，必填。支持路径参数"`
	Header          map[string]interface{} `json:"header" jsonschema:"API请求头，可选。键值对格式，通常不需要传递该字段，除非用户要求"`
	Body            map[string]interface{} `json:"body" jsonschema:"存在body参数时传入"`
	Query           map[string]interface{} `json:"query" jsonschema:"存在query参数时传入"`
	PathParams      map[string]interface{} `json:"path_params" jsonschema:"存在path参数时传入"`
	ResponseBizCode string                 `json:"response_biz_code" jsonschema:"response 业务码，存在时传入，通常不需要传递该字段，除非用户要求"`
	Tags            []string               `json:"tags" jsonschema:"文档标签，可选。用于分类和搜索，用于方便用户快速定位和筛选文档，需要总结出用户接口的特点并提取关键字"`
}) (*mcp.CallToolResult, struct {
	Success  bool            `json:"success"`
	Document *types.Document `json:"document"`
	Message  string          `json:"message"`
}, error) {
	userId, err := utils.ExtractUserID(req)
	if err != nil {
		return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{Text: fmt.Sprintf("获取用户ID失败: %v", err)},
				},
				IsError: true,
			}, struct {
				Success  bool            `json:"success"`
				Document *types.Document `json:"document"`
				Message  string          `json:"message"`
			}{Success: false, Message: fmt.Sprintf("获取用户ID失败: %v", err)}, err
	}

	projectId, err := utils.ExtractProjectID(req)
	if err != nil {
		return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{Text: fmt.Sprintf("获取项目ID失败: %v", err)},
				},
				IsError: true,
			}, struct {
				Success  bool            `json:"success"`
				Document *types.Document `json:"document"`
				Message  string          `json:"message"`
			}{Success: false, Message: fmt.Sprintf("获取项目ID失败: %v", err)}, err
	}

	// 验证必需参数
	if input.Name == "" {
		return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{Text: "文档名称不能为空"},
				},
				IsError: true,
			}, struct {
				Success  bool            `json:"success"`
				Document *types.Document `json:"document"`
				Message  string          `json:"message"`
			}{Success: false, Message: "文档名称不能为空"}, fmt.Errorf("文档名称不能为空")
	}

	m := md5.New()
	m.Write([]byte(uuid.NewString()))

	now := time.Now()

	// 创建API文档对象
	document := &types.Document{
		ID:          hex.EncodeToString(m.Sum(nil)),
		ProjectID:   projectId,
		CreatorID:   userId,
		Type:        types.DocumentTypeAPI,
		Name:        input.Name,
		Description: input.Description,
		Tags:        input.Tags,
		APIContent: types.APIDocumentContent{
			Method:          input.Method,
			Path:            input.Path,
			Header:          input.Header,
			Query:           input.Query,
			PathParams:      input.PathParams,
			ResponseBizCode: input.ResponseBizCode,
			Body:            input.Body,
		},
		CreatedAt: &now,
		UpdatedAt: &now,
	}
	// 保存文档
	_, err = es.ESClient.SaveDocument(document)
	if err != nil {
		return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{Text: fmt.Sprintf("保存API文档失败: %v", err)},
				},
				IsError: true,
			}, struct {
				Success  bool            `json:"success"`
				Document *types.Document `json:"document"`
				Message  string          `json:"message"`
			}{Success: false, Message: fmt.Sprintf("保存失败: %v", err)}, err
	}

	// 构建成功响应
	resultText := fmt.Sprintf("# API文档创建成功\n\n")
	resultText += fmt.Sprintf("**文档名称:** %s\n", document.Name)
	resultText += fmt.Sprintf("**文档ID:** %s\n", document.ID)
	resultText += fmt.Sprintf("**HTTP方法:** %s\n", input.Method)
	resultText += fmt.Sprintf("**API路径:** %s\n", input.Path)
	resultText += fmt.Sprintf("**项目ID:** %s\n", document.ProjectID)
	resultText += fmt.Sprintf("**创建时间:** %s\n", document.CreatedAt)
	if len(document.Tags) > 0 {
		resultText += fmt.Sprintf("**标签:** %v\n", document.Tags)
	}
	resultText += "\n**说明:** API文档已保存，可通过ID或搜索功能访问"

	return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: resultText},
			},
		}, struct {
			Success  bool            `json:"success"`
			Document *types.Document `json:"document"`
			Message  string          `json:"message"`
		}{Success: true, Document: document, Message: "API文档创建成功"}, nil
}

// EditAPIDocument 编辑API文档
func EditAPIDocument(_ context.Context, req *mcp.CallToolRequest, input struct {
	ID              string                 `json:"id" jsonschema:"API文档ID"`
	Name            string                 `json:"name" jsonschema:"文档名称，必填。可使用中文或英文，比如 userLogin、用户登录接口、用户微信登录接口等"`
	Description     string                 `json:"description" jsonschema:"文档详细描述，必填。清晰说明API的功能、业务场景和使用方法"`
	Method          string                 `json:"method" jsonschema:"HTTP方法，必填。标准HTTP方法之一：GET、POST、PUT、DELETE、PATCH、HEAD、OPTIONS、TRACE等"`
	Path            string                 `json:"path" jsonschema:"API相对路径，必填。支持路径参数"`
	Header          map[string]interface{} `json:"header" jsonschema:"API请求头，可选。键值对格式，通常不需要传递该字段，除非用户要求"`
	Body            map[string]interface{} `json:"body" jsonschema:"存在body参数时传入"`
	Query           map[string]interface{} `json:"query" jsonschema:"存在query参数时传入"`
	PathParams      map[string]interface{} `json:"path_params" jsonschema:"存在path参数时传入"`
	ResponseBizCode string                 `json:"response_biz_code" jsonschema:"response 业务码，存在时传入，通常不需要传递该字段，除非用户要求"`
	Tags            []string               `json:"tags" jsonschema:"文档标签，可选。用于分类和搜索，用于方便用户快速定位和筛选文档，需要总结出用户接口的特点并提取关键字"`
}) (*mcp.CallToolResult, struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}, error) {
	doc := map[string]interface{}{
		"api_content": types.APIDocumentContent{
			Method:          input.Method,
			Path:            input.Path,
			Header:          input.Header,
			Query:           input.Query,
			PathParams:      input.PathParams,
			ResponseBizCode: input.ResponseBizCode,
			Body:            input.Body,
		},
		"updated_at": time.Now(),
	}
	if input.Name != "" {
		doc["name"] = input.Name
	}
	if input.Description != "" {
		doc["description"] = input.Description
	}
	if input.Tags != nil {
		doc["tags"] = input.Tags
	}

	// 保存文档
	err := es.ESClient.EditDocument(input.ID, doc)
	if err != nil {
		return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{Text: fmt.Sprintf("保存API文档失败: %v", err)},
				},
				IsError: true,
			}, struct {
				Success bool   `json:"success"`
				Message string `json:"message"`
			}{Success: false, Message: fmt.Sprintf("保存失败: %v", err)}, err
	}

	// 构建成功响应
	resultText := fmt.Sprintf("API文档编辑成功")

	return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: resultText},
			},
		}, struct {
			Success bool   `json:"success"`
			Message string `json:"message"`
		}{Success: true, Message: "API文档编辑成功"}, nil
}

// DelAPIDocument 删除API文档
func DelAPIDocument(_ context.Context, req *mcp.CallToolRequest, input struct {
	ID string `json:"id" jsonschema:"API文档ID"`
}) (*mcp.CallToolResult, struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}, error) {
	_, err := es.ESClient.Client().Delete().
		Index(config.Config.DocumentIndex).
		Id(input.ID).
		Refresh("true").
		Do(context.Background())
	if err != nil {
		return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{Text: fmt.Sprintf("删除API文档失败: %v", err)},
				},
				IsError: true,
			}, struct {
				Success bool   `json:"success"`
				Message string `json:"message"`
			}{Success: false, Message: fmt.Sprintf("删除失败: %v", err)}, err
	}

	return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("API文档删除成功，ID: %s", input.ID)},
			},
		}, struct {
			Success bool   `json:"success"`
			Message string `json:"message"`
		}{Success: true, Message: "API文档删除成功"}, nil
}

// GetAPIDocument 获取API文档
func GetAPIDocument(_ context.Context, req *mcp.CallToolRequest, input struct {
	ID string `json:"id" jsonschema:"API文档ID"`
}) (*mcp.CallToolResult, *types.Document, error) {
	// 获取文档
	document, err := es.ESClient.GetDocument(input.ID)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("获取API文档失败: %v", err)},
			},
			IsError: true,
		}, nil, err
	}

	// 验证文档类型
	if document.Type != types.DocumentTypeAPI {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: "指定的文档不是API文档类型"},
			},
			IsError: true,
		}, nil, fmt.Errorf("文档类型不匹配")
	}

	// 验证权限 - 使用token提取工具函数
	projectId, err := utils.ExtractProjectID(req)
	if err == nil && document.ProjectID != projectId {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("无权访问该API文档，文档属于项目 %s，但当前项目是 %s", document.ProjectID, projectId)},
			},
			IsError: true,
		}, nil, fmt.Errorf("无权访问")
	}

	// 格式化输出
	docJSON, err := json.MarshalIndent(document, "", "  ")
	if err != nil {
		docJSON = []byte(fmt.Sprintf("格式化文档失败: %v", err))
	}

	resultText := fmt.Sprintf("# API文档详情\n\n")
	resultText += fmt.Sprintf("```json\n%s\n```", string(docJSON))

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: resultText},
		},
	}, document, nil
}

// SearchAPIDocuments 搜索API文档
func SearchAPIDocuments(_ context.Context, req *mcp.CallToolRequest, input struct {
	Query  string   `json:"query" jsonschema:"搜索关键词"`
	Method string   `json:"method,omitempty" jsonschema:"HTTP方法筛选"`
	Tags   []string `json:"tags,omitempty" jsonschema:"标签筛选"`
	Limit  int      `json:"limit,omitempty" jsonschema:"返回数量限制，默认20"`
	Offset int      `json:"offset,omitempty" jsonschema:"偏移量，默认0"`
}) (*mcp.CallToolResult, struct {
	Success   bool             `json:"success"`
	Documents []types.Document `json:"documents"`
	Count     int              `json:"count"`
	Message   string           `json:"message"`
}, error) {
	projectId, err := utils.ExtractProjectID(req)
	if err != nil {
		return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{Text: fmt.Sprintf("获取项目ID失败: %v", err)},
				},
				IsError: true,
			}, struct {
				Success   bool             `json:"success"`
				Documents []types.Document `json:"documents"`
				Count     int              `json:"count"`
				Message   string           `json:"message"`
			}{Success: false, Count: 0, Message: fmt.Sprintf("获取项目ID失败: %v", err)}, err
	}

	// 构建搜索请求
	apiType := types.DocumentTypeAPI
	searchReq := &types.SearchRequest{
		Query:  input.Query,
		Limit:  input.Limit,
		Offset: input.Offset,
		Tags:   input.Tags,
		Type:   &apiType,
	}

	// 执行搜索
	result, err := es.ESClient.SearchDocuments(projectId, searchReq)
	if err != nil {
		return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{Text: fmt.Sprintf("搜索API文档失败: %v", err)},
				},
				IsError: true,
			}, struct {
				Success   bool             `json:"success"`
				Documents []types.Document `json:"documents"`
				Count     int              `json:"count"`
				Message   string           `json:"message"`
			}{Success: false, Count: 0, Message: fmt.Sprintf("搜索失败: %v", err)}, err
	}

	// 过滤结果 - 只返回API文档
	var apiDocs []types.Document
	for _, doc := range result.Documents {
		if doc.Type == types.DocumentTypeAPI {
			// 如果指定了方法，进行筛选
			if input.Method != "" && doc.APIContent.Method != input.Method {
				continue
			}
			apiDocs = append(apiDocs, doc)
		}
	}

	// 构建结果文本
	resultText := fmt.Sprintf("# API文档搜索结果: \"%s\"\n\n", input.Query)
	resultText += fmt.Sprintf("找到 **%d** 个匹配的API文档（共 %d 个）\n\n", len(apiDocs), result.Total)

	for i, doc := range apiDocs {
		resultText += fmt.Sprintf("**%d. %s**\n", i+1, doc.Name)
		resultText += fmt.Sprintf("   - ID: `%s`\n", doc.ID)
		resultText += fmt.Sprintf("   - 方法: %s\n", doc.APIContent.Method)
		resultText += fmt.Sprintf("   - 路径: %s\n", doc.APIContent.Path)
		resultText += fmt.Sprintf("   - 描述: %s\n", doc.Description)
		// 请求头
		if doc.APIContent.Header != nil {
			resultText += fmt.Sprintf("   - 请求头: %v\n", doc.APIContent.Header)
		}
		// 查询参数
		if doc.APIContent.Query != nil {
			resultText += fmt.Sprintf("   - 查询参数: %v\n", doc.APIContent.Query)
		}
		// 路径参数
		if doc.APIContent.PathParams != nil {
			resultText += fmt.Sprintf("   - 路径参数: %v\n", doc.APIContent.PathParams)
		}
		// Body参数
		if doc.APIContent.Body != nil {
			resultText += fmt.Sprintf("   - Body参数: %v\n", doc.APIContent.Body)
		}
		// 响应体
		if doc.APIContent.ResponseBizCode != "" {
			resultText += fmt.Sprintf("   - 业务状态码: %v\n", doc.APIContent.ResponseBizCode)
		}
		if len(doc.Tags) > 0 {
			resultText += fmt.Sprintf("   - 标签: %v\n", doc.Tags)
		}
		resultText += fmt.Sprintf("   - 更新时间: %s\n", doc.UpdatedAt)
		resultText += "\n"
	}

	resultText += "---\n\n"
	resultText += "**搜索建议：**\n"
	resultText += "- 使用更具体的关键词获得更精确的结果\n"
	resultText += "- 可以在API路径中使用通配符\n"
	resultText += "- 通过标签分类管理不同类型的API文档"

	return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: resultText},
			},
		}, struct {
			Success   bool             `json:"success"`
			Documents []types.Document `json:"documents"`
			Count     int              `json:"count"`
			Message   string           `json:"message"`
		}{Success: true, Documents: apiDocs, Count: len(apiDocs), Message: "搜索成功"}, nil
}

// ListAPIDocuments 列出API文档
func ListAPIDocuments(_ context.Context, req *mcp.CallToolRequest, input struct {
	Method string `json:"method,omitempty" jsonschema:"HTTP方法筛选"`
	Limit  int    `json:"limit,omitempty" jsonschema:"返回数量限制，默认20"`
	Offset int    `json:"offset,omitempty" jsonschema:"偏移量，默认0"`
}) (*mcp.CallToolResult, struct {
	Success   bool             `json:"success"`
	Documents []types.Document `json:"documents"`
	Count     int              `json:"count"`
	Message   string           `json:"message"`
}, error) {
	// 获取用户ID和项目ID
	userId, err := utils.ExtractUserID(req)
	if err != nil {
		return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{Text: fmt.Sprintf("获取用户ID失败: %v", err)},
				},
				IsError: true,
			}, struct {
				Success   bool             `json:"success"`
				Documents []types.Document `json:"documents"`
				Count     int              `json:"count"`
				Message   string           `json:"message"`
			}{Success: false, Count: 0, Message: fmt.Sprintf("获取用户ID失败: %v", err)}, err
	}

	projectId, err := utils.ExtractProjectID(req)
	if err != nil {
		return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{Text: fmt.Sprintf("获取项目ID失败: %v", err)},
				},
				IsError: true,
			}, struct {
				Success   bool             `json:"success"`
				Documents []types.Document `json:"documents"`
				Count     int              `json:"count"`
				Message   string           `json:"message"`
			}{Success: false, Count: 0, Message: fmt.Sprintf("获取项目ID失败: %v", err)}, err
	}

	fmt.Printf("🔐 用户 %s 正在项目 %s 中列出API文档\n", userId, projectId)

	// 设置默认值
	limit := input.Limit
	if limit <= 0 {
		limit = 20
	}
	offset := input.Offset
	if offset < 0 {
		offset = 0
	}

	// 执行列表
	apiType := types.DocumentTypeAPI
	docType := &apiType
	result, err := es.ESClient.ListDocuments(projectId, docType, limit, offset)
	if err != nil {
		return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{Text: fmt.Sprintf("列出API文档失败: %v", err)},
				},
				IsError: true,
			}, struct {
				Success   bool             `json:"success"`
				Documents []types.Document `json:"documents"`
				Count     int              `json:"count"`
				Message   string           `json:"message"`
			}{Success: false, Count: 0, Message: fmt.Sprintf("列出失败: %v", err)}, err
	}

	// 如果指定了方法筛选
	if input.Method != "" {
		var filteredDocs []types.Document
		for _, doc := range result.Documents {
			if doc.Type == types.DocumentTypeAPI && doc.APIContent.Method == input.Method {
				filteredDocs = append(filteredDocs, doc)
			}
		}
		result.Documents = filteredDocs
		result.Total = int64(len(filteredDocs))
	}

	// 构建结果文本
	methodFilter := "所有方法"
	if input.Method != "" {
		methodFilter = input.Method
	}

	resultText := fmt.Sprintf("# API文档列表 (%s)\n\n", methodFilter)
	resultText += fmt.Sprintf("显示 **%d** 个API文档（共 %d 个）\n\n", len(result.Documents), result.Total)

	for i, doc := range result.Documents {
		resultText += fmt.Sprintf("**%d. %s**\n", i+1, doc.Name)
		resultText += fmt.Sprintf("   - ID: `%s`\n", doc.ID)
		resultText += fmt.Sprintf("   - 方法: %s\n", doc.APIContent.Method)
		resultText += fmt.Sprintf("   - 路径: %s\n", doc.APIContent.Path)
		resultText += fmt.Sprintf("   - 描述: %s\n", doc.Description)
		if len(doc.Tags) > 0 {
			resultText += fmt.Sprintf("   - 标签: %v\n", doc.Tags)
		}
		resultText += fmt.Sprintf("   - 更新时间: %s\n", doc.UpdatedAt)
		resultText += "\n"
	}

	resultText += "---\n\n"
	resultText += "**API文档管理提示：**\n"
	resultText += "- 使用 `search_api_documents` 进行详细搜索\n"
	resultText += "- 使用 `get_api_document` 获取完整文档结构\n"
	resultText += "- 通过标签分类管理不同类型的API文档"

	return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: resultText},
			},
		}, struct {
			Success   bool             `json:"success"`
			Documents []types.Document `json:"documents"`
			Count     int              `json:"count"`
			Message   string           `json:"message"`
		}{Success: true, Documents: result.Documents, Count: len(result.Documents), Message: "列出成功"}, nil
}

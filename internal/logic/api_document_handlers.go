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
	ProjectID       string                 `json:"project_id,omitempty" jsonschema:"项目ID，可选。要列出API文档的项目ID。必须是有效的项目标识符，不传将查询默认项目，可从"list_projects"或"get_project_info"工具中获取。"`
	Name            string                 `json:"name" jsonschema:"文档名称，必填。用于标识和搜索API文档。建议使用有意义的名称，如：'userLogin'、'getUserInfo'、'createOrder'等。支持中英文，长度2-50字符。"`
	Description     string                 `json:"description" jsonschema:"文档描述，必填。详细说明API的功能、业务场景和使用方法。建议包括：API用途、输入输出说明、使用示例、注意事项等。"`
	Method          string                 `json:"method" jsonschema:"HTTP方法，必填。API的请求方法。必须是标准HTTP方法：GET(查询)、POST(创建)、PUT(更新)、DELETE(删除)、PATCH(部分更新)、HEAD、OPTIONS、TRACE。请根据API的实际功能选择。"`
	Path            string                 `json:"path" jsonschema:"API路径，必填。API的访问路径，如：'/api/users/{id}'、'/api/orders'。支持路径参数，使用'{paramName}'格式。必须以'/'开头，避免使用特殊字符。"`
	Header          map[string]interface{} `json:"header" jsonschema:"请求头，可选。API请求所需的HTTP头部信息，如：{'Authorization': 'Bearer token', 'Content-Type': 'application/json'}。通常不需要传递，除非有特殊认证要求。"`
	Body            map[string]interface{} `json:"body" jsonschema:"请求体，可选。POST/PUT请求的JSON数据，如：{'name': '张三', 'age': 25}。GET/DELETE请求通常不需要。请提供完整的请求体结构和字段说明。"`
	Query           map[string]interface{} `json:"query" jsonschema:"查询参数，可选。URL查询参数，如：{'page': 1, 'size': 10}。用于GET请求的参数传递，格式为键值对。"`
	PathParams      map[string]interface{} `json:"path_params" jsonschema:"路径参数，可选。URL路径中的参数，如：{'id': 123}。对应路径中的'{paramName}'占位符，用于动态路径构建。"`
	ResponseBizCode string                 `json:"response_biz_code" jsonschema:"业务状态码，可选。自定义的业务响应码，如：'SUCCESS'、'USER_NOT_FOUND'。用于业务逻辑判断，通常与HTTP状态码配合使用。"`
	Tags            []string               `json:"tags" jsonschema:"文档标签，可选。用于分类和搜索API文档。建议使用：['user', 'auth', 'v1']等有意义的标签。支持多个标签，便于文档管理和筛选。"`
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

	projectID := input.ProjectID
	if projectID == "" {
		projectId, err := utils.ExtractProjectID(req)
		if err == nil {
			projectID = projectId
		}
	}

	is, err := es.ESClient.CheckUserInProject(projectID, userId)
	if err != nil {
		return nil, struct {
			Success  bool            `json:"success"`
			Document *types.Document `json:"document"`
			Message  string          `json:"message"`
		}{Success: false, Message: fmt.Sprintf("系统错误: %v", err)}, err
	}
	if !is {
		return nil, struct {
			Success  bool            `json:"success"`
			Document *types.Document `json:"document"`
			Message  string          `json:"message"`
		}{Success: false, Message: fmt.Sprintf("用户不在项目中: %v", err)}, err
	}

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
		ProjectID:   projectID,
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
	ID              string                 `json:"id" jsonschema:"文档ID，必填。要编辑的API文档的唯一标识符。必须是已存在的文档ID，可通过搜索或列表功能获取。"`
	Name            string                 `json:"name" jsonschema:"文档名称，可选。更新后的文档名称，如：'userLogin'、'getUserInfo'。支持中英文，长度2-50字符。不传则保持原名称不变。"`
	Description     string                 `json:"description" jsonschema:"文档描述，可选。更新后的文档详细说明，包括API功能、业务场景和使用方法。建议包含具体的使用示例和注意事项。不传则保持原描述不变。"`
	Method          string                 `json:"method" jsonschema:"HTTP方法，可选。更新后的API请求方法。必须是标准HTTP方法：GET(查询)、POST(创建)、PUT(更新)、DELETE(删除)、PATCH(部分更新)等。不传则保持原方法不变。"`
	Path            string                 `json:"path" jsonschema:"API路径，可选。更新后的API访问路径，如：'/api/users/{id}'。支持路径参数，必须以'/'开头。不传则保持原路径不变。"`
	Header          map[string]interface{} `json:"header" jsonschema:"请求头，可选。更新后的API请求头部信息，如：{'Authorization': 'Bearer token'}。通常不需要修改，除非有特殊认证要求变更。"`
	Body            map[string]interface{} `json:"body" jsonschema:"请求体，可选。更新后的POST/PUT请求JSON数据，如：{'name': '张三', 'age': 25}。请提供完整的请求体结构。不传则保持原请求体不变。"`
	Query           map[string]interface{} `json:"query" jsonschema:"查询参数，可选。更新后的URL查询参数，如：{'page': 1, 'size': 10}。用于GET请求的参数传递。不传则保持原查询参数不变。"`
	PathParams      map[string]interface{} `json:"path_params" jsonschema:"路径参数，可选。更新后的URL路径参数，如：{'id': 123}。对应路径中的'{paramName}'占位符。不传则保持原路径参数不变。"`
	ResponseBizCode string                 `json:"response_biz_code" jsonschema:"业务状态码，可选。更新后的自定义业务响应码，如：'SUCCESS'、'USER_NOT_FOUND'。用于业务逻辑判断。不传则保持原业务码不变。"`
	Tags            []string               `json:"tags" jsonschema:"文档标签，可选。更新后的文档标签数组，如：['user', 'auth', 'v1']。用于分类和搜索。不传则保持原标签不变。"`
}) (*mcp.CallToolResult, struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}, error) {
	if input.ID == "" {
		return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{Text: "文档ID不能为空"},
				},
				IsError: true,
			}, struct {
				Success bool   `json:"success"`
				Message string `json:"message"`
			}{Success: false, Message: "文档ID不能为空"}, fmt.Errorf("文档ID不能为空")
	}

	userId, err := utils.ExtractUserID(req)
	if err != nil {
		return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{Text: fmt.Sprintf("获取用户ID失败: %v", err)},
				},
				IsError: true,
			}, struct {
				Success bool   `json:"success"`
				Message string `json:"message"`
			}{Success: false, Message: fmt.Sprintf("获取用户ID失败: %v", err)}, err
	}

	projectID, err := es.ESClient.GetDocumentProjectID(input.ID)
	if err != nil {
		return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{Text: fmt.Sprintf("系统错误: %v", err)},
				},
				IsError: true,
			}, struct {
				Success bool   `json:"success"`
				Message string `json:"message"`
			}{Success: false, Message: fmt.Sprintf("系统错误: %v", err)}, err
	}

	is, err := es.ESClient.CheckUserInProject(projectID, userId)
	if err != nil {
		return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{Text: fmt.Sprintf("系统错误: %v", err)},
				},
				IsError: true,
			}, struct {
				Success bool   `json:"success"`
				Message string `json:"message"`
			}{Success: false, Message: fmt.Sprintf("系统错误: %v", err)}, err
	}
	if !is {
		return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{Text: fmt.Sprintf("用户不在项目中: %v", userId)},
				},
				IsError: true,
			}, struct {
				Success bool   `json:"success"`
				Message string `json:"message"`
			}{Success: false, Message: fmt.Sprintf("用户不在项目中: %v", userId)}, err
	}

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

	err = es.ESClient.EditDocument(input.ID, doc)
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
	ID string `json:"id" jsonschema:"文档ID，必填。要删除的API文档的唯一标识符。必须是已存在的文档ID，删除操作不可恢复，请谨慎操作。建议在删除前先备份重要文档。"`
}) (*mcp.CallToolResult, struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}, error) {
	userId, err := utils.ExtractUserID(req)
	if err != nil {
		return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{Text: fmt.Sprintf("获取用户ID失败: %v", err)},
				},
				IsError: true,
			}, struct {
				Success bool   `json:"success"`
				Message string `json:"message"`
			}{Success: false, Message: fmt.Sprintf("获取用户ID失败: %v", err)}, err
	}

	projectID, err := es.ESClient.GetDocumentProjectID(input.ID)
	if err != nil {
		return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{Text: fmt.Sprintf("系统错误: %v", err)},
				},
				IsError: true,
			}, struct {
				Success bool   `json:"success"`
				Message string `json:"message"`
			}{Success: false, Message: fmt.Sprintf("系统错误: %v", err)}, err
	}

	is, err := es.ESClient.CheckUserInProject(projectID, userId)
	if err != nil {
		return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{Text: fmt.Sprintf("系统错误: %v", err)},
				},
				IsError: true,
			}, struct {
				Success bool   `json:"success"`
				Message string `json:"message"`
			}{Success: false, Message: fmt.Sprintf("系统错误: %v", err)}, err
	}
	if !is {
		return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{Text: fmt.Sprintf("用户不在项目中: %v", userId)},
				},
				IsError: true,
			}, struct {
				Success bool   `json:"success"`
				Message string `json:"message"`
			}{Success: false, Message: fmt.Sprintf("用户不在项目中: %v", userId)}, err
	}

	_, err = es.ESClient.Client().Delete().
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
	ID string `json:"id" jsonschema:"文档ID，必填。要获取的API文档的唯一标识符。必须是已存在的文档ID，可通过搜索或列表功能获取。返回完整的文档结构和内容。"`
}) (*mcp.CallToolResult, *types.Document, error) {
	userId, err := utils.ExtractUserID(req)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("获取用户ID失败: %v", err)},
			},
			IsError: true,
		}, nil, err
	}

	projectID, err := es.ESClient.GetDocumentProjectID(input.ID)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("获取文档项目ID失败: %v", err)},
			},
			IsError: true,
		}, nil, err
	}

	is, err := es.ESClient.CheckUserInProject(projectID, userId)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("检查用户在项目中失败: %v", err)},
			},
			IsError: true,
		}, nil, err
	}
	if !is {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("用户不在项目中: %v", userId)},
			},
			IsError: true,
		}, nil, err
	}

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
	ProjectID string   `json:"project_id,omitempty" jsonschema:"项目ID，可选。要列出API文档的项目ID。必须是有效的项目标识符，不传将查询默认项目，可从"list_projects"或"get_project_info"工具中获取。"`
	Query     string   `json:"query" jsonschema:"搜索关键词，必填。用于搜索API文档的关键词，可以是文档名称、路径、描述中的任意词汇。如：'user'、'login'、'/api/users'。支持模糊搜索。"`
	Method    string   `json:"method,omitempty" jsonschema:"HTTP方法筛选，可选。按HTTP方法过滤结果，如：'GET'、'POST'等。必须是标准HTTP方法。不传则返回所有方法的文档。"`
	Tags      []string `json:"tags,omitempty" jsonschema:"标签筛选，可选。按标签过滤文档，如：['user', 'auth', 'v1']。支持多个标签组合筛选。不传则忽略标签过滤。"`
	Limit     int      `json:"limit,omitempty" jsonschema:"返回数量限制，可选。控制返回结果数量，默认20，最大100。用于分页浏览。"`
	Offset    int      `json:"offset,omitempty" jsonschema:"偏移量，可选。分页查询的起始位置，默认0。用于跳过前面的结果，获取后续数据。"`
}) (*mcp.CallToolResult, struct {
	Success   bool             `json:"success"`
	Documents []types.Document `json:"documents"`
	Count     int              `json:"count"`
	Message   string           `json:"message"`
}, error) {
	userId, err := utils.ExtractUserID(req)
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
			}{Success: false, Count: 0, Documents: make([]types.Document, 0), Message: fmt.Sprintf("搜索API文档失败: %v", err)}, err
	}

	projectID := input.ProjectID
	if projectID == "" {
		projectId, err := utils.ExtractProjectID(req)
		if err == nil {
			projectID = projectId
		}
	}

	is, err := es.ESClient.CheckUserInProject(projectID, userId)
	if err != nil {
		return nil, struct {
			Success   bool             `json:"success"`
			Documents []types.Document `json:"documents"`
			Count     int              `json:"count"`
			Message   string           `json:"message"`
		}{Success: false, Count: 0, Documents: make([]types.Document, 0), Message: fmt.Sprintf("检查用户在项目中失败: %v", err)}, err
	}
	if !is {
		return nil, struct {
			Success   bool             `json:"success"`
			Documents []types.Document `json:"documents"`
			Count     int              `json:"count"`
			Message   string           `json:"message"`
		}{Success: false, Count: 0, Documents: make([]types.Document, 0), Message: fmt.Sprintf("用户不在项目中: %v", userId)}, err
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
	result, err := es.ESClient.SearchDocuments(projectID, searchReq)
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
			}{Success: false, Count: 0, Documents: make([]types.Document, 0), Message: fmt.Sprintf("搜索失败: %v", err)}, err
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

	if len(apiDocs) == 0 {
		// 没有找到文档时的友好提示
		resultText += fmt.Sprintf("🔍 **未找到匹配的API文档**\n\n")
		resultText += fmt.Sprintf("搜索关键词 \"%s\" 没有找到任何匹配的API文档。\n\n", input.Query)
		resultText += "**建议：**\n"
		resultText += "- 尝试使用更通用的关键词，如 'user'、'api'、'system'\n"
		resultText += "- 检查关键词拼写是否正确\n"
		resultText += "- 使用不同的HTTP方法筛选（如 'GET'、'POST'）\n"
		resultText += "- 尝试使用标签筛选，如 'auth'、'user'、'admin'\n"
		resultText += "- 如果这是新项目，可以先使用 `create_api_document` 创建第一个API文档\n\n"
		resultText += "**可用操作：**\n"
		resultText += "- 使用 `list_api_documents` 查看所有API文档\n"
		resultText += "- 使用 `create_api_document` 创建新的API文档\n"
		resultText += "- 使用 `list_projects` 查看可用项目"
	} else {
		// 找到文档时的正常显示
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
	ProjectID string `json:"project_id,omitempty" jsonschema:"项目ID，可选。要列出API文档的项目ID。必须是有效的项目标识符，不传将查询默认项目，可从"list_projects"或"get_project_info"工具中获取。"`
	Query     string `json:"query,omitempty" jsonschema:"搜索关键词，可选。用于在列表中进一步搜索的过滤关键词，如：'user'、'login'。支持模糊搜索。不传则返回所有文档。"`
	Method    string `json:"method,omitempty" jsonschema:"HTTP方法筛选，可选。按HTTP方法过滤结果，如：'GET'、'POST'等。必须是标准HTTP方法。不传则返回所有方法的文档。"`
	Limit     int    `json:"limit,omitempty" jsonschema:"返回数量限制，可选。控制返回结果数量，默认20，最大100。用于分页浏览。"`
	Offset    int    `json:"offset,omitempty" jsonschema:"偏移量，可选。分页查询的起始位置，默认0。用于跳过前面的结果，获取后续数据。"`
}) (*mcp.CallToolResult, struct {
	Success   bool             `json:"success"`
	Documents []types.Document `json:"documents,omitempty"`
	Count     int              `json:"count"`
	Message   string           `json:"message"`
}, error) {
	userId, err := utils.ExtractUserID(req)
	if err != nil {
		return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{Text: fmt.Sprintf("搜索API文档失败: %v", err)},
				},
				IsError: true,
			}, struct {
				Success   bool             `json:"success"`
				Documents []types.Document `json:"documents,omitempty"`
				Count     int              `json:"count"`
				Message   string           `json:"message"`
			}{Success: false, Count: 0, Message: fmt.Sprintf("搜索API文档失败: %v", err)}, err
	}

	projectID := input.ProjectID
	if projectID == "" {
		projectId, err := utils.ExtractProjectID(req)
		if err == nil {
			projectID = projectId
		}
	}

	is, err := es.ESClient.CheckUserInProject(projectID, userId)
	if err != nil {
		return nil, struct {
			Success   bool             `json:"success"`
			Documents []types.Document `json:"documents,omitempty"`
			Count     int              `json:"count"`
			Message   string           `json:"message"`
		}{Success: false, Count: 0, Message: fmt.Sprintf("检查用户在项目中失败: %v", err)}, err
	}
	if !is {
		return nil, struct {
			Success   bool             `json:"success"`
			Documents []types.Document `json:"documents,omitempty"`
			Count     int              `json:"count"`
			Message   string           `json:"message"`
		}{Success: false, Count: 0, Message: fmt.Sprintf("用户不在项目中: %v", userId)}, err
	}

	fmt.Printf("🔐 用户 %s 正在项目 %s 中列出API文档\n", userId, projectID)

	// 设置默认值
	limit := input.Limit
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	offset := input.Offset
	if offset < 0 {
		offset = 0
	}

	// 执行列表
	apiType := types.DocumentTypeAPI
	docType := &apiType
	result, err := es.ESClient.ListDocuments(input.ProjectID, docType, limit, offset)
	if err != nil {
		return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{Text: fmt.Sprintf("列出API文档失败: %v", err)},
				},
				IsError: true,
			}, struct {
				Success   bool             `json:"success"`
				Documents []types.Document `json:"documents,omitempty"`
				Count     int              `json:"count"`
				Message   string           `json:"message"`
			}{Success: false, Count: 0, Message: fmt.Sprintf("列出API文档失败: %v", err)}, err
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

	if len(result.Documents) == 0 {
		// 没有找到文档时的友好提示
		resultText += fmt.Sprintf("📝 **暂无API文档**\n\n")
		resultText += fmt.Sprintf("当前项目中还没有任何API文档。\n\n")
		resultText += "**建议：**\n"
		resultText += "- 使用 `create_api_document` 创建第一个API文档\n"
		resultText += "- 可以先定义常见的API，如用户登录、数据查询等\n"
		resultText += "- 为API文档添加清晰的描述和标签便于管理\n"
		resultText += "- 使用 `list_text_documents` 查看是否有其他类型的文档\n\n"
		resultText += "**创建API文档示例：**\n"
		resultText += "```bash\n"
		resultText += "# 创建用户登录API文档\n"
		resultText += "create_api_document(\n"
		resultText += "  name='userLogin',\n"
		resultText += "  description='用户登录接口',\n"
		resultText += "  method='POST',\n"
		resultText += "  path='/api/auth/login',\n"
		resultText += "  body={'username': 'string', 'password': 'string'},\n"
		resultText += "  tags=['auth', 'user']\n"
		resultText += ")\n```"
	} else {
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
	}

	resultText += "---\n\n"
	resultText += "**API文档管理提示：**\n"
	resultText += "- 使用 `search_api_documents` 进行详细搜索\n"
	resultText += "- 使用 `get_api_document` 获取完整文档结构\n"
	resultText += "- 通过标签分类管理不同类型的API文档"
	if len(result.Documents) == 0 {
		resultText += "\n- 现在就创建第一个API文档开始管理你的API吧！"
	}

	return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: resultText},
			},
		}, struct {
			Success   bool             `json:"success"`
			Documents []types.Document `json:"documents,omitempty"`
			Count     int              `json:"count"`
			Message   string           `json:"message"`
		}{Success: true, Documents: result.Documents, Count: len(result.Documents), Message: "列出API文档成功"}, nil
}

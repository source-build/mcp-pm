package logic

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/source-build/mcp-pm/internal/es"
	"github.com/source-build/mcp-pm/internal/types"
	"github.com/source-build/mcp-pm/internal/utils"
)

/*
文本文档管理处理器

注意：这些工具专门用于管理文本文档，包括：
- 技术文档和说明文档
- 提示词模板和脚本
- 配置文件和说明
- 规范文档和设计文档

适用场景：
- 项目文档的集中管理
- 知看和管理技术文档
- 维护API的使用指南
- 保存团队知识库和最佳实践

数据结构说明：
- DocumentType: "text" 标识这是文本文档
- content.Text: 包含content_type和content字段等文本特定字段
- 标签：用于分类不同类型或状态的文档（如README, README, note, config等）
*/

// CreateTextDocument 创建文本文档
func CreateTextDocument(_ context.Context, req *mcp.CallToolRequest, input struct {
	ProjectID   string            `json:"project_id,omitempty" jsonschema:"项目ID，可选。要列出文本文档的项目ID。必须是有效的项目标识符，不传将查询默认项目，可从"list_projects"或"get_project_info"工具中获取。"`
	Name        string            `json:"name" jsonschema:"文档名称，必填。文本文档的显示名称，用于搜索和识别。建议使用描述性名称，如：'project-readme'、'server-config'、'api-guide'等。支持中英文，长度2-100字符。"`
	Description string            `json:"description" jsonschema:"文档描述，必填。详细说明文本文档的用途、内容和适用场景，帮助AI和团队成员理解文档的价值。建议包括：文档目的、使用方法、维护说明、更新频率等。"`
	ContentType string            `json:"content_type" jsonschema:"内容类型，必填。文本文档的具体类型，支持的值：'readme'(README文档)、'prompt'(提示词模板)、'config'(配置文件)、'note'(普通笔记)、'spec'(技术规范)。必须为有效枚举值。"`
	Variables   map[string]string `json:"variables" jsonschema:"模板变量，可选。用于提示词模板的变量键值对，如：{'username': '用户名', 'api_key': 'API密钥'}。仅在content_type为'prompt'时有效，用于模板参数替换。"`
	Content     interface{}       `json:"content" jsonschema:"文档内容，必填。文本文档的实际内容，可以是字符串、JSON对象或任意可转换为字符串的内容。根据content_type自动格式化存储。如：README文档使用Markdown格式，配置文件使用JSON格式，提示词模板使用纯文本。"`
	Tags        []string          `json:"tags,omitempty" jsonschema:"文档标签，可选。用于分类和搜索文本文档。建议使用：['readme', 'config', 'guide', 'template', 'docs']等有意义的标签。支持多个标签，便于文档管理和团队协作。"`
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

	fmt.Printf("🔐 用户 %s 正在项目 %s 中创建文本文档: %s\n", userId, projectID, input.Name)

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

	// 创建文本文档对象
	document := &types.Document{
		ID:          hex.EncodeToString(m.Sum(nil)),
		ProjectID:   projectID,
		Type:        types.DocumentTypeText,
		Name:        input.Name,
		Description: input.Description,
		Tags:        input.Tags,
		TextContent: types.TextDocumentContent{
			ContentType: input.ContentType,
			Content:     input.Content,
			Variables:   input.Variables,
		},
	}

	// 保存文档
	_, err = es.ESClient.SaveDocument(document)
	if err != nil {
		return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{Text: fmt.Sprintf("保存文本文档失败: %v", err)},
				},
				IsError: true,
			}, struct {
				Success  bool            `json:"success"`
				Document *types.Document `json:"document"`
				Message  string          `json:"message"`
			}{Success: false, Message: fmt.Sprintf("保存失败: %v", err)}, err
	}

	// 构建成功响应
	resultText := fmt.Sprintf("# 文本文档创建成功\n\n")
	resultText += fmt.Sprintf("**文档名称:** %s\n", document.Name)
	resultText += fmt.Sprintf("**文档ID:** %s\n", document.ID)
	resultText += fmt.Sprintf("**内容类型:** %s\n", document.TextContent.ContentType)
	resultText += fmt.Sprintf("**项目ID:** %s\n", document.ProjectID)
	resultText += fmt.Sprintf("**创建时间:** %s\n", document.CreatedAt)
	if len(document.Tags) > 0 {
		resultText += fmt.Sprintf("**标签:** %v\n", document.Tags)
	}
	resultText += "\n**内容预览:**\n"

	// 根据内容类型显示预览
	contentStr := fmt.Sprintf("%v", document.TextContent.Content)
	if len(contentStr) > 200 {
		contentStr = contentStr[:200] + "..."
	}
	resultText += fmt.Sprintf("```\n%s\n```", contentStr)
	resultText += "\n**说明:** 文本文档已保存，可通过ID或搜索功能访问"

	return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: resultText},
			},
		}, struct {
			Success  bool            `json:"success"`
			Document *types.Document `json:"document"`
			Message  string          `json:"message"`
		}{Success: true, Document: document, Message: "文本文档创建成功"}, nil
}

// GetTextDocument 获取文本文档
func GetTextDocument(_ context.Context, req *mcp.CallToolRequest, input struct {
	ID string `json:"id" jsonschema:"文档ID，必填。要获取的文本文档的唯一标识符。必须是已存在的文档ID，可通过搜索或列表功能获取。返回完整的文档结构和内容。"`
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
				&mcp.TextContent{Text: fmt.Sprintf("获取文本文档失败: %v", err)},
			},
			IsError: true,
		}, nil, err
	}

	// 验证文档类型
	if document.Type != types.DocumentTypeText {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: "指定的文档不是文本文档类型"},
			},
			IsError: true,
		}, nil, fmt.Errorf("文档类型不匹配")
	}

	// 格式化输出
	docJSON, err := json.MarshalIndent(document, "", "  ")
	if err != nil {
		docJSON = []byte(fmt.Sprintf("格式化文档失败: %v", err))
	}

	resultText := fmt.Sprintf("# 文本文档详情\n\n")
	resultText += fmt.Sprintf("```json\n%s\n```", string(docJSON))

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: resultText},
		},
	}, document, nil
}

// SearchTextDocuments 搜索文本文档
func SearchTextDocuments(_ context.Context, req *mcp.CallToolRequest, input struct {
	ProjectID   string   `json:"project_id,omitempty" jsonschema:"项目ID，可选。要列出文本文档的项目ID。必须是有效的项目标识符，不传将查询默认项目，可从"list_projects"或"get_project_info"工具中获取。"`
	Query       string   `json:"query" jsonschema:"搜索关键词，必填。用于搜索文本文档的关键词，可以是文档名称、描述或内容中的任意词汇。如：'readme'、'config'、'guide'、'tutorial'。支持模糊搜索和内容全文检索。"`
	ContentType string   `json:"content_type,omitempty" jsonschema:"内容类型筛选，可选。按内容类型过滤结果，支持的值：'readme'、'prompt'、'config'、'note'、'spec'。不传则返回所有类型的文档。"`
	Tags        []string `json:"tags,omitempty" jsonschema:"标签筛选，可选。按标签过滤文档，如：['guide', 'tutorial', 'docs']。支持多个标签组合筛选。不传则忽略标签过滤。"`
	Limit       int      `json:"limit,omitempty" jsonschema:"返回数量限制，可选。控制返回结果数量，默认20，最大100。用于分页浏览和性能优化。"`
	Offset      int      `json:"offset,omitempty" jsonschema:"偏移量，可选。分页查询的起始位置，默认0。用于跳过前面的结果，获取后续数据。支持大量文档的分页浏览。"`
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
		}{Success: false, Count: 0, Documents: make([]types.Document, 0), Message: fmt.Sprintf("检查用户在项目中失败: %v", err)}, err
	}
	if !is {
		return nil, struct {
			Success   bool             `json:"success"`
			Documents []types.Document `json:"documents,omitempty"`
			Count     int              `json:"count"`
			Message   string           `json:"message"`
		}{Success: false, Count: 0, Message: fmt.Sprintf("用户不在项目中: %v", userId)}, err
	}

	// 构建搜索请求
	textType := types.DocumentTypeText
	searchReq := &types.SearchRequest{
		Query:  input.Query,
		Limit:  input.Limit,
		Offset: input.Offset,
		Tags:   input.Tags,
		Type:   &textType,
	}

	// 执行搜索
	result, err := es.ESClient.SearchDocuments(projectID, searchReq)
	if err != nil {
		return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{Text: fmt.Sprintf("搜索文本文档失败: %v", err)},
				},
				IsError: true,
			}, struct {
				Success   bool             `json:"success"`
				Documents []types.Document `json:"documents,omitempty"`
				Count     int              `json:"count"`
				Message   string           `json:"message"`
			}{Success: false, Count: 0, Message: fmt.Sprintf("搜索失败: %v", err)}, err
	}

	// 过滤结果 - 只返回文本文档
	var textDocs []types.Document
	for _, doc := range result.Documents {
		if doc.Type == types.DocumentTypeText {
			// 如果指定了内容类型，进行筛选
			if input.ContentType != "" && doc.TextContent.ContentType != input.ContentType {
				continue
			}
			textDocs = append(textDocs, doc)
		}
	}

	// 构建结果文本
	resultText := fmt.Sprintf("# 文本文档搜索结果: \"%s\"\n\n", input.Query)

	if len(textDocs) == 0 {
		// 没有找到文档时的友好提示
		resultText += fmt.Sprintf("🔍 **未找到匹配的文本文档**\n\n")
		resultText += fmt.Sprintf("搜索关键词 \"%s\" 没有找到任何匹配的文本文档。\n\n", input.Query)
		resultText += "**建议：**\n"
		resultText += "- 尝试使用更通用的关键词，如 'readme'、'config'、'guide'\n"
		resultText += "- 检查关键词拼写是否正确\n"
		resultText += "- 尝试使用不同的内容类型筛选（如 'readme'、'prompt'）\n"
		resultText += "- 尝试使用标签筛选，如 'docs'、'guide'、'tutorial'\n"
		resultText += "- 如果这是新项目，可以先使用 `create_text_document` 创建第一个文本文档\n\n"
		resultText += "**可用操作：**\n"
		resultText += "- 使用 `list_text_documents` 查看所有文本文档\n"
		resultText += "- 使用 `create_text_document` 创建新的文本文档\n"
		resultText += "- 使用 `list_projects` 查看可用项目"
	} else {
		// 找到文档时的正常显示
		resultText += fmt.Sprintf("找到 **%d** 个匹配的文本文档（共 %d 个）\n\n", len(textDocs), result.Total)

		for i, doc := range textDocs {
			resultText += fmt.Sprintf("**%d. %s**\n", i+1, doc.Name)
			resultText += fmt.Sprintf("   - ID: `%s`\n", doc.ID)
			resultText += fmt.Sprintf("   - 类型: %s\n", doc.TextContent.ContentType)
			resultText += fmt.Sprintf("   - 描述: %s\n", doc.Description)
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
	resultText += "- 可以搜索文档内容中的任意文本\n"
	resultText += "- 通过标签分类管理不同类型的文本文档\n"
	resultText += "- 重点注意，首次请把content字段内容分析总结并展示给用户，如果用户不需要，那么后续可以不用继续展示该部分"

	return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: resultText},
			},
		}, struct {
			Success   bool             `json:"success"`
			Documents []types.Document `json:"documents,omitempty"`
			Count     int              `json:"count"`
			Message   string           `json:"message"`
		}{Success: true, Documents: textDocs, Count: len(textDocs), Message: "搜索成功"}, nil
}

// ListTextDocuments 列出文本文档
func ListTextDocuments(_ context.Context, req *mcp.CallToolRequest, input struct {
	ProjectID   string `json:"project_id,omitempty" jsonschema:"项目ID，可选。要列出文本文档的项目ID。必须是有效的项目标识符，不传将查询默认项目，可从"list_projects"或"get_project_info"工具中获取。"`
	ContentType string `json:"content_type,omitempty" jsonschema:"内容类型筛选，可选。按内容类型过滤结果，支持的值：'readme'、'prompt'、'config'、'note'、'spec'。不传则返回所有类型的文档。用于快速定位特定类型的文档。"`
	Limit       int    `json:"limit,omitempty" jsonschema:"返回数量限制，可选。控制返回结果数量，默认20，最大100。用于分页浏览和性能优化，避免一次性返回过多数据。"`
	Offset      int    `json:"offset,omitempty" jsonschema:"偏移量，可选。分页查询的起始位置，默认0。用于跳过前面的结果，获取后续数据。支持大量文档的分页浏览。"`
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

	fmt.Printf("🔐 用户 %s 正在项目 %s 中列出文本文档\n", userId, projectID)

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
	textType := types.DocumentTypeText
	docType := &textType
	result, err := es.ESClient.ListDocuments(projectID, docType, limit, offset)
	if err != nil {
		return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{Text: fmt.Sprintf("列出文本文档失败: %v", err)},
				},
				IsError: true,
			}, struct {
				Success   bool             `json:"success"`
				Documents []types.Document `json:"documents,omitempty"`
				Count     int              `json:"count"`
				Message   string           `json:"message"`
			}{Success: false, Count: 0, Message: fmt.Sprintf("列出失败: %v", err)}, err
	}

	// 如果指定了内容类型筛选
	if input.ContentType != "" {
		var filteredDocs []types.Document
		for _, doc := range result.Documents {
			if doc.Type == types.DocumentTypeText && doc.TextContent.ContentType == input.ContentType {
				filteredDocs = append(filteredDocs, doc)
			}
		}
		result.Documents = filteredDocs
		result.Total = int64(len(filteredDocs))
	}

	// 构建结果文本
	typeFilter := "所有类型"
	if input.ContentType != "" {
		typeFilter = input.ContentType
	}

	resultText := fmt.Sprintf("# 文本文档列表 (%s)\n\n", typeFilter)

	if len(result.Documents) == 0 {
		// 没有找到文档时的友好提示
		resultText += fmt.Sprintf("📝 **暂无文本文档**\n\n")
		resultText += fmt.Sprintf("当前项目中还没有任何文本文档。\n\n")
		resultText += "**建议：**\n"
		resultText += "- 使用 `create_text_document` 创建第一个文本文档\n"
		resultText += "- 可以先创建项目README文档，介绍项目概况\n"
		resultText += "- 创建配置文件模板，便于团队使用\n"
		resultText += "- 添加提示词模板，提高AI助手效率\n"
		resultText += "- 使用 `list_api_documents` 查看是否有API文档\n\n"
		resultText += "**创建文本文档示例：**\n"
		resultText += "```bash\n"
		resultText += "# 创建项目README文档\n"
		resultText += "create_text_document(\n"
		resultText += "  name='project-readme',\n"
		resultText += "  description='项目介绍和使用指南',\n"
		resultText += "  content_type='readme',\n"
		resultText += "  content='# 项目标题\\n\\n项目简介...',\n"
		resultText += "  tags=['readme', 'guide']\n"
		resultText += ")\n```"
	} else {
		// 找到文档时的正常显示
		resultText += fmt.Sprintf("显示 **%d** 个文本文档（共 %d 个）\n\n", len(result.Documents), result.Total)

		for i, doc := range result.Documents {
			resultText += fmt.Sprintf("**%d. %s**\n", i+1, doc.Name)
			resultText += fmt.Sprintf("   - ID: `%s`\n", doc.ID)
			resultText += fmt.Sprintf("   - 类型: %s\n", doc.TextContent.ContentType)
			resultText += fmt.Sprintf("   - 描述: %s\n", doc.Description)
			if len(doc.Tags) > 0 {
				resultText += fmt.Sprintf("   - 标签: %v\n", doc.Tags)
			}
			resultText += fmt.Sprintf("   - 更新时间: %s\n", doc.UpdatedAt)
			resultText += "\n"
		}
	}

	resultText += "---\n\n"
	resultText += "**文本文档管理提示：**\n"
	resultText += "- 使用 `search_text_documents` 进行详细搜索\n"
	resultText += "- 使用 `get_text_document` 获取完整文档内容\n"
	resultText += "- 通过标签分类管理不同类型的文本文档"
	resultText += "- 内容类型可用于区分文档用途（如README、配置、提示词等）"
	if len(result.Documents) == 0 {
		resultText += "\n- 现在就创建第一个文本文档开始管理你的知识库吧！"
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
		}{Success: true, Documents: result.Documents, Count: len(result.Documents), Message: "列出成功"}, nil
}

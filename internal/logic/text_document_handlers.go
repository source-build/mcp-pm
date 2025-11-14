package logic

import (
	"context"
	"encoding/json"
	"fmt"

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
//
// 🎯 **AI使用指南**
// 这个工具专门用于创建和管理各种类型的文本文档。AI可以通过以下方式使用：
//
// 1. **创建README文档**：
//   - contentType: "readme"
//   - content: 使用Markdown格式的项目说明文档
//   - 适用于：项目说明、使用指南、安装说明等
//
// 2. **创建配置文档**：
//   - contentType: "config"
//   - content: JSON、YAML、TOML等配置文件内容
//   - 适用于：应用配置、环境配置、服务配置等
//
// 3. **创建提示词模板**：
//   - contentType: "prompt"
//   - content: 包含变量占位符的提示词，如"请分析以下代码：{{.code}}"
//   - 适用于：代码审查模板、分析模板、生成模板等
//
// 4. **创建技术规范**：
//   - contentType: "spec"
//   - content: 结构化的技术规范文档
//   - 适用于：接口规范、数据结构定义、编码规范等
//
// 5. **创建普通笔记**：
//   - contentType: "note"
//   - content: 纯文本内容，可包含简单的Markdown格式
//   - 适用于：会议记录、技术笔记、待办事项等
//
// ⚠️ **AI调用注意事项**
// - **必需参数**: name, description, content_type, content必须提供
// - **命名规范**: 建议使用"用途-类型"格式，如"project-readme"或"server-config"
// - **内容格式**: 根据contentType选择合适的内容格式
// - **变量占位符**: 在prompt类型中，使用{{.variable}}格式
// - **标签分类**: 使用标签进行分类管理，如['docs', 'config', 'template']
//
// 🔍 **搜索和发现**: AI可以通过以下方式找到创建的文档
// - 使用`search_text_documents`按内容、类型或标签搜索
// - 使用`list_text_documents`按内容类型分类浏览
// - 通过content_type参数快速筛选特定类型文档
//
// 📋 **最佳实践**
// - 为README文档添加清晰的结构和目录
// - 配置文档包含注释说明每个配置项
// - 提示词模板使用有意义的变量名
// - 技术规范遵循标准格式和术语
//
// 🛡️ **参数详细说明**
func CreateTextDocument(ctx context.Context, req *mcp.CallToolRequest, input struct {
	Name        string            `json:"name" jsonschema:"文档名称，必填。建议使用描述性名称，如'project-readme'或'server-config'"`
	Description string            `json:"description" jsonschema:"文档详细描述，必填。清晰说明文档的用途、内容和适用场景，帮助AI理解这个文档的作用"`
	ContentType string            `json:"content_type" jsonschema:"内容类型，必填。支持的类型：'readme'(README文档), 'prompt'(提示词模板), 'config'(配置文件), 'note'(普通笔记), 'spec'(技术规范)"`
	Variables   map[string]string `json:"variables" jsonschema:"模板变量，用于提示词模板，可选"`
	Content     interface{}       `json:"content" jsonschema:"文档内容，必填。可以是字符串、JSON对象或任意可转换为字符串的内容。根据content_type自动格式化存储"`
	Tags        []string          `json:"tags,omitempty" jsonschema:"文档标签，可选。用于分类和搜索，建议使用如：['readme', 'config', 'guide', 'template', 'docs']等标签"`
}) (*mcp.CallToolResult, struct {
	Success  bool            `json:"success"`
	Document *types.Document `json:"document"`
	Message  string          `json:"message"`
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

	fmt.Printf("🔐 用户 %s 正在项目 %s 中创建文本文档: %s\n", userId, projectId, input.Name)

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

	// 创建文本文档对象
	document := &types.Document{
		ProjectID:   projectId,
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
	document.ID, err = es.ESClient.SaveDocument(document)
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
//
// 根据文档ID获取详细的文本文档信息，包含完整内容结构。
// 主要用于：
// - 查看文档完整内容
// - 获取特定格式的内容
// - 复制文档用于其他工具
//
// 注意事项：
//   - 需要提供确切的文档ID
//   - 返回的JSON格式保持原始结构
//   - 只能访问当前用户有权限的项目文档
func GetTextDocument(ctx context.Context, req *mcp.CallToolRequest, input struct {
	ID string `json:"id" jsonschema:"文本文档ID"`
}) (*mcp.CallToolResult, *types.Document, error) {
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

	// 验证权限 - 使用token提取工具函数
	projectId, err := utils.ExtractProjectID(req)
	if err == nil && document.ProjectID != projectId {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("无权访问该文本文档，文档属于项目 %s，但当前项目是 %s", document.ProjectID, projectId)},
			},
			IsError: true,
		}, nil, fmt.Errorf("无权访问")
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
//
// 在当前项目的所有文本文档中搜索，支持多种筛选条件：
// - 关键词搜索：在文档名称、描述、内容等中搜索
// - 内容类型筛选：只返回指定内容类型的文档
// - 标签筛选：根据标签分类搜索特定类型的文档
//
// 搜索逻辑：
//   - 支持模糊匹配和精确匹配
//   - 搜索范围包括：文档名称、描述、文档内容
//   - 支持组合搜索条件
//
// 使用技巧：
//   - 搜索配置相关文档：关键词使用"config"
//   - 查找特定类型：contentType参数设为"readme"
//   - 按标签筛选：tags参数设为["setup", "guide"]
//   - 搜索文档内容：可使用文档中的任意关键词
func SearchTextDocuments(ctx context.Context, req *mcp.CallToolRequest, input struct {
	Query       string   `json:"query" jsonschema:"搜索关键词"`
	ContentType string   `json:"content_type,omitempty" jsonschema:"内容类型筛选"`
	Tags        []string `json:"tags,omitempty" jsonschema:"标签筛选"`
	Limit       int      `json:"limit,omitempty" jsonschema:"返回数量限制，默认20"`
	Offset      int      `json:"offset,omitempty" jsonschema:"偏移量，默认0"`
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
	textType := types.DocumentTypeText
	searchReq := &types.SearchRequest{
		Query:  input.Query,
		Limit:  input.Limit,
		Offset: input.Offset,
		Tags:   input.Tags,
		Type:   &textType,
	}

	// 执行搜索
	result, err := es.ESClient.SearchDocuments(projectId, searchReq)
	if err != nil {
		return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{Text: fmt.Sprintf("搜索文本文档失败: %v", err)},
				},
				IsError: true,
			}, struct {
				Success   bool             `json:"success"`
				Documents []types.Document `json:"documents"`
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
			Documents []types.Document `json:"documents"`
			Count     int              `json:"count"`
			Message   string           `json:"message"`
		}{Success: true, Documents: textDocs, Count: len(textDocs), Message: "搜索成功"}, nil
}

// ListTextDocuments 列出文本文档
//
// 列出当前项目的所有文本文档，支持按内容类型分类筛选。
// 用于：
// - 查看项目中的所有文档
// - 按文档类型分类浏览
// - 快速定位特定类型的文档
//
// 分类说明：
//   - 不指定类型：列出所有文本文档
//   - 指定类型（如readme）：只返回该类型的文档
//   - 支持的类型：readme, prompt, config, note, spec等
//
// 注意事项：
//   - 返回按更新时间倒序排列
//   - 支持分页浏览大量文档
//   - 只列出用户有权限访问的项目文档
func ListTextDocuments(ctx context.Context, req *mcp.CallToolRequest, input struct {
	ContentType string `json:"content_type,omitempty" jsonschema:"内容类型筛选"`
	Limit       int    `json:"limit,omitempty" jsonschema:"返回数量限制，默认20"`
	Offset      int    `json:"offset,omitempty" jsonschema:"偏移量，默认0"`
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

	fmt.Printf("🔐 用户 %s 正在项目 %s 中列出文本文档\n", userId, projectId)

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
	result, err := es.ESClient.ListDocuments(projectId, docType, limit, offset)
	if err != nil {
		return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{Text: fmt.Sprintf("列出文本文档失败: %v", err)},
				},
				IsError: true,
			}, struct {
				Success   bool             `json:"success"`
				Documents []types.Document `json:"documents"`
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

	resultText += "---\n\n"
	resultText += "**文本文档管理提示：**\n"
	resultText += "- 使用 `search_text_documents` 进行详细搜索\n"
	resultText += "- 使用 `get_text_document` 获取完整文档内容\n"
	resultText += "- 通过标签分类管理不同类型的文本文档"
	resultText += "- 内容类型可用于区分文档用途（如README、配置、提示词等）"

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

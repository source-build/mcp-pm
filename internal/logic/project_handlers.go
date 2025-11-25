package logic

import (
	"context"
	"fmt"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/source-build/mcp-pm/internal/pm"
	"github.com/source-build/mcp-pm/internal/types"
	"github.com/source-build/mcp-pm/internal/utils"
)

/*
项目管理处理器

注意：这些工具专门用于项目上下文管理，包括：
- 列出用户可访问的项目
- 切换当前工作项目
- 获取项目详细信息

适用场景：
- 多项目环境下的项目切换
- 查看项目基本信息和权限
- 管理当前工作上下文

重要说明：
- 这些工具不包含项目创建、删除等管理功能
- 项目管理应通过其他系统进行
- 这里只提供项目查看和切换功能

数据结构说明：
- Project: 包含项目基本信息、用户权限列表、元数据等
- UserContext: 包含当前用户Token、选中的项目信息、过期时间等
- 标签：用于分类不同类型的项目（如production, staging, development等）
*/

// ListProjects 列出项目处理器
func ListProjects(_ context.Context, req *mcp.CallToolRequest, input struct{}) (*mcp.CallToolResult, struct {
	Success  bool             `json:"success"`
	Projects []*types.Project `json:"projects"`
	Count    int              `json:"count"`
	Message  string           `json:"message"`
}, error) {
	// 获取项目管理器
	p := pm.NewProjectManager()

	// 获取用户的项目列表
	projects, err := p.ListAvailableProjects(req)
	if err != nil {
		return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{Text: fmt.Sprintf("获取项目列表失败: %v", err)},
				},
				IsError: true,
			}, struct {
				Success  bool             `json:"success"`
				Projects []*types.Project `json:"projects"`
				Count    int              `json:"count"`
				Message  string           `json:"message"`
			}{Success: false, Count: 0, Message: fmt.Sprintf("获取失败: %v", err)}, err
	}

	// 获取当前项目ID（如果有的话）
	currentProjectID, err := utils.ExtractProjectID(req)
	if err != nil {
		// 如果无法获取项目ID，说明token中没有当前项目信息
		currentProjectID = ""
	}

	// 构建结果文本
	resultText := fmt.Sprintf("# 可访问的项目列表 (共 %d 个)\n\n", len(projects))

	if currentProjectID != "" {
		resultText += fmt.Sprintf("**Token中的项目ID:** %s\n\n", currentProjectID)
	}

	for i, project := range projects {
		status := ""
		if project.ID == currentProjectID {
			status = " 📍 **当前项目**"
		}

		resultText += fmt.Sprintf("**%d. %s**%s\n", i+1, project.Name, status)
		resultText += fmt.Sprintf("   - ID: `%s`\n", project.ID)
		resultText += fmt.Sprintf("   - 描述: %s\n", project.Description)
		resultText += fmt.Sprintf("   - 用户数: %d\n", len(project.UserIds))
		resultText += fmt.Sprintf("   - 创建时间: %s\n", project.CreatedAt)

		if len(project.Metadata) > 0 {
			resultText += fmt.Sprintf("   - 元数据: %d 项\n", len(project.Metadata))
		}
		resultText += "\n"
	}

	resultText += "---\n\n"
	resultText += "**使用说明：**\n"
	resultText += "2. 使用 `get_project_info` 获取项目详细信息\n"
	resultText += "3. 切换项目后，所有文档操作将在当前项目上下文中进行"

	return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: resultText},
			},
		}, struct {
			Success  bool             `json:"success"`
			Projects []*types.Project `json:"projects"`
			Count    int              `json:"count"`
			Message  string           `json:"message"`
		}{Success: true, Projects: projects, Count: len(projects), Message: "获取项目列表成功"}, nil
}

// GetProjectInfo 获取项目信息处理器
func GetProjectInfo(_ context.Context, req *mcp.CallToolRequest, input struct {
	ProjectID string `json:"project_id,omitempty" jsonschema:"项目ID，可选。要查询的项目标识符，如果不提供则使用Token中的当前项目。当用户说'查询当前项目'、'查看项目信息'时可不传，系统会自动使用当前项目。传参时必须是有效的项目ID。"`
}) (*mcp.CallToolResult, struct {
	Success bool           `json:"success"`
	Project *types.Project `json:"project"`
	Message string         `json:"message"`
}, error) {
	// 获取项目管理器
	p := pm.NewProjectManager()

	// 获取用户ID和项目ID
	userId, err := utils.ExtractUserID(req)
	if err != nil {
		return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{Text: fmt.Sprintf("获取用户ID失败: %v", err)},
				},
				IsError: true,
			}, struct {
				Success bool           `json:"success"`
				Project *types.Project `json:"project"`
				Message string         `json:"message"`
			}{Success: false, Message: fmt.Sprintf("获取用户ID失败: %v", err)}, err
	}

	// 确定要查询的项目ID
	projectID := input.ProjectID
	if projectID == "" {
		// 尝试从Token中获取项目ID
		projectID, err = utils.ExtractProjectID(req)
		if err != nil {
			return &mcp.CallToolResult{
					Content: []mcp.Content{
						&mcp.TextContent{Text: "没有指定项目ID，且Token中没有项目信息"},
					},
					IsError: true,
				}, struct {
					Success bool           `json:"success"`
					Project *types.Project `json:"project"`
					Message string         `json:"message"`
				}{Success: false, Message: "无法确定项目ID"}, nil
		}
	}

	if projectID == "" {
		return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{Text: "没有指定项目ID，且Token中没有项目信息"},
				},
				IsError: true,
			}, struct {
				Success bool           `json:"success"`
				Project *types.Project `json:"project"`
				Message string         `json:"message"`
			}{Success: false, Message: "无法确定项目ID"}, nil
	}

	fmt.Printf("🔐 用户 %s 正在查询项目 %s 的信息\n", userId, projectID)

	// 获取项目信息
	project, err := p.GetProjectInfo(req, projectID)
	if err != nil {
		return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{Text: fmt.Sprintf("获取项目信息失败: %v", err)},
				},
				IsError: true,
			}, struct {
				Success bool           `json:"success"`
				Project *types.Project `json:"project"`
				Message string         `json:"message"`
			}{Success: false, Message: fmt.Sprintf("获取失败: %v", err)}, err
	}

	// 构建结果文本
	resultText := fmt.Sprintf("# 项目详细信息\n\n")
	resultText += fmt.Sprintf("**项目名称:** %s\n", project.Name)
	resultText += fmt.Sprintf("**项目ID:** %s\n", project.ID)
	resultText += fmt.Sprintf("**项目描述:** %s\n", project.Description)
	resultText += fmt.Sprintf("**创建时间:** %s\n", project.CreatedAt)
	resultText += fmt.Sprintf("**更新时间:** %s\n", project.UpdatedAt)
	resultText += fmt.Sprintf("**用户数量:** %d\n", len(project.UserIds))

	if len(project.UserIds) > 0 {
		resultText += fmt.Sprintf("**用户IDs:** `%s`\n", strings.Join(project.UserIds, "`, `"))
	}

	if len(project.Metadata) > 0 {
		resultText += "\n**项目元数据:**\n"
		for key, value := range project.Metadata {
			resultText += fmt.Sprintf("- %s: %s\n", key, value)
		}
	}

	resultText += "\n\n---\n\n"
	resultText += "**项目操作提示：**\n"
	resultText += "- 使用 `list_api_documents` 查看API文档\n"
	resultText += "- 使用 `list_text_documents` 查看文本文档"

	return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: resultText},
			},
		}, struct {
			Success bool           `json:"success"`
			Project *types.Project `json:"project"`
			Message string         `json:"message"`
		}{Success: true, Project: project, Message: "获取项目信息成功"}, nil
}

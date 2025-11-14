package pm

import (
	"fmt"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/source-build/mcp-pm/internal/es"
	"github.com/source-build/mcp-pm/internal/types"
	"github.com/source-build/mcp-pm/internal/utils"
)

// ProjectManager 项目上下文管理器
type ProjectManager struct {
	esClient *es.Client
}

// NewProjectManager 创建项目管理器
func NewProjectManager() *ProjectManager {
	return &ProjectManager{
		esClient: es.ESClient,
	}
}

// hasProjectAccess 检查用户是否有项目访问权限
func (pm *ProjectManager) hasProjectAccess(userId string, project *types.Project) bool {
	for _, id := range project.UserIds {
		if id == userId {
			return true
		}
	}
	return false
}

// ListAvailableProjects 列出用户可以访问的项目
func (pm *ProjectManager) ListAvailableProjects(req *mcp.CallToolRequest) ([]*types.Project, error) {
	userId, err := utils.ExtractUserID(req)
	if err != nil {
		return nil, fmt.Errorf("提取用户ID失败: %v", err)
	}

	// 获取用户的所有项目
	projects, err := pm.esClient.FindProjectsByToken(userId)
	if err != nil {
		return nil, fmt.Errorf("查找项目失败: %v", err)
	}

	return projects, nil
}

// GetProjectInfo 获取项目信息
func (pm *ProjectManager) GetProjectInfo(req *mcp.CallToolRequest, projectID string) (*types.Project, error) {
	userId, err := utils.ExtractUserID(req)
	if err != nil {
		return nil, fmt.Errorf("提取用户ID失败: %v", err)
	}

	// 获取项目
	project, err := pm.esClient.GetProject(projectID)
	if err != nil {
		return nil, fmt.Errorf("获取项目失败: %v", err)
	}

	// 检查访问权限
	if !pm.hasProjectAccess(userId, project) {
		return nil, fmt.Errorf("用户无权访问项目 %s", projectID)
	}

	return project, nil
}

// CreateSampleProject 创建示例项目（仅用于演示）
// 这个方法展示了如何创建项目，但在实际MCP项目中不应该调用
func CreateSampleProject() *types.Project {
	t := time.Now()
	return &types.Project{
		ID:          "mxd",
		Name:        "秒选达",
		Description: "秒选达项目",
		UserIds:     []string{"1"},
		Metadata:    map[string]string{"type": "mxd"},
		CreatedAt:   &t,
	}
}

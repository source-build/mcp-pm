package utils

import (
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// ExtractUserID 提取用户ID
func ExtractUserID(req *mcp.CallToolRequest) (string, error) {
	userId, ok := req.Extra.TokenInfo.Extra["user_id"].(string)
	if !ok {
		return "", fmt.Errorf("user_id not found in token info")
	}
	return userId, nil
}

// ExtractProjectID 提取项目ID
func ExtractProjectID(req *mcp.CallToolRequest) (string, error) {
	projectID, ok := req.Extra.TokenInfo.Extra["project_id"].(string)
	if !ok {
		return "", fmt.Errorf("project_id not found in token info")
	}
	return projectID, nil
}

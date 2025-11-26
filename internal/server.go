package internal

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/modelcontextprotocol/go-sdk/auth"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/source-build/mcp-pm/internal/config"
	"github.com/source-build/mcp-pm/internal/logic"
	"github.com/source-build/mcp-pm/internal/middleware"
	"github.com/source-build/mcp-pm/internal/token"
)

type GetTimeParams struct {
	City string `json:"city" jsonschema:"City to get time for (nyc, sf, or boston)"`
}

// createMCPServer 创建并配置MCP服务器
func createMCPServer() *mcp.Server {
	// 创建 MCP 服务器实例
	server := mcp.NewServer(&mcp.Implementation{
		Name:    config.Config.ServerName,
		Version: config.Config.ServerVersion,
	}, nil)

	// ========================================
	// 项目管理工具
	// ========================================
	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_default_project",
		Description: "获取默认项目。用于获取当前用户默认设置的项目信息。",
	}, logic.GetDefaultProject)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_projects",
		Description: "列出用户可以访问的项目。用于查看当前用户有权访问的所有项目列表。",
	}, logic.ListProjects)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_project_info",
		Description: "获取项目详细信息。用于查看当前或指定项目的详细信息和元数据。",
	}, logic.GetProjectInfo)

	// ========================================
	// API文档管理工具
	// ========================================

	mcp.AddTool(server, &mcp.Tool{
		Name:        "create_api_document",
		Description: "创建API文档。用于创建和管理各种API接口文档，支持REST类型。",
	}, logic.CreateAPIDocument)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "edit_api_document",
		Description: "编辑API文档。根据文档ID编辑指定的API文档内容。",
	}, logic.EditAPIDocument)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "del_api_document",
		Description: "删除API文档。根据文档ID删除指定的API文档。",
	}, logic.DelAPIDocument)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_api_document",
		Description: "获取API文档详情。根据文档ID获取完整的API文档信息，包含请求/响应结构。",
	}, logic.GetAPIDocument)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "search_api_documents",
		Description: "搜索API文档。在当前项目中搜索API文档，支持关键词、HTTP方法、标签等多种筛选条件。",
	}, logic.SearchAPIDocuments)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_api_documents",
		Description: "列出API文档。按HTTP方法分类列出当前项目的所有API文档。",
	}, logic.ListAPIDocuments)

	// ========================================
	// 文本文档管理工具
	// ========================================

	mcp.AddTool(server, &mcp.Tool{
		Name:        "create_text_document",
		Description: "创建文本文档。用于创建和管理各种文本文档，支持README、配置、提示词、规范等类型。",
	}, logic.CreateTextDocument)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_text_document",
		Description: "获取文本文档详情。根据文档ID获取完整的文本文档信息和内容。",
	}, logic.GetTextDocument)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "search_text_documents",
		Description: "搜索文本文档。在当前项目中搜索文本文档，支持关键词、内容类型、标签等筛选条件。",
	}, logic.SearchTextDocuments)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_text_documents",
		Description: "列出文本文档。按内容类型分类列出当前项目的所有文本文档。",
	}, logic.ListTextDocuments)

	return server
}

func verifyJWT(ctx context.Context, tokenString string, req *http.Request) (*auth.TokenInfo, error) {
	claims, err := token.VerifyJWT(tokenString)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", auth.ErrInvalidToken, err)
	}

	return &auth.TokenInfo{
		Scopes:     []string{"read"},      // User permissions
		Expiration: claims.ExpiresAt.Time, // Token expiration time
		Extra: map[string]interface{}{
			"user_id":    claims.UserID,
			"project_id": claims.ProjectID,
		},
	}, nil
}

func Server() {
	// 创建 MCP 服务器实例
	server := createMCPServer()

	// Create authentication middleware.
	jwtAuth := auth.RequireBearerToken(verifyJWT, &auth.RequireBearerTokenOptions{
		Scopes: []string{"read"}, // Require "read" permission
	})

	// 添加日志记录中间件
	// 用于记录所有请求和响应，便于调试和监控
	server.AddReceivingMiddleware(middleware.CreateLoggingMiddleware())

	// 创建可流式 HTTP 处理函数
	// 将MCP协议转换为HTTP接口
	handler := mcp.NewStreamableHTTPHandler(func(req *http.Request) *mcp.Server { return server }, nil)

	authenticatedHandler := jwtAuth(handler)

	// jwt: 处理JWT认证请求
	http.HandleFunc("/mcp", authenticatedHandler.ServeHTTP)

	url := fmt.Sprintf("%s:%s", config.Config.HTTPAddr, config.Config.HTTPPort)

	log.Printf("=======================================")
	log.Printf("API文档管理MCP服务器启动")
	log.Printf("服务地址: %s", url)
	log.Printf("=======================================")

	// 启动 HTTP 服务器
	if err := http.ListenAndServe(url, nil); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}

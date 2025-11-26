package main

import (
	"log"

	"github.com/source-build/mcp-pm/internal"
	"github.com/source-build/mcp-pm/internal/config"
	"github.com/source-build/mcp-pm/internal/es"
)

func main() {
	// 加载配置
	err := config.LoadConfig()
	if err != nil {
		panic("无法加载配置" + err.Error())
	}

	// 初始化Elasticsearch
	err = es.InitESClient()
	if err != nil {
		log.Fatalf("初始化Elasticsearch失败: %v", err)
	}
	defer es.ESClient.Close()

	// 启动MCP服务器
	internal.Server()
}

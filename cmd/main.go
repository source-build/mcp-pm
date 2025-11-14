package main

import (
	"log"

	"github.com/source-build/mcp-pm/internal"
	"github.com/source-build/mcp-pm/internal/config"
	"github.com/source-build/mcp-pm/internal/es"
)

//func generateToken() {
//	t, err := token.GenerateToken("123456", "mxd", time.Hour*24)
//	if err != nil {
//		panic("生成token失败" + err.Error())
//	}
//	fmt.Println(t)
//}

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

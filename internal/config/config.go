package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
)

// Options 应用配置结构
type Options struct {
	ServerName    string
	ServerVersion string
	HTTPAddr      string
	HTTPPort      string

	// Elasticsearch配置
	EsURL      string
	EsUsername string
	EsPassword string
	EsSniff    bool

	// 认证配置
	ServerJwtSecret string

	// 索引配置
	ProjectIndex  string
	DocumentIndex string
}

var Config *Options

// LoadConfig 加载配置
func LoadConfig() error {
	// 加载.env文件
	if err := godotenv.Load(); err != nil {
		log.Printf("警告: 无法加载.env文件: %v", err)
		return err
	}

	// 获取服务器配置
	serverName := os.Getenv("SERVER_NAME")
	if serverName == "" {
		serverName = "api-docs-manager"
	}

	serverVersion := os.Getenv("SERVER_VERSION")
	if serverVersion == "" {
		serverVersion = "1.0.0"
	}

	// 获取HTTP配置
	httpAddr := os.Getenv("HTTP_ADDR")
	if httpAddr == "" {
		httpAddr = "0.0.0.0"
	}

	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = "8999"
	}

	// 获取Elasticsearch配置
	esURL := os.Getenv("ES_URL")
	if esURL == "" {
		esURL = "http://localhost:9200"
	}

	esUsername := os.Getenv("ES_USERNAME")
	if esUsername == "" {
		esUsername = "elastic"
	}

	esPassword := os.Getenv("ES_PASSWORD")
	if esPassword == "" {
		esPassword = "changeme"
	}

	// 获取索引配置
	projectIndex := os.Getenv("PROJECT_INDEX")
	if projectIndex == "" {
		projectIndex = "projects"
	}

	documentIndex := os.Getenv("DOCUMENT_INDEX")
	if documentIndex == "" {
		documentIndex = "documents"
	}

	// 获取认证配置
	serverJwtSecret := os.Getenv("JWT_SECRET")
	if serverJwtSecret == "" {
		serverJwtSecret = "udHFYTnV68auQjeGrYbKDYxwVufhpzcG"
	}

	Config = &Options{
		ServerName:      serverName,
		ServerVersion:   serverVersion,
		HTTPAddr:        httpAddr,
		HTTPPort:        httpPort,
		EsURL:           esURL,
		EsUsername:      esUsername,
		EsPassword:      esPassword,
		EsSniff:         false,
		ProjectIndex:    projectIndex,
		DocumentIndex:   documentIndex,
		ServerJwtSecret: serverJwtSecret,
	}

	return nil
}

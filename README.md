# 源构互联 - 项目管理MCP工具

[![Go Version](https://img.shields.io/badge/Go-1.25.0-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Version](https://img.shields.io/badge/Version-2.0.0-red.svg)](CHANGELOG.md)

## 🏢 关于源构互联

源构互联是一家胆子非常大的公司。本MCP工具是我们推出的开源项目，旨在帮助开发团队更好地管理API文档和技术资料。

## 📋 项目概述

这是一个基于Model Context Protocol (MCP) 的专业项目管理工具，支持完整的文档生命周期管理，包括API文档、技术文档、配置文件等多种类型的内容管理。

### 🎯 核心价值

- **🏗️ 企业级多租户架构** - 支持多项目隔离，确保数据安全
- **📚 全方位文档管理** - 统一管理API文档、技术文档、配置文件
- **🔍 智能搜索引擎** - 基于Elasticsearch的高性能全文搜索
- **🔄 灵活的协作模式** - 支持团队协作和权限管理
- **🚀 云原生设计** - 容器化部署，支持微服务架构

## ✨ 核心特性

### 🏢️ 多租户项目管理

| 功能特性 | 说明 | 优势 |
|---------|------|------|
| **项目隔离** | 每个项目数据完全独立 | 数据安全，权限隔离 |
| **用户认证** | 基于JWT Token的安全认证机制 | 高安全性，防止未授权访问 |
| **权限控制** | 细粒度的项目访问权限 | 精确控制用户操作范围 |
| **上下文管理** | 自动项目上下文切换和管理 | 提高工作效率，支持多项目管理 |

### 📚 多类型文档管理

#### 📡 API文档管理
- **支持REST API**: 完整的REST接口文档管理
- **完整的请求/响应定义**: 包含参数、示例、业务码
- **灵活的参数类型**: 支持body、query、path参数
- **标签分类**: 按功能模块、状态、团队进行分类

#### 📄 文本文档管理
- **README文档**: 项目说明、安装指南、使用文档
- **配置文件**: JSON等配置模板
- **提示词模板**: AI辅助的代码审查、分析模板
- **技术规范**: 接口规范、编码标准、设计文档
- **团队笔记**: 会议记录、技术笔记、待办事项

### 🔍 智能搜索系统

| 搜索维度 | 支持条件 | 应用场景 |
|---------|---------|---------|
| **关键词搜索** | 模糊匹配、相关性排序 | 快速定位相关内容 |
| **类型筛选** | 文档类型、API方法、内容分类 | 精确查找特定类型文档 |
| **标签搜索** | 多标签组合、层级筛选 | 按业务场景分类查找 |
| **全文搜索** | 内容深度搜索、高亮显示 | 搜索文档内部内容 |

## 🛠️ 工具生态

### 🏢️ 项目管理工具集

| 工具名称 | 功能描述 | 使用场景 |
|---------|---------|---------|
| `list_projects` | 列出用户可访问的所有项目 | 项目概览，资源查看 |
| `get_project_info` | 获取项目的详细配置信息 | 项目配置查看，权限确认 |

**使用示例:**
```bash
# 查看所有可用项目
list_projects

# 查看当前项目详细信息
get_project_info

# 获取指定项目信息
get_project_info --project_id="web-backend-001"
```

### 📡 API文档管理工具集

| 工具名称 | 功能描述 | 支持类型 | 典型应用 |
|---------|---------|---------|---------|
| `create_api_document` | 创建标准化API接口文档 | REST API | 用户认证API、支付接口文档 |
| `get_api_document` | 获取API文档完整结构 | 所有API类型 | 查看接口详情，复制文档模板 |
| `search_api_documents` | 多维度API文档搜索 | 关键词、HTTP方法、标签 | 查找所有用户相关接口 |
| `list_api_documents` | 按分类列出API文档 | HTTP方法分类 | 浏览所有POST类型接口 |

**使用示例:**
```bash
# 创建用户管理API文档
create_api_document \
  --name="user-auth-api" \
  --method="POST" \
  --path="/api/auth/login" \
  --description="用户认证接口" \
  --tags="auth,user,public" \
  --request='{"username": "string", "password": "string"}' \
  --response='{"code": 200, "result": {"token": "jwt_token"}, "msg": "登录成功"}'

# 搜索所有认证相关API
search_api_documents --query="认证" --tags="auth"

# 查看所有GET方式的接口
list_api_documents --method="GET"
```

### 📄 文本文档管理工具集

| 工具名称 | 功能描述 | 内容类型 | 业务场景 |
|---------|---------|---------|---------|
| `create_text_document` | 创建各类文本文档 | README、配置、提示词、笔记、规范 | 项目说明文档、配置模板 |
| `get_text_document` | 获取文本文档完整内容 | 所有文本类型 | 查看配置详情、获取文档内容 |
| `search_text_documents` | 智能文本文档搜索 | 关键词、内容类型、标签 | 查找特定配置文件 |
| `list_text_documents` | 按类型列出文本文档 | 内容类型分类 | 浏览所有README文档 |

**使用示例:**
```bash
# 创建项目README文档
create_text_document \
  --name="project-readme" \
  --content_type="readme" \
  --description="项目说明文档" \
  --tags="readme,documentation,onboarding" \
  --content="# 项目说明\n\n这是一个企业级Web应用项目..."

# 创建应用配置模板
create_text_document \
  --name="production-config" \
  --content_type="config" \
  --description="生产环境配置模板" \
  --tags="config,production" \
  --content='{"debug": false, "port": 8080, "database": {"host": "prod-db"}}'

# 创建代码审查提示词
create_text_document \
  --name="code-review-prompt" \
  --content_type="prompt" \
  --description="代码审查提示词模板" \
  --content="请审查以下{{.language}}代码，重点关注{{.focus_points}}"

# 查找所有配置相关文档
search_text_documents --query="配置" --content_type="config"
```

## 🎯 典型应用场景

### 🏢 企业级项目协作

**场景: 多团队协作开发大型企业应用**
```bash
# 1. 项目经理查看所有项目
list_projects

# 2. API团队管理接口文档
create_api_document --name="user-registration-api" ...
search_api_documents --query="用户注册" ...

# 3. 运维团队管理配置文档
create_text_document --name="deployment-config" ...
search_text_documents --content_type="config" ...
```

### 🚀 微服务架构管理

**场景: 管理复杂的微服务生态系统**
```bash
# 为每个微服务创建独立项目
# user-service, order-service, payment-service, notification-service

# 管理服务间API文档
create_api_document --name="user-order-integration" ...
create_api_document --name="order-payment-webhook" ...

# 统一管理配置和部署文档
create_text_document --name="kubernetes-deployment" --content_type="config"
create_text_document --name="service-monitoring-guide" --content_type="readme"
```

### 📚 技术知识库建设

**场景: 构建企业技术知识库**
```bash
# 创建编码规范文档
create_text_document --name="go-coding-standards" --content_type="spec"

# 创建最佳实践指南
create_text_document --name="database-design-patterns" --content_type="spec"

# 创建故障排查手册
create_text_document --name="troubleshooting-guide" --content_type="note"

# 创建代码审查清单
create_text_document --name="review-checklist" --content_type="prompt"
```

## 🏗️ 技术架构

### 系统架构图
```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   MCP Client    │    │   AI Assistant  │    │  Web Dashboard  │
└─────────┬───────┘    └─────────┬───────┘    └─────────┬───────┘
          │                      │                      │
          └──────────────────────┼──────────────────────┘
                                 │
                    ┌─────────────▼─────────────┐
                    │      MCP Server           │
                    │  (源构互联项目管理工具)    │
                    └─────────────┬─────────────┘
                                 │
          ┌──────────────────────┼──────────────────────┐
          │                      │                      │
┌─────────▼─────────┐  ┌─────────▼─────────┐  ┌─────────▼─────────┐
│   Authentication  │  │   Project Mgmt    │  │  Document Mgmt    │
│     Service       │  │     Service       │  │     Service       │
└─────────┬─────────┘  └─────────┬─────────┘  └─────────┬─────────┘
          │                      │                      │
          └──────────────────────┼──────────────────────┘
                                 │
                    ┌─────────────▼─────────────┐
                    │     Elasticsearch         │
                    │   (搜索引擎与数据存储)     │
                    └───────────────────────────┘
```

### 核心组件

#### 🔐 认证与授权系统
- **JWT Token认证**: 基于JWT的安全认证机制
- **权限隔离**: 多租户环境下的数据隔离
- **上下文管理**: 项目上下文的自动切换和管理

#### 📊 数据存储层
- **Elasticsearch**: 主要的搜索和数据存储引擎
- **文档索引**: 高性能的全文搜索和索引管理
- **数据建模**: 灵活的文档结构支持

#### 🚀 应用服务层
- **项目管理器**: 处理项目创建、切换、权限验证
- **文档处理器**: 分类处理不同类型的文档操作
- **搜索引擎**: 智能搜索和结果排序

### 数据模型设计

#### 项目模型 (Project)
```json
{
  "id": "web-prod-001",
  "name": "Web生产环境项目",
  "description": "企业级Web应用生产环境",
  "user_ids": ["token-admin-001", "token-dev-002", "token-ops-003"],
  "created_at": "2024-01-15T09:00:00Z",
  "updated_at": "2024-01-20T14:30:00Z",
  "metadata": {
    "environment": "production",
    "team": "backend-team",
    "domain": "user-management",
    "criticality": "high"
  }
}
```

#### API文档模型 (APIDocument)
```json
{
  "id": "api-user-auth-001",
  "project_id": "web-prod-001",
  "type": "api",
  "name": "用户认证API",
  "description": "用户登录和身份验证接口",
  "tags": ["auth", "security", "public-api"],
  "api_content": {
    "method": "POST",
    "path": "/api/v1/auth/login",
    "body": {
      "username": {"type": "string", "required": true},
      "password": {"type": "string", "required": true}
    },
    "response_biz_code": "200",
    "header": {
      "Content-Type": "application/json"
    }
  },
  "created_at": "2024-01-15T10:00:00Z",
  "updated_at": "2024-01-18T16:45:00Z"
}
```

#### 文本文档模型 (TextDocument)
```json
{
  "id": "doc-readme-001",
  "project_id": "web-prod-001",
  "type": "text",
  "name": "项目README文档",
  "description": "项目说明和快速开始指南",
  "tags": ["readme", "documentation", "onboarding"],
  "text_content": {
    "content_type": "readme",
    "content": "# Web应用项目\n\n## 项目概述\n这是一个企业级Web应用程序...",
    "variables": {
      "project_name": "web-app",
      "version": "2.0.0"
    }
  },
  "created_at": "2024-01-15T11:00:00Z",
  "updated_at": "2024-01-19T09:30:00Z"
}
```

## 🚀 快速开始

### 环境要求

- **Go**: 1.25.0 或更高版本
- **Elasticsearch**: 7.0+ (推荐 8.0+)
- **Docker**: 可选，用于容器化部署
- **Git**: 版本控制管理

### 安装部署

#### 方式一: 直接编译安装
```bash
# 1. 克隆代码库
git clone https://github.com/source-build/mcp-pm
cd mcp-pm

# 2. 安装依赖
go mod tidy

# 3. 编译应用
go build -o mcp-pm-mcp ./cmd

# 4. 配置环境变量
cp .env.example .env
# 编辑 .env 文件配置Elasticsearch连接信息

# 5. 启动服务
./mcp-pm-mcp
```

#### 方式二: Docker容器部署
```bash
# 1. 使用Docker Compose一键部署
docker-compose up -d

# 2. 查看服务状态
docker-compose ps

# 3. 查看日志
docker-compose logs -f mcp-pm-mcp
```

### 配置说明

#### 环境变量配置
```bash
# Elasticsearch配置
ES_URL=http://localhost:9200
ES_USERNAME=elastic
ES_PASSWORD=changeme
ES_SNIFF=false

# 服务器配置
SERVER_NAME=Sourcebuild Pm
SERVER_VERSION=2.0.0
HTTP_ADDR=localhost
HTTP_PORT=9999

# 索引配置
PROJECT_INDEX=mcp_pm_projects
DOCUMENT_INDEX=mcp_pm_documents
```

#### Elasticsearch配置
```json
{
  "settings": {
    "number_of_shards": 1,
    "number_of_replicas": 1,
    "refresh_interval": "5s"
  },
  "mappings": {
    "properties": {
      "project_id": {"type": "keyword"},
      "name": {
        "type": "text",
        "analyzer": "standard",
        "fields": {
          "keyword": {"type": "keyword"}
        }
      },
      "description": {
        "type": "text",
        "analyzer": "standard"
      },
      "api_content": {
        "type": "object",
        "dynamic": true
      },
      "text_content": {
        "type": "object",
        "dynamic": true
      },
      "tags": {"type": "keyword"},
      "created_at": {"type": "date"},
      "updated_at": {"type": "date"}
    }
  }
}
```

## 📚 使用指南

### 基础操作流程

#### 1. 项目初始化
```bash
# 查看所有可访问的项目
list_projects

# 确认当前项目信息
get_project_info
```

#### 2. API文档管理
```bash
# 创建用户管理API文档
create_api_document \
  --name="user-management-api" \
  --method="POST" \
  --path="/api/v1/users" \
  --description="用户创建和管理接口" \
  --tags="user,management,public" \
  --request='{"username": "string", "email": "string", "password": "string"}' \
  --response='{"code": 200, "result": {"id": "user123"}, "msg": "创建成功"}'

# 搜索特定API文档
search_api_documents --query="用户管理" --method="POST"

# 列出所有GET方式接口
list_api_documents --method="GET"
```

#### 3. 文本文档管理
```bash
# 创建项目README
create_text_document \
  --name="project-readme" \
  --content_type="readme" \
  --description="项目说明文档" \
  --tags="readme,documentation,onboarding" \
  --content="# 项目说明\n\n这是一个企业级项目..."

# 创建配置模板
create_text_document \
  --name="production-config" \
  --content_type="config" \
  --description="生产环境配置模板" \
  --tags="config,production" \
  --content='{"debug": false, "port": 8080, "database": {"host": "prod-db.company.com"}}'

# 搜索配置文档
search_text_documents --content_type="config"
```

## 🔧 开发指南

### 项目结构详解
```
mcp-pm/
├── cmd/                          # 应用程序入口
│   └── main.go                   # 主程序启动逻辑
├── internal/                     # 内部包（不对外暴露）
│   ├── server.go                 # MCP服务器配置和路由
│   ├── config/                   # 配置管理
│   │   └── config.go             # 环境变量和配置结构
│   ├── middleware/               # 中间件层
│   │   └── logging_middleware.go # 日志记录中间件
│   ├── pm/                       # 项目管理器
│   │   └── project_manager.go    # 项目管理核心逻辑
│   ├── logic/                    # 业务逻辑处理器
│   │   ├── project_handlers.go      # 项目相关操作
│   │   ├── api_document_handlers.go # API文档操作
│   │   └── text_document_handlers.go # 文本文档操作
│   ├── es/                       # Elasticsearch客户端
│   │   ├── client.go             # ES客户端封装
│   │   └── models.go             # ES数据模型
│   ├── token/                    # Token管理
│   │   └── jwt.go                # JWT Token处理
│   ├── types/                    # 数据类型定义
│   │   ├── document.go           # 文档相关类型
│   │   └── models.go             # 业务模型
│   ├── utils/                    # 工具函数
│   │   └── extra.go              # 通用工具函数
│   └── cache/                    # 缓存层（预留）
│       └── redis.go              # Redis缓存实现
├── .env                          # 环境变量配置
├── .env.example                 # 环境变量示例
├── go.mod                       # Go模块依赖
├── go.sum                       # 依赖版本锁定
└── README.md                    # 项目说明文档
```

### 技术栈

#### 核心依赖
- **github.com/modelcontextprotocol/go-sdk v1.1.0**: MCP协议实现
- **github.com/olivere/elastic/v7 v7.0.32**: Elasticsearch客户端
- **github.com/golang-jwt/jwt/v5 v5.2.2**: JWT认证
- **github.com/joho/godotenv v1.5.1**: 环境变量管理
- **github.com/redis/go-redis/v9 v9.16.0**: Redis客户端（预留）

#### 开发工具
- **Go 1.25.0**: 主要编程语言
- **Elasticsearch 7.0+**: 搜索引擎和数据存储
- **Redis**: 缓存层（功能预留）

## 🐛 故障排除

### 常见问题及解决方案

#### 连接问题
**问题**: `failed to connect to Elasticsearch`
```bash
# 检查ES服务状态
curl -X GET "localhost:9200/_cluster/health"

# 检查网络连通性
telnet localhost 9200

# 验证配置
echo $ES_URL
```

#### 认证问题
**问题**: `authentication failed`
```bash
# 检查Token格式
echo $AUTH_TOKEN | jq .

# 验证Token有效性
curl -H "X-API-Token: $AUTH_TOKEN" http://localhost:9999/mcp
```

#### 搜索问题
**问题**: `search_phase_execution_exception`
```bash
# 检查ES索引映射
curl -X GET "localhost:9200/mcp_pm_documents/_mapping"

# 重建索引（如果需要）
curl -X DELETE "localhost:9200/mcp_pm_documents"
```

### 调试模式

#### 启用详细日志
```bash
# 开发环境调试
LOG_LEVEL=debug ./mcp-pm-mcp

# 生产环境监控
./mcp-pm-mcp 2>&1 | tee mcp-pm.log
```

#### 健康检查
```bash
# 检查服务状态
curl http://localhost:9999/mcp

# 检查ES连接
curl http://localhost:9200/_cluster/health
```

## 📊 监控与运维

### 关键指标监控

#### 应用指标
- **请求处理时间**: 平均响应时间和P99延迟
- **错误率**: 4xx和5xx错误比例
- **并发连接数**: 当前活跃连接数
- **内存使用**: 应用内存占用情况

#### 业务指标
- **文档创建数量**: 按类型统计的文档创建趋势
- **搜索请求量**: 搜索操作的QPS和成功率
- **用户活跃度**: 活跃用户和项目数量

#### 基础设施指标
- **Elasticsearch性能**: 索引大小、查询延迟、节点状态
- **系统资源**: CPU、内存、磁盘、网络使用率
- **日志分析**: 错误日志、访问日志、审计日志

### 运维建议

#### 性能优化
1. **Elasticsearch优化**: 合理设置分片数量和副本配置
2. **缓存策略**: 合理使用Redis缓存热点数据
3. **连接池**: 优化数据库和ES连接池配置
4. **索引优化**: 定期重建索引，提高查询性能

#### 安全加固
1. **网络安全**: 使用HTTPS和VPN保护数据传输
2. **访问控制**: 实施最小权限原则
3. **数据加密**: 敏感数据存储加密
4. **审计日志**: 记录所有关键操作日志

#### 备份策略
1. **数据备份**: 定期备份Elasticsearch数据
2. **配置备份**: 备份应用配置和环境变量
3. **灾难恢复**: 制定详细的灾难恢复计划

## 🤝 贡献指南

### 开发环境搭建
```bash
# 1. Fork并克隆代码库
git clone https://github.com/your-username/mcp-pm
cd mcp-pm

# 2. 创建开发分支
git checkout -b feature/your-feature-name

# 3. 安装依赖
go mod tidy

# 4. 启动开发环境
# 配置本地ES实例
# 设置环境变量
go run ./cmd
```

### 代码规范
- 遵循Go官方代码规范
- 使用gofmt格式化代码
- 编写单元测试和集成测试
- 添加必要的注释和文档

### 提交流程
1. 创建功能分支
2. 编写代码和测试
3. 提交Pull Request
4. 代码审查
5. 合并到主分支

## 📄 许可证

本项目采用 [MIT 许可证](LICENSE)。

## 🙋‍♂️ 支持

如果您在使用过程中遇到问题，可以通过以下方式获取帮助：

- 📧 邮箱: support@sourcebuild.cn
- 🐛 问题反馈: [GitHub Issues](https://github.com/source-build/mcp-pm/issues)
- 📖 文档: [项目文档](https://docs.sourcebuild.cn)

---

**源构互联** - 让开发更高效 🚀

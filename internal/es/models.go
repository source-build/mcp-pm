package es

import (
	"context"
	"fmt"
	"time"

	"github.com/source-build/mcp-pm/internal/config"
	"github.com/source-build/mcp-pm/internal/types"
)

// ProjectES 项目ES模型
type ProjectES struct {
	// 项目ID
	ID string `json:"id,omitempty"`
	// 项目名称
	Name string `json:"name,omitempty"`
	// 项目描述
	Description string `json:"description,omitempty"`
	// 用户ID列表
	UserIds []string `json:"user_ids,omitempty"`
	// 创建时间
	CreatedAt *time.Time `json:"created_at,omitempty"`
	// 更新时间
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	// 项目元数据
	Metadata map[string]string `json:"metadata,omitempty"`
}

// DocumentES 文档ES模型
type DocumentES struct {
	// 文档ID
	ID string `json:"id,omitempty"`
	// 项目ID
	ProjectID string `json:"project_id,omitempty"`
	// 创建用户ID
	CreatorID string `json:"creator_id,omitempty"`
	// 文档类型
	Type string `json:"type,omitempty"`
	// 文档名称
	Name string `json:"name,omitempty"`
	// 文档描述
	Description string `json:"description,omitempty"`
	// API文档内容
	ApiContent types.APIDocumentContent `json:"api_content,omitempty"`
	// 文本文档内容
	TextContent interface{} `json:"text_content,omitempty"`
	// 文档标签
	Tags []string `json:"tags,omitempty"`
	// 创建时间
	CreatedAt *time.Time `json:"created_at,omitempty"`
	// 更新时间
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	// 文档元数据
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// InitProjectESMappingJson 初始化项目ES索引结构
func InitProjectESMappingJson() error {
	mappingJson := map[string]interface{}{
		"settings": map[string]interface{}{
			"number_of_shards":   1,
			"number_of_replicas": 0,
			"refresh_interval":   "5s",
			"analysis": map[string]interface{}{
				"analyzer": map[string]interface{}{
					"ik_max_word": map[string]interface{}{
						"type": "ik_max_word",
					},
					"ik_smart": map[string]interface{}{
						"type": "ik_smart",
					},
				},
			},
		},
		"mappings": map[string]interface{}{
			"properties": map[string]interface{}{
				"id": map[string]interface{}{
					"type": "keyword",
				},
				"name": map[string]interface{}{
					"type":            "text",
					"analyzer":        "ik_max_word",
					"search_analyzer": "ik_smart",
					"fields": map[string]interface{}{
						"keyword": map[string]interface{}{
							"type": "keyword",
						},
					},
				},
				"description": map[string]interface{}{
					"type":     "text",
					"analyzer": "ik_max_word",
					"index":    false, // 描述不参与搜索，减少索引大小
				},
				"user_ids": map[string]interface{}{
					"type": "keyword",
				},
				"created_at": map[string]interface{}{
					"type":   "date",
					"format": "strict_date_optional_time||epoch_millis",
				},
				"updated_at": map[string]interface{}{
					"type":   "date",
					"format": "strict_date_optional_time||epoch_millis",
				},
				"metadata": map[string]interface{}{
					"type":    "object",
					"dynamic": true,
					"enabled": false, // 元数据不参与搜索，只存储
				},
			},
		},
	}

	_, err := ESClient.client.CreateIndex(config.Config.ProjectIndex).BodyJson(mappingJson).Do(context.Background())
	if err != nil {
		return fmt.Errorf("初始化项目ES索引失败: %v", err)
	}

	fmt.Println("初始化项目ES索引成功")
	return nil
}

// InitDocumentESMappingJson 初始化文档ES索引结构
func InitDocumentESMappingJson() error {
	mappingJson := map[string]interface{}{
		"settings": map[string]interface{}{
			"number_of_shards":   1,
			"number_of_replicas": 0,
			"refresh_interval":   "5s",
			"analysis": map[string]interface{}{
				"analyzer": map[string]interface{}{
					"ik_max_word": map[string]interface{}{
						"type": "ik_max_word",
					},
					"ik_smart": map[string]interface{}{
						"type": "ik_smart",
					},
				},
			},
		},
		"mappings": map[string]interface{}{
			"properties": map[string]interface{}{
				"id": map[string]interface{}{
					"type": "keyword",
				},
				"project_id": map[string]interface{}{
					"type": "keyword",
				},
				"creator_id": map[string]interface{}{
					"type": "keyword",
				},
				"type": map[string]interface{}{
					"type": "keyword",
				},
				"name": map[string]interface{}{
					"type":            "text",
					"analyzer":        "ik_max_word",
					"search_analyzer": "ik_smart",
					"fields": map[string]interface{}{
						"keyword": map[string]interface{}{
							"type": "keyword",
						},
					},
				},
				"description": map[string]interface{}{
					"type":            "text",
					"analyzer":        "ik_max_word",
					"search_analyzer": "ik_smart",
				},
				"api_content": map[string]interface{}{
					"type": "nested",
					"properties": map[string]interface{}{
						"method": map[string]interface{}{
							"type": "keyword",
						},
						"path": map[string]interface{}{
							"type": "keyword",
						},
						"body": map[string]interface{}{
							"type":    "object",
							"dynamic": true,
						},
						"query": map[string]interface{}{
							"type":    "object",
							"dynamic": true,
						},
						"path_params": map[string]interface{}{
							"type":    "object",
							"dynamic": true,
						},
						"response_biz_code": map[string]interface{}{
							"type": "keyword",
						},
						"header": map[string]interface{}{
							"type":    "object",
							"dynamic": true,
						},
					},
				},
				"text_content": map[string]interface{}{
					"type":    "object",
					"dynamic": true,
				},
				"tags": map[string]interface{}{
					"type": "keyword",
				},
				"created_at": map[string]interface{}{
					"type":   "date",
					"format": "strict_date_optional_time||epoch_millis",
				},
				"updated_at": map[string]interface{}{
					"type":   "date",
					"format": "strict_date_optional_time||epoch_millis",
				},
				"metadata": map[string]interface{}{
					"type":    "object",
					"dynamic": true,
					"enabled": false, // 元数据不参与搜索，只存储
				},
			},
		},
	}

	_, err := ESClient.client.CreateIndex(config.Config.DocumentIndex).BodyJson(mappingJson).Do(context.Background())
	if err != nil {
		return fmt.Errorf("初始化文档ES索引失败: %v", err)
	}

	fmt.Println("初始化文档ES索引成功")
	return nil
}

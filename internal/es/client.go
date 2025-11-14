package es

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/olivere/elastic/v7"
	"github.com/source-build/mcp-pm/internal/config"
	"github.com/source-build/mcp-pm/internal/types"
)

var ESClient *Client

// Client Elasticsearch客户端封装
type Client struct {
	client *elastic.Client
	ctx    context.Context
}

// InitESClient 初始化Elasticsearch客户端
func InitESClient() error {
	client, err := elastic.NewClient(
		elastic.SetURL(config.Config.EsURL),
		elastic.SetBasicAuth(config.Config.EsUsername, config.Config.EsPassword),
		elastic.SetSniff(config.Config.EsSniff),
		elastic.SetHealthcheckInterval(10*time.Second),
		elastic.SetErrorLog(log.New(log.Writer(), "ELASTIC_ERROR ", log.LstdFlags)),
		elastic.SetInfoLog(log.New(log.Writer(), "ELASTIC_INFO ", log.LstdFlags)),
	)
	if err != nil {
		return fmt.Errorf("创建Elasticsearch客户端失败: %v", err)
	}

	// 测试连接
	ctx := context.Background()
	info, code, err := client.Ping(config.Config.EsURL).Do(ctx)
	if err != nil {
		return fmt.Errorf("elasticsearch连接测试失败: %v", err)
	}

	log.Printf("Elasticsearch连接成功: 版本=%s, 状态码=%d", info.Version.Number, code)

	ESClient = &Client{
		client: client,
		ctx:    ctx,
	}

	//初始化索引
	//if err = ESClient.InitIndices(); err != nil {
	//	return fmt.Errorf("初始化索引失败: %v", err)
	//}

	return nil
}

// InitIndices 初始化索引
func (c *Client) InitIndices() error {
	// 创建项目索引
	//if err := InitProjectESMappingJson(); err != nil {
	//	return fmt.Errorf("创建项目索引失败: %v", err)
	//}

	// 创建文档索引
	if err := InitDocumentESMappingJson(); err != nil {
		return fmt.Errorf("创建文档索引失败: %v", err)
	}

	return nil
}

// createProjectIndex 创建项目索引
func (c *Client) createProjectIndex() error {
	indexName := config.Config.ProjectIndex

	// 检查索引是否存在
	exists, err := c.client.IndexExists(indexName).Do(c.ctx)
	if err != nil {
		return fmt.Errorf("检查项目索引失败: %v", err)
	}

	if exists {
		log.Printf("项目索引 %s 已存在", indexName)
		return nil
	}

	// 定义索引映射
	mapping := map[string]interface{}{
		"settings": map[string]interface{}{
			"number_of_shards":   1,
			"number_of_replicas": 0,
		},
		"mappings": map[string]interface{}{
			"properties": map[string]interface{}{
				"id": map[string]interface{}{
					"type": "keyword",
				},
				"name": map[string]interface{}{
					"type": "text",
					"fields": map[string]interface{}{
						"keyword": map[string]interface{}{
							"type": "keyword",
						},
					},
				},
				"description": map[string]interface{}{
					"type": "text",
				},
				"user_ids": map[string]interface{}{
					"type": "keyword",
				},
				"created_at": map[string]interface{}{
					"type": "date",
				},
				"updated_at": map[string]interface{}{
					"type": "date",
				},
				"metadata": map[string]interface{}{
					"type":    "object",
					"dynamic": true,
				},
			},
		},
	}

	// 创建索引
	_, err = c.client.CreateIndex(indexName).BodyJson(mapping).Do(c.ctx)
	if err != nil {
		return fmt.Errorf("创建项目索引失败: %v", err)
	}

	log.Printf("项目索引 %s 创建成功", indexName)
	return nil
}

// GetProject 获取项目
func (c *Client) GetProject(id string) (*types.Project, error) {
	get, err := c.client.Get().
		Index(config.Config.ProjectIndex).
		Id(id).
		Do(c.ctx)

	if err != nil {
		if elastic.IsNotFound(err) {
			return nil, fmt.Errorf("项目不存在: %s", id)
		}
		return nil, fmt.Errorf("获取项目失败: %v", err)
	}

	var project types.Project
	if err = json.Unmarshal(get.Source, &project); err != nil {
		return nil, fmt.Errorf("反序列化项目失败: %v", err)
	}

	return &project, nil
}

// SaveProject 保存项目
func (c *Client) SaveProject(project *types.Project) error {
	now := time.Now()
	project.UpdatedAt = &now
	if project.CreatedAt.IsZero() {
		project.CreatedAt = &now
	}

	_, err := c.client.Index().
		Index(config.Config.ProjectIndex).
		Id(project.ID).
		BodyJson(project).
		Refresh("wait_for").
		Do(c.ctx)

	if err != nil {
		return fmt.Errorf("保存项目失败: %v", err)
	}

	log.Printf("项目保存成功: %s", project.ID)
	return nil
}

// FindProjectsByToken 根据用户ID查找项目
func (c *Client) FindProjectsByToken(userIds string) ([]*types.Project, error) {
	query := elastic.NewTermQuery("user_ids", userIds)
	searchResult, err := c.client.Search().
		Index(config.Config.ProjectIndex).
		Query(query).
		Sort("updated_at", false).
		Do(c.ctx)

	if err != nil {
		return nil, fmt.Errorf("搜索项目失败: %v", err)
	}

	var projects []*types.Project
	for _, hit := range searchResult.Hits.Hits {
		var project types.Project
		if err = json.Unmarshal(hit.Source, &project); err != nil {
			log.Printf("反序列化项目失败: %v", err)
			continue
		}
		projects = append(projects, &project)
	}

	return projects, nil
}

// SaveDocument 保存文档
func (c *Client) SaveDocument(document *types.Document) (id string, err error) {
	now := time.Now()
	document.UpdatedAt = &now
	if document.CreatedAt == nil {
		document.CreatedAt = &now
	}

	doc := DocumentES{
		ProjectID:   document.ProjectID,
		CreatorID:   document.CreatorID,
		Type:        string(document.Type),
		Name:        document.Name,
		Description: document.Description,
		ApiContent:  document.APIContent,
		TextContent: document.TextContent,
		Tags:        document.Tags,
		CreatedAt:   document.CreatedAt,
		UpdatedAt:   document.UpdatedAt,
		Metadata:    document.Metadata,
	}
	resp, err := c.client.Index().
		Index(config.Config.DocumentIndex).
		Id(document.ID).
		BodyJson(doc).
		Refresh("wait_for").
		Do(c.ctx)
	if err != nil {
		return "", fmt.Errorf("保存文档失败: %v", err)
	}

	return resp.Id, nil
}

// GetDocument 获取文档
func (c *Client) GetDocument(id string) (*types.Document, error) {
	get, err := c.client.Get().
		Index(config.Config.DocumentIndex).
		Id(id).
		Do(c.ctx)

	if err != nil {
		if elastic.IsNotFound(err) {
			return nil, fmt.Errorf("文档不存在: %s", id)
		}
		return nil, fmt.Errorf("获取文档失败: %v", err)
	}

	var document types.Document
	if err = json.Unmarshal(get.Source, &document); err != nil {
		return nil, fmt.Errorf("反序列化文档失败: %v", err)
	}

	return &document, nil
}

// SearchDocuments 搜索文档
func (c *Client) SearchDocuments(projectID string, req *types.SearchRequest) (*types.SearchResult, error) {
	boolQuery := elastic.NewBoolQuery()

	// 项目过滤
	boolQuery.Must(elastic.NewTermQuery("project_id", projectID))

	// 文档类型过滤
	if req.Type != nil {
		boolQuery.Must(elastic.NewTermQuery("type", string(*req.Type)))
	}

	// 标签过滤
	if len(req.Tags) > 0 {
		boolQuery.Must(elastic.NewTermsQuery("tags", req.Tags))
	}

	// 关键词搜索
	if req.Query != "" {
		boolQuery.Must(elastic.NewMultiMatchQuery(req.Query, "name", "description"))
	}

	// 设置默认值
	limit := req.Limit
	if limit <= 0 {
		limit = 20
	}
	offset := req.Offset
	if offset < 0 {
		offset = 0
	}

	// 执行搜索
	searchResult, err := c.client.Search().
		Index(config.Config.DocumentIndex).
		Query(boolQuery).
		From(offset).
		Size(limit).
		Sort("updated_at", false).
		Do(c.ctx)
	if err != nil {
		return nil, fmt.Errorf("搜索文档失败: %v", err)
	}
	var documents []types.Document
	for _, hit := range searchResult.Hits.Hits {
		var document types.Document
		if err = json.Unmarshal(hit.Source, &document); err != nil {
			log.Printf("反序列化文档失败: %v", err)
			continue
		}
		documents = append(documents, document)
	}

	return &types.SearchResult{
		Documents: documents,
		Total:     searchResult.TotalHits(),
		Offset:    offset,
		Limit:     limit,
	}, nil
}

// ListDocuments 列出项目文档
func (c *Client) ListDocuments(projectID string, docType *types.DocumentType, limit, offset int) (*types.SearchResult, error) {
	var queries []elastic.Query

	// 项目过滤
	projectQuery := elastic.NewTermQuery("project_id", projectID)
	queries = append(queries, projectQuery)

	// 文档类型过滤
	if docType != nil {
		typeQuery := elastic.NewTermQuery("type", string(*docType))
		queries = append(queries, typeQuery)
	}

	// 组合查询
	var query elastic.Query
	if len(queries) == 1 {
		query = queries[0]
	} else {
		query = elastic.NewBoolQuery().Must(queries...)
	}

	// 设置默认值
	if limit <= 0 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	// 执行搜索
	searchResult, err := c.client.Search().
		Index(config.Config.DocumentIndex).
		Query(query).
		From(offset).
		Size(limit).
		Sort("updated_at", false).
		Do(c.ctx)

	if err != nil {
		return nil, fmt.Errorf("列出文档失败: %v", err)
	}

	var documents []types.Document
	for _, hit := range searchResult.Hits.Hits {
		var document types.Document
		if err := json.Unmarshal(hit.Source, &document); err != nil {
			log.Printf("反序列化文档失败: %v", err)
			continue
		}
		documents = append(documents, document)
	}

	return &types.SearchResult{
		Documents: documents,
		Total:     searchResult.TotalHits(),
		Offset:    offset,
		Limit:     limit,
	}, nil
}

// Close 关闭客户端
func (c *Client) Close() error {
	if c.client != nil {
		c.client.Stop()
	}
	return nil
}

package es

import (
	"Mail-Achive/pkg/config"
	"Mail-Achive/pkg/log"
	"Mail-Achive/pkg/model"
	"Mail-Achive/pkg/utils"
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/olivere/elastic/v7"
)

const (
	// defaultElasticHealthCheck - health check time for elastic
	defaultElasticHealthCheck = 60

	// he index mapping string
	mapping = `{
		"settings":{
			"number_of_shards":1,
			"number_of_replicas":0
		},
		"mappings":{
			"properties":{
				"name":{
					"type":"text"
				},
				"from":{
					"type":"text"
				},
				"created_at":{
					"type":"date"
				},
				"title":{
					"type":"text"
				},
				"content": {
					"type":"text"
				}
			}
		}
	}`
)

// Handler for elastic search
type Handler struct {
	document string
	url      string
	fields   []string
	client   *elastic.Client
	ctx      context.Context
}

// NewHandler returns a handler for elastic
func NewHandler(settings *config.Config) *Handler {
	s := &Handler{
		document: settings.Server.DocumentName,
		fields:   settings.Server.MatchFields,
		url:      settings.Server.ElasticURL,
		ctx:      context.Background(),
	}

	// init the elastic client
	client, err := elastic.NewClient(
		elastic.SetURL(s.url),
		elastic.SetSniff(false),
		elastic.SetHealthcheck(true),
		elastic.SetHealthcheckInterval(time.Second*defaultElasticHealthCheck),
	)
	if err != nil {
		panic(err)
	}
	s.client = client

	// getting the ES version number is quite common, so there's a shortcut
	version, err := client.ElasticsearchVersion(s.url)
	if err != nil {
		panic(err)
	}
	log.Infof("Elasticsearch version: %v", version)

	// create the index if the specified index is not existed
	if err := s.CreateIndex(); err != nil {
		panic(err)
	}

	return s
}

// Close - release the connections
func (s *Handler) Close() {
	return
}

// Bulk the emails to elastic
func (s *Handler) Bulk(emails []*model.Email) error {
	if len(emails) == 0 {
		return nil
	}

	// upset the email content into elastic
	bulk := s.client.Bulk().Index(s.document)
	for _, email := range emails {
		// md5 the content as the id
		id := utils.MD5Str(email.Content)
		bulk.Add(elastic.NewBulkIndexRequest().Id(id).Doc(email))
	}

	if _, err := bulk.Do(s.ctx); err != nil {
		log.Errorf("upset the elastic document: %v", err)
	}

	return nil
}

// Search the email by multiple match fields
func (s *Handler) Search(ctx context.Context, query string, skip int, take int) ([]*model.Email, error) {
	// new a multiple match query
	match := elastic.NewMultiMatchQuery(query, s.fields...)

	// do the search with the match query
	result, err := s.client.Search().
		Index().
		Query(match).
		Sort("created_at", true).
		From(skip).Size(take).
		Pretty(true).
		Do(ctx)
	if err != nil {
		return nil, err
	}
	log.Infof("Query tooks %d milliseconds", result.TookInMillis)

	emails := []*model.Email{}
	// unmarshal the search result
	for _, hit := range result.Hits.Hits {
		var email model.Email
		if err := json.Unmarshal(hit.Source, &email); err != nil {
			return nil, err
		}
		emails = append(emails, &email)
	}

	return emails, nil
}

// CreateIndex checks if the index is existed, if not, create a new index
func (s *Handler) CreateIndex() error {
	// check if the specified index exists
	exists, err := s.client.IndexExists(s.document).Do(s.ctx)
	if err != nil {
		return err
	}
	if !exists {
		// create the index
		createIndex, err := s.client.CreateIndex(s.document).BodyString(mapping).Do(s.ctx)
		if err != nil {
			return err
		}
		if !createIndex.Acknowledged {
			return errors.New("Not acknowledged for CreateIndex")
		}
		log.Infof("Elastic index %s is created", s.document)
	}

	return nil
}

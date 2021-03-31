// Create file in v.1.0.0
// srvcheck_elasticsearch_repo.go is file that define implement elasticsearch history repository using elasticsearch
// this elasticsearch repository struct embed esRepositoryRequiredComponent struct in ./srvcheck.go file

package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/DMS-SMS/v1-health-check/domain"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/pkg/errors"
	"log"
	"net/http"
	"time"
)

// esElasticsearchCheckHistoryRepository is to handle ElasticsearchCheckHistoryRepository model using elasticsearch as data store
type esElasticsearchCheckHistoryRepository struct {
	// myCfg is used for get elasticsearch history repository config about elasticsearch
	myCfg esElasticsearchCheckHistoryRepoConfig

	// esCli is elasticsearch client connection injected from the outside package
	esCli *elasticsearch.Client

	// reqBodyWriter is implementation of reqBodyWriter interface to write []byte for request body
	reqBodyWriter reqBodyWriter
}

// esElasticsearchCheckHistoryRepoConfig is the config for elasticsearch check history repository using elasticsearch
type esElasticsearchCheckHistoryRepoConfig interface {
	// get common method from embedding esRepositoryComponentConfig
	esRepositoryComponentConfig
}

// NewESElasticsearchCheckHistoryRepository return new object that implement ElasticsearchCheckHistoryRepository interface
func NewESElasticsearchCheckHistoryRepository(
	cfg esElasticsearchCheckHistoryRepoConfig,
	cli *elasticsearch.Client,
	w reqBodyWriter,
) domain.ElasticsearchCheckHistoryRepository {
	repo := &esElasticsearchCheckHistoryRepository{
		myCfg:         cfg,
		esCli:         cli,
		reqBodyWriter: w,
	}

	if err := repo.Migrate(); err != nil {
		log.Fatal(errors.Wrap(err, "could not migrate repository").Error())
	}

	return repo
}

// Implement Migrate method of ElasticsearchCheckHistoryRepository interface
func (eer *esElasticsearchCheckHistoryRepository) Migrate() error {
	resp, err := (esapi.IndicesExistsRequest{
		Index: []string{eer.myCfg.IndexName()},
	}).Do(context.Background(), eer.esCli)

	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to call IndicesExists, resp: %+v", resp))
	}

	if resp.StatusCode == http.StatusNotFound {
		if err := eer.createIndex(); err != nil {
			return errors.Wrap(err, "failed to create index")
		}
	}

	return nil
}

// createIndex method create index with name, share number in esElasticsearchCheckHistoryRepository
func (eer *esElasticsearchCheckHistoryRepository) createIndex() error {
	body := map[string]interface{}{}
	body["settings.number_of_shards"] = eer.myCfg.IndexShardNum()
	body["settings.number_of_replicas"] = eer.myCfg.IndexReplicaNum()

	b, _ := json.Marshal(body)
	if _, err := eer.reqBodyWriter.Write(b); err != nil {
		return errors.Wrap(err, "failed to write map to body writer")
	}

	buf := &bytes.Buffer{}
	if _, err := eer.reqBodyWriter.WriteTo(buf); err != nil {
		return errors.Wrap(err, "failed to body writer WriteTo method")
	}

	resp, err := (esapi.IndicesCreateRequest{
		Index:         eer.myCfg.IndexName(),
		Body:          bytes.NewReader(buf.Bytes()),
		MasterTimeout: time.Second * 5,
		Timeout:       time.Second * 5,
	}).Do(context.Background(), eer.esCli)

	return errors.Wrap(err, fmt.Sprintf("failed to call IndicesCreate, resp: %+v", resp))
}

// Implement Store method of ElasticsearchCheckHistoryRepository interface
func (eer *esElasticsearchCheckHistoryRepository) Store(history *domain.ElasticsearchCheckHistory) (b []byte, err error) {
	body, _ := json.Marshal(history.DottedMapWithPrefix(""))
	if _, err = eer.reqBodyWriter.Write(body); err != nil {
		err = errors.Wrap(err, "failed to write map to body writer")
		return
	}

	buf := &bytes.Buffer{}
	if _, err = eer.reqBodyWriter.WriteTo(buf); err != nil {
		err = errors.Wrap(err, "failed to body writer WriteTo method")
		return
	}

	resp, err := (esapi.IndexRequest{
		Index:        eer.myCfg.IndexName(),
		Body:         bytes.NewReader(buf.Bytes()),
		Timeout:      time.Second * 5,
	}).Do(context.Background(), eer.esCli)

	if err != nil {
		err = errors.Wrap(err, fmt.Sprintf("failed to call IndexRequest, resp: %+v", resp))
		return
	} else if resp.IsError() {
		err = errors.Errorf("IndexRequest return error code, resp: %+v", resp)
		return
	}

	result := map[string]interface{}{}
	_ = json.NewDecoder(resp.Body).Decode(&result)
	b, _ = json.Marshal(result)
	return
}

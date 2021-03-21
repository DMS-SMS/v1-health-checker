// Create file in v.1.0.0
// syscheck_cpu_repo.go is file that define implement cpu history repository using elasticsearch
// this cpu repository struct embed esRepositoryRequiredComponent struct in ./syscheck.go file

package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/pkg/errors"
	"log"
	"net/http"
	"time"

	"github.com/DMS-SMS/v1-health-check/domain"
)

// esCPUCheckHistoryRepository is to handle CPUCheckHistory model using elasticsearch as data store
type esCPUCheckHistoryRepository struct {
	// myCfg is used for get cpu check history repository config about elasticsearch
	myCfg esCPUCheckHistoryRepoConfig

	// esCli is elasticsearch client connection injected from the outside package
	esCli *elasticsearch.Client

	// bodyWriter is implementation of reqBodyWriter interface to write []byte for request body
	bodyWriter reqBodyWriter
}

// esCPUCheckHistoryRepoConfig is the config for cpu check history repository using elasticsearch
type esCPUCheckHistoryRepoConfig interface {
	// get common method from embedding esRepositoryComponentConfig
	esRepositoryComponentConfig
}

// NewESCPUCheckHistoryRepository return new object that implement CPUCheckHistoryRepository interface
func NewESCPUCheckHistoryRepository(cfg esCPUCheckHistoryRepoConfig, cli *elasticsearch.Client, w reqBodyWriter) domain.CPUCheckHistoryRepository {
	repo := &esCPUCheckHistoryRepository{
		myCfg:      cfg,
		esCli:      cli,
		bodyWriter: w,
	}

	if err := repo.Migrate(); err != nil {
		log.Fatal(errors.Wrap(err, "could not migrate repository").Error())
	}

	return repo
}

// Implement Migrate method of CPUCheckHistoryRepository interface
func (esr *esCPUCheckHistoryRepository) Migrate() error {
	resp, err := (esapi.IndicesExistsRequest{
		Index: []string{esr.myCfg.IndexName()},
	}).Do(context.Background(), esr.esCli)

	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to call IndicesExists, resp: %+v", resp))
	}

	if resp.StatusCode == http.StatusNotFound {
		if err := esr.createIndex(); err != nil {
			return errors.Wrap(err, "failed to create index")
		}
	}

	return nil
}

// createIndex method create index with name, share number in esCPUCheckHistoryRepository
func (esr *esCPUCheckHistoryRepository) createIndex() error {
	body := map[string]interface{}{}
	body["settings.number_of_shards"] = esr.myCfg.IndexShardNum()
	body["settings.number_of_replicas"] = esr.myCfg.IndexReplicaNum()

	b, _ := json.Marshal(body)
	if _, err := esr.bodyWriter.Write(b); err != nil {
		return errors.Wrap(err, "failed to write map to body writer")
	}

	buf := &bytes.Buffer{}
	if _, err := esr.bodyWriter.WriteTo(buf); err != nil {
		return errors.Wrap(err, "failed to body writer WriteTo method")
	}

	resp, err := (esapi.IndicesCreateRequest{
		Index:         esr.myCfg.IndexName(),
		Body:          bytes.NewReader(buf.Bytes()),
		MasterTimeout: time.Second * 5,
		Timeout:       time.Second * 5,
	}).Do(context.Background(), esr.esCli)

	return errors.Wrap(err, fmt.Sprintf("failed to call IndicesCreate, resp: %+v", resp))
}

// Implement Store method of CPUCheckHistoryRepository interface
func (esr *esCPUCheckHistoryRepository) Store(history *domain.CPUCheckHistory) (b []byte, err error) {
	body, _ := json.Marshal(history.DottedMapWithPrefix(""))
	if _, err = esr.bodyWriter.Write(body); err != nil {
		err = errors.Wrap(err, "failed to write map to body writer")
		return
	}

	buf := &bytes.Buffer{}
	if _, err = esr.bodyWriter.WriteTo(buf); err != nil {
		err = errors.Wrap(err, "failed to body writer WriteTo method")
		return
	}

	resp, err := (esapi.IndexRequest{
		Index:        esr.myCfg.IndexName(),
		Body:         bytes.NewReader(buf.Bytes()),
		Timeout:      time.Second * 5,
	}).Do(context.Background(), esr.esCli)

	if err != nil || resp.IsError() {
		// 단순 err인지 status code상 오류인지 확인
		err = errors.Wrap(err, fmt.Sprintf("failed to call IndexRequest, resp: %+v", resp))
		return
	}

	result := map[string]interface{}{}
	_ = json.NewDecoder(resp.Body).Decode(&result)
	b, _ = json.Marshal(result)
	return
}

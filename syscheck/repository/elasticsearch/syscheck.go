// Create package in v.1.0.0
// elasticsearch package is for implementations of syscheck domain repository using elasticsearch
// In practice, repository struct declaration and implementation is in this syscheck.go file

package elasticsearch

import (
	"github.com/elastic/go-elasticsearch/v7"
)

// esDiskCheckHistoryRepository is to handle DiskCheckHistory model using elasticsearch as data store
type esDiskCheckHistoryRepository struct {
	cli *elasticsearch.Client
}

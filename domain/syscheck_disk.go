// Create file in v.1.0.0
// syscheck_disk.go is file that declare model struct & repo interface about disk in syscheck domain.
// also, additional method of model struct is declared in this file, too.

package domain

import "context"

// DiskCheckHistory model is used for record disk health check history and result
type DiskCheckHistory struct {
	// get required component by embedding systemCheckHistoryComponent
	systemCheckHistoryComponent
}

// DiskCheckHistoryRepository is abstract method used in business layer
// Repository is implemented with elastic search in v.1.0.0
type DiskCheckHistoryRepository interface {
	// get required component by embedding systemCheckHistoryRepositoryComponent
	systemCheckHistoryRepositoryComponent

	// Store method save DiskCheckHistory model in repository
	// b in return represents bytes of response body(map[string]interface{})
	Store(*DiskCheckHistory) (b []byte, err error)
}

// DiskCheckUseCase is interface used as business process handler about disk check
type DiskCheckUseCase interface {
	// CheckDisk method check disk capacity status and store disk check history using repository
	CheckDisk(ctx context.Context) error
}

// MapWithPrefixKey convert DiskCheckHistory to dotted map and return using MapWithPrefixKey of upper struct
// all key value of Map start with prefix received from parameter
func (dh *DiskCheckHistory) MapWithPrefixKey(prefix string) (m map[string]interface{}) {
	m = dh.systemCheckHistoryComponent.MapWithPrefixKey(prefix)

	if prefix != "" {
		prefix += "."
	}

	// add field in map

	return
}
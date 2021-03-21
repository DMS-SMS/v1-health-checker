// Create file in v.1.0.0
// syscheck_cpu.go is file that declare model struct & repo interface about cpu health check in syscheck domain.
// also, additional method of model struct is declared in this file, too.

package domain

import "github.com/inhies/go-bytesize"

// CPUCheckHistory model is used for record cpu health check history and result
type CPUCheckHistory struct {
	// get required component by embedding systemCheckHistoryComponent
	systemCheckHistoryComponent

	// Usage specifies current cpu usage of runtime system looked in cpu check
	Usage bytesize.ByteSize

	// Free specifies freed cpu size while recovering cpu health
	Free bytesize.ByteSize
}
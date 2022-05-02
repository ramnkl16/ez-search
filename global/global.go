package global

import (
	"fmt"
	"os"
)

var (
	WorkingDir string
)

func init() {
	WorkingDir, _ = os.Getwd()
}

type StatusEnum int

const (
	EVENT_TYPE_HYBRIS_PRODUCT           = "HybrisProduct"
	EVENT_TYPE_HYBRIS_PRODUCT_SUMMARY   = "HybrisProductSummary"
	EVENT_TYPE_OCC_IMPORT               = "OccImport"
	EVENT_TYPE_OCC_IMPORT_STATUS_UPDATE = "OccImportStatusUpdate"
	EVENT_TYPE_OCC_IMPORT_LOGS_UPDATE   = "OccImportLogsUpdate"
	EVENT_TYPE_SC_SYNC                  = "scsync"
	EVENT_TYPE_OCC_SYNC                 = "occsync"
	EVENT_TYPE_DETETE_LOG               = "dellogs"

	INT_HYBRIS     = "hybris"
	INT_OCC        = "occ"
	INT_SC_MASTER  = "scm"
	INT_SC_WEB     = "scw"
	INT_SC_LIVEWEB = "sclw"

	STATUS_INACTIVE   StatusEnum = 0
	STATUS_ACTIVE     StatusEnum = 1
	STATUS_QUEUED     StatusEnum = 2
	STATUS_INPROGRESS StatusEnum = 3
	STATUS_SUSPEND    StatusEnum = 4
	STATUS_COMPLETED  StatusEnum = 5
	STATUS_ERROR      StatusEnum = 6
	STATUS_DELETED    StatusEnum = 7
)

func Find(slice []string, val string) (int, bool) {
	//fmt.Println("before for loop in attr", val)
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}

	return -1, false
}

func (e StatusEnum) String() string {
	switch e {
	case STATUS_INACTIVE:
		return "InActive"
	case STATUS_ACTIVE:
		return "Active"
	case STATUS_QUEUED:
		return "Queued"
	case STATUS_INPROGRESS:
		return "Inprogress"
	case STATUS_SUSPEND:
		return "Suspend"
	case STATUS_COMPLETED:
		return "Completed"
	case STATUS_ERROR:
		return "Error"
	case STATUS_DELETED:
		return "Deleted"
	default:
		return fmt.Sprintf("%d", int(e))
	}
}

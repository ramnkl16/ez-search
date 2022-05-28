package global

import (
	"fmt"
	"os"
)

var (
	WorkingDir        string
	CsvWatcherPath    string
	CsvFileExt        string
	MaxIndexbatchSize int
)

func init() {
	WorkingDir, _ = os.Getwd()
	CsvWatcherPath = "\\csvwatcher"
	CsvFileExt = ".csv"
}

type StatusEnum int

const (
	EVENT_TYPE_INDEXFROMCSV        = "indexfromcsv"
	EVENT_TYPE_MYSQL_SYNC          = "mysqlsync"
	EVENT_TYPE_MSSQL_SYNC          = "mssqlsync"
	EVENT_TYPE_AS400_SYNC          = "as400sync"
	EVENT_TYPE_WINEVENTVIEWER_SYNC = "wevsync"
	EVENT_TYPE_DETETE_LOG          = "dellogs"

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

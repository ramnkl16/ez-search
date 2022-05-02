package ezsearch

type Config struct {
	IndexBatchSize    int    `json:"indexBatchSize"`
	MaxResultsInaPage int    `json:"MaxResultsInaPage"`
	IndexBasePath     string `json:"indexBasePath"`
	IndexTablesPath   string `json:"indexTablesPath"`
	ApplogIndexPath   string `json:"applogIndexPath"`
	BoltDbBucketName  string
}

var (
	Conf Config
)

func SetConfig(c Config) {
	Conf = c
}

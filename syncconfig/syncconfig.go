package syncconfig

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/ramnkl16/ez-search/datasources/catalogboltdb"
	"github.com/ramnkl16/ez-search/logger"

	"github.com/ramnkl16/ez-search/ezsearch"
)

var (
	Gconfig *Config
)

type (
	Config struct {
		BoltDBSettings     catalogboltdb.Config `json:"boltdbSettings"`
		EzsearchSettings   ezsearch.Config      `json:"searchSettings"`
		LoggerSettings     logger.Config        `json:"loggerSettings"`
		CSVFileWatcherpath string               `json:"csvWatcher"`
		CSVFileExtension   string               `json:"csvFileExt"`
		DbName             string               `json:"dbName"`
		Port               string               `json:"port"`
		HostName           string               `json:"hostName"`
		WorkingDir         string               `json:"workingDir"`
		RunScheduler       bool
	}
)

func NewCofig() *Config {
	c := &Config{
		DbName:   "product-int-logs",
		Port:     "8015",
		HostName: "localhost",
	}
	c.WorkingDir, _ = os.Getwd()

	return c
}
func OpenConfig(configFile string) (*Config, error) {
	fmt.Println(fmt.Sprintf("config File: %v", configFile))
	c := NewCofig()
	f, err := os.Open(configFile)
	if err != nil {
		fmt.Println("unable to open config file", err)
		return nil, err
	}
	d := json.NewDecoder(f)
	err = d.Decode((c))
	if err != nil {
		fmt.Println("unable to parse config file", err)

	}
	return c, nil

}

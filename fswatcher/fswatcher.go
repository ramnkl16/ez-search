package fswatcher

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/adriansr/fsnotify"
	"github.com/ramnkl16/ez-search/ezeventqueue"
	"github.com/ramnkl16/ez-search/global"
	"github.com/ramnkl16/ez-search/logger"
	"github.com/ramnkl16/ez-search/models"
	"github.com/ramnkl16/ez-search/services"
	"github.com/ramnkl16/ez-search/utils/date_utils"
	"go.uber.org/zap/zapcore"
)

var watcher *fsnotify.Watcher

func watchDir(path string, fi os.FileInfo, err error) error {
	if fi.Mode().IsDir() {
		return watcher.Add(path)
	}

	return nil
}
func WatchCSVFiles(csvFilePath string) {

	watcher, _ = fsnotify.NewWatcher()
	defer watcher.Close()
	wpath := filepath.Join(global.WorkingDir, global.CsvWatcherPath)
	logger.Info("csvwatcher path", zapcore.Field{Key: "watcher", Type: zapcore.StringType, String: wpath})
	if err := filepath.Walk(wpath, watchDir); err != nil {
		logger.Error(err.Error(), err)
	}
	done := make(chan bool)

	go func() {
		for {
			select {
			case event := <-watcher.Events:
				processEvent(event)
			case err := <-watcher.Errors:
				logger.Error(err.Error(), err)
			}
		}
	}()
	<-done
}

///register event as csv index
func processEvent(event fsnotify.Event) {

	if strings.HasSuffix(event.Name, global.CsvFileExt) {

		fileName := filepath.Join(global.WorkingDir, global.CsvWatcherPath, event.Name)
		cd := ezeventqueue.CSVImportCustomData{
			FileName:  fileName,
			IndexName: event.Name,
		}
		cdByte, err := json.Marshal(cd)
		if err != nil {
			logger.Error("Failed while unmarshal|processEvent", err, zapcore.Field{Key: "fileName", Type: zapcore.StringType, String: fileName})
		}
		services.EventQueueCustomService.CreateWithIndex(models.EventQueue{EventType: global.EVENT_TYPE_INDEXFROMCSV,
			EventData: string(cdByte), Status: int(global.STATUS_ACTIVE), RetryCount: 0, StartAt: date_utils.GetNowSearchFormat(), IsActive: "t"})
	} else {
		logger.Warn("Invalid file extension", zapcore.Field{Key: "ext", Type: zapcore.StringType, String: event.Name})
	}

}

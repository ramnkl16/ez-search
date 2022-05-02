package logger

import (
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/natefinch/lumberjack"
	"github.com/ramnkl16/ez-search/global"
	"github.com/ramnkl16/ez-search/utils/uid_utils"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	envLogLevel   = "LOG_LEVEL"
	envLogOutput  = "LOG_OUTPUT"
	UTCDateLayout = "2006-01-02T15:04:05Z"
)

var (
	Log                 logger
	hasIndexDatePattern = true
	loggerPatternName   = ""
	sugarLogger         *zap.SugaredLogger
)

type ezLogger interface {
	Print(v ...interface{})
	Printf(format string, v ...interface{})
}

type logger struct {
	log *zap.Logger
}
type errorLog struct {
	Level string `json:"l"`
	Msg   string `json:"m"`
	Time  string `json:"t"`
}

func InitLogger() {
	if len(Conf.ApplogIndexPath) == 0 {
		Conf.ApplogIndexPath = "indexes/applogs"
	}

	//index pattern parsing
	grp := global.RegexParseDate.FindAllSubmatch([]byte(Conf.ApplogIndexPath), -1)
	if grp != nil {
		hasIndexDatePattern = true
		loggerPatternName = string(grp[0][1])
	}

	var core zapcore.Core
	conEncoder := getEncoder()
	if len(Conf.LogOutput) > 0 && Conf.EnableConsoleLog {
		writerSyncer := getLogWriter()
		conEncoder := getEncoder()
		core = zapcore.NewTee(
			zapcore.NewCore(conEncoder, writerSyncer, getLevel()),
			zapcore.NewCore(conEncoder, zapcore.AddSync(os.Stdout), getLevel()),
		)
	} else if len(Conf.LogOutput) > 0 && !Conf.EnableConsoleLog {
		writerSyncer := getLogWriter()
		core = zapcore.NewTee(
			zapcore.NewCore(conEncoder, writerSyncer, getLevel()),
		)
	} else {
		conEncoder := getEncoder()
		core = zapcore.NewTee(
			zapcore.NewCore(conEncoder, zapcore.AddSync(os.Stdout), getLevel()),
		)
	}

	Log.log = zap.New(core, zap.AddCaller())

	//sugarLogger = logger.Sugar()
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(UTCDateLayout)
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getLogWriter() zapcore.WriteSyncer {
	var confFile string
	if strings.ToLower(Conf.LogOutput) != "stdout" || !strings.Contains(Conf.LogOutput, ":") {
		confFile = path.Join(global.WorkingDir, Conf.LogOutput)
	}
	lumberJackLogger := &lumberjack.Logger{
		Filename:   confFile,
		MaxSize:    1,
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   false,
	}
	return zapcore.AddSync(lumberJackLogger)
}

// func Initialize() {
// 	// initialize the rotator
// 	logFile := "/var/log/app-%Y-%m-%d-%H.log"
// 	rotator, err := rotatelogs.New(
// 		logFile,
// 		rotatelogs.WithMaxAge(60*24*time.Hour),
// 		rotatelogs.WithRotationTime(time.Hour))
// 	if err != nil {
// 		panic(err)
// 	}

// 	// initialize the JSON encoding config
// 	encoderConfig := map[string]string{
// 		"levelEncoder": "l",
// 		"timeKey":      "m",
// 		"timeEncoder":  "d",
// 	}
// 	data, _ := json.Marshal(encoderConfig)
// 	var encCfg zapcore.EncoderConfig
// 	if err := json.Unmarshal(data, &encCfg); err != nil {
// 		panic(err)
// 	}

// 	// add the encoder config and rotator to create a new zap logger
// 	w := zapcore.AddSync(rotator)
// 	core := zapcore.NewCore(
// 		zapcore.NewJSONEncoder(encCfg),
// 		w,
// 		zap.InfoLevel)
// 	logger := zap.New(core)

// 	logger.Info("Now logging in a rotated file")
// }

func getLevel() zapcore.Level {
	//fmt.Println("log level", Conf.LogLevel)
	switch strings.ToLower(Conf.LogLevel) {
	case "debug":
		return zap.DebugLevel
	case "info":
		return zap.InfoLevel
	case "warn":
		return zap.WarnLevel
	case "error":
		return zap.ErrorLevel
	default:
		return zap.DebugLevel
	}
}

func getOutput() string {
	output := strings.TrimSpace(Conf.LogOutput)
	//fmt.Println("logger|output", output)
	if output == "" {
		return "stdout"
	}
	return output
}

func (l logger) Print(v ...interface{}) {
	logMsg("debug", fmt.Sprintf("%v", v))
}

func Info(msg string, tags ...zap.Field) {
	logMsg("info", msg, tags...)
}

func Debug(msg string, tags ...zap.Field) {
	logMsg("debug", msg, tags...)
}

func Warn(msg string, tags ...zap.Field) {
	logMsg("warn", msg, tags...)
}
func createBleveIndex(msg string, level string, tags []zapcore.Field) {
	var l errorLog
	if len(tags) > 0 {
		l = errorLog{Time: time.Now().UTC().Format(UTCDateLayout), Level: level, Msg: fmt.Sprintf("%s%v", msg, tags)}
	} else {
		l = errorLog{Time: time.Now().UTC().Format(UTCDateLayout), Level: level, Msg: msg}
	}
	i, _ := getIndex()
	//fmt.Println("createBleveIndex", l)
	err := i.Index(uid_utils.GetUid("al", true), l)
	if err != nil {
		Error(err.Error(), err)
	}
}

func Error(msg string, err error, tags ...zap.Field) {
	tags = append(tags, zap.NamedError("error", err))
	logMsg("error", msg, tags...)
}
func logMsg(l string, msg string, tags ...zap.Field) {
	if Conf.EnableConsoleLog {
		switch l {
		case "error":
			Log.log.Error(msg, tags...)
		case "warn":
			Log.log.Warn(msg, tags...)
		case "info":
			Log.log.Info(msg, tags...)
		default:
			Log.log.Debug(msg, tags...)
		}
		Log.log.Sync()
	}
	if Conf.EnableTextIndexLog {
		createBleveIndex(msg, l, tags)
	}
}

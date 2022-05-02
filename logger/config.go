package logger

type Config struct {
	ApplogIndexPath    string `json:"applogIndexPath"`
	EnableTextIndexLog bool   `json:"enableTextIndexLog"`
	EnableConsoleLog   bool   `json:"EnableConsoleLog"` //set true would write log on console
	LogOutput          string `json:"logOutput"`        //setting file name would create log file and append
	LogLevel           string `json:"logLevel"`
}

var (
	Conf Config
)

func SetConfig(c Config) {
	Conf = c
}

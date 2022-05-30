package main

import (
	"flag"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/kardianos/service"
	"github.com/ramnkl16/ez-search/app"
	"github.com/ramnkl16/ez-search/global"
	"github.com/ramnkl16/ez-search/syncconfig"
)

var (
	router     *gin.Engine
	config     *syncconfig.Config
	configFile string
	productSku string
	language   string
	workingDir string

	trace            bool
	debug            bool
	info             bool
	warn             bool
	errlvl           bool
	fatal            bool
	paniclvl         bool
	fullSync         bool //enable full sync
	flagPort         string
	flagForce        bool //avoid check sum validation
	flagRunSchedule  bool //run only rest api
	flagLocalStorage bool
)

type refTest struct {
	A  int
	B  string
	A1 float64
	B1 string
}

type program struct{}

func (p *program) Start(s service.Service) error {
	// Start should not block. Do the actual work async.
	go p.run()
	return nil
}
func (p *program) run() {
	c, err := syncconfig.OpenConfig(configFile)
	if err != nil {
		fmt.Println("unable to open config file ", err)
		return
	}

	config = c
	// if len(c.WorkingDir) > 0 {
	// 	syn.WorkingDir = c.WorkingDir
	// }
	if len(workingDir) > 0 {

		global.WorkingDir = workingDir
		c.WorkingDir = workingDir
		fmt.Println("env workdir", workingDir, c.WorkingDir)
		//	coredb.Initialize(workingDir)

	}
	global.CsvFileExt = c.CSVFileExtension
	global.CsvWatcherPath = c.CSVFileWatcherpath
	global.MaxIndexbatchSize = c.EzsearchSettings.IndexBatchSize

	syncconfig.Gconfig = c
	//config.HybrisSettings.Force = flagForce

	config.RunScheduler = flagRunSchedule
	//	logger.Info(fmt.Sprintf("hybris setting: %v", config.HybrisSettings.CatalogName))
	router = gin.Default()
	router.Use(serverHeader)
	app.StartApplication(config, router, productSku, language)
}
func (p *program) Stop(s service.Service) error {
	// Stop should not block. Return with a few seconds.
	return nil
}

func main() {

	// sEnc, _ := b64.StdEncoding.DecodeString("eyJ1IjoicmFtQHNiZGluYy5jb20iLCAicCI6IndlbGNvbWUxMjMifQ==")
	// fmt.Println(string(sEnc))
	// return
	flag.StringVar(&configFile, `c`, ``, `Location of configuration file`)
	flag.StringVar(&productSku, `s`, ``, `product sku`)
	flag.StringVar(&language, "l", "", "language code")
	flag.StringVar(&workingDir, "wd", "", "working dir")
	// flag.StringVar(&log_level, "ll", "stdout", "log standard output")
	flag.BoolVar(&fullSync, `fullsync`, false, ``)
	flag.BoolVar(&flagForce, `force`, false, ``)
	flag.BoolVar(&flagRunSchedule, `rs`, false, ``)
	flag.BoolVar(&flagLocalStorage, `localstorage`, false, ``)
	flag.Parse()
	if configFile == `` {
		fmt.Println("No config file specified (-c)", nil)
	}
	// if log_output != `` {
	// 	os.Setenv("LOG_OUTPUT", log_output)
	// }
	// if log_level == `` {
	// 	os.Setenv("LOG_LEVEL", "debug")
	// } else {
	// 	os.Setenv("LOG_LEVEL", log_level)
	// }

	// s := strings.Split("sicne  30 days  ago", " ")
	// for idx, i := range s {
	// 	if len(i) > 0 {
	// 		fmt.Println(idx, len(i), i)
	// 	}
	// }

	//getParsedQueryByKeyword()
	//return
	svcConfig := &service.Config{
		Name:        "Ezsearch",
		DisplayName: "Ez service",
		Description: "Ez service",
	}

	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		fmt.Println(err.Error())
	}
	// //logger, err = s.Logger(nil)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	err = s.Run()
	if err != nil {
		fmt.Println(err)
	}

}

func serverHeader(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*") //TODO refactor later.
	c.Header("Access-Control-Allow-Headers", "Access-Control-Allow-Origin, Accept, x-auth, x-ns, Content-Type, Authorization")
	c.Header("Access-Control-Allow-Private-Network", "true")
	c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS, HEAD, DELETE")
}

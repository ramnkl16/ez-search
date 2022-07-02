package app

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jasonlvhit/gocron"
	"github.com/ramnkl16/ez-search/abstractimpl"
	"github.com/ramnkl16/ez-search/api/ping"
	"github.com/ramnkl16/ez-search/auth"
	"github.com/ramnkl16/ez-search/datasources/catalogboltdb"
	"github.com/ramnkl16/ez-search/ezeventqueue"
	"github.com/ramnkl16/ez-search/ezsmtp"

	//"github.com/ramnkl16/ez-search/ezeventqueue"
	"github.com/ramnkl16/ez-search/ezsearch"
	"github.com/ramnkl16/ez-search/global"
	"github.com/ramnkl16/ez-search/logger"
	"github.com/ramnkl16/ez-search/models"
	"github.com/ramnkl16/ez-search/rest_errors"
	"github.com/ramnkl16/ez-search/services"
	"github.com/ramnkl16/ez-search/syncconfig"
	"github.com/ramnkl16/ez-search/utils/cache_utils"
	"github.com/ramnkl16/ez-search/utils/date_utils"
)

func StartApplication(config *syncconfig.Config, router *gin.Engine, productSku string, language string) {

	cache_utils.Initialize(2 * 60 * 60)            //two hours caching enabled. //must be initialized befor any logger cache
	cache_utils.InitializeCredential(24 * 60 * 60) //24 hours credential cache
	cache_utils.Cache.SetCheckExpirationCallback(ezsearch.CheckExpirationCallback)
	logger.SetConfig(config.LoggerSettings)
	//logger.Initialize()
	os.Chdir(config.WorkingDir)
	logger.InitLogger()
	logger.BuildAppIndexSchema()

	catalogboltdb.SetConfig(config.BoltDBSettings)
	config.EzsearchSettings.BoltDbBucketName = config.BoltDBSettings.BuketName
	ezsearch.SetConfig(config.EzsearchSettings)
	ezsmtp.SetConfig(config.EzsmptSettings)
	abstractimpl.IndexTablesPath = config.EzsearchSettings.IndexTablesPath
	abstractimpl.IndexBasePath = config.EzsearchSettings.IndexBasePath
	abstractimpl.CreateTables()
	abstractimpl.DefaultmetaData()
	router.Use(gin.Recovery())
	router.Use(AuthTokenValidation())
	scheduleJobForLogCleanup()
	// ezsmtp.SendEmail([]string{"ram.duraisamy@sbdinc.com"}, "", nil)
	// //ezsmtp.SendEmailUsinggoogle([]string{"ram.duraisamy@sbdinc.com"})
	// return

	//ezsearch.TriggerIndex()
	//fswatcher.WatchCSVFiles(config.CSVFileWatcherpath)
	//fmt.Println("runScheudleJob", config.RunScheduler)
	if config.RunScheduler {
		s := gocron.NewScheduler()
		s.Every(1).Minute().Do(ezeventqueue.ProcessEventqueue)
		s.Start()
	}

	//return

	logger.Info(fmt.Sprintf("started logging %s %s", config.WorkingDir, config.Port))
	router.GET("/ping", ping.Ping)
	mapUrls(router)
	//fmt.Printf())

	fullPath := path.Join(global.WorkingDir, "/swagger-ui")
	logger.Info(fmt.Sprintf("staticFS fullpath %s", fullPath))
	router.StaticFS("/swagger-ui", http.Dir(fullPath))
	fullwebUi := path.Join(global.WorkingDir, "/web-ui")
	router.StaticFS("/web-ui", http.Dir(fullwebUi))
	//abstractimpl.Delete(abstractimpl.EventQueueTable, "")
	// //ezcsv.GetJsonFromCsv("C:\\go-prj\\ez-search\\uploads\\Userinformation.csv", 0)

	// //cdCsv:= "{\"fileName\":\"C:\\\\go-prj\\\\ez-search\\\\uploads\\\\Userinformation.csv\",\"ignoreEmpty\":true,\"indexName\":\"macindex/new/customer\",\"uniqueIndexColIndex\":1}"
	// //cdstr = strings.ReplaceAll(cdstr, `\"`, `"`)
	// // bytes, err := rawIn.MarshalJSON()
	// //fmt.Println("rawin", cdstr)
	// // if err != nil {
	// // 	panic(err)
	// // }
	// cdstr := "{\"host\":\"TOW-P-SQLHA07\\\\sqlha07\",\"dbName\":\"MBAWEB\",\"lastSyncAt\":\"2000-01-01\",\"goScriptBlock\":{\"goScriptKey\":\"mac.mba75.goscript.distnumLoop\",\"queryKey\":\"mac.mba75.customer.query\",\"distnumKey\":\"mac.mba75.distnum.list\",\"params\":[\"2000-01-01\"]},\"docIdColName\":\"\",\"indexName\":\"indexes/mac/mba75/customers\",\"userName\":\"macuser\",\"password\":\"T001sM@C\"}"
	// var cd ezmssqlconn.MsSqlEventCustomData
	// cd.GoScriptBlock.Params = []string{"2000-01-01"}
	// err := json.Unmarshal([]byte(cdstr), &cd)
	// if err != nil {
	// 	logger.Error("Failed|ExecuteMsSqlScript|Unmarshal", err)
	// }
	// s, _ := json.Marshal(cd)
	// fmt.Println(string(s))
	// return
	//ezeventqueue.ProcessEventqueue()
	// schemas := ezcsv.GenerateIndexSchema("DistNum,TID,YearNum,WeekNum,ReportNum,StartDate,EndDate,StartDisbNum,EndDisbNum,StartAcctTransNum,EndAcctTransNum,StartTransNum,EndTransNum,CountSales,CountTransactions,ApMacStartBalance,ApMacEndBalance,ApMacInvoice,ApMacMACredit,ApMacNationalAccountCredit,ApMacOther", false)
	// fmt.Println("schema", schemas)
	//return
	router.Run(config.Port)
	// // if len(language) > 0 {
	// // 	hybris.UpdateSchedulerJobStatusAndTime(language)
	// // }
	// if len(productSku) == 0 {
	// 	hybris.ExecuteSchudleJob()
	// } else {
	// 	logger.Info("singl sku execution")
	// 	hybris.ExecuteSyncJobBySku(productSku, language)
	// }
}

func scheduleJobForLogCleanup() {
	logger.Debug("ScheduleJobForLogCleanup")
	ed := fmt.Sprintf("{\"noDays\": 30, \"indexNameKey\":\"schedulejob.delete_logs.key\"}")
	//edStr, _ := json.Marshal(ed)
	e := models.EventQueue{EventType: global.EVENT_TYPE_DETETE_LOG, EventData: ed,
		StartAt: date_utils.GetNowSearchFormat(), IsActive: "t", Status: int(global.STATUS_ACTIVE), RetryCount: 0, RetryMax: 5, RecurringInSeconds: 24 * 60 * 60}
	e.ID = "dellogs"
	err := abstractimpl.CreateOrUpdate(e, abstractimpl.EventQueueTable, "dellogs")
	if err != nil {
		logger.Error("Failed to add event queue", err)
	}
}
func AuthTokenValidation() gin.HandlerFunc {
	return func(c *gin.Context) {
		authToken := auth.GetXauthToken(c.Request)
		v, exist := cache_utils.GetFromCredentialCache(authToken)
		println("auth token", authToken, auth.GetNamespace(c.Request), v, exist)
		//fmt.Println(authToken)
		// for a := range cache_utils.Cache.GetKeys() {
		// 	fmt.Println("Avaialbl Keys: ")
		// 	fmt.Println(a)
		// }
		//fmt.Println("url.path", c.Request.URL.Path)
		//fmt.Println("cach key", cache_utils.Cache.GetKeys())
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)

		} else if value, exists := cache_utils.GetFromCredentialCache(authToken); exists != nil && c.Request.URL.Path != "/api/auth/login" &&
			c.Request.URL.Path != "/api/user" &&
			c.Request.URL.Path != "/ping" &&
			// c.Request.URL.Path != "/api/addorupdate" &&
			//c.Request.URL.Path != "/api/search" &&
			// c.Request.URL.Path != "/api/getindexes" &&
			// c.Request.URL.Path != "/api/getfields" &&
			// c.Request.URL.Path != "/api/createschema" &&
			!strings.Contains(c.Request.URL.Path, "/web-ui/") &&
			!strings.Contains(c.Request.URL.Path, "/swagger-ui/") {
			//fmt.Printf("Got value: %v\n", value)
			// fmt.Println("not authorized")
			// fmt.Println(exists)
			logger.Debug(fmt.Sprintf("AuthTokenValidation|credential cache value %s not authorized %v", value, exists))
			//has token then try login decode base64 user credential as it is cached
			if len(authToken) > 0 {
				//fmt.Println("authtoken", authToken)
				uEnc, err := b64.StdEncoding.DecodeString(authToken)
				if err != nil {
					c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"errCode": rest_errors.InvalidAuthToken, "status": http.StatusUnauthorized, "errorDesc": "Invalid auth token while parse from base 64", "datetime": date_utils.GetNowSearchFormat()})
					return
				}
				var bu models.UserBase64
				err = json.Unmarshal(uEnc, &bu)
				if err != nil {
					c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"errCode": rest_errors.InvalidAuthToken, "status": http.StatusUnauthorized, "errorDesc": "Invalid auth token while unmarshal from base 64", "datetime": date_utils.GetNowSearchFormat()})
					return
				}
				_, err = services.UserService.Login(bu.UserName, bu.Password, bu.Namespace)
				if err != nil {
					c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"errCode": rest_errors.InvalidAuthToken, "status": http.StatusUnauthorized, "errorDesc": "Invalid auth token", "datetime": date_utils.GetNowSearchFormat()})
					return
				}
				c.Next()

			} else {

				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"errCode": rest_errors.InvalidAuthToken, "status": http.StatusUnauthorized, "errorDesc": "Invalid auth token", "datetime": date_utils.GetNowSearchFormat()})
			}
			//return
		} else {
			logger.Debug(fmt.Sprintf("AuthTokenValidation|authorized %s url=%s", value, c.Request.URL))
			c.Next()
		}
	}
}

//}

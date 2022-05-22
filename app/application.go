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
	abstractimpl.IndexTablesPath = config.EzsearchSettings.IndexTablesPath
	abstractimpl.IndexBasePath = config.EzsearchSettings.IndexBasePath
	abstractimpl.CreateTables()
	router.Use(gin.Recovery())
	router.Use(AuthTokenValidation())
	//ezsearch.TriggerIndex()

	if !config.Restapi {
		s := gocron.NewScheduler()
		s.Start()
	}

	logger.Info(fmt.Sprintf("started logging %s %s", config.WorkingDir, config.Port))
	router.GET("/ping", ping.Ping)
	mapUrls(router)
	//fmt.Printf())
	fullPath := path.Join(global.WorkingDir, "/swagger-ui")
	logger.Info(fmt.Sprintf("staticFS fullpath %s", fullPath))
	router.StaticFS("/swagger-ui", http.Dir(fullPath))
	fullwebUi := path.Join(global.WorkingDir, "/web-ui")
	router.StaticFS("/web-ui", http.Dir(fullwebUi))

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

func AuthTokenValidation() gin.HandlerFunc {
	return func(c *gin.Context) {
		authToken := auth.GetXauthToken(c.Request)
		println("auth token", authToken, auth.GetNamespace(c.Request))
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

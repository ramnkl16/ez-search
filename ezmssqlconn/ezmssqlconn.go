package ezmssqlconn

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/matryer/goscript"
	"github.com/ramnkl16/ez-search/abstractimpl"
	"github.com/ramnkl16/ez-search/coredb"
	"github.com/ramnkl16/ez-search/global"
	"github.com/ramnkl16/ez-search/logger"
	"github.com/ramnkl16/ez-search/models"
	"github.com/ramnkl16/ez-search/rest_errors"
	"github.com/ramnkl16/ez-search/utils/date_utils"
	"go.uber.org/zap"
)

var goscript1 = `import (
	"fmt"
	"strings"
)

func goscript(dists string, dt string) ([]string, error) {

	list := make([]string, 0)
	for _, s := range strings.Split(dists, ",") {
		list = append(list, fmt.Sprintf("select top 1  c.CustNum cn,c.LastName ln,c.FirstName fn,	c.Address a,c.City c,c.State s,c.PostalCode pc,c.HomePhone hp,c.WorkPhone wp,cast(c.DiscountPct as int) dp,c.TaxID ti,c.TaxExemptFlag tef,c.Status s1,c.StatusDate sd,c.EMail e, c.CreatedBy cb,c.CreatedDate cd,c.ModifiedBy mb,c.ModifiedDate md ,c.StopNum sn,CustomerTypeOID cto,PaymentSource ps ,CashCustomer cc,cast (ShopOID as char(36)) so,ReceiptByEmail rbe,ReceiptByPrinted rbp,CustomerNumber cn2,ShopNum sn3, d.distnum d from customer c inner join distributor d on c.[DistributorOID] = d.distributorOid where d.distnum = %s and c.createdDate>='%s'", s, dt))
	}

	return list, nil
}`

type GoScriptBlock struct {
	GoScriptKey string   `json:"goScriptKey"`
	QueryKey    string   `json:"queryKey"`
	DistnumKey  string   `json:"distnumKey"`
	Params      []string `json:"params"`
}
type MsSqlEventCustomData struct {
	Host          string            `json:"host"`   // id
	DbName        string            `json:"dbName"` // id
	LastSyncAt    string            `json:"lastSyncAt,omitempty"`
	GoScriptBlock GoScriptBlock     `json:"goScriptBlock"`
	IsStoredProc  bool              `json:"isStoredProc,omitempty"`
	DocIdColName  string            `json:"docIdColName"`
	IndexName     string            `json:"indexName"`
	FieldsMap     map[string]string `json:"fieldsMap,omitempty"`
	UserName      string            `json:"userName"`
	Password      string            `json:"password"`
	SaveOnLocal   string            `json:"saveOnLocal"`
}

func ExecuteMsSqlScript(eq *models.EventQueue) rest_errors.RestErr {
	//fmt.Println("ExecuteMsSqlScript", eq.EventData)
	var cd MsSqlEventCustomData
	err := json.Unmarshal([]byte(eq.EventData), &cd)
	if err != nil {
		logger.Error("Failed|ExecuteMsSqlScript|Unmarshal", err, zap.String("ref1", eq.EventData))
		return rest_errors.NewInternalServerError("Failed|ExecuteMsSqlScript", err)
	}

	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;database=%s;",
		cd.Host, cd.UserName, cd.Password, cd.DbName)
	//fmt.Println("msssql conn string", connString)
	db, errdb := sql.Open("mssql", connString)
	if errdb != nil {
		logger.Error("Failed open db:", errdb, zap.String("ref1", cd.Host), zap.String("ref2", cd.DbName))
	}
	ctx := context.Background()

	// Check if database is alive.
	err = db.PingContext(ctx)
	if err != nil {
		logger.Error("Failed while ping db", err, zap.String("ref1", cd.Host), zap.String("ref2", cd.DbName))
		return rest_errors.NewInternalServerError("Failed while ping db", err)
	}
	gosciptGen, err := coredb.GetKey(cd.GoScriptBlock.GoScriptKey)
	if err != nil {
		restErr := rest_errors.NewInternalServerError(fmt.Sprintf("Please try register goscript,missing key %s", cd.GoScriptBlock.GoScriptKey), err)
		return restErr
	}
	query, err := coredb.GetKey(cd.GoScriptBlock.QueryKey)
	if err != nil {
		restErr := rest_errors.NewInternalServerError(fmt.Sprintf("Please try register sql script,missing key %s", cd.GoScriptBlock.QueryKey), err)
		return restErr
	}
	distNums, err := coredb.GetKey(cd.GoScriptBlock.DistnumKey)
	if err != nil {
		restErr := rest_errors.NewInternalServerError(fmt.Sprintf("Please try register sql script,missing key %s", cd.GoScriptBlock.DistnumKey), err)
		return restErr
	}

	script := goscript.New(string(gosciptGen))
	//fmt.Println("goscriptblock", string(distNums), string(query), cd.GoScriptBlock.Params)
	list, err := script.Execute(string(distNums), string(query), cd.GoScriptBlock.Params) //,8605,8651,8750,8780,8782,8986", "2000-01-01")
	logger.Info(fmt.Sprintf("query list|%v", list))

	if err != nil {
		logger.Error("Failed at execuge go script, Check customdata query defintion", err)
	}
	for _, q := range list.([]string) {
		results, err := getDataset(q, cd.DocIdColName, db, ctx)
		if err != nil {
			logger.Error("Failed while getdataset", err)
			continue
		}
		logger.Warn(fmt.Sprintf("distnum %s(%d)", q, len(results)))
		abstractimpl.BatchCreateOrUpdate(cd.IndexName, results)
		if cd.SaveOnLocal == "t" {

			saveAsJsonFile(results, cd.IndexName)
		}
	}
	return nil
}
func saveAsJsonFile(jsonData map[string]interface{}, entityName string) {
	wd := filepath.Join(global.WorkingDir, "downloads")
	_, err := os.Stat(wd)
	if os.IsNotExist(err) {
		os.MkdirAll(wd, os.ModeDir)
	}
	fileName := fmt.Sprintf(`%s/%s_%s.json`, wd, strings.ReplaceAll(entityName, "/", "_"), date_utils.GetNowFileLayout())
	logger.Warn(fmt.Sprintf("saveasjsonFIle filename %s", fileName))
	f, err := os.Create(fileName)
	defer f.Close()
	fileBytes, err := json.Marshal(jsonData)
	if err != nil {
		logger.Error("Failed whil marshal", err)
	}
	f.Write(fileBytes)

}

func getDataset(q, docIdColName string, db *sql.DB, ctx context.Context) (map[string]interface{}, error) {

	rows, err := db.QueryContext(ctx, q)
	if err != nil {
		logger.Error("Failed while execute", err)
		return nil, rest_errors.NewRestError("Failed while execute", http.StatusInternalServerError, err.Error(), nil)
	}
	defer rows.Close()
	columns, err := rows.Columns()

	// for each database row / record, a map with the column names and row values is added to the allMaps slice
	allMaps := make(map[string]interface{}, 0)

	for rows.Next() {
		values := make([]interface{}, len(columns))
		pointers := make([]interface{}, len(columns))
		for i := range values {
			pointers[i] = &values[i]
		}
		err := rows.Scan(pointers...)
		if err != nil {
			logger.Error("Failed while fetch", err)
			continue
		}

		resultMap := make(map[string]interface{})
		for i, val := range values {
			//fmt.Printf("Field=%s val=%v\n", columns[i], val)
			if val != nil {
				resultMap[columns[i]] = val
			}
		}
		//fmt.Println(resultMap)
		if len(docIdColName) == 0 {
			docIdColName = columns[0]
		}
		val := resultMap[docIdColName]
		//fmt.Println(docIdColName, val, fmt.Sprintf("%v", val))
		allMaps[fmt.Sprintf("%v", val)] = resultMap
	}
	return allMaps, nil
}

func QueueNewEventForMsSqlSync(cd *MsSqlEventCustomData, eq *models.EventQueue) string {
	return fmt.Sprintf(` {
	"createdAt": "%s",
	"eventData": "{
	"host":"%s",  
	"dbName":"%s",
	"userName":"%s",
	"password":"%s",
	"query":"%s",
	"isStoredProc":"%t",
	"fieldsMap":"%s",
	"indexName:"%s",
	"docIdColName:"%s"
	 }",
	"eventType": "mssqlsync",
	"id": "newid",
	"isActive": "t",
	"message": "",
	"retryCount":0,
	"startAt": "%s",
	"status":1,
	"updatedAt": "%s",
	"RecurringInSeconds":0
  }`, eq.CreatedAt, cd.Host, cd.DbName, cd.UserName, cd.Password, cd.GoScriptBlock, cd.IsStoredProc, cd.FieldsMap, cd.IndexName, cd.DocIdColName, eq.StartAt, eq.StartAt)
}

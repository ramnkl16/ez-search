package abstractimpl

import (
	"encoding/json"
	"errors"
	"fmt"
	"path"
	"strings"
	"time"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/analysis/lang/en"
	"github.com/blevesearch/bleve/v2/analysis/token/keyword"
	"github.com/ramnkl16/ez-search/common"
	"github.com/ramnkl16/ez-search/coredb"
	"github.com/ramnkl16/ez-search/ezsearch"

	"github.com/ramnkl16/ez-search/global"
	"github.com/ramnkl16/ez-search/logger"
	"github.com/ramnkl16/ez-search/rest_errors"
	"github.com/ramnkl16/ez-search/utils/cache_utils"
	"github.com/ramnkl16/ez-search/utils/crypto_utils"
	"github.com/ramnkl16/ez-search/utils/uid_utils"
	"go.uber.org/zap/zapcore"
)

var (
	IndexTablesPath string
	IndexBasePath   string
)

const (
	EventQueueTable    string = "tables/tables.eventqueue"
	EventQueueHisTable string = "tables/tables.eventqueuehistory"
	QueryMetaTable     string = "tables/tables.querymeta"
	UserTable          string = "tables/tables.user"
	MenuTable          string = "tables/tables.menu"
	UserGroupTable     string = "tables/tables.usergroup"
	NamespaceTable     string = "tables/tables.namespace"
	UserMenuTable      string = "tables/tables.usermenu"
)

//create tables for eventqueue eventqueuehistory and querydef
func CreateTables() {
	tables := common.GetIndexes(true)
	var isExist bool
	//fmt.Println("createtable", tables)
	isExist = tables[EventQueueTable]
	if !isExist {
		BuildIndexSchema(EventQueueTable, eqSchema, "tables")
	}
	isExist = tables[EventQueueHisTable]
	if !isExist {
		BuildIndexSchema(EventQueueHisTable, eqSchema, "tables")
	}
	isExist = tables[QueryMetaTable]
	if !isExist {
		BuildIndexSchema(QueryMetaTable, queryMetaSchema, "tables")
	}
	isExist = tables[UserTable]
	if !isExist {
		BuildIndexSchema(UserTable, usrSchema, "tables")
		var m interface{}
		id := uid_utils.GetUid("us", true)
		json.Unmarshal([]byte(fmt.Sprintf(`{	"username": "admin",
			"id":"%s",
			"namespaceId": "platform",
			"userGroupId": "grp1",
			"email": "admin@gost.com",
			"token": "%s",
			"mobile": "12345627890",
			"firstName": "Admin",
			"lastName": "",
			"userRoleId": "role1",
			"isActive": "t",
			"emailVerified": "2022-03-28T01:24:27Z",
			"passwordUpdatedAt": "2022-03-28T01:24:27Z"
		}`, id, crypto_utils.GetMd5("welcome@123"))), &m)
		CreateOrUpdate(m, UserTable, id)
	}
	isExist = tables[MenuTable]
	if !isExist {
		BuildIndexSchema(MenuTable, menuSchema, "tables")
		var m1 interface{}
		json.Unmarshal([]byte(`{"id":"root", "name":"Root","parentId":"", "link":"/", "isActive":"t","updatedAt":"2022-03-27T00:00:00Z", "createdAt":"2022-03-27T00:00:00Z" }`), &m1)
		CreateOrUpdate(m1, MenuTable, "root")
		json.Unmarshal([]byte(`{"id":"user", "name":"User","parentId":"root","link":"user", "isActive":"t","updatedAt":"2022-03-27T00:00:00Z", "createdAt":"2022-03-27T00:00:00Z" }`), &m1)
		CreateOrUpdate(m1, MenuTable, "user")
		// json.Unmarshal([]byte(`{"id":"usrGrp", "name":"User Group","parentId":"root","link":"usergroup", "isActive":"t","updatedAt":"2022-03-27T00:00:00Z", "createdAt":"2022-03-27T00:00:00Z" }`), &m1)
		// CreateOrUpdate(m1, MenuTable, "usrGrp")
		// json.Unmarshal([]byte(`{"id":"usrmnu", "name":"User menu","parentId":"root", "link":"usermenu", "isActive":"t","updatedAt":"2022-03-27T00:00:00Z", "createdAt":"2022-03-27T00:00:00Z" }`), &m1)
		// CreateOrUpdate(m1, MenuTable, "usrmnu")
		// json.Unmarshal([]byte(`{"id":"ns", "name":"Namespace","parentId":"root", "link":"namespace", "isActive":"t","updatedAt":"2022-03-27T00:00:00Z", "createdAt":"2022-03-27T00:00:00Z" }`), &m1)
		// CreateOrUpdate(m1, MenuTable, "ns")
		// json.Unmarshal([]byte(`{"id":"menu", "name":"Menu","parentId":"root", "link":"menu", "isActive":"t","updatedAt":"2022-03-27T00:00:00Z", "createdAt":"2022-03-27T00:00:00Z" }`), &m1)
		// CreateOrUpdate(m1, MenuTable, "menu")
		json.Unmarshal([]byte(`{"id":"qryDef", "name":"Query Defintion","parentId":"root", "link":"querydef", "isActive":"t","updatedAt":"2022-03-27T00:00:00Z", "createdAt":"2022-03-27T00:00:00Z" }`), &m1)
		CreateOrUpdate(m1, MenuTable, "querydef")
		json.Unmarshal([]byte(`{"id":"indexList", "name":"Show all indexes","parentId":"root", "link":"indexlist", "isActive":"t","updatedAt":"2022-03-27T00:00:00Z", "createdAt":"2022-03-27T00:00:00Z" }`), &m1)
		CreateOrUpdate(m1, MenuTable, "indexlist")
		json.Unmarshal([]byte(`{"id":"indexFields", "name":"Show Fields","parentId":"root", "link":"indexlist", "isActive":"t","updatedAt":"2022-03-27T00:00:00Z", "createdAt":"2022-03-27T00:00:00Z" }`), &m1)
		CreateOrUpdate(m1, MenuTable, "indexFields")

	}
	isExist = tables[UserGroupTable]
	if !isExist {
		BuildIndexSchema(UserGroupTable, grpSchema, "tables")
		var m interface{}
		json.Unmarshal([]byte(`{	
			"id":"grp1",
			"namespaceId": "platform",
			"desc": "Main group for platform",
			"name": "main group",
			"isActive": "t",
			"createdAt": "2022-03-28T01:24:27Z",
			"updatedAt": "2022-03-28T01:24:27Z"
		}`), &m)
		CreateOrUpdate(m, UserGroupTable, "grp1")
	}
	isExist = tables[NamespaceTable]
	if !isExist {
		BuildIndexSchema(NamespaceTable, nsSchema, "tables")
		var m1 interface{}
		json.Unmarshal([]byte(`{"id":"platform", "name":"platform","code":"platform", "isActive":"t","updatedAt":"2022-03-27T00:00:00Z", "createdAt":"2022-03-27T00:00:00Z" }`), &m1)
		CreateOrUpdate(m1, NamespaceTable, "platform")
	}
	isExist = tables[UserMenuTable]
	if !isExist {
		BuildIndexSchema(UserMenuTable, usrMenuSchema, "tables")
		var m1 interface{}
		json.Unmarshal([]byte(`{"id":"platform", "nsId":"platform","menuId":"root", "cd":"customdata", "refId":"admin", "refType":"NS", "privilege":31, "isActive":"t","updatedAt":"2022-03-27T00:00:00Z", "createdAt":"2022-03-27T00:00:00Z" }`), &m1)
		CreateOrUpdate(m1, UserMenuTable, "platform")
	}
}

func BuildIndexSchema(indexName string, fields []common.BleveFieldDef, indexFolderName string) rest_errors.RestErr {
	// a generic reusable mapping for english text
	englishTextFieldMapping := bleve.NewTextFieldMapping()
	englishTextFieldMapping.Analyzer = en.AnalyzerName
	keywordFieldMapping := bleve.NewTextFieldMapping()
	keywordFieldMapping.Analyzer = keyword.Name

	indexmapping := bleve.NewDocumentMapping()
	for _, f := range fields {
		switch strings.ToLower(f.Type) {
		case "bool":

			indexmapping.AddFieldMappingsAt(f.Name, bleve.NewBooleanFieldMapping())
		case "date":
			indexmapping.AddFieldMappingsAt(f.Name, bleve.NewDateTimeFieldMapping())
		case "numeric":
			indexmapping.AddFieldMappingsAt(f.Name, bleve.NewNumericFieldMapping())
		case "geopoint":
			indexmapping.AddFieldMappingsAt(f.Name, bleve.NewGeoPointFieldMapping())
		default:
			indexmapping.AddFieldMappingsAt(f.Name, bleve.NewTextFieldMapping())
		}
	}
	indexPath := indexName
	if !strings.Contains(indexName, "/") {
		//if !isTable {
		indexPath = fmt.Sprintf("%s%c%s", indexFolderName, '/', indexName)
		// } else {
		// 	indexPath = fmt.Sprintf("%s%c%s", IndexTablesPath, '/', indexName)
		// }
	}
	patternIndexName := indexName

	grp := global.RegexParseDate.FindAllSubmatch([]byte(indexPath), -1)
	if grp != nil {
		dt := time.Now().UTC()

		dtFormat := string(grp[0][1])
		//dtVal := time.Now().UTC().Format(dtFormat)
		dtVal := dt.Format(dtFormat)
		indexPath = strings.Replace(indexPath, fmt.Sprintf("{%s}", dtFormat), dtVal, -1)
		//fmt.Println("formated index Name", indexPath, dtVal)

	}
	indexMapping := bleve.NewIndexMapping()
	docMapName := "docs" //strings.ReplaceAll(indexName, ".", "")
	//fmt.Println("buildindexschema|", indexMapping)
	indexMapping.AddDocumentMapping(docMapName, indexmapping)

	index, err := ezsearch.GetIndex(indexPath)
	//fmt.Println("afterGetIndex(indexName)")
	if err != nil || index == nil {
		logger.Info("BuildIndexSchema|Creating  new index ... ", zapcore.Field{Type: zapcore.StringType, Key: "indexname", String: indexName})
		//fmt.Println("BuildIndexSchema|Creating  new index ... ", indexPath)
		// create a mapping
		//fmt.Println("index base path", IndexBasePath)

		//fmt.Println("index path", indexPath)
		index, err := bleve.New(indexPath, indexMapping)
		if err != nil {
			//fmt.Println("BuildIndexSchema|Failed schemabuild mapping", err)
			logger.Info("BuildIndexSchema|Failed schemabuild mapping", zapcore.Field{Type: zapcore.StringType, String: err.Error(), Key: "msg"})
			return rest_errors.NewBadRequestError(err.Error())
		}
		cache_utils.AddOrUpdateCache(indexPath, index)
	} else {
		msg := fmt.Sprintf("Index %s is already exist, provide new index name", indexPath)
		key := fmt.Sprintf("%s.schema", patternIndexName)
		bytes, _ := coredb.GetKey(key)
		//fmt.Println("before if not found key in core db", key, string(bytes))
		if bytes == nil || len(bytes) == 0 {
			logger.Debug(fmt.Sprintf("not found key in core db key:%s data:%s", key, string(bytes)))
			bytes, _ := json.Marshal(fields)
			coredb.AddKey(key, bytes)
		}
		return rest_errors.NewInternalServerError(msg, errors.New(msg))
	}
	bytes, _ := json.Marshal(fields)
	coredb.AddKey(fmt.Sprintf("%s.schema", patternIndexName), bytes)
	return nil
}
func getIndexNameWithTable(indexName string) string {
	if strings.Contains(indexName, "/") {
		return indexName
	}
	s := fmt.Sprintf("%s%c%s", IndexTablesPath, '/', indexName)
	//fmt.Println("index path", s)
	return s
}

//index name start with tables then that should be created under tables
// the table name pattern should be  tables.eventqueue.bleve
func GetTable(tableName string) (bleve.Index, rest_errors.RestErr) {
	var index bleve.Index
	var err error
	//fmt.Println("Get Index called", indexName, "Count ", cache_utils.Cache.Count())
	i, err := cache_utils.Cache.Get(tableName)
	//fmt.Println("after cache get")
	if err != nil {
		logger.Error("GetTable|Failed while cache_utils.cache", err, zapcore.Field{String: tableName, Key: "p1", Type: zapcore.StringType})
		i, err = bleve.Open(path.Join(global.WorkingDir, getIndexNameWithTable(tableName)))

		if err != nil {
			logger.Error("GetTable|Failed while open index", err, zapcore.Field{String: tableName, Key: "p1", Type: zapcore.StringType})
			return nil, rest_errors.NewBadRequestError(err.Error())
		}
		//fmt.Println("GetIndex|index has found under data folder", indexName, err)
		cache_utils.AddOrUpdateCache(tableName, i)
	}

	index = i.(bleve.Index)
	//fmt.Println("before return ", index.Name())
	return index, nil
}

// func GetIndexes(true bool) {
// 	panic("unimplemented")
// }

var (
	eqSchema = []common.BleveFieldDef{
		{Name: "id", Type: "text"},
		{Name: "eventType", Type: "text"},
		{Name: "eventData", Type: "text"},
		{Name: "status", Type: "numeric"},
		{Name: "retryCount", Type: "numeric"},
		{Name: "startAt", Type: "date"},
		{Name: "isActive", Type: "text"},
		{Name: "createdAt", Type: "date"},
		{Name: "updatedAt", Type: "date"},
	}
	eqData = `{"eventType":"product", "eventData":"{customjsondata}","status":1, "retryCount":0, 
	"startAt":"2022-01-22T02:04:00z", "isActive":"t", "createdAt":"2022-01-22T02:04:00z", "updatedAt":"2022-01-22T02:04:00z"}`
	eqHisSchema = []common.BleveFieldDef{
		{Name: "id", Type: "text"},
		{Name: "eventQueueId", Type: "text"},
		{Name: "eventType", Type: "text"},
		{Name: "eventData", Type: "text"},
		{Name: "status", Type: "numeric"},
		{Name: "retryCount", Type: "numeric"},
		{Name: "startAt", Type: "date"},
		{Name: "isActive", Type: "text"},
		{Name: "createdAt", Type: "date"},
		{Name: "updatedAt", Type: "date"}}
	eqHisData = `{"eventQueueId":"id1", "eventType":"product", "eventData":"{customjsondata}","status":1, "retryCount":0, 
	"startAt":"2022-01-22T02:04:00z", "isActive":'t', "createdAt":"2022-01-22T02:04:00z", "updatedAt":"2022-01-22T02:04:00z"}`
	queryMetaSchema = []common.BleveFieldDef{
		{Name: "id", Type: "text"},
		{Name: "name", Type: "text"},
		{Name: "division", Type: "text"},
		{Name: "module", Type: "text"},
		{Name: "page", Type: "text"},
		{Name: "cd", Type: "text"},
		{Name: "isActive", Type: "text"},
		{Name: "createdAt", Type: "date"},
		{Name: "updatedAt", Type: "date"}}

	data = `{"name":"product query","division":"","module":"","page":"search","data":"{"q":"select * from index where date>2022-01-22T02:04:00z"}",
 "isActive":true, "createdAt":"2022-01-22T02:04:00z", "updatedAt":"2022-01-22T02:04:00z"}`
	usrSchema = []common.BleveFieldDef{
		{Name: "id", Type: "text"},
		{Name: "usernName", Type: "text"},
		{Name: "namespaceId", Type: "text"},
		{Name: "token", Type: "text"},
		{Name: "email", Type: "text"},
		{Name: "mobile", Type: "text"},
		{Name: "firstName", Type: "text"},
		{Name: "lastName", Type: "text"},
		{Name: "roleId", Type: "text"},
		{Name: "isActive", Type: "text"},
		{Name: "emailVerified", Type: "date"},
		{Name: "passwordUpdatedAt", Type: "date"},
		{Name: "createdAt", Type: "date"},
		{Name: "updatedAt", Type: "date"},
	}
	grpSchema = []common.BleveFieldDef{
		{Name: "id", Type: "text"},
		{Name: "name", Type: "text"},
		{Name: "namespaceId", Type: "text"},
		{Name: "desc", Type: "text"},
		{Name: "isActive", Type: "text"},
		{Name: "createdAt", Type: "date"},
		{Name: "updatedAt", Type: "date"},
	}
	menuSchema = []common.BleveFieldDef{
		{Name: "id", Type: "text"},
		{Name: "name", Type: "text"},
		{Name: "link", Type: "text"},
		{Name: "parentId", Type: "text"},
		{Name: "isActive", Type: "bool"},
		{Name: "createdAt", Type: "date"},
		{Name: "updatedAt", Type: "date"},
	}

	nsSchema = []common.BleveFieldDef{
		{Name: "id", Type: "text"},
		{Name: "name", Type: "text"},
		{Name: "code", Type: "text"},
		{Name: "customJson", Type: "text"},
		{Name: "cotextToken", Type: "text"},
		{Name: "isActive", Type: "text"},
		{Name: "createdAt", Type: "date"},
		{Name: "updatedAt", Type: "date"},
	}
	usrMenuSchema = []common.BleveFieldDef{
		{Name: "id", Type: "text"},
		{Name: "nsId", Type: "text"},
		{Name: "menuId", Type: "text"},
		{Name: "refId", Type: "text"},
		{Name: "cd", Type: "text"}, //custom data
		{Name: "refType", Type: "text"},
		{Name: "privilege", Type: "numeric"},
		{Name: "isActive", Type: "text"},
		{Name: "createdAt", Type: "date"},
		{Name: "updatedAt", Type: "date"},
	}
)

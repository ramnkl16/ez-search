# Ez-search
Ez search is powered by [bleve](http://blevesearch.com/docs/Home/) text index. Explose as Rest api to ingest data from any external system. Sql Parse engine helps to explore indexed data using SQL query  Select * from indexes/indexName 

You can use the go build command to distribute to any os after change the go env variable. for example if you want to deploy on windows machine then set "GOOS=windows" and run go build that would generate executable file ez-search.exe

## Local dev env 
after check out run the app from your local machine using the below command.
local rest api setup  [go run .\main.go -c config.json -wd c:\go-prj\ez-search] you can provide any port number
for in the config.json should be available under root folder.

Post man collections is available in the root folder. 
to ingest index data use the below rest api 
    
   {{baseUrl}}/api/addorupdate?indexName=indexes/winevent-{2006-01-02}&indexTranDate=2022-05-10

It supports date based index pattern. in case want to create daily separate index then you should IndexName{2006-01-02}, in case yearly then IndexName{2006}

get the indexed data use the below api call.
[post] http://localhost:8015/api/search

## Schema build
By default schema creates automatically when calls the ingest api. Recommend to create schema first which helps to query date time types using since clause on sql query. 

schema field type [bool|text|date|numeric|geopoint]
sample schema defintion sample json [{"name":"name", "type":"text"},{"name":"startDt", "type":"date"}, {"name":"age", "type":"numeric"}]
last 10 years date range against launched date field  [10*360*24*60]

## bleve index search query  
Field Scoping 
You can qualify the field for these searches by prefixing them with the name of the field separated by a colon.
[name:ram] parsing field logic is upto [:] "name" field name and "ram" should match in the index document. Would apply as match query
[select id,name,age from indexName where name:ram,age:>40,+age:<=50,startDt>2022-01-01T01:01:00Z facets name limit 1, 10]

Terms query In where condition if the filed name missed then automatically construct the term query in the below query "ram" will searched any document using term query which mean find the "ram" any where in the document on all text fields
[select id,name,age from indexName where ram,age:>40,+age:<=50,startDt>2022-01-01T01:01:00Z facets name limit 1, 10]

Regular Expressions
You can use regular expressions in addition to using terms by wrapping the expression in forward slashes (/).
[name:/r*/] in the value part starts with forward slash then apply regex query
[select id,name,age from indexName where name:/r*/,age:>40,+age:<=50,startDt>2022-01-01T01:01:00Z facets name limit 1, 10]

Required, Optional, and Exclusion
When your query string includes multiple items, by default these are placed into the SHOULD clause of a Boolean Query.
You can change this by prefixing your items with a "+" or "-". The "+" Prefixing with plus places that item in the MUST portion of the boolean query. The "-" Prefixing with a minus places that item in the MUST NOT portion of the boolean query.
[select id,name,age from indexName where name:ram,age:>40,+age:<=50,startDt>2022-01-01T01:01:00Z facets name limit 1, 10]

Numeric / Date Ranges
You can perform ranges by using the >, >=, <, and <= operators, followed by a valid numeric/datetime value.

Escaping
The following quoted string enumerates the characters which may be escaped:

[+-=&|><!(){}[]^\"~*?:\\/]
NOTE: this list contains the space character.

In order to escape these characters, they are prefixed with the \ (backslash) character. In all cases, using the escaped version produces the character itself and is not interpreted by the lexer.

Example: "my\ name" will be interpreted as a single argument to a match query with the value “my name”.

Example: "contains {a\" character} will be interpreted as a single argument to a phrase query with the value contains {a " character}.

Date field is formated and converted to UTC time zone. 
Examaple 2022-02-19T20:49:03Z  golang format is [2006-01-02T15:04:05Z] which is equalant [yyyy-MM-ddThh:mm:ssZ] 
while searching must follow the same format.

Log settings 
    "loggerSettings":{
        "applogIndexPath":"indexes/applogs-{2006-01-02}", index document creation
        "enableConsoleLog":true,  --> set true writes logs into console 
        "enableTextIndexLog":true,--> set true writes logs into bleve search 
        "logOutput":"logs.txt", --> set as empty file log would be disabled otherwise logs writes on specified  file name under root of  logs folder
        "logLevel":"debug"
    },

swagger json generator help to build rest
swagger document generation new path `swag.exe init .`
and copy the json file into swagger-ui folder copy `.\docs\swagger.json .\swagger-ui\`
- code generation from xml document table schema and including dto,dao,service and controllers along with swagger tags
    [`.\codege.exe`] make sure that folders require xml defintion files are available under codedef folder. code generator always looking for codedef folder 
--Build release mode `go build -ldflags "-s -w"`

--event queue for delete old index documents
 {
        "id": "dellogs",
        "eventType": "dellogs",
        "eventData": "{\"noDays\":15,\"indexNameKey\":\"schedulejob.delete_logs.key\"}",
        "status": 1,
        "startAt": "2022-07-02T11:04:29Z",
        "retryCount": 0,
        "retryMax": 5,
        "retryDuraition": 0,
        "message": "",
        "isActive": "t",
        "createdAt": "",
        "updatedAt": "",
        "recurringInSeconds": 86400,
        "lastSyncAt": "",
        "timeTaken": 0
    }

--event queue fetch data from mssql ezmssqlconn\readme.txt script to register the events
ref
eventdata json {"host":"TOW-P-SQLHA07\\sqlha07","dbName":"MBAWEB","lastSyncAt":"2000-01-01","goScriptBlock":{"goScriptKey":"mac.mba75.goscript.distnumLoop","queryKey":"mac.mba75.customer.query","distnumKey":"mac.mba75.distnum.list","params":["2000-01-01"]},"docIdColName":"","indexName":"indexes/mac/mba75/customers","userName":"macuser","password":"T001sM@C","saveOnLocal":""}
{
        "id": "mac.mba75.customer",
        "eventType": "mssql",
        "eventData": "{\"host\":\"TOW-P-SQLHA07\\\\sqlha07\",\"dbName\":\"MBAWEB\",\"lastSyncAt\":\"2000-01-01\",\"goScriptBlock\":{\"goScriptKey\":\"mac.mba75.goscript.distnumLoop\",\"queryKey\":\"mac.mba75.customer.query\",\"distnumKey\":\"mac.mba75.distnum.list\",\"params\":[\"2000-01-01\"]},\"docIdColName\":\"\",\"indexName\":\"indexes/mac/mba75/customers\",\"userName\":\"macuser\",\"password\":\"T001sM@C\",\"saveOnLocal\":\"\"}",
        "status": 1,
        "startAt": "2022-07-02T11:04:29Z",
        "retryCount": 0,
        "retryMax": 5,
        "retryDuraition": 0,
        "message": "",
        "isActive": "t",
        "createdAt": "",
        "updatedAt": "",
        "recurringInSeconds": 86400,
        "lastSyncAt": "",
        "timeTaken": 0
    }
 
    
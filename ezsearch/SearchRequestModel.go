package ezsearch

type SearchRequestModel struct {
	DateField string `json:"dtField"`
	DateRange int32  `json:"dtRange"` //in mins for ex, 30->last 30mins, 60-> last 60 mins, 24*60-> a day ago, 360*24*60 a year ago

	Fields         []string `json:"fields"` //Fields that are generated out put
	Facets         []string `json:"facets"` //list of facet name that needs to be provided
	BoolQueries    []string `json:"bools"`
	TermQueries    []string `json:"terms"`
	PharseQueries  []string `json:"pharses"`  // used to build additional query filter.
	DateQueries    []string `json:"dates"`    //You can perform date range searches by using the >, >=, <, and <= operators, followed by a date value in quotes.
	NumericQueries []string `json:"numerics"` //You can perform date range searches by using the >, >=, <, and <= operators, followed by a date value in quotes.
	GeoQueries     []string `json:"geos"`     //You can perform date range searches by using the >, >=, <, and <= operators, followed by a date value in quotes.

	//Example: created:>"2016-09-21" will perform an Date Range Query on the created field for values after September 21, 2016.
	IndexName    string   `json:"indexName"`    //index that needs to be searched
	LocalStorage bool     `json:"localStorage"` //if it is set as tru, then will generate csv file and stored under local server
	From         int      `json:"from"`
	Size         int      `json:"size"`
	SortBy       []string `json:"sortBy"`
}
type SearchRequestQuery struct {
	Query string `json:"q"`
}

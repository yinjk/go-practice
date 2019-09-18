package main

type SearchRes struct {
	Took         int
	TimeOut      bool `json:"time_out"`
	Shards       Shards
	Hits         Hits
	Aggregations map[string]interface{}
}
type Shards struct {
	Total      int
	Successful int
	Skipped    int
	Failed     int
}
type Hits struct {
	Total    interface{}
	MaxScore float64 `json:"max_score"`
	Hits     interface{}
}

////////////////// Used to decode the return value of esapi get method   /////////////////
type GetRes struct {
	Id     string      `json:"_id"`
	Index  string      `json:"_index"`
	Source interface{} `json:"_source"`
}

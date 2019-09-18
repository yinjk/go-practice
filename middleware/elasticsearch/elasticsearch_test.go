package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"testing"
)

func TestEasyAPI_Search(t *testing.T) {
	client := New([]string{"http://172.21.82.11:9200"})
	dsl := `{
    "query": {
        "bool": {
            "must": [
                {
                    "match": {
                        "kind": "SERVER"
                    }
                }
            ]
        }
    }
}`
	var result []map[string]interface{}
	if res, err := client.Search(&result, "zipkin:span*", dsl, SearchWithSize(0)); err != nil {
		panic(err)
	} else {
		fmt.Println(fmt.Sprintf("took: %d, total: %v, current: %d", res.Took, res.Hits.Total, len(result)))
	}
}

func TestEasyAPI_Search2(t *testing.T) {
	client := New([]string{"http://localhost:9200"})
	dsl := `{
    "query": {
        "bool": {
            "filter": {
                "range": {
                    "timestamp": {
                        "gte": 1567586321608992
                    }
                }
            },
            "must": [
                {
                    "match": {
                        "localEndpoint.serviceName": "withdraw"
                    }
                },
                {
                    "match": {
                        "kind": "SERVER"
                    }
                }
            ]
        }
    },
    "aggs": {
        "tps": {
            "date_histogram": {
                "field": "timestamp_millis",
                "interval": "15s",
                "format": "yyyy-MM-dd HH:mm:ss",
                "min_doc_count": 0
            }
        }
    }
}`
	var result []map[string]interface{}
	if res, err := client.Search(&result, "zipkin*span-2019-09-09", dsl, SearchWithSize(100)); err != nil {
		panic(err)
	} else {
		fmt.Println(fmt.Sprintf("took: %d, total: %s, current: %d", res.Took, res.Hits.Total, len(result)))
		for _, span := range result {
			span = span["_source"].(map[string]interface{})
			//fmt.Println(span)
			fmt.Println(fmt.Sprintf("kind: %s,  serviceName: %s, serviceInstance: %s, timestamp: %s", span["kind"], span["localEndpoint"].(map[string]interface{})["serviceName"], span["tags"].(map[string]interface{})["serverInstance"], strconv.FormatFloat(span["timestamp"].(float64), 'f', 0, 64)))
		}
		for _, value := range res.Aggregations["tps"].(map[string]interface{})["buckets"].([]interface{}) {
			valueMap := value.(map[string]interface{})
			fmt.Println(fmt.Sprintf("%s ==> %v", valueMap["key_as_string"], valueMap["doc_count"]))
		}
	}
}

func Test95time(t *testing.T) {
	client := New([]string{"http://localhost:9200"})
	dsl := `{
    "query": {
        "bool": {
            "filter": {
                "bool": {
                    "must": [
                        {
                            "range": {
                                "timestamp": {
                                    "gte": 1567586321608992
                                }
                            }
                        },
                        {
                            "term": {
                                "_q": "RTCode=0"
                            }
                        }
                    ]
                }
            },
            "must": [
                {
                    "match": {
                        "localEndpoint.serviceName": "withdraw"
                    }
                },
                {
                    "match": {
                        "kind": "SERVER"
                    }
                }
            ]
        }
    },
    "aggs": {
        "time": {
            "date_histogram": {
                "field": "timestamp_millis",
                "interval": "15s",
                "format": "yyyy-MM-dd hh:mm:ss",
                "min_doc_count": 0
            },
            "aggs": {
                "duration": {
                    "percentiles": {
                        "field": "duration",
                        "percents": [
                            95,
                            99
                        ]
                    }
                },
                "duration_avg": {
                    "avg": {
                        "field": "duration"
                    }
                }
            }
        }
    }
}`
	var result []map[string]interface{}
	if res, err := client.Search(&result, "zipkin*span-2019-09-09", dsl, SearchWithSize(100)); err != nil {
		panic(err)
	} else {
		fmt.Println(fmt.Sprintf("took: %d, total: %s, current: %d", res.Took, res.Hits.Total, len(result)))
		for _, span := range result {
			span = span["_source"].(map[string]interface{})
			//fmt.Println(span)
			fmt.Println(fmt.Sprintf("kind: %s,  serviceName: %s, serviceInstance: %s, timestamp: %s, duration: %v", span["kind"], span["localEndpoint"].(map[string]interface{})["serviceName"], span["tags"].(map[string]interface{})["serverInstance"], strconv.FormatFloat(span["timestamp"].(float64), 'f', 0, 64), span["duration"]))
		}
		for _, value := range res.Aggregations["time"].(map[string]interface{})["buckets"].([]interface{}) {
			valueMap := value.(map[string]interface{})
			fmt.Println(fmt.Sprintf("%s ==> %v", valueMap["key_as_string"], valueMap["doc_count"]))
			fmt.Println(valueMap)
		}
	}
}

func TestEasyAPI_CreateIndex(t *testing.T) {
	client := New([]string{"http://172.21.82.11:9200"})
	var r = make(map[string]string)
	r["title"] = "hello"
	r["name"] = "jack"
	if b, err := json.Marshal(&r); err != nil {
		panic(err)
	} else {
		fmt.Println(string(b))
	}

	if err := client.Index("my-test-index", "", "", &r); err != nil {
		panic(err)
	}
}

func TestEasyAPI_Delete(t *testing.T) {
	client := New([]string{"http://172.21.82.11:9200"})
	if err := client.Delete("my-test-index", "ydSBBG0Bs5IKAp01iTEi"); err != nil {
		panic(err)
	}
}

func TestEasyAPI_Update(t *testing.T) {
	client := New([]string{"http://172.21.82.11:9200"})
	var r = make(map[string]string)
	r["title"] = "wa kakak"
	r["name"] = "love"
	if err := client.Update("my-test-index", "H9RlBG0Bs5IKAp01GxlA", &r); err != nil {
		panic(err)
	}
}

func TestEasyAPI_Get(t *testing.T) {
	client := New([]string{"http://172.21.82.11:9200"})
	var result map[string]interface{}
	if err := client.Get("my-test-index", "H9RlBG0Bs5IKAp01GxlA", &result); err != nil {
		panic(err)
	}
	fmt.Println(result)
}

func TestEasyAPI_Search_test(t *testing.T) {
	client := New([]string{"http://172.21.108.131:9200"})
	dsl := `{
    "query": {
        "bool": {
            "must": [
				{
					"term": {
						"_q": "originalTraceId=e4b8a734-fd42-bd7f-8ee5-6d21e3d179e8"
					}
				}
            ]
        }
    }
}`
	var result []map[string]interface{}
	if res, err := client.Search(&result, "zipkin:span*", dsl, SearchWithSize(20)); err != nil {
		panic(err)
	} else {
		fmt.Println(fmt.Sprintf("took: %d, total: %v, current: %d", res.Took, res.Hits.Total, len(result)))
	}
}

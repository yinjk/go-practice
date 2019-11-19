package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/olivere/elastic"
	"reflect"
	"strconv"
)

func main() {
	// Create a client
	client, err := elastic.NewClient(elastic.SetURL("http://172.21.108.131:9200"), elastic.SetSniff(false))
	if err != nil {
		panic(err)
	}

	// filters
	timeFilter := elastic.NewRangeQuery("timestamp").Gte(1567996860000000)
	rtCodeFilter := elastic.NewTermQuery("_q", "RTCode=0")
	orgFilter := elastic.NewTermQuery("_q", "serverOrg=ORG001")
	azFilter := elastic.NewTermQuery("_q", "serverAz=AZ0001")
	dcnFilter := elastic.NewTermQuery("_q", "serverDcn=DCN001")
	instanceFilter := elastic.NewTermQuery("_q", "serverInstance=INS001")

	//querys
	serviceNameQuery := elastic.NewTermQuery("localEndpoint.serviceName", "dtsagent")
	kindQuery := elastic.NewTermQuery("kind", "SERVER")

	//aggregations
	dateAgg := elastic.NewDateHistogramAggregation().Field("timestamp_millis").Interval("15s").Format("yyyy-MM-dd HH:mm:ss").MinDocCount(0).ExtendedBounds("2019-09-09 02:41:00", "2019-09-09 02:46:00")

	query := elastic.NewBoolQuery().Filter(timeFilter, rtCodeFilter, orgFilter, azFilter, dcnFilter, instanceFilter).Must(serviceNameQuery, kindQuery)
	result, err := client.Search("zipkin*span-2019-09-09").Query(query).Aggregation("tps", dateAgg).Size(100).RestTotalHitsAsInt(true).Do(context.Background())
	if err != nil {
		panic(err)
	}
	if result.Hits.TotalHits != 0 {
		fmt.Println(fmt.Sprintf("took: %d, total: %d, current: %d", result.TookInMillis, result.Hits.TotalHits, len(result.Hits.Hits)))
		var res map[string]interface{}
		for _, item := range result.Each(reflect.TypeOf(res)) {
			span := item.(map[string]interface{})
			fmt.Println(fmt.Sprintf("kind: %s,  serviceName: %s, serviceInstance: %s, timestamp: %s", span["kind"], span["localEndpoint"].(map[string]interface{})["serviceName"], span["tags"].(map[string]interface{})["serverInstance"], strconv.FormatFloat(span["timestamp"].(float64), 'f', 0, 64)))
		}
	}
	items, found := result.Aggregations.DateHistogram("tps")
	if found {
		for _, v := range items.Buckets {
			fmt.Println(fmt.Sprintf("%s ==> %v", *v.KeyAsString, v.DocCount))
		}
	}
	bytes, err := json.Marshal(result.Aggregations)
	if err != nil {
		panic(err)
	}
	var res map[string]interface{}
	if err = json.Unmarshal(bytes, &res); err != nil {
		panic(err)
	}
	fmt.Println(res)
}

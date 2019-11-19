package main

import (
	"encoding/json"
	"fmt"
	"github.com/olivere/elastic"
	"testing"
)

func Test_test(t *testing.T) {
	dateAgg := elastic.NewDateHistogramAggregation().Field("timestamp_millis").Interval("15s").Format("2006-01-02 15:04:05").MinDocCount(0).ExtendedBounds("2016-08-19", "2016-10-19")
	rtPercentile := elastic.NewPercentilesAggregation().Field("duration").Percentiles(95, 99)
	rtAvg := elastic.NewAvgAggregation().Field("duration")
	dateAgg = dateAgg.SubAggregation("rt_percentile", rtPercentile)
	dateAgg = dateAgg.SubAggregation("rt_avg", rtAvg)
	source, _ := dateAgg.Source()
	bytes, e := json.Marshal(source)
	if e != nil {
		panic(e)
	}
	fmt.Println(string(bytes))
}

func Test_test1(t *testing.T) {

}

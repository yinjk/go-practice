package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	es7 "github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"log"
	"strings"
	//es8 "github.com/elastic/go-elasticsearch/v8"
)

func main() {

}

type easyAPI struct {
	client *es7.Client
}

func New(esAddr []string) *easyAPI {
	config := es7.Config{
		Addresses: esAddr,
	}
	client, err := es7.NewClient(config)
	if err != nil {
		panic(err)
	}
	return &easyAPI{client: client}
}

//Index 索引文档，在指定index内创建document
func (e *easyAPI) Index(index, documentType, documentId string, v interface{}) (err error) {
	var b []byte
	if b, err = json.Marshal(v); err != nil {
		log.Printf("EsayApi Index() method unmarshal struct failed, errors(%s)\n", err.Error())
		return err
	}

	// Set up the request object.
	req := esapi.IndexRequest{
		Index:        index,
		DocumentID:   documentId,
		DocumentType: documentType,
		Body:         strings.NewReader(string(b)),
		Refresh:      "true",
	}

	// Perform the request with the client.
	res, err := req.Do(context.Background(), e.client)
	if err != nil {
		log.Printf("Es Index() method req.Do failed, errors(%s)\n", err.Error())
		return
	}
	defer func() { _ = res.Body.Close() }()

	if res.IsError() {
		log.Printf("[%s] Error indexing document", res.Status())
		var e map[string]interface{}
		if err = json.NewDecoder(res.Body).Decode(&e); err != nil {
			return err
		}
		reason := e["error"].(map[string]interface{})["reason"].(string)
		return errors.New(fmt.Sprintf("[%s] index document failed, the reason is: %s \n", res.Status(), reason))
	}
	return
}

//Delete delete by documentId
func (e *easyAPI) Delete(index, documentId string) (err error) {
	res, err := e.client.Delete(index, documentId,
		e.client.Delete.WithContext(context.Background()),
	)
	if err != nil {
		log.Printf("esayAPI Delete() method error, errors(%s)", err.Error())
		return
	}
	if res.IsError() {
		log.Printf("[%s] Error delete document", res.Status())
		var e map[string]interface{}
		if err = json.NewDecoder(res.Body).Decode(&e); err != nil {
			return err
		}
		var reason string
		if e["error"] != nil {
			reason = e["error"].(string)
		}
		return errors.New(fmt.Sprintf("[%s] delete document failed, the reason is: %s \n", res.Status(), reason))
	}
	return
}

//Update update document, used doc to update
func (e *easyAPI) Update(index, documentId string, v interface{}) (err error) {
	var doc = make(map[string]interface{})
	//update需要在数据外面包一个doc或script，见：https://www.elastic.co/guide/en/elasticsearch/reference/current/docs-update.html#_updates_with_a_partial_document
	doc["doc"] = v
	var b bytes.Buffer
	if err = json.NewEncoder(&b).Encode(&doc); err != nil {
		log.Printf("EsayApi Index() method unmarshal struct failed, errors(%s)\n", err.Error())
		return err
	}
	res, err := e.client.Update(index, documentId, &b,
		e.client.Update.WithContext(context.Background()))
	if err != nil {
		log.Printf("esayAPI Delete() method error, errors(%s)", err.Error())
		return
	}
	if res.IsError() {
		log.Printf("[%s] Error update document", res.Status())
		var e map[string]interface{}
		if err = json.NewDecoder(res.Body).Decode(&e); err != nil {
			return err
		}
		var reason string
		if e["error"] != nil {
			reason = e["error"].(map[string]interface{})["reason"].(string)
		}
		return errors.New(fmt.Sprintf("[%s] update document failed, the reason is: %s \n", res.Status(), reason))
	}
	return
}

//Get get by documentId
func (e *easyAPI) Get(index, documentId string, result interface{}) (err error) {
	res, err := e.client.Get(index, documentId,
		e.client.Get.WithRefresh(true),
		e.client.Get.WithContext(context.Background()))
	if err != nil {
		log.Printf("esayAPI Get() method error, errors(%s)", err.Error())
		return
	}
	if res.IsError() {
		log.Printf("[%s] Error get document", res.Status())
		var e map[string]interface{}
		if err = json.NewDecoder(res.Body).Decode(&e); err != nil {
			return err
		}
		var reason string
		if e["error"] != nil {
			reason = e["error"].(map[string]interface{})["reason"].(string)
		}
		return errors.New(fmt.Sprintf("[%s] get document failed, the reason is: %s \n", res.Status(), reason))
	}
	var getRes GetRes
	getRes.Source = result
	if err = json.NewDecoder(res.Body).Decode(&getRes); err != nil {
		log.Printf("esayAPI Get() method error, decode response errors(%v)", err)
		return err
	}
	return
}

func SearchWithSize(size int) func(*esapi.SearchRequest) {
	return func(r *esapi.SearchRequest) {
		r.Size = &size
	}
}
func SearchWithFrom(from int) func(*esapi.SearchRequest) {
	return func(r *esapi.SearchRequest) {
		r.From = &from
	}
}

//Search search with dsl
func (e *easyAPI) Search(result interface{}, index, dsl string, options ...func(*esapi.SearchRequest)) (searchRes *SearchRes, err error) {
	var buf bytes.Buffer
	buf.WriteString(dsl)
	// Perform the search request.
	var searchOptions = []func(*esapi.SearchRequest){
		e.client.Search.WithContext(context.Background()),
		e.client.Search.WithIndex(index),
		e.client.Search.WithBody(&buf),
		e.client.Search.WithTrackTotalHits(true),
		e.client.Search.WithPretty(),
	}
	for _, option := range options {
		searchOptions = append(searchOptions, option)
	}
	res, err := e.client.Search(
		searchOptions...,
	)
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			log.Fatalf("Error parsing the response body: %s", err)
			return nil, err
		} else {
			// Print the response status and error information.
			errMsg := fmt.Sprintf("[%s] %s: %s",
				res.Status(),
				e["error"].(map[string]interface{})["type"],
				e["error"].(map[string]interface{})["reason"],
			)
			return nil, errors.New(fmt.Sprintf("ES search() method error, errors(%s)\n", errMsg))
		}
	}
	searchRes = &SearchRes{}
	if result != nil {
		searchRes.Hits.Hits = result
	}
	if err := json.NewDecoder(res.Body).Decode(&searchRes); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}
	return
}

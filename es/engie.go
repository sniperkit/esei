package es

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"os"

	elastic "gopkg.in/olivere/elastic.v5"
)

type EsInfo struct {
	EsURL   string
	EsUser  string
	EsPass  string
	EsIndex string
	EsType  string
	EsSize  int
	IsDebug bool
	client  *elastic.Client
	ctx     context.Context
}

// Do Elasticsearch入口函数
func (e *EsInfo) Do(file string) error {
	var err error
	e.client, err = e.clientInit()
	if err != nil {
		return err
	}

	e.ctx = context.Background()
	data, hit, err := e.search()
	if err != nil {
		return err
	}

	log.Printf("Total Hit [%d] Once Receive [%d]\n", hit, len(data))

	fs, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}

	for _, d := range data {
		fs.WriteString(d + "\n")

	}
	return nil
}

// clientInit 使用用户提供的参数初始化ElasticSearch Client
func (e *EsInfo) clientInit() (*elastic.Client, error) {
	var client *elastic.Client
	var err error

	if (e.EsUser != "") && (e.EsPass != "") {
		if e.IsDebug {
			log.Println(e.EsURL, e.EsUser, e.EsPass)
			client, err = elastic.NewClient(elastic.SetTraceLog(log.New(os.Stdout, "", 0)), elastic.SetSniff(false), elastic.SetURL(e.EsURL), elastic.SetBasicAuth(e.EsUser, e.EsPass))
		} else {
			client, err = elastic.NewClient(elastic.SetSniff(false), elastic.SetURL(e.EsURL), elastic.SetBasicAuth(e.EsUser, e.EsPass))
		}

	} else {
		if e.IsDebug {
			log.Println(e.EsURL, e.EsUser, e.EsPass)
			client, err = elastic.NewClient(elastic.SetTraceLog(log.New(os.Stdout, "", 0)), elastic.SetSniff(false), elastic.SetURL(e.EsURL))
		} else {
			client, err = elastic.NewClient(elastic.SetSniff(false), elastic.SetURL(e.EsURL))
		}
	}

	return client, err
}

// search 设定检索条件，返回查询的数据,检索的总记录条数
func (e *EsInfo) search() ([]string, int64, error) {
	var result []string
	searchResult, err := e.client.Search().
		Index(e.EsIndex).
		Query(nil).
		From(0).Size(e.EsSize).
		Pretty(true).
		Do(e.ctx)
	if err != nil {
		return nil, 0, errors.New("Search ElasticSearch Error. " + err.Error())
	}

	for _, hit := range searchResult.Hits.Hits {
		// var content []byte
		content, err := json.Marshal(hit.Source)
		if err != nil {
			log.Println(err.Error())
			break
		}

		result = append(result, string(content))
	}
	return result, searchResult.Hits.TotalHits, nil
}

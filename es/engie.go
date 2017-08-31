package es

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	elastic "gopkg.in/olivere/elastic.v5"
)

// EsInfo 封装的ElasticSearch操作数据
type EsInfo struct {
	EsURL   string
	EsUser  string
	EsPass  string
	EsIndex string
	EsType  string
	EsSize  int
	IsDebug bool
	// Mode 运行模式 0 导出/ 1 导入
	Mode   int
	client *elastic.Client
	ctx    context.Context
}

// Do Elasticsearch入口函数
func (e *EsInfo) Do(file string) error {
	var err error
	e.client, err = e.clientInit()
	if err != nil {
		return err
	}

	e.ctx = context.Background()

	switch e.Mode {
	case 0:
		data, err := e.export()
		if err != nil {
			return err
		}

		log.Printf("Total Hit [%d] \n", len(data))

		fs, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE, 0755)
		if err != nil {
			return err
		}

		defer func() {
			fs.Close()
		}()
		for _, d := range data {
			fs.WriteString(d + "\n")

		}
	case 1:
		fs, err := os.OpenFile(file, os.O_RDONLY, 0755)
		if err != nil {
			return err
		}

		defer func() {
			fs.Close()
		}()
		var data []string

		scanner := bufio.NewScanner(fs)
		for scanner.Scan() {
			data = append(data, scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			return err
		}

		ticker := make(chan int, 1)
		go func() {
			for {
				now := time.Now()
				next := now.Add(time.Second * time.Duration(30))
				next = time.Date(next.Year(), next.Month(), next.Day(), next.Hour(), next.Minute(), next.Second(), 0, next.Location())
				ti := time.NewTimer(next.Sub(now))
				select {
				case <-ti.C:
					ticker <- 1
				}
			}
		}()
		return e.esimport(data, ticker)

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
func (e *EsInfo) export() ([]string, error) {
	var result []string
	searchResult, err := e.client.Search().
		Index(e.EsIndex).
		Query(nil).
		From(0).Size(e.EsSize).
		Pretty(true).
		Do(e.ctx)
	if err != nil {
		return nil, errors.New("Search ElasticSearch Error. " + err.Error())
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
	return result, nil
}

// esimport 加载指定数据到ElasticSearch中
// data 从文件中加载的数据,必须为json格式
func (e *EsInfo) esimport(data []string, ticker chan int) error {
	fmt.Printf("ESEI Will Load [%d] Records \n", len(data))
	retry := false
	hasRetry := 0
	breakPoint := 0
	totalLoad := 0
	tda := data[breakPoint:]
	for hasRetry < 10 {

		if retry {
			fmt.Printf("ESEI waitting %ds and will try to reimport.", hasRetry*5)
			fmt.Printf("ESEI Has Retry [%d] \n", hasRetry)
			time.Sleep(time.Duration(5*hasRetry) * time.Second)
			tda = tda[breakPoint:]
		}

		retry = false
		for i, d := range tda {
			_, err := e.client.Index().
				Index(e.EsIndex).
				Type(e.EsType).
				BodyString(d).
				Do(e.ctx)

			if err != nil {
				fmt.Println("Insert ElasticSearch Error. " + err.Error())
				retry = true
				hasRetry++
				breakPoint = i
				break
			}
			// 恢复计数器
			hasRetry = 0

			select {
			case _, ok := <-ticker:
				if ok {
					totalLoad = totalLoad + i
					fmt.Printf("ESEI Has Load [%d] Records \n", totalLoad)
				}
			default:

			}

		}

		if !retry {
			break
		}
	}

	fmt.Println("ESEI Load Complete!")
	return nil
}

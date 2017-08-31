package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/andy-zhangtao/esei/es"
)

const (
	// ESURLEMPTY url不得为空
	ESURLEMPTY = "esurl can't be empty"
	// ESIDXEMPTY index不得为空
	ESIDXEMPTY = "index can't be empty"
)

// esurl Elasticsearch访问地址
var esurl string

// esuser ElasticSearch用户名,如果启用X-Pack，则必填
var esuser string

// espasswd ElasticSearch口令,如果启用X-Pack，则必填
var espasswd string

// esindex ElasticSearch索引名称
var esindex string

// estype ElasticSearch类型名称
var estype string

// debug 是否显示所有debug信息
var debug bool

// file 导出文件名称
var file string

// size 处理条数
var size int

func init() {
	flag.StringVar(&esurl, "esurl", "", "The URL of Elasticsearch")
	flag.StringVar(&esuser, "user", "", "The user name of Elasticsearch. If you enable X-Pack, Maybe you should tell me this value")
	flag.StringVar(&espasswd, "passwd", "", "The user password of Elasticsearch. If you enable X-Pack, Maybe you should tell me this value")
	flag.StringVar(&esindex, "index", "", "The index name you want to export/import")
	flag.StringVar(&estype, "type", "", "The type name you want to export/import")
	flag.BoolVar(&debug, "debug", false, "If true, ESEI will show all ElasticSearch info. Default is false")
	flag.StringVar(&file, "out", "out.json", "The output file name, default is out.json")

	flag.IntVar(&size, "size", 10, "Export/Import the number of records")
}

func isParaValid() error {
	if esurl == "" {
		return errors.New(ESURLEMPTY)
	}

	if esindex == "" {
		return errors.New(ESIDXEMPTY)
	}

	return nil
}
func main() {
	flag.Parse()

	err := isParaValid()
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("------------")
		fmt.Println("[esei -h] get more info")
		os.Exit(-1)
	}

	ei := es.EsInfo{
		EsURL:   esurl,
		EsUser:  esuser,
		EsPass:  espasswd,
		EsIndex: esindex,
		EsType:  estype,
		EsSize:  size,
		IsDebug: debug,
	}

	err = ei.Do(file)
	if err != nil {
		log.Println(err.Error())
	}
}

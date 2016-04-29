package main

import (
    "fmt"

	"github.com/stefanoschrs/proxymeister/sqlite"
	"github.com/stefanoschrs/proxymeister/crawler"
)

func fetchAndInsertProxies(){
	const dbpath = "./data/proxy.db"
	db := sqlite.InitDB(dbpath)
	defer db.Close()
	sqlite.CreateTable(db)

	for _, proxy := range crawler.FetchProxies(){
		sqlite.InsertProxy(db, proxy.Ip, proxy.Port)
	}

	fmt.Println(sqlite.SelectAllProxies(db))
}

func main(){
	fetchAndInsertProxies()
}

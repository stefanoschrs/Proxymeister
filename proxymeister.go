package main

import (
    "fmt"

	"github.com/stefanoschrs/proxymeister/sqlite"
	"github.com/stefanoschrs/proxymeister/crawler"
	"github.com/stefanoschrs/proxymeister/validator"
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

func validateProxies(){
	myIp := validator.GetMyIp()
	fmt.Printf("My IP: %s\n", myIp)

	const dbpath = "./data/proxy.db"
	db := sqlite.InitDB(dbpath)
	defer db.Close()

	for _, proxy := range sqlite.SelectRecentProxies(db){
		proxyUrl := fmt.Sprintf("http://%s:%d", proxy.Ip, proxy.Port)
		fmt.Printf("%-30s", proxyUrl)
		isValid := validator.Validate(myIp, proxyUrl)

		if isValid {
			fmt.Printf(" Valid\n")
			sqlite.UpdateProxy(db, proxy.Ip, proxy.Port, 1)
		} else {
			fmt.Printf(" Not Valid\n")
			sqlite.UpdateProxy(db, proxy.Ip, proxy.Port, 2)
		}
	}
}

func main(){
	// fetchAndInsertProxies()
	validateProxies()
}

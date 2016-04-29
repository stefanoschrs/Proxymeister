package main

import (
    "fmt"
	"time"

	"github.com/stefanoschrs/proxymeister/sqlite"
	"github.com/stefanoschrs/proxymeister/crawler"
	"github.com/stefanoschrs/proxymeister/validator"
)

const FETCH_INTERVAL = time.Hour * 1
const VALIDATE_INTERVAL = time.Minute* 15
const DB_PATH = "./data/proxy.db"

func fetchAndInsertProxies(){
	fmt.Printf("%-25s Interval:%s\n", "Proxy Fetcher Started!", FETCH_INTERVAL)
	for range time.Tick(FETCH_INTERVAL){
		db := sqlite.InitDB(DB_PATH)
		defer db.Close()
		sqlite.CreateTable(db)

		for _, proxy := range crawler.FetchProxies(){
			sqlite.InsertProxy(db, proxy.Ip, proxy.Port)
		}

		fmt.Println(sqlite.SelectAllProxies(db))
	}
}

func validateProxies(){
	fmt.Printf("%-25s Interval:%s\n", "Proxy Validator Started!", VALIDATE_INTERVAL)
	for range time.Tick(VALIDATE_INTERVAL){
		myIp := validator.GetMyIp()
		fmt.Printf("My IP: %s\n", myIp)

		db := sqlite.InitDB(DB_PATH)
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
}

func main(){
	fmt.Println(" _______                                                                  __              __                                ")
	fmt.Println("/       \\                                                                /  |            /  |                              ")
	fmt.Println("$$$$$$$  | ______    ______   __    __  __    __  _____  ____    ______  $$/   _______  _$$ |_     ______    ______         ")
	fmt.Println("$$ |__$$ |/      \\  /      \\ /  \\  /  |/  |  /  |/     \\/    \\  /      \\ /  | /       |/ $$   |   /      \\  /      \\ ")
	fmt.Println("$$    $$//$$$$$$  |/$$$$$$  |$$  \\/$$/ $$ |  $$ |$$$$$$ $$$$  |/$$$$$$  |$$ |/$$$$$$$/ $$$$$$/   /$$$$$$  |/$$$$$$  |      ")
	fmt.Println("$$$$$$$/ $$ |  $$/ $$ |  $$ | $$  $$<  $$ |  $$ |$$ | $$ | $$ |$$    $$ |$$ |$$      \\   $$ | __ $$    $$ |$$ |  $$/       ")
	fmt.Println("$$ |     $$ |      $$ \\__$$ | /$$$$  \\ $$ \\__$$ |$$ | $$ | $$ |$$$$$$$$/ $$ | $$$$$$  |  $$ |/  |$$$$$$$$/ $$ |          ")
	fmt.Println("$$ |     $$ |      $$    $$/ /$$/ $$  |$$    $$ |$$ | $$ | $$ |$$       |$$ |/     $$/   $$  $$/ $$       |$$ |             ")
	fmt.Println("$$/      $$/        $$$$$$/  $$/   $$/  $$$$$$$ |$$/  $$/  $$/  $$$$$$$/ $$/ $$$$$$$/     $$$$/   $$$$$$$/ $$/              ")
	fmt.Println("                                       /  \\__$$ |                                                                          ")
	fmt.Println("                                       $$    $$/                                                                            ")
	fmt.Println("                                        $$$$$$/                                                               v.1.0.0       ")
	fmt.Println()

	go fetchAndInsertProxies()
	go validateProxies()

	time.Sleep(time.Hour * 2)
}

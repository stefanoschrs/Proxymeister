package main

import (
    "fmt"
	"time"
    "log"
	"strconv"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/stefanoschrs/proxymeister/sqlite"
	"github.com/stefanoschrs/proxymeister/crawler"
	"github.com/stefanoschrs/proxymeister/validator"
)

const FETCH_INTERVAL = time.Minute * 30
const VALIDATE_INTERVAL = time.Minute * 10
const DB_PATH = "./data/proxy.db"
const PORT = 5000

func fetchAndInsertProxies(){
	fmt.Printf("> %-25s%-10s %s\n", "Proxy Fetcher Started!", "Interval:", FETCH_INTERVAL)
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
	fmt.Printf("> %-25s%-10s %s\n", "Proxy Validator Started!", "Interval:", VALIDATE_INTERVAL)
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

func IndexHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "Yo! Welcome to Proxymeister's API")
}

func ProxyHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var count int

	if s, err := strconv.ParseInt(vars["count"], 10, 32); err == nil {
		count = int(s)
	} else{
		count = 1
	}

	db := sqlite.InitDB(DB_PATH)
	defer db.Close()

	var proxyList []string
	for _, proxy := range sqlite.SelectRandomProxies(db, count){
		proxyList = append(proxyList, fmt.Sprintf("%s:%d", proxy.Ip, proxy.Port))
	}

	json.NewEncoder(w).Encode(proxyList)
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

	router := mux.NewRouter().StrictSlash(true)
    router.HandleFunc("/", IndexHandler)
    router.HandleFunc("/proxy", ProxyHandler)
    router.HandleFunc("/proxy/{count}", ProxyHandler)

	fmt.Printf("> %-25s%-10s %d\n", "API Server Started", "Port:", PORT)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", PORT), router))
}

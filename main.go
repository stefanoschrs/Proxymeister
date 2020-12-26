package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/stefanoschrs/proxymeister/internal/cron"
	"github.com/stefanoschrs/proxymeister/internal/database"
	"github.com/stefanoschrs/proxymeister/internal/webserver"

	"github.com/stefanoschrs/proxymeister/pkg/utils"

	"github.com/gin-gonic/gin"
)

func printIntro() {
	i := "\033[36;1m"
	i2 := "\033[34;1m"
	d := "\033[0m"
	fmt.Println(i + "______                    " + i2 + "                _     _" + d)
	fmt.Println(i + "| ___ \\                   " + i2 + "               (_)   | |" + d)
	fmt.Println(i + "| |_/ / __ _____  ___   _ " + i2 + " _ __ ___   ___ _ ___| |_ ___ _ __" + d)
	fmt.Println(i + "|  __/ '__/ _ \\ \\/ / | | |" + i2 + "| '_ ` _ \\ / _ \\ / __| __/ _ \\ '__|" + d)
	fmt.Println(i + "| |  | | | (_) >  <| |_| |" + i2 + "| | | | | |  __/ \\__ \\ ||  __/ |" + d)
	fmt.Println(i + "\\_|  |_|  \\___/_/\\_\\\\__, |" + i2 + "|_| |_| |_|\\___|_|___/\\__\\___|_|" + d)
	fmt.Println(i + "                     __/ |" + d)
	fmt.Println(i + "	            |___/ " + d)
	fmt.Println()
}

func main() {
	printIntro()

	// ------------------------- Initialize Config ------------------------- //
	err := utils.LoadConfig()
	if err != nil {
		log.Fatal("config.LoadConfig", err)
	}

	// ------------------------ Initialize Database ------------------------ //
	db, err := database.Init()
	if err != nil {
		log.Fatal("database.Init", err)
	}

	// ----------------------- Initialize Random Seed ---------------------- //
	rand.Seed(time.Now().UTC().UnixNano())

	// ------------------------------- Cron -------------------------------- //
	// TODO: Create cron service, give access from API to re-run job etc
	_, err = cron.Init(db)
	if err != nil {
		log.Fatal("cron.Init", err)
	}

	// -------------------------------- Gin -------------------------------- //
	router := webserver.Init()

	webserver.SetMiddleware(router)
	router.Use(func(c *gin.Context) {
		c.Set("db", db)
	})

	webserver.SetRoutes(router)

	if !gin.IsDebugging() {
		port := os.Getenv("PORT")
		if port == "" {
			port = "8080"
		}

		log.Printf("Listening on :%s...\n", port)
	}

	err = router.Run()
	if err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"github.com/stefanoschrs/proxymeister/pkg/utils"
	"log"

	"github.com/stefanoschrs/proxymeister/internal/database"
	"github.com/stefanoschrs/proxymeister/pkg/types"
)

func main() {
	err := utils.LoadEnv()
	if err != nil {
		log.Fatal(err)
	}

	db, err := database.Init()
	if err != nil {
		log.Fatal(err)
	}

	// Migrate
	err = db.AutoMigrate(
		&types.Proxy{},
	)
	if err != nil {
		log.Fatal(err)
	}
}

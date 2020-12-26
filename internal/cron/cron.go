package cron

import (
	"fmt"
	"log"

	"github.com/stefanoschrs/proxymeister/internal/database"
	"github.com/stefanoschrs/proxymeister/pkg/crawler"

	"github.com/robfig/cron/v3"
)

func fetchProxies(db database.DB) {
	proxies, crawlerErr := crawler.FetchProxies()
	if crawlerErr != nil {
		log.Println(crawlerErr)
		return
	}

	fmt.Printf("Found %d proxies\n", len(proxies))
	for _, p := range proxies {
		proxy, created, createErr := db.CreateProxy(p)
		if createErr != nil {
			log.Println(crawlerErr)
			continue
		}
		if created {
			fmt.Printf("Proxy %s:%d added!\n", proxy.Ip, proxy.Port)
		}
	}
}

func Init(db database.DB) (c *cron.Cron, err error) {
	c = cron.New()

	// fetchProxies
	_, err = c.AddFunc("@midnight", func() {
		fetchProxies(db)
	})
	if err != nil {
		err = fmt.Errorf("failed to add fetchProxies func. %w", err)
		return
	}

	//log.Println("Fetching proxies..")
	//fetchProxies(db)

	c.Start()

	return
}

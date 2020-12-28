package cron

import (
	"fmt"
	"github.com/stefanoschrs/proxymeister/internal/logging"
	"log"
	"strings"

	"github.com/stefanoschrs/proxymeister/internal/database"

	"github.com/stefanoschrs/proxymeister/pkg/crawler"
	"github.com/stefanoschrs/proxymeister/pkg/types"
	"github.com/stefanoschrs/proxymeister/pkg/utils"
	"github.com/stefanoschrs/proxymeister/pkg/validator"

	"github.com/robfig/cron/v3"
)

func fetchProxies(db database.DB) {
	logging.Debug("Fetching proxies..")

	proxies, err := crawler.FetchProxies()
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Printf("Found %d proxies\n", len(proxies))
	for _, p := range proxies {
		proxy, created, createErr := db.CreateProxy(p)
		if createErr != nil {
			log.Println(err)
			continue
		}
		if created {
			fmt.Printf("Proxy %s:%d added!\n", proxy.Ip, proxy.Port)
		}
	}
}

func checkProxies(db database.DB) {
	logging.Debug("Checking proxies..")

	type empty struct{}
	type resultTry struct {
		Latency int64
		Error   error
	}
	type result struct {
		Proxy types.Proxy
		Tries []resultTry
	}

	workerCount := 10
	validationTries := 2

	workerData := make(chan types.Proxy, workerCount)
	gather := make(chan result)
	tracker := make(chan empty)

	var results []result

	myIp, err := utils.GetMyIP()
	if err != nil {
		logging.Error("utils.GetMyIP", err)
		return
	}

	// Initialize workers
	for i := 0; i < workerCount; i++ {
		go func(t chan empty, w chan types.Proxy, g chan result) {
			for proxy := range w {
				logging.Debugf("Processing %s:%d..", proxy.Ip, proxy.Port)

				var tries []resultTry
				for j := 0; j < validationTries; j++ {
					latency, validationErr := validator.Validate(myIp, proxy.Ip, proxy.Port, true)

					tries = append(tries, resultTry{
						Latency: latency,
						Error:   validationErr,
					})
				}

				g <- result{
					proxy,
					tries,
				}
			}
			tracker <- empty{}
		}(tracker, workerData, gather)
	}

	// Gather results
	go func() {
		for r := range gather {
			results = append(results, r)
		}
		tracker <- empty{}
	}()

	// Add data to be processed
	proxies, err := db.GetProxies()
	if err != nil {
		logging.Error("db.GetProxies", err)
		return
	}

	logging.Debugf("Found %d proxies", len(proxies))
	for _, proxy := range proxies {
		workerData <- proxy
	}
	close(workerData)

	// Track
	for i := 0; i < workerCount; i++ {
		<-tracker
	}
	close(gather)

	<-tracker
	close(tracker)

	for _, r := range results {
		logging.Debugf("Updating %s:%d", r.Proxy.Ip, r.Proxy.Port)

		var sumLatency int64
		var successfulTries int64
		for _, try := range r.Tries {
			if try.Error != nil {
				processValidationError(try.Error)
				continue
			}
			sumLatency += try.Latency
			successfulTries += 1
		}
		if sumLatency > 0 {
			r.Proxy.Status = types.ProxyStatusActive
			r.Proxy.Latency = sumLatency / successfulTries
			r.Proxy.FailedChecks = 0
		} else {
			r.Proxy.Status = types.ProxyStatusInactive
			r.Proxy.Latency = 0
			r.Proxy.FailedChecks += 1
		}

		err = db.UpdateProxy(r.Proxy)
		if err != nil {
			logging.Error("failed to update proxy",
				"id", r.Proxy.ID,
				"err", err)
		}
	}
}

func processValidationError(err error) {
	// TODO: Detect all errors and handle accordingly

	if strings.Contains(err.Error(), "context deadline exceeded") {
		return
	} else if strings.Contains(err.Error(), "unexpected EOF") {
		return
	} else if strings.Contains(err.Error(), "Too many open connections") {
		return
	} else if strings.Contains(err.Error(), "Proxy Authentication Required") {
		return
	}

	logging.Error(err)
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
	//fetchProxies(db)

	// checkProxies
	_, err = c.AddFunc("0 */2 * * *", func() {
		checkProxies(db)
	})
	if err != nil {
		err = fmt.Errorf("failed to add checkProxies func. %w", err)
		return
	}
	//checkProxies(db)

	c.Start()

	return
}

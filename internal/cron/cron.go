package cron

import (
	"fmt"
	"log"
	"strings"

	"github.com/stefanoschrs/proxymeister/internal/database"
	"github.com/stefanoschrs/proxymeister/internal/logging"

	"github.com/stefanoschrs/proxymeister/pkg/crawler"
	"github.com/stefanoschrs/proxymeister/pkg/types"
	"github.com/stefanoschrs/proxymeister/pkg/utils"
	"github.com/stefanoschrs/proxymeister/pkg/validator"

	"github.com/robfig/cron/v3"
)

func FetchProxies(db database.DB) {
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

func CheckProxies(db database.DB) {
	logging.Debug("Checking proxies..")

	type empty struct{}

	// TODO: Retrieve from .env
	workerCount := 10
	validationTries := 3

	workerData := make(chan types.Proxy, workerCount)
	tracker := make(chan empty)

	myIp, err := utils.GetMyIP()
	if err != nil {
		logging.Error("utils.GetMyIP", err)
		return
	}

	var workingProxies uint

	// Initialize workers
	for i := 0; i < workerCount; i++ {
		// TODO: Move func
		go func(t chan empty, w chan types.Proxy) {
			for proxy := range w {
				logging.Debugf("Processing %s:%d ..", proxy.Ip, proxy.Port)

				var sumLatency int64
				var successfulTries int64

				for j := 0; j < validationTries; j++ {
					latency, validationErr := validator.Validate(myIp, proxy.Ip, proxy.Port, true)
					if validationErr != nil {
						//logging.Debugf("[%s] %v", proxy.Ip, validationErr.Error())
						processValidationError(validationErr)
						continue
					}
					//logging.Debugf("[%s] Success", proxy.Ip)

					sumLatency += latency
					successfulTries += 1
				}

				if successfulTries > 1 {
					proxy.Status = types.ProxyStatusActive
					proxy.Latency = sumLatency / successfulTries
					proxy.FailedChecks = 0

					// TODO: Move somewhere else
					workingProxies++
				} else {
					proxy.Status = types.ProxyStatusInactive
					proxy.Latency = 0
					proxy.FailedChecks += 1
				}

				err = db.UpdateProxy(proxy)
				if err != nil {
					logging.Error("failed to update proxy",
						"id", proxy.ID,
						"err", err)
				}
			}
			tracker <- empty{}
		}(tracker, workerData)
	}

	// Add data to be processed
	proxies, err := db.GetProxiesInternal()
	if err != nil {
		logging.Error("db.GetProxies", err)
		return
	}
	//proxies = proxies[:10]

	logging.Debugf("Found %d proxies", len(proxies))
	for _, proxy := range proxies {
		workerData <- proxy
	}
	close(workerData)

	// Track
	for i := 0; i < workerCount; i++ {
		<-tracker
	}
	close(tracker)

	logging.Debugf("Checking finished! Working: %d/%d", workingProxies, len(proxies))
}

func Init(db database.DB) (c *cron.Cron, err error) {
	c = cron.New()

	// fetchProxies
	_, err = c.AddFunc("@midnight", func() {
		FetchProxies(db)
	})
	if err != nil {
		err = fmt.Errorf("failed to add fetchProxies func. %w", err)
		return
	}

	// checkProxies
	_, err = c.AddFunc("0 */2 * * *", func() {
		CheckProxies(db)
	})
	if err != nil {
		err = fmt.Errorf("failed to add checkProxies func. %w", err)
		return
	}

	// TODO: Clean up proxies with X failed tries

	c.Start()

	return
}

func processValidationError(err error) {
	// TODO: Detect all errors and handle accordingly

	if strings.Contains(err.Error(), "context deadline exceeded") {
		return
	} else if strings.Contains(err.Error(), "unexpected EOF") {
		return
	} else if strings.Contains(err.Error(), "certificate signed by unknown authority") {
		return
	} else if strings.Contains(err.Error(), "connect: connection refused") {
		return
	} else if strings.Contains(err.Error(), "connect: no route to host") {
		return
	} else if strings.Contains(err.Error(), "i/o timeout") {
		return
	} else if strings.Contains(err.Error(), "read: connection reset by peer") {
		return
	} else if strings.Contains(err.Error(), "tls: server chose an unconfigured cipher suite") {
		return
	} else if strings.Contains(err.Error(), "malformed HTTP status code") {
		return
	} else if strings.Contains(err.Error(), "Too many open connections") {
		return
	} else if strings.Contains(err.Error(), "Proxy Authentication Required") {
		return
	} else if strings.Contains(err.Error(), "Forbidden") {
		return
	} else if strings.Contains(err.Error(), "Forwarding failure") {
		return
	} else if strings.Contains(err.Error(), "Bad Gateway") {
		return
	} else if strings.Contains(err.Error(), "Method Not Allowed") {
		return
	} else if strings.Contains(err.Error(), "Service Unavailable") {
		return
	}

	// Custom
	if strings.Contains(err.Error(), types.ErrEmptyResponse) {
		return
	}

	logging.Error(err)
}

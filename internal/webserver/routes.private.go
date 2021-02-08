package webserver

import (
	"net/http"

	"github.com/stefanoschrs/proxymeister/internal/cron"
	"github.com/stefanoschrs/proxymeister/internal/utils"

	"github.com/gin-gonic/gin"
)

// Trigger fetching from sources
func postAdminFetch(c *gin.Context) {
	db := utils.ExtractDB(c)

	go cron.FetchProxies(db)

	c.Status(http.StatusOK)
}

// Trigger checking of proxies
func postAdminCheck(c *gin.Context) {
	db := utils.ExtractDB(c)

	go cron.CheckProxies(db)

	c.Status(http.StatusOK)
}

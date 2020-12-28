package webserver

import (
	"github.com/stefanoschrs/proxymeister/internal/cron"
	"net/http"

	"github.com/stefanoschrs/proxymeister/internal/utils"

	"github.com/gin-gonic/gin"
)

func postAdminFetch(c *gin.Context) {
	db := utils.ExtractDB(c)

	go cron.FetchProxies(db)

	c.Status(http.StatusOK)
}

func postAdminCheck(c *gin.Context) {
	db := utils.ExtractDB(c)

	go cron.CheckProxies(db)

	c.Status(http.StatusOK)
}

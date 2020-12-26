package webserver

import (
	"log"
	"net/http"

	"github.com/stefanoschrs/proxymeister/internal/utils"

	"github.com/gin-gonic/gin"
)

func headHealth(c *gin.Context) {
	c.Status(http.StatusOK)
}

func getProxies(c *gin.Context) {
	db := utils.ExtractDB(c)

	proxies, err := db.GetProxies()
	if err != nil {
		log.Println("db.GetProxies", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, proxies)
}

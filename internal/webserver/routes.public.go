package webserver

import (
	"fmt"
	"github.com/stefanoschrs/proxymeister/pkg/types"
	"log"
	"net/http"
	"strings"

	"github.com/stefanoschrs/proxymeister/internal/utils"

	"github.com/gin-gonic/gin"
)

func headHealth(c *gin.Context) {
	c.Status(http.StatusOK)
}

// @Param limit		query	int		false	-1		"Limit results to value"
// @Param latency	query	int		false	-1		"Filter by latency lower than value"
// @Param fields	query	string	false	""		"Return fields: address,ip,port,latency,updatedAt"
// @Param format	query	string	false	"json"	"Output format: json,csv"
func getProxies(c *gin.Context) {
	db := utils.ExtractDB(c)

	dbProxies, err := db.GetProxies(map[string]interface{}{
		"status":  types.ProxyStatusActive,
		"limit":   c.Query("limit"),
		"latency": c.Query("latency"),
	})
	if err != nil {
		log.Println("db.GetProxies", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	fields := c.Query("fields")
	format := c.Query("format")

	var proxies []map[string]interface{}
	for _, dbPproxy := range dbProxies {
		proxy := make(map[string]interface{})
		proxy["address"] = fmt.Sprintf("%s:%d", dbPproxy.Ip, dbPproxy.Port)
		proxy["ip"] = dbPproxy.Ip
		proxy["port"] = dbPproxy.Port
		proxy["latency"] = dbPproxy.Latency
		proxy["updatedAt"] = dbPproxy.UpdatedAt

		for key := range proxy {
			if fields != "" && !strings.Contains(fields, key) {
				delete(proxy, key)
			}
		}

		proxies = append(proxies, proxy)
	}

	if format == "csv" {
		var result string

		if len(proxies) > 0 {
			for key := range proxies[0] {
				result += fmt.Sprintf("%v,", key)
			}
			result = result[:len(result)-1]
			result += "\n"
		}

		for _, proxy := range proxies {
			for _, val := range proxy {
				result += fmt.Sprintf("%v,", val)
			}
			result = result[:len(result)-1]
			result += "\n"
		}

		c.String(http.StatusOK, result)
		return
	}

	c.JSON(http.StatusOK, proxies)
}

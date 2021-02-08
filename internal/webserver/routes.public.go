package webserver

import (
	"fmt"
	"github.com/stefanoschrs/proxymeister/pkg/types"
	"log"
	"net/http"
	"reflect"
	"strings"

	"github.com/stefanoschrs/proxymeister/internal/utils"

	"github.com/gin-gonic/gin"
)

// Health check
func headHealth(c *gin.Context) {
	c.Status(http.StatusOK)
}

// Return active proxies
//
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

	var proxies []types.ProxyFE
	for _, dbPproxy := range dbProxies {
		proxy := types.ProxyFE{}
		if fields == "" || strings.Contains(fields, "address") {
			address := fmt.Sprintf("%s:%d", dbPproxy.Ip, dbPproxy.Port)
			proxy.Address = &address
		}
		if fields == "" || strings.Contains(fields, "ip") {
			ip := dbPproxy.Ip
			proxy.Ip = &ip
		}
		if fields == "" || strings.Contains(fields, "port") {
			port := dbPproxy.Port
			proxy.Port = &port
		}
		if fields == "" || strings.Contains(fields, "latency") {
			latency := dbPproxy.Latency
			proxy.Latency = &latency
		}
		if fields == "" || strings.Contains(fields, "updatedAt") {
			updatedAt := dbPproxy.UpdatedAt
			proxy.UpdatedAt = &updatedAt
		}

		proxies = append(proxies, proxy)
	}

	if format == "csv" {
		var result string

		if len(proxies) == 0 {
			c.String(http.StatusOK, result)
			return
		}

		// Create header
		v := reflect.ValueOf(proxies[0])
		for i := 0; i < v.NumField(); i++ {
			result += fmt.Sprintf("%v,", v.Type().Field(i).Name)
		}
		result = result[:len(result)-1]
		result += "\n"

		// Add rows
		for _, proxy := range proxies {
			v = reflect.ValueOf(proxy)
			for i := 0; i < v.NumField(); i++ {
				result += fmt.Sprintf("%v,", v.Field(i).Interface())
			}
			result = result[:len(result)-1]
			result += "\n"
		}

		c.String(http.StatusOK, result)
		return
	}

	c.JSON(http.StatusOK, proxies)
}

package webserver

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func corsMiddleware(c *gin.Context) {
	if c.Request.Method != "OPTIONS" {
		c.Next()
		return
	}

	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
	c.Header("Access-Control-Allow-Headers", "authorization, origin, content-type, accept")
	c.Header("Allow", "HEAD,GET,POST,PUT,PATCH,DELETE,OPTIONS")
	c.Header("Content-Type", "application/json")
	c.AbortWithStatus(http.StatusOK)
}

func loggerMiddleware() gin.HandlerFunc {
	if gin.IsDebugging() {
		return gin.Logger()
	}

	return gin.LoggerWithFormatter(func (param gin.LogFormatterParams) string {
		if param.StatusCode >= 200 && param.StatusCode < 300 || param.StatusCode == 304 {
			return ""
		}

		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	})
}

package utils

import (
	"github.com/stefanoschrs/proxymeister/internal/database"

	"github.com/gin-gonic/gin"
)

func ExtractDB(c *gin.Context) database.DB {
	dbContext, _ := c.Get("db")
	return dbContext.(database.DB)
}

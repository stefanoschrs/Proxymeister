package webserver

import (
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Init() *gin.Engine {
	router := gin.New()
	if os.Getenv("GIN_MODE") != "" {
		gin.SetMode(os.Getenv("GIN_MODE"))
	}

	return router
}

func SetMiddleware(router *gin.Engine) {
	router.Use(gin.Recovery())
	router.Use(loggerMiddleware())
	router.Use(corsMiddleware)
}

func SetRoutes(router *gin.Engine) {
	// CORS
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowCredentials = true
	corsConfig.AllowAllOrigins = true

	// Public routes: /
	apiPublic := router.Group("/")
	apiPublic.Use(cors.New(corsConfig))
	{
		apiPublic.HEAD("/health", headHealth)
		apiPublic.GET("/proxies", getProxies)
	}

	// Private routes: /admin
	apiPrivate := router.Group("/admin")
	// TODO: Read from file e.g .htaccess, config.yml
	adminUser := os.Getenv("ADMIN_USER")
	if adminUser == "" {
		adminUser = "admin"
	}
	apiPrivate.Use(gin.BasicAuth(map[string]string{
		adminUser: os.Getenv("ADMIN_PASS"),
	}))
	apiPrivate.Use(cors.New(corsConfig))
	{
		apiPrivate.POST("/fetch", postAdminFetch)
		apiPrivate.POST("/check", postAdminCheck)
	}
}

package api

import (
	"github.com/anotherhadi/eleakxir/leak"
	"github.com/gin-gonic/gin"
)

type API struct {
	Router    *gin.Engine
	Dataleaks *leak.Dataleaks
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().
			Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func NewAPI(dataleaks *leak.Dataleaks, dev bool) *API {
	router := gin.Default()
	if dev == false {
		gin.SetMode(gin.ReleaseMode)
	}
	router.Use(CORSMiddleware())
	api := &API{
		Router:    router,
		Dataleaks: dataleaks,
	}

	api.SetupRoutes()

	return api
}

func (api *API) Run(addr string) error {
	return api.Router.Run(addr)
}

package middleware

import (
	"github.com/Depado/ginprom"
	"github.com/gin-gonic/gin"
	"log"
)

func GinpromMiddleware(p *ginprom.Prometheus) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("Prometheus: ")
		log.Println(p)
		c.Set("ginprom", p)
		c.Next()
	}
}

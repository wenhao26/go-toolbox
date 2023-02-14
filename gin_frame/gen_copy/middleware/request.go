package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func Request() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()
		c.Set("_t", t)
		c.Next()

		// 执行时间值
		latency := time.Since(t).String()
		log.Println("Latency=", latency)
	}
}

package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
	"time"
)

func RequestIdMiddleware(methodName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestId := uuid.New().String()
		c.Header("X-Request-ID", requestId)

		startTime := time.Now()
		log.Printf("Start %s request with ID: %s at %s", methodName, requestId, startTime.Format(time.RFC3339))

		c.Next()

		endTime := time.Now()
		log.Printf("Finished %s request with ID: %s at %s", methodName, requestId, endTime.Format(time.RFC3339))
	}
}

func SetCors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type")
	}
}

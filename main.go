package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.Default()
	router.GET("/hello", withRequestContext(), withTimeout(), controller)
	router.Run(":8001")
}

// RequestIDKey ... context key
type RequestIDKey struct{}

var requestIDKey = RequestIDKey{}

func withRequestContext() gin.HandlerFunc {
	return func(c *gin.Context) {
		reqID := c.GetHeader("X-Request-ID")
		ctx := context.WithValue(c, requestIDKey, reqID)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func withTimeout() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			ctx    context.Context
			cancel context.CancelFunc
		)
		timeout, err := time.ParseDuration("1s")
		if err == nil {
			ctx, cancel = context.WithTimeout(c.Request.Context(), timeout)
		} else {
			ctx, cancel = context.WithCancel(c.Request.Context())
		}
		defer cancel()

		c.Request = c.Request.WithContext(ctx)

		c.Next()

		c.String(http.StatusOK, "Hello ธนพ\n")
	}
}

func controller(c *gin.Context) {
	ctx := c.Request.Context()
	done := make(chan bool)
	go func() {
		time.Sleep(3 * time.Second)
		fmt.Printf("Value: %s\n", ctx.Value(requestIDKey))
		done <- true
	}()

	select {
	case <-ctx.Done():
		fmt.Println(ctx.Err())
		return
	case <-done:
		return
	}
}

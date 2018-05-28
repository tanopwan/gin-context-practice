package main

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.Default()
	router.GET("/hello", func(c *gin.Context) {
		var (
			ctx    context.Context
			cancel context.CancelFunc
		)
		timeout, err := time.ParseDuration("1s")
		if err == nil {
			ctx, cancel = context.WithTimeout(c, timeout)
		} else {
			ctx, cancel = context.WithCancel(c)
		}
		defer cancel()

		c.Request = c.Request.WithContext(ctx)

		doWork(ctx)

		c.String(http.StatusOK, "Hello ธนพ\n")
	})
	router.Run(":8001")
}

func doWork(ctx context.Context) error {
	done := make(chan bool)
	go func() {
		time.Sleep(5 * time.Second)
		done <- true
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-done:
		return nil
	}
}

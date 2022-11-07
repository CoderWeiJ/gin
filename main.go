package main

import (
	"github.com/CoderWeiJ/web/gin"
	"net/http"
)

func main() {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello ginktutu\n")
	})
	// index out of range for testing Recovery()
	r.GET("/panic", func(c *gin.Context) {
		names := []string{"ginktutu"}
		c.String(http.StatusOK, names[100])
	})

	r.Run(":9999")
}

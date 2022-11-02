package main

import (
	"github.com/CoderWeiJ/web/gin"
	"log"
	"net/http"
)

func main() {
	engine := gin.New()
	engine.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"success": true})
	})
	log.Fatal(http.ListenAndServe(":9999", engine))
}

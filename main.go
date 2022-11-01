package main

import (
	"github.com/CoderWeiJ/web/gin"
	"log"
	"net/http"
)

func main() {
	engine := gin.New()
	log.Fatal(http.ListenAndServe(":9999", engine))
}

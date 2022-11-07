package main

import (
	"github.com/CoderWeiJ/web/gin"
	"net/http"
)

func main() {
	r := gin.New()
	r.GET("/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "<h1>Index Page</h1>")
	})
	v1 := r.Group("/v1")
	{
		v1.GET("/", func(c *gin.Context) {
			c.HTML(http.StatusOK, "<h1>Hello gin</h1>")
		})

		v1.GET("/hello", func(c *gin.Context) {
			// expect /hello?name=ginktutu
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
		})
	}
	v2 := r.Group("/v2")
	{
		v2.GET("/hello/:name", func(c *gin.Context) {
			// expect /hello/ginktutu
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
		})
		v2.POST("/login", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"username": c.PostForm("username"),
				"password": c.PostForm("password"),
			})
		})

	}

	r.Run(":9999")
}

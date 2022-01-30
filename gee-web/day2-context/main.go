package main

import (
	"gee"
	"net/http"
)

func main() {
	r := gee.New()

	r.GET("/hello", func(c *gee.Context) {
		//expect /hello?name=crea
		c.String(http.StatusOK, "hello %s, you are at %s\n", c.Query("name"), c.Path)
	})
	r.POST("/login", func(c *gee.Context) {
		c.JSON(http.StatusOK, gee.H{
			"username": c.PostForm("username"),
			"password": c.PostForm("password"),
		})
	})

	r.Run(":8080")
}

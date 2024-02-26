package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewRouter() {
	r := gin.Default()
	r.LoadHTMLGlob("../api/templates/**/*")

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.GET("/id", func(c *gin.Context) {
		c.HTML(http.StatusOK, "home.html", gin.H{"lol": "lol"})
	})
	r.GET("/order", func(c *gin.Context) {
		c.HTML(http.StatusOK, "order.html", gin.H{"lol": "lol"})
	})
	r.Run(":8080")

}

/*
order_id := c.PostForm("id")
c.JSON(200, gin.H{
            "status":  "posted to login",
            "message": "whoo",
            "form": formContent})
    })

	router.GET("/posts/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "posts/index.tmpl", gin.H{
			"title": "Posts",
		})
	})

*/

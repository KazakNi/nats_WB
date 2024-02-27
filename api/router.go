package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Storage interface {
	Get(string) ([]byte, bool)
	Set(string, []byte)
}

func NewRouter(cache Storage) {

	r := gin.Default()
	r.LoadHTMLGlob("../api/templates/**/*")

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.GET("/id", func(c *gin.Context) {
		c.HTML(http.StatusOK, "home.html", gin.H{})
	})
	r.POST("/order", GetOrder(cache))
	r.Run(":8080")

}

package api

import (
	"encoding/json"
	"fmt"

	"net/http"

	"github.com/gin-gonic/gin"
)

func GetOrder(cache Storage) gin.HandlerFunc {
	fn := func(c *gin.Context) {

		var order Order
		order_id := c.PostForm("order_id")

		fmt.Printf("Incoming id %s\n", order_id)

		bytes, found := cache.Get(order_id)
		json.Unmarshal(bytes, &order)

		if found {
			c.HTML(http.StatusOK, "order.html", gin.H{"order": order})
		} else {
			c.HTML(http.StatusNotFound, "404.html", gin.H{"data": "Order not found!"})
		}

	}

	return gin.HandlerFunc(fn)
}

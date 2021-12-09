package main

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func main() {
	cpt := 1
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "vous avez vu cette page " + strconv.Itoa(cpt) + " fois",
		})
		cpt += 1
	})
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

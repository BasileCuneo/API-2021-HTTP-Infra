package main

import (
	"encoding/base64"
	"math/rand"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func genId() string {
	token := make([]byte, 64)
	rand.Seed(time.Now().UnixNano())
	if _, err := rand.Read(token); err != nil {
		// from the documentation : It always returns len(p) and a nil error.
		return base64.StdEncoding.EncodeToString(token)
	}
	return base64.StdEncoding.EncodeToString(token)
}
func main() {
	r := gin.Default()
	idCpt := make(map[string]int)
	r.GET("/", func(c *gin.Context) {
		userId, err := c.Cookie("id")
		//if  cookie is not set
		if err != nil {
			//gen an id
			strId := string(genId())
			//set the cookie
			c.SetCookie("id", strId, 120, "/", "", false, false)
			idCpt[strId] = 0
			c.JSON(200, gin.H{
				"message": "first time uh ?",
			})
			return
		}
		timeViewed, ok := idCpt[userId]
		//the user's cookie is set
		if ok {
			c.JSON(200, gin.H{
				"message": "vous avez vu cette page " + strconv.Itoa(timeViewed+1) + " fois",
			})
			idCpt[userId]++
		} else {
			c.JSON(200, gin.H{
				"message": "Really ?",
			})
		}
	})
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

// package main

// import (
// 	"encoding/base64"
// 	"math/rand"
// 	"strconv"
// 	"time"

// 	"github.com/gin-gonic/gin"
// )

// func genId() string {
// 	token := make([]byte, 64)
// 	rand.Seed(time.Now().UnixNano())
// 	if _, err := rand.Read(token); err != nil {
// 		// from the documentation : It always returns len(p) and a nil error.
// 		return base64.StdEncoding.EncodeToString(token)
// 	}
// 	return base64.StdEncoding.EncodeToString(token)
// }
// func main() {
// 	r := gin.Default()
// 	idCpt := make(map[string]int)
// 	r.GET("/", func(c *gin.Context) {
// 		userId, err := c.Cookie("id")
// 		//if  cookie is not set
// 		if err != nil {
// 			//gen an id
// 			strId := string(genId())
// 			//set the cookie
// 			c.SetCookie("id", strId, 120, "/", "", false, false)
// 			idCpt[strId] = 0
// 			c.JSON(200, gin.H{
// 				"message": "first time uh ?",
// 			})
// 			return
// 		}
// 		timeViewed, ok := idCpt[userId]
// 		//the user's cookie is set
// 		if ok {
// 			c.JSON(200, gin.H{
// 				"message": "vous avez vu cette page " + strconv.Itoa(timeViewed+1) + " fois",
// 			})
// 			idCpt[userId]++
// 		} else {
// 			c.JSON(200, gin.H{
// 				"message": "Really ?",
// 			})
// 		}
// 	})
// 	r.Run(":80") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
// }
package main

import (
	"math/rand"

	"github.com/gin-gonic/gin"
)

type Person struct {
	Firstname string   `json:"firstname"` //json binding for the struct
	Lastname  string   `json:"lastname"`
	Hobby     []string `json:"hobby"`
	Interest  []string `json:"interest"`
	Skills    []string `json:"skills"`
}

func genPerson() Person {
	firstname := []string{"John", "Jane", "Jack", "Jill", "Joe", "Joris", "Basile", "Bastien", "Adrien", "Kevin", "George", "Alfred", "Zoe", "Yandira", "Alice", "Camille"}
	lastname := []string{"Schaller", "Cueno", "Doe", "Simplon", "Simpson", "Smith", "Smythe", "Wichoud", "Lemaire", "Lemaitre", "Lemaitre", "Muller", "Nguyen"}
	hobby := []string{"Football", "Cinema", "Music", "Reading", "Traveling", "Cooking", "Sleeping", "Swimming", "Dancing", "Skiing", "Running", "Skiing"}
	interest := []string{"Functional language", "NLP", "ML", "DevOps", "movie", "reverse engeenering", "Data analysis", "Cloud infra", "github action"}
	skills := []string{"Go", "Python", "Java", "C++", "C#", "PHP", "Ruby", "Latex", "HTML", "CSS", "SQL", "NoSQL", "Linux", "Windows", "MacOS", "Android", "iOS"}
	//shuflling the arrays
	for i := range skills {
		j := rand.Intn(i + 1)
		skills[i], skills[j] = skills[j], skills[i]
	}
	for i := range hobby {
		j := rand.Intn(i + 1)
		hobby[i], hobby[j] = hobby[j], hobby[i]
	}
	for i := range interest {
		j := rand.Intn(i + 1)
		interest[i], interest[j] = interest[j], interest[i]
	}

	return Person{
		Firstname: firstname[rand.Intn(len(firstname))],
		Lastname:  lastname[rand.Intn(len(lastname))],
		Hobby:     hobby[0:rand.Intn(len(hobby))],
		Interest:  interest[0:rand.Intn(len(interest)/2)],
		Skills:    skills[0:rand.Intn(len(skills)/2)]}

}
func main() {

	//default router for the web app
	r := gin.Default()

	//register a route with a handler (lambda)
	r.GET("/", func(c *gin.Context) {
		users := []Person{}
		n := rand.Intn(10)
		for i := 0; i < n; i++ {
			users = append(users, genPerson())
		}
		c.JSON(200, users)
	})

	r.Run(":80")
}

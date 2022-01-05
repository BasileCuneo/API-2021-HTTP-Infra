package main

import (
	"math/rand"

	"github.com/gin-gonic/gin"
)

const MAX_PERSON = 10

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
		j := rand.Intn(i + 1) // + 1 because i+1 is not inclusive
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
		n := rand.Intn(MAX_PERSON)
		for i := 0; i < n; i++ {
			users = append(users, genPerson())
		}
		c.JSON(200, users)
	})

	r.Run(":80")
}

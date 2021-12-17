package main

//import gin-gonic
import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/gin-gonic/gin"
)

func getContainerData() []types.Container {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}
	return containers
}
func main() {

	r := gin.Default()
	r.LoadHTMLGlob("templates/*.tmpl")
	r.GET("/", func(c *gin.Context) {
		data := make(map[string]interface{})
		containers := getContainerData()
		data["nRoutes"] = 15
		data["nContainers"] = len(containers)
		data["containers"] = containers
		c.HTML(200, "index.tmpl", data)
	})
	r.Run(":9090")
}

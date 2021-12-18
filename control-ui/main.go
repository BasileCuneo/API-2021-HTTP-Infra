package main

//import gin-gonic
import (
	"context"
	"net/http"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/gin-gonic/gin"
)

type Img struct {
	ID         string   `json:"ID"`
	Containers int64    `json:"Containers"`
	RepoTags   []string `json:"rRepoTags"`
}

var cli *client.Client

func getContainerData() ([]types.Container, []Img) {
	ctx := context.Background()
	var err error
	cli, err = client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	info := make(map[string]int64)
	for _, container := range containers {
		info[container.Image]++
	}
	imgList, err := cli.ImageList(ctx, types.ImageListOptions{})
	if err != nil {
		panic(err)
	}
	images := make([]Img, 0)
	for k, v := range info {
		images = append(images, Img{ID: k, Containers: v, RepoTags: []string{}})
	}
	for _, imgL := range imgList {
		for idx, img := range images {
			if img.ID == imgL.ID {
				images[idx].RepoTags = imgL.RepoTags
			}
		}
	}

	return containers, images
}

//var urlPrefix = os.Getenv("PREFIX")
func indexHandler(c *gin.Context) {
	data := make(map[string]interface{})
	containers, images := getContainerData()
	data["nRoutes"] = 15
	data["nContainers"] = len(containers)
	data["containers"] = containers
	data["nImages"] = len(images)
	data["images"] = images
	c.HTML(200, "index.tmpl", data)
}
func stopHandler(c *gin.Context) {
	id := c.Param(":id")
	println("Stopping " + id)
	//cli.ContainerStop(context.Background(), id, nil)
	c.JSON(http.StatusOK, gin.H{})
}
func scaleUpHandler(c *gin.Context) {
	id := c.Param(":id")
	println("Scaling up " + id)
	c.Next()
}
func scaleDownHandler(c *gin.Context) {
	id := c.Param(":id")
	println("Scaling up " + id)
	c.JSON(http.StatusOK, gin.H{})
}
func main() {

	r := gin.Default()
	r.LoadHTMLGlob("templates/*.tmpl")
	r.GET("/", indexHandler)
	r.POST("/panel/stop/:id", stopHandler)
	r.POST("/panel/scaleup/:id", scaleUpHandler)
	r.POST("/panel/scaledown/:id", scaleDownHandler)
	r.Run(":9090")
}

package main

//import gin-gonic
import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/gin-gonic/gin"
)

func remove(s []types.ImageSummary, i int) []types.ImageSummary {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

var cli *client.Client

func getData() ([]types.Container, []types.ImageSummary) {
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
	//because Docker Engine is shitty
	for i, img := range imgList {
		if img.RepoTags[0] == "<none>:<none>" {
			imgList = remove(imgList, i)
		}
	}
	//because Docker Engine is shitty and has bugs and doesn't count container of new images --'
	for _, container := range containers {
		for _, img := range imgList {
			if container.ImageID == img.ID {
				img.Containers++
				//because we found the images container, we can stop the loop
				break
			}
		}
	}

	return containers, imgList
}

//var urlPrefix = os.Getenv("PREFIX")
func indexHandler(c *gin.Context) {
	data := make(map[string]interface{})
	containers, images := getData()
	data["nRoutes"] = 15
	data["nContainers"] = len(containers)
	data["containers"] = containers
	data["nImages"] = len(images)
	data["images"] = images
	data["prefix"] = ""
	c.HTML(200, "index.tmpl", data)
}
func stopHandler(c *gin.Context) {
	id := c.Param("id")
	fmt.Printf("Stopping %q", id)
	err := cli.ContainerStop(context.Background(), id, nil)
	if err != nil {
		c.HTML(http.StatusOK, "apiFailure.tmpl", gin.H{
			"errorMessage": err.Error(),
		})

	} else {
		c.HTML(http.StatusOK, "apiSuccess.tmpl", gin.H{
			"action": "Arret",
			"what":   id,
		})
	}
}
func scaleUpHandler(c *gin.Context) {
	id := c.Param("id")
	println("Scaling up " + id)
	err := errors.New("not implemented")
	if err != nil {
		c.HTML(http.StatusOK, "apiFailure.tmpl", gin.H{
			"errorMessage": err.Error(),
		})

	} else {
		c.HTML(http.StatusOK, "apiSuccess.tmpl", gin.H{
			"action": "Scale up",
			"what":   id,
		})
	}
}
func scaleDownHandler(c *gin.Context) {
	id := c.Param("id")
	println("Scaling down " + id)
	err := errors.New("not implemented scaling up the image ")
	if err != nil {
		c.HTML(http.StatusOK, "apiFailure.tmpl", gin.H{
			"errorMessage": err.Error(),
		})
	} else {
		c.HTML(http.StatusOK, "apiSuccess.tmpl", gin.H{
			"action": "Scale down",
			"what":   id,
		})
	}
}

func main() {

	r := gin.Default()
	r.LoadHTMLGlob("templates/*.tmpl")
	r.GET("/", indexHandler)
	r.GET("/panel/stop/:id", stopHandler)
	r.GET("/panel/scaleup/:id", scaleUpHandler)
	r.GET("/panel/scaledown/:id", scaleDownHandler)
	r.Run(":9090")
}

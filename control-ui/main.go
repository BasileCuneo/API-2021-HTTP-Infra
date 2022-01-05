package main

//import gin-gonic
import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/docker/docker/api/types"

	"github.com/docker/docker/client"
	"github.com/gin-gonic/gin"
)

var urlPrefix = os.Getenv("PREFIX")

//swap and trim for efficiency purpose
func remove(s []types.ImageSummary, i int) []types.ImageSummary {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

var cli *client.Client

func getData() ([]types.Container, []types.ImageSummary) {
	//_ = getStat()
	ctx := context.Background()
	var err error
	cli, err = client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{All: true})
	if err != nil {
		panic(err)
	}
	imgList, err := cli.ImageList(ctx, types.ImageListOptions{})
	if err != nil {
		panic(err)
	}
	var curatedList []types.ImageSummary
	// Because Docker Engine is shitty and show
	// some random <none>:<none> images and i don't know how to filter them except like this
	for i, img := range imgList {
		if img.RepoTags[0] != "<none>:<none>" {
			curatedList = append(curatedList, imgList[i])
			//because Docker Engine is shitty and has bugs and doesn't correctly count new container of images --'
			curatedList[len(curatedList)-1].Containers = 0
		}
	}

	for _, container := range containers {
		for i := range curatedList {
			if container.ImageID == curatedList[i].ID {
				curatedList[i].Containers++
				//because we found the container's image, we can stop the loop
				//(Containers can't have more than one base image)
				break
			}
		}
	}

	return containers, curatedList
}

func indexHandler(c *gin.Context) {
	data := make(map[string]interface{})
	containers, images := getData()
	data["nRoutes"] = 4 //who gonna check anyway ?
	data["nContainers"] = len(containers)
	data["containers"] = containers
	data["nImages"] = len(images)
	data["images"] = images
	data["prefix"] = urlPrefix
	c.HTML(200, "index.tmpl", data)
}
func startHandler(c *gin.Context) {
	id := c.Param("id")
	fmt.Printf("Starting %q", id)
	err := cli.ContainerStart(context.Background(), id, types.ContainerStartOptions{})
	if err != nil {
		c.HTML(http.StatusOK, "apiFailure.tmpl", gin.H{
			"errorMessage": err.Error(),
			"prefix":       urlPrefix,
		})
	} else {
		c.HTML(http.StatusOK, "apiSuccess.tmpl", gin.H{
			"action": "Start",
			"what":   id,
			"prefix": urlPrefix,
		})
	}
}
func stopHandler(c *gin.Context) {
	id := c.Param("id")
	fmt.Printf("Stopping %q", id)
	err := cli.ContainerStop(context.Background(), id, nil)
	if err != nil {
		c.HTML(http.StatusOK, "apiFailure.tmpl", gin.H{
			"errorMessage": err.Error(),
			"prefix":       urlPrefix,
		})

	} else {
		c.HTML(http.StatusOK, "apiSuccess.tmpl", gin.H{
			"action": "Arret",
			"what":   id,
			"prefix": urlPrefix,
		})
	}
}
func removeHandler(c *gin.Context) {
	id := c.Param("id")
	fmt.Printf("Removing %q\n", id)
	err := cli.ContainerRemove(context.Background(), id, types.ContainerRemoveOptions{RemoveVolumes: true, RemoveLinks: false, Force: true})
	if err != nil {
		c.HTML(http.StatusOK, "apiFailure.tmpl", gin.H{
			"errorMessage": err.Error(),
			"prefix":       urlPrefix,
		})
	} else {
		c.HTML(http.StatusOK, "apiSuccess.tmpl", gin.H{
			"action": "Remove",
			"what":   id,
			"prefix": urlPrefix,
		})
	}

}
func increment(name string, imgId string) string {
	_, resp := getData()

	var n int64
	for _, img := range resp {
		if img.ID == imgId {
			n = img.Containers + 1
		}
	}
	index := strings.LastIndex(name, "_")

	return fmt.Sprintf("%s_%d", name[:index], n)
}

func cloneContainer(original types.Container) (string, error) {
	resp, _ := cli.ContainerInspect(context.Background(), original.ID)
	newCtnConfig := resp.Config
	newCtnHostConfig := resp.HostConfig
	newCtnName := increment(original.Names[0], original.ImageID)
	newCtnBody, err := cli.ContainerCreate(context.Background(), newCtnConfig, newCtnHostConfig, nil, nil, newCtnName)
	if err != nil {
		return "", err
	}
	err2 := cli.ContainerStart(context.Background(), newCtnBody.ID, types.ContainerStartOptions{})
	return newCtnBody.ID, err2
}

func scaleUpHandler(c *gin.Context) {
	id := c.Param("id")
	//find in the containers list an image with the ImageID similar to the id
	ctn, _ := getData()
	//the container we want to scale up
	var sample types.Container
	var i int
	for i = range ctn {
		if ctn[i].ImageID == id {
			sample = ctn[i]
		}
	}

	if sample.ImageID == "" {
		c.HTML(http.StatusOK, "apiFailure.tmpl", gin.H{
			"errorMessage": "No such container, to scale up we need a container to copy from. Please run docker-compose up -d first",
			"prefix":       urlPrefix,
		})
		return
	}
	id, err := cloneContainer(sample)

	if err != nil {
		c.HTML(http.StatusOK, "apiFailure.tmpl", gin.H{
			"errorMessage": err.Error(),
			"prefix":       urlPrefix,
		})
	} else {
		c.HTML(http.StatusOK, "apiSuccess.tmpl", gin.H{
			"action": "ScaleUp",
			"what":   id,
			"prefix": urlPrefix,
		})
	}
}

func scaleDownHandler(c *gin.Context) {
	imageId := c.Param("id")
	println("Scaling down " + imageId)
	ctn, _ := getData()
	var sample types.Container
	for i := len(ctn) - 1; i > 0; i-- {
		if ctn[i].ImageID == imageId {
			sample = ctn[i]
		}
	}
	fmt.Printf("\n\n%+v\n\n", sample)
	fmt.Printf("id : %s\n", sample.ID)
	fmt.Printf("Imageid : %s\n", sample.ImageID)
	fmt.Printf("Image : %s\n", sample.Image)

	err := cli.ContainerStop(context.Background(), sample.ID, nil)
	if err != nil {
		c.HTML(http.StatusOK, "apiFailure.tmpl", gin.H{
			"errorMessage": err.Error(),
			"prefix":       urlPrefix,
		})
		return

	}

	cli.ContainerRemove(context.Background(), sample.ID, types.ContainerRemoveOptions{RemoveVolumes: true, RemoveLinks: false, Force: true})
	if err != nil {
		c.HTML(http.StatusOK, "apiFailure.tmpl", gin.H{
			"errorMessage": err.Error(),
			"prefix":       urlPrefix,
		})
		return
	}
	c.HTML(http.StatusOK, "apiSuccess.tmpl", gin.H{
		"action": "Scale down",
		"what":   sample.Names[0],
		"prefix": urlPrefix,
	})
}

func lost(c *gin.Context) {
	c.HTML(http.StatusOK, "lost.tmpl", gin.H{
		"title": "lost ?",
	})
}
func main() {
	fmt.Printf("Starting the server behind url %q\n", os.Getenv("PREFIX"))
	r := gin.Default()
	r.LoadHTMLGlob("templates/*.tmpl")
	r.GET("/", indexHandler)
	r.GET("/panel/start/:id", startHandler)
	r.GET("/panel/stop/:id", stopHandler)
	r.GET("/panel/remove/:id", removeHandler)
	r.GET("/panel/scaleup/:id", scaleUpHandler)
	r.GET("/panel/scaledown/:id", scaleDownHandler)
	r.NoRoute(lost)
	r.Run(":80")
}

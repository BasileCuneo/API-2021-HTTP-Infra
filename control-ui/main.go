package main

//import gin-gonic
import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"

	"github.com/docker/docker/client"
	"github.com/gin-gonic/gin"
)

func getStat() int {

	// resp, err := http.Get("http://localhost:8080/api/overview")
	// if err != nil {
	// 	panic(err)
	// }
	// defer resp.Body.Close()
	// ct := make([]byte, 1024) //big enough for stats
	// resp.Body.Read(ct)
	// fmt.Println("Stats: " + string(ct))
	// var data interface{}
	// json.Unmarshal(ct, &data)

	return 0 //seconde["total"].(int)
}

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
	// Because Docker Engine is shitty and show
	// some random <none>:<none> and has a default value of -1
	// for the number of containers --'
	for i, img := range imgList {
		if img.RepoTags[0] == "<none>:<none>" {
			imgList = remove(imgList, i)
		} else {
			imgList[i].Containers = 0
		}
	}
	//because Docker Engine is shitty and has bugs and doesn't correctly count new container of images --'
	for _, container := range containers {
		for i := range imgList {
			if container.ImageID == imgList[i].ID {
				imgList[i].Containers++
				//because we found the container's image, we can stop the loop
				//(Containers can't have more than one base image)
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
func startHandler(c *gin.Context) {
	id := c.Param("id")
	fmt.Printf("Starting %q", id)
	err := cli.ContainerStart(context.Background(), id, types.ContainerStartOptions{})
	if err != nil {
		c.HTML(http.StatusOK, "apiFailure.tmpl", gin.H{
			"errorMessage": err.Error(),
		})
	} else {
		c.HTML(http.StatusOK, "apiSuccess.tmpl", gin.H{
			"action": "Start",
			"what":   id,
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
		})

	} else {
		c.HTML(http.StatusOK, "apiSuccess.tmpl", gin.H{
			"action": "Arret",
			"what":   id,
		})
	}
}
func removeHandler(c *gin.Context) {
	id := c.Param("id")
	fmt.Printf("Removing %q", id)
	err := cli.ContainerRemove(context.Background(), id, types.ContainerRemoveOptions{RemoveVolumes: true, RemoveLinks: false, Force: true})
	if err != nil {
		c.HTML(http.StatusOK, "apiFailure.tmpl", gin.H{
			"errorMessage": err.Error(),
		})
	} else {
		c.HTML(http.StatusOK, "apiSuccess.tmpl", gin.H{
			"action": "Remove",
			"what":   id,
		})
	}
}
func scaleUpHandler(c *gin.Context) {
	id := c.Param("id")
	//find in the containers list an image with the ImageID similar to the id
	ctn, _ := getData()
	//the container we want to scale up
	var sample types.Container
	for i := range ctn {
		if ctn[i].ImageID == id {
			sample = ctn[i]
		}
	}
	//	ctx context.Context, config *container.Config, hostConfig *container.HostConfig, networkingConfig *network.NetworkingConfig, platform *v1.Platform, containerName string

	_, err := cli.ContainerCreate(context.Background(), &container.Config{}, &container.HostConfig{}, nil, nil, sample.ID)
	if err != nil {
		c.HTML(http.StatusOK, "apiFailure.tmpl", gin.H{
			"errorMessage": err.Error(),
		})
	} else {
		c.HTML(http.StatusOK, "apiSuccess.tmpl", gin.H{
			"action": "ScaleUp",
			"what":   id,
		})
	}

	// println("Scaling up " + id)
	// err := errors.New("not implemented : scaling up the image ")
	// if err != nil {
	// 	c.HTML(http.StatusOK, "apiFailure.tmpl", gin.H{
	// 		"errorMessage": err.Error(),
	// 	})

	// } else {
	// 	c.HTML(http.StatusOK, "apiSuccess.tmpl", gin.H{
	// 		"action": "Scale up",
	// 		"what":   id,
	// 	})
	// }
}
func scaleDownHandler(c *gin.Context) {
	id := c.Param("id")
	println("Scaling down " + id)
	err := errors.New("not implemented scaling down the image ")
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
	r.GET("/panel/start/:id", startHandler)
	r.GET("/panel/stop/:id", stopHandler)
	r.GET("/panel/remove/:id", removeHandler)
	r.GET("/panel/scaleup/:id", scaleUpHandler)
	r.GET("/panel/scaledown/:id", scaleDownHandler)

	r.Run(":9090")
}

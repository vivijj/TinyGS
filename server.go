package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

// the max time the client wait the response for
const WAITPERIOD int = 5
const picDir string = "/home/ljj/pic"
const textDir string = "/home/ljj/text"

func exist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func getPicPath(name string) string {
	return picDir + name + ".jpg"
}

func getTextPath(name string) string {
	return textDir + name + ".txt"
}

func getPic(name string) []byte {
	path := "./pic" + name + ".jpg"
	var file []byte
	if exist(path) {
		file, _ = ioutil.ReadFile(path)
	}
	return file
}
func main() {
	router := gin.Default() // add log file
	f, _ := os.Create("server.log")
	gin.DefaultWriter = io.MultiWriter(f)
	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"1": "1",
		})
	})

	router.POST("/picmaker", func(c *gin.Context) {
		// TODO: decide get post data from json file or post data
		id := c.Query("id")
		message := c.PostForm("message")
		filePath1 := getTextPath(id)
		msg := []byte(message)
		err := ioutil.WriteFile(filePath1, msg, 0644)
		if err != nil {
			log.Fatal(err)
		}
		// c.FileAttachment("/pic/"+id+".jpg", "1.jpg")
		picPath := getPicPath(id)
		for i := 0; i < WAITPERIOD; i++ {
			if exist(picPath) {
				c.File(picPath)
				fmt.Println("remove the " + picPath)
				err := os.Remove(picPath)
				if err != nil {
					fmt.Println("fail remove the pic ")
				}
				return
			}
			time.Sleep(1 * time.Second)
		}
		fmt.Println("can not get the pic")
		// if can't get the pic after 5s, may be something wrong with the model or server
		c.String(http.StatusInternalServerError, "something wrong with the server model")
	})
	router.Run(":8001")
}

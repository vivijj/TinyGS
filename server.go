package main

import (
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
const picDir string = "/home/ljj/pic/"
const textDir string = "/home/ljj/text/"

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

func main() {
	gin.SetMode(gin.ReleaseMode)
	f, _ := os.Create("serverGin.log")
	gin.DefaultWriter = io.MultiWriter(f)

	router := gin.Default()

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
		getPicSuc := false
		for i := 0; i < WAITPERIOD; i++ {
			if exist(picPath) {
				c.File(picPath)
				getPicSuc = true
				break
			}
			time.Sleep(1 * time.Second)
		}
		if getPicSuc {
			err := os.Remove(picPath)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			// if can't get the pic after 5s, may be something wrong with the model or server
			c.String(http.StatusInternalServerError, "something wrong with the server model")
		}
	})
	router.Run(":8081")
}

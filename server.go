package main

import (
	"fmt"
	"github.com/TinyGS/config"
	"github.com/TinyGS/logger"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

// the max time the client wait the response for
const WAITPERIOD int = 5
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

// getKey set the key use in the filename of text&picture
func getKey(id string) (k string){
	t := time.Now().Format("20060102150405")
	k = id + "_" + t
	return
}

func generateText(id, message string) {
	textPath := config.Conf.CacheConfig.TextFolder + getKey(id) + ".txt"
	fmt.Println("textpath is:", textPath)
	if exist(textPath) {
		panic("some thing wrong with the AI model batch")
	}
	err := ioutil.WriteFile(textPath, []byte(message), 0644)
	if err != nil {
		zap.L().Error("can't create the text file", zap.String("textPath", textPath))
	}
}

func generatePic(id string) string{
	picPath := config.Conf.CacheConfig.PicFolder + getKey(id) + ".jpg"
	for i:= 0; i < 4; i++ {
		if exist(picPath) {
			return picPath
		}
		time.Sleep(1 * time.Second)
	}
	return ""
}

func picMaker(c * gin.Context) {
	id := c.Query("id")
	msg := c.PostForm("message")

	generateText(id, msg)
	picPath := generatePic(id)
	fmt.Println("picpath is :", picPath)
	if picPath == "" {
		fmt.Println("pic path is nil")
		panic("pic path is null")
		//c.String(http.StatusInternalServerError, "something wrong with the server model")
	} else {
		fmt.Println("picpath is not nil")
		//c.File(picPath)
	}
}


func main() {
	if len(os.Args) <= 1 {
		return
	}
	if err := config.Init(os.Args[1]); err != nil {
		panic(err)
	}
	// init logger
	if err := logger.InitLogger(config.Conf.LogConfig); err != nil {
		fmt.Printf("init logger failed, err:%v\n", err)
		return
	}
	gin.SetMode(config.Conf.Mode)
	router := gin.New()
	router.Use(logger.GinLogger(),logger.GinRecovery(false))

	router.GET("/", func(C *gin.Context) {
		C.String(http.StatusOK, "Restful-API")
	})
	router.POST("/picmaker", picMaker)

	addr := fmt.Sprintf(":%v", config.Conf.Port)
	_ = router.Run(addr)
}

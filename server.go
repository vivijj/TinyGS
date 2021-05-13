package main

import (
	"fmt"
	"github.com/TinyGS/config"
	"github.com/TinyGS/logger"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// the max time the client wait the response for
const WAITPERIOD int = 5
const CONN = 100
const DELETEDURATION = 1 // hours
// use to delete the pic after serve it
type ExpirePic struct {
	sync.RWMutex
	expirePic []string
}
var expirePicture = &ExpirePic{expirePic: make([]string, CONN)}

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
	for i:= 0; i < WAITPERIOD; i++ {
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
		c.String(http.StatusInternalServerError, "something wrong with the server model")
	} else {
		fmt.Println("picpath is not nil")
		c.File(picPath)
		// indicate that the pic could be safe delete
		expirePicture.expirePic = append(expirePicture.expirePic, picPath)
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

	go func() {
		t:= time.NewTicker(DELETEDURATION *time.Hour)
		for {
			<- t.C
			// add the lock to ensure that no goroutine append the slice when delete picture
			expirePicture.Lock()
			for _, v := range expirePicture.expirePic {
				err := os.Remove(v)
				if err != nil {
					zap.L().Error("fail to remove the file",zap.String("path",v),zap.Any("error",err))
				}
			}
			// clean up the expirePic after delete the picture
			expirePicture.expirePic = expirePicture.expirePic[:0]
			zap.L().Info("finish delete the expire pic")
			expirePicture.Unlock()
		}
	}()

	router.GET("/", func(C *gin.Context) {
		C.String(http.StatusOK, "Restful-API")
	})
	router.POST("/picmaker", picMaker)

	addr := fmt.Sprintf(":%v", config.Conf.Port)
	_ = router.Run(addr)
}

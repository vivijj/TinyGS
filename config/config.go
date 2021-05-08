package config

import (
	"encoding/json"
	"io/ioutil"
)

// Config : Global Config
type Config struct {
	Mode       string `json:"mode"`
	Port       int    `json:"port"`
	*LogConfig `json:"log"`
	*CacheConfig `json:"cache"`
}

// CacheConfig indicate where the tmp text & picture is use as share folder with the AI model
type CacheConfig struct {
	PicFolder string `json:"pic_folder"`
	TextFolder string `json:"text_folder"`
}
// LogConfig : config of logger
type LogConfig struct {
	Level      string `json:"level"`
	Filename   string `json:"filename"`
	MaxSize    int    `json:"max_size"`
	MaxAge     int    `json:"max_age"`
	MaxBackups int    `json:"max_backups"`
}

// Conf : Global Conf variable
var Conf = new(Config)

func Init(filePath string) error {
	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, Conf)
}

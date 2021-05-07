package app

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

var AppConfig Config

type Config struct {
	Telegram Telegram `yaml:"telegram"`
	Exchange Exchange `yaml:"exchange"`
	Mqtt     Mqtt     `yaml:"mqtt"`
	Market   []Market `yaml:"market"`
}

type Exchange struct {
	Ftx Ftx `yaml:"ftx"`
}

type Ftx struct {
	APIKey     string `yaml:"apiKey"`
	APISecret  string `yaml:"apiSecret"`
	SubAccount string `yaml:"subAccount"`
}

type Market struct {
	Symbol string `yaml:"symbol"`
	Market string `yaml:"market"`
}

type Mqtt struct {
	Server   string `yaml:"server"`
	Port     int64  `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	CAFile   string `yaml:"caFile"`
	Topic    string `yaml:"topic"`
	ClientId string `yaml:"clientId"`
}

type Telegram struct {
	Token string `yaml:"token"`
}

func ParseConfig() {
	data, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Fatalln("failed to read config.yaml")
	}

	err = yaml.Unmarshal(data, &AppConfig)
	if err != nil {
		log.Fatalln("failed to parse config.yaml")
	}
}

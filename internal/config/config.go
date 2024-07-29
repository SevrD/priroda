package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

const pathToConfig = "config.yaml"

type Config struct {
	DatabaseDNS   string `yaml:"DATABASE_URL"`
	BotToken      string `yaml:"botToken"`
	BotAdminToken string `yaml:"botAdminToken"`
	Rules         string `yaml:"rules"`
	LoginAdmin    string `yaml:"loginAdmin"`
	ChannelID     int64  `yaml:"channelID"`

	// Services    struct {
	// 	Loms           string `yaml:"loms"`
	// 	Checkout       string `yaml:"checkout"`
	// 	ProductService string `yaml:"productService"`
	// } `yaml:"services"`
}

var AppConfig = Config{}

func Init() error {
	rawYaml, err := os.ReadFile(pathToConfig)
	if err != nil {
		log.Println("read config file:", err)

		return err
	}

	err = yaml.Unmarshal(rawYaml, &AppConfig)
	if err != nil {
		log.Println("parse config file: %w", err)
		return err
	}

	return nil
}
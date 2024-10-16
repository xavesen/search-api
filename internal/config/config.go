package config

import (
	"encoding/base64"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	ListenAddr 				string		`mapstructure:"LISTEN_ADDR"`
	
	ElasticSearchURLsStr	string		`mapstructure:"ELASTIC_SEARCH_URL"`
	ElasticSearchURLs		[]string
	ElasticSearchKey		string		`mapstructure:"ELASTIC_SEARCH_KEY"`

	KafkaAddrsStr			string		`mapstructure:"KAFKA_ADDR"`
	KafkaAddrs 				[]string
	KafkaTopic				string		`mapstructure:"KAFKA_TOPIC"`

	DbAddr					string		`mapstructure:"DB_ADDR"`
	Db						string		`mapstructure:"DB"`
	DbUser					string		`mapstructure:"DB_USER"`
	DbPass					string		`mapstructure:"DB_PASSWORD"`

	LogLevel 				log.Level	`mapstructure:"LOG_LEVEL"`

	JwtAccessTTL			int			`mapstructure:"JWT_ACCESS_TOKEN_TTL"`
	JwtRefreshTTL			int			`mapstructure:"JWT_REFRESH_TOKEN_TTL"`
	JwtSalt        	 		string 		`mapstructure:"JWT_TOKEN_SALT"`
	JwtKeyStr				string		`mapstructure:"JWT_KEY"`
	JwtKey					[]byte
	TokenHeaderName			string		`mapstructure:"TOKEN_HEADER_NAME"`
}

func LoadConfig() (*Config, error) {
	log.Info("Loading config from environment")
	var config Config

	viper.AutomaticEnv()

	log.Info("Parsing environment variables to config struct")
	if err := viper.Unmarshal(&config); err != nil {
		log.Errorf("Error parsing environment variables to config struct: %s", err.Error())
		return nil, err
	}

	config.KafkaAddrs = strings.Split(config.KafkaAddrsStr, ";")
	config.ElasticSearchURLs = strings.Split(config.ElasticSearchURLsStr, ";")
	jwtKey, err := base64.StdEncoding.DecodeString(config.JwtKeyStr)
	if err != nil {
		log.Errorf("Error decoding JWT_KEY from base64: %s", err)
		os.Exit(1)
	}
	config.JwtKey = jwtKey

	log.Infof("Setting log level to %s", config.LogLevel.String())
	log.SetLevel(config.LogLevel)

	log.Info("Successfully loaded config from environment")
	return &config, nil
}

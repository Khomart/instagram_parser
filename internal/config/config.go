package config

import (
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	OpenAIApiKey  string `mapstructure:"OPENAI_API_KEY"`
	AllowedEmails []string
	Host          string
	Port          int
}

func NewConfig() (*Config, error) {
	c := Config{}
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {
		return nil, err
	}
	err = viper.Unmarshal(&c)
	if err != nil {
		return nil, err
	}
	allowedEmailsEnv := viper.GetString("ALLOWED_EMAILS")
	emailList := strings.Split(allowedEmailsEnv, ", ")
	c.AllowedEmails = emailList

	return &c, nil
}

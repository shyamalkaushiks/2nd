package config

import (
	"errors"
	"runtime"

	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type ConfigStruct struct {
	SERVICE_IP   string
	SERVICE_PORT string
	DATABASE     string
	DB_SERVER    string
	DB_PORT      string
	DB_USER      string
	DB_PASSWORD  string
	DB_DATABASE  string
	API_SECRET   string
}

const windowfile string = "F:\\Dean-ai\\users\\config\\user.toml"
const defaultconfig string = "F:\\Dean-ai\\users\\config\\user.toml"

// config file
var showVersion = false
var Config ConfigStruct
var configFile string

func LoadConfig() error {
	if runtime.GOOS == "windows" {
		flag.StringVarP(&configFile, "config", "c", windowfile, "assets configuration file.")
	} else {
		flag.StringVarP(&configFile, "config", "c", defaultconfig, "assets configuration file.")
	}
	flag.BoolVarP(&showVersion, "version", "v", false, "Display version information and exit")
	flag.Parse()
	viper.SetConfigFile(configFile)
	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigParseError); ok {
			// Log.Error().Err(err).Msg("error in loading config file, viper.ConfigParseError")
			return errors.New("error in loading config file, viper.ConfigParseError")
		} else {
			// Log.Error().Err(err).Msg("error in loading config file")
			return errors.New("error in loading config file")
		}
	}
	viper.Unmarshal(&Config)
	return nil
}

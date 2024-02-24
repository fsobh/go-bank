package util

import "github.com/spf13/viper"

// Config Stores all configuration env variables loaded from app.env using viper
// The annotations tell viper what the name of each value is in the .env file (uses map structure)
type Config struct {
	DBDriver      string `mapstructure:"DB_DRIVER"`
	DBSource      string `mapstructure:"DB_SOURCE"`
	ServerAddress string `mapstructure:"SERVER_ADDRESS"`
}

func LoadConfig(path string) (config Config, err error) {

	viper.AddConfigPath(path)  // tell it where the file is
	viper.SetConfigName("app") // tell it the file name
	viper.SetConfigType("env") // tell it the file type (could also be json, xml,...)

	//Override existing Env variable values with what's read from the file
	viper.AutomaticEnv()

	//Read in the config
	err = viper.ReadInConfig()

	if err != nil {
		return
	}

	//unmarshal what's read from the env file into the config struct were returning
	err = viper.Unmarshal(&config)
	return
}

// you can override some at runtime like this : `SERVER_ADDRESS=0.0.0.0:8081 make server`

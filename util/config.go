package util

import (
	"crypto/ed25519"
	"encoding/base64"
	"fmt"
	"github.com/spf13/viper"
	"os"
	"time"
)

// Config Stores all configuration env variables loaded from app.env using viper
// The annotations tell viper what the name of each value is in the .env file (uses map structure)
type Config struct {
	DBDriver            string `mapstructure:"DB_DRIVER"`
	DBSource            string `mapstructure:"DB_SOURCE"`
	ServerAddress       string `mapstructure:"SERVER_ADDRESS"`
	PasetoPublicKeyStr  string `mapstructure:"PASETO_PUBLIC_KEY"`
	PasetoPrivateKeyStr string `mapstructure:"PASETO_PRIVATE_KEY"`
	PasetoPublicKey     ed25519.PublicKey
	PasetoPrivateKey    ed25519.PrivateKey
	PasetoSymmetricKey  string        `mapstructure:"PASETO_SYM_KEY"`
	AccessTokenDuration time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
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

	err = viper.Unmarshal(&config)
	// The keys are currently base64-encoded strings, we need to decode them and convert to ed25519 keys
	publicKey := base64StringToBytes(config.PasetoPublicKeyStr)

	privateKey := base64StringToBytes(config.PasetoPrivateKeyStr)

	// Set the properly typed keys in the config
	config.PasetoPublicKey = publicKey
	config.PasetoPrivateKey = privateKey
	return
}
func base64StringToBytes(keyString string) []byte {
	key, err := base64.StdEncoding.DecodeString(keyString)
	if err != nil {
		_ = fmt.Errorf("cannot decode Keys: %w", err)
		os.Exit(1)
	}
	return key
}

// you can override some at runtime like this : `SERVER_ADDRESS=0.0.0.0:8081 make server`

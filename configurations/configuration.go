package configurations

import (
	"github.com/stripe/stripe-go"
	"os"
	"sync"
)

// Configuration is a struct that describes all of the application configuration values
type Configuration struct {
	Port string
	MongoURI string
	StripeKey string
}

func getEnv(key string, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// buildConfiguration returns back the hydrated configuration file which is based on various environment variables
func buildConfiguration() *Configuration {

	return &Configuration{
		Port: getEnv("PORT", "3000"),
		MongoURI: "mongodb://localhost:27017",
		// See api keys here: https://dashboard.stripe.com/account/apikeys
		StripeKey: "sk_test_zuIs2qpDAX9EtFlu5qxeHPtC00yrXBp5Ni",
	}
}

// Thread-safe singleton with lazy instantiation
var configuration *Configuration
var once sync.Once

func GetConfiguration() *Configuration {
	once.Do(func() {
		configuration = buildConfiguration()
		stripe.Key = configuration.StripeKey
	})

	return configuration
}

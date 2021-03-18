package main

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/MaxPolarfox/goTools/mongoDB"
	"github.com/MaxPolarfox/tasks/pkg/tasks"
	"github.com/MaxPolarfox/tasks/pkg/types"
)

const ServiceName = "tasks"
const EnvironmentVariable = "APP_ENV"

func main() {
	// Load current environment
	env := os.Getenv(EnvironmentVariable)

	// load config options
	options := loadEnvironmentConfig(env)

	// db collections
	tasksCollection := mongoDB.NewMongo(options.DB.Tasks)

	// create the instance of service
	s := tasks.NewService(options, tasksCollection)

	// start the server
	s.Start()
}

// loadEnvironmentConfig will use the environment string and concatenate to a proper config file to use
func loadEnvironmentConfig(env string) types.Options {
	configFile := "config/" + ServiceName + "/" + env + ".json"
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		panic(err)
	}
	return parseConfigFile(configFile)
}

func parseConfigFile(configFile string) types.Options {
	var opts types.Options
	byts, err := ioutil.ReadFile(configFile)

	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(byts, &opts)
	if err != nil {
		panic(err)
	}

	return opts
}

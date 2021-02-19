package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"

	"github.com/MaxPolarfox/goTools/mongoDB"
	"github.com/MaxPolarfox/tasks/pkg/grpc/service"
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

	netListener := createNetListener(options.Port)

	// db collections
	tasksCollection := mongoDB.NewMongo(options.DB.Tasks)

	// create the instance of service implementation
	taskServiceImpl := tasks.NewService(options, tasksCollection)

	// create gRPC server
	gRPCServer := grpc.NewServer()

	// register service implementation to gRPC server
	service.RegisterTaskServiceServer(gRPCServer, taskServiceImpl)

	// start the server
	if err := gRPCServer.Serve(netListener); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}

// createNetListener creates listener
func createNetListener (port int) net.Listener {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
		panic(fmt.Sprintf("failed to listen: %v", err))
	}

	return lis
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


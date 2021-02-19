package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"github.com/MaxPolarfox/tasks/pkg/types"
	"os"

	"google.golang.org/grpc"

	"github.com/MaxPolarfox/tasks/pkg/grpc/messages"
	"github.com/MaxPolarfox/tasks/pkg/grpc/service"
)

const ServiceName = "tasks"
const EnvironmentVariable = "APP_ENV"

type Client interface {
	AddTask(ctx context.Context, data string) (*types.CreatedTaskResponse, error)
}

type TasksClientImpl struct {
	client service.TaskServiceClient
}

func NewTasksClient() Client {
	var conn *grpc.ClientConn

	serverAddress := fmt.Sprintf("localhost:%d", 3005)

	conn, err := grpc.Dial(serverAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("conn failed: %v", err)
		panic(fmt.Sprintf("conn failed: %v", err))
	}

	client := service.NewTaskServiceClient(conn)

	return &TasksClientImpl{
		client,
	}
}

func (i *TasksClientImpl) AddTask(ctx context.Context, data string) (*types.CreatedTaskResponse, error) {

	metricName := "TasksClientImpl.AddTask"

	taskMSG := messages.Task{
		Data: data,
	}

	responseMSG, err := i.client.Add(ctx, &taskMSG)
	if err != nil {
		log.Printf(metricName, "err", err)
		return nil, err
	}

	response := &types.CreatedTaskResponse{
		ID: responseMSG.AddedTask.Id,
	}

	return response, nil
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
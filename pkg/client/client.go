package client

import (
	"context"
	"fmt"
	"log"
	"github.com/MaxPolarfox/tasks/pkg/types"

	"google.golang.org/grpc"

	"github.com/MaxPolarfox/tasks/pkg/grpc/messages"
	"github.com/MaxPolarfox/tasks/pkg/grpc/service"
)

type Client interface {
	AddTask(ctx context.Context, data string) (*types.CreatedTaskResponse, error)
}

type TasksClientImpl struct {
	client service.TaskServiceClient
}

func NewTasksClient(options types.Options) Client {

	var conn *grpc.ClientConn

	serverAddress := fmt.Sprintf("localhost:%d", options.Port)

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
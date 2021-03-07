package grpc_client

import (
	"context"
	"errors"
	"fmt"
	"github.com/MaxPolarfox/tasks/pkg/types"
	"google.golang.org/grpc"
	"log"

	"github.com/MaxPolarfox/tasks/pkg/grpc/messages"
	"github.com/MaxPolarfox/tasks/pkg/grpc/service"
)

type Client interface {
	AddTask(ctx context.Context, data string) (*types.CreatedTaskResponse, error)
	GetAllTasks(ctx context.Context) (*[]types.Task, error)
	DeleteTask(ctx context.Context, taskID string) error
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
		log.Printf(metricName+".Add", "err", err)
		return nil, err
	}

	response := &types.CreatedTaskResponse{
		ID: responseMSG.AddedTask.Id,
	}

	return response, nil
}

func (i *TasksClientImpl) GetAllTasks(ctx context.Context) (*[]types.Task, error) {
	metricName := "TasksClientImpl.AddTask"

	msg := messages.Action{Data: "delete"}

	responseMSG, err := i.client.GetAll(ctx, &msg)
	if err != nil {
		log.Printf(metricName, "err", err)
		return nil, err
	}

	tasks := make([]types.Task, len(responseMSG.Tasks))

	for i, taskMsg := range responseMSG.Tasks {
		tasks[i] = types.Task{ID: taskMsg.Id, Data: taskMsg.Data}
	}

	return &tasks, nil
}

func (i *TasksClientImpl) DeleteTask(ctx context.Context, taskID string) error {

	metricName := "TasksClientImpl.DeleteTask"

	msg := messages.DeleteRequest{Id: taskID}

	responseMSG, err := i.client.Delete(ctx, &msg)
	if err != nil {
		log.Printf(metricName+".Delete", "err", err)
		return err
	}

	if responseMSG.Error != nil {
		return errors.New(responseMSG.Error.Message)
	}
	return nil
}

package client

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"

	goToolsClient "github.com/MaxPolarfox/goTools/client"
	"github.com/MaxPolarfox/tasks/pkg/taskspb"
	"github.com/MaxPolarfox/tasks/pkg/types"
)

type Client interface {
	CreateTask(ctx context.Context, data string) (*types.CreatedTaskResponse, error)
	GetTasks(ctx context.Context) (*[]types.Task, error)
	DeleteTask(ctx context.Context, taskID string) error
}

type TasksClientImpl struct {
	client taskspb.TaskServiceClient
}

func NewTasksClient(options goToolsClient.Options) Client {
	var conn *grpc.ClientConn

	serverAddress := options.URL

	conn, err := grpc.Dial(serverAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("conn failed: %v", err)
		panic(fmt.Sprintf("conn failed: %v", err))
	}

	client := taskspb.NewTaskServiceClient(conn)

	return &TasksClientImpl{
		client,
	}
}

func (i *TasksClientImpl) CreateTask(ctx context.Context, data string) (*types.CreatedTaskResponse, error) {
	metricName := "Client.CreateTask"

	res := &taskspb.CreateTaskReq{
		Data: data,
	}

	createRes, err := i.client.CreateTask(ctx, res)
	if err != nil {
		log.Printf(metricName+".client.CreateTask", "err", err)
		return nil, err
	}

	response := &types.CreatedTaskResponse{
		ID: createRes.GetId(),
	}

	return response, nil
}

func (i *TasksClientImpl) GetTasks(ctx context.Context) (*[]types.Task, error) {
	metricName := "Client.GetTasks"

	getTasksRes, err := i.client.GetTasks(ctx, &taskspb.GetTasksReq{})
	if err != nil {
		log.Printf(metricName+".client.GetTasks", "err", err)
		return nil, err
	}

	tasks := make([]types.Task, len(getTasksRes.Tasks))

	for i, taskMsg := range getTasksRes.GetTasks() {
		tasks[i] = types.Task{ID: taskMsg.Id, Data: taskMsg.Data}
	}

	return &tasks, nil
}

func (i *TasksClientImpl) DeleteTask(ctx context.Context, taskID string) error {
	metricName := "Client.DeleteTask"

	req := &taskspb.DeleteTaskReq{Id: taskID}

	_, err := i.client.DeleteTask(ctx, req)
	if err != nil {
		log.Printf(metricName+".client.DeleteTask", "err", err)
		return err
	}

	return nil
}

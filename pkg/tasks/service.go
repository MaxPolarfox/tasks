package tasks

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"log"
	"net"
	"os"
	"os/signal"

	uuid "github.com/satori/go.uuid"

	"github.com/MaxPolarfox/goTools/mongoDB"
	"github.com/MaxPolarfox/tasks/pkg/taskspb"
	"github.com/MaxPolarfox/tasks/pkg/types"
)

// Service is a implementation of TasksService Grpc Service
type Service struct {
	options types.Options
	db      DB
}

type DB struct {
	tasks mongoDB.Mongo
}

//NewService returns the pointer to the Service.
func NewService(options types.Options, tasksCollection mongoDB.Mongo) *Service {
	return &Service{
		options: options,
		db:      DB{tasksCollection},
	}
}

func (s *Service) Start() {
	// listen to the appropriate signals, and notify a channel
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.options.Port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	opts := []grpc.ServerOption{}
	server := grpc.NewServer(opts...)
	taskspb.RegisterTaskServiceServer(server, s)
	reflection.Register(server)
	go func() {
		if err := server.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()
	<-stopChan // wait for a signal to exit
	log.Println("shutting down the server")
	server.Stop()
	log.Println("Stopping listener")
	lis.Close()
	log.Println("End of program")
}

// CreateTask creates new task
func (s *Service) CreateTask(ctx context.Context, req *taskspb.CreateTaskReq) (*taskspb.CreateTaskRes, error) {
	taskID := uuid.NewV4().String()
	data := req.GetData()

	newTask := types.Task{
		ID:   taskID,
		Data: data,
	}

	// insert newTask to DB
	_, err := s.db.tasks.InsertOne(ctx, newTask)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("unexpected error: %v", err),
		)
	}

	res := &taskspb.CreateTaskRes{
		Id: taskID,
	}

	return res, nil
}

// GetTasks returns all the tasks
func (s *Service) GetTasks(ctx context.Context, req *taskspb.GetTasksReq) (*taskspb.GetTasksRes, error) {
	res := taskspb.GetTasksRes{Tasks: []*taskspb.Task{}}

	filter := bson.M{}
	cursor, err := s.db.tasks.Find(ctx, filter)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("unexpected error: %v", err),
		)
	}

	err = cursor.All(ctx, &res.Tasks)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("unexpected error: %v", err),
		)
	}

	return &res, nil
}

// DeleteTask function implementation of gRPC Service.
func (s *Service) DeleteTask(ctx context.Context, req *taskspb.DeleteTaskReq) (*taskspb.DeleteTaskRes, error) {
	taskID := req.GetId()

	filter := bson.M{"id": taskID}
	deleteRes, err := s.db.tasks.DeleteOne(ctx, filter)
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("undexpected error: %v", err))
	}

	if deleteRes.DeletedCount == 0 {
		return nil, status.Errorf(codes.NotFound, "")
	}

	return &taskspb.DeleteTaskRes{}, nil
}

package tasks

import (
	"context"
	"log"

	uuid "github.com/satori/go.uuid"

	"github.com/MaxPolarfox/goTools/mongoDB"
	"github.com/MaxPolarfox/tasks/pkg/grpc/messages"
	"github.com/MaxPolarfox/tasks/pkg/grpc/service"
	"github.com/MaxPolarfox/tasks/pkg/types"
)

//TasksServiceGrpcImpl is a implementation of TasksService Grpc Service.
type Service struct {
	options types.Options
	db DB
}

type DB struct {
	tasks mongoDB.Mongo
}

//NewService returns the pointer to the Service.
func NewService(options types.Options, tasksCollection mongoDB.Mongo) *Service {
	return &Service{
		options: options,
		db: DB{tasksCollection},
	}
}

// Add function implementation of gRPC Service.
func (s *Service) Add(ctx context.Context, newTaskMSG *messages.Task) (*service.AddRepositoryResponse, error) {
	log.Println("Received request for adding task")

	metricName := "Service.Add"

	taskID := uuid.NewV4().String()

	// add id to the taskMSG
	newTaskMSG.Id = taskID

	// create new task in DB
	newTask := types.Task{
		taskID,
		newTaskMSG.Data,
	}

	// insert newTask to DB
	_, err := s.db.tasks.InsertOne(ctx, newTask)
	if err != nil {
		log.Println(metricName, "err", err)
		return &service.AddRepositoryResponse{
			AddedTask: newTaskMSG,
			Error: &service.Error{Message: err.Error()},
		}, err
	}

	return &service.AddRepositoryResponse{ AddedTask: newTaskMSG, Error: nil}, nil
}
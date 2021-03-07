package tasks

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
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

// Add function implementation of gRPC Service.
func (s *Service) Add(ctx context.Context, newTaskMSG *messages.Task) (*service.AddRepositoryResponse, error) {
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
			Error:     &service.Error{Message: err.Error()},
		}, err
	}

	return &service.AddRepositoryResponse{AddedTask: newTaskMSG, Error: nil}, nil
}

// GetAll function implementation of gRPC Service.
func (s *Service) GetAll(ctx context.Context, msg *messages.Action) (*service.GetAllTasksResponse, error) {
	metricName := "Service.GetAll"

	response := service.GetAllTasksResponse{Error: nil, Tasks: []*messages.Task{}}

	filter := bson.M{}
	cursor, err := s.db.tasks.Find(ctx, filter)
	if err != nil {
		log.Println(metricName+".Find", "err", "err")
		response.Error = &service.Error{Message: err.Error()}
		return &response, err
	}

	err = cursor.All(ctx, &response.Tasks)
	if err != nil {
		log.Println(metricName+".All", "err", "err")
		response.Error = &service.Error{Message: err.Error()}
		return &response, err
	}

	return &response, nil
}

// Delete function implementation of gRPC Service.
func (s *Service) Delete(ctx context.Context, msg *messages.DeleteRequest) (*service.DeleteResponse, error) {
	metricName := "Service.Delete"

	response := service.DeleteResponse{Error: nil}

	filter := bson.M{"id": msg.Id}
	_, err := s.db.tasks.DeleteOne(ctx, filter)
	if err != nil {
		log.Println(metricName+".DeleteOne", "err", "err")
		response.Error = &service.Error{Message: err.Error()}
		return &response, err
	}

	response.Task = &messages.Task{Id: msg.Id}

	return &response, nil
}

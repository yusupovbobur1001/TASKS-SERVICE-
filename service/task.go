package service

import (
	"context"
	pb "task_service/genproto/task_service"
	"task_service/storage"

	"google.golang.org/protobuf/types/known/emptypb"
)

type TaskService struct {
	storage storage.ITaskStorage
	pb.UnimplementedTaskServiceServer
}

func NewTaskService(storage storage.ITaskStorage) *TaskService {
	return &TaskService{
		storage: storage,
	}
}
func (service *TaskService) CreateTask(ctx context.Context, req *pb.TaskRequest) (*pb.Task, error) {
	return service.storage.CreateTask(ctx, req)
}
func (service *TaskService) UpdateTask(ctx context.Context, req *pb.Task) (*pb.Task, error) {
	return service.storage.UpdateTask(ctx, req)
}
func (service *TaskService) GetTask(ctx context.Context, req *pb.PrimaryKey) (*pb.Task, error) {
	return service.storage.GetTask(ctx, req)
}
func (service *TaskService) GetAllTasks(ctx context.Context, req *pb.GetListRequest) (*pb.TasksResponse, error) {
	return service.storage.GetAllTasks(ctx, req)
}
func (service *TaskService) DeleteTask(ctx context.Context, req *pb.PrimaryKey) (*emptypb.Empty, error) {
	return service.storage.DeleteTask(ctx, req)
}

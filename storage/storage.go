package storage

import (
	"context"
	pb "task_service/genproto/task_service"

	"google.golang.org/protobuf/types/known/emptypb"
)

type IStorage interface {
	Details() IDetailStorage
	Tasks() ITaskStorage
}
type IStorageRepo struct {
	detail IDetailStorage
	task   ITaskStorage
}

func (I IStorageRepo) Details() IDetailStorage {
	return I.detail
}

func (I IStorageRepo) Tasks() ITaskStorage {
	return I.task
}

func NewIStorageRepo(detail IDetailStorage, task ITaskStorage) IStorage {
	return &IStorageRepo{
		detail: detail,
		task:   task,
	}
}

type IDetailStorage interface {
	CreateDetail(ctx context.Context, request *pb.DetailRequest) (*pb.Detail, error)
	UpdateDetail(ctx context.Context, request *pb.Detail) (*pb.Detail, error)
	GetDetail(ctx context.Context, request *pb.PrimaryKey) (*pb.Detail, error)
	GetAllDetails(ctx context.Context, request *pb.GetListRequest) (*pb.DetailResponse, error)
	DeleteDetail(ctx context.Context, request *pb.PrimaryKey) (*emptypb.Empty, error)
}
type ITaskStorage interface {
	CreateTask(ctx context.Context, request *pb.TaskRequest) (*pb.Task, error)
	UpdateTask(ctx context.Context, request *pb.Task) (*pb.Task, error)
	GetTask(ctx context.Context, request *pb.PrimaryKey) (*pb.Task, error)
	DeleteTask(ctx context.Context, request *pb.PrimaryKey) (*emptypb.Empty, error)
	GetAllTasks(ctx context.Context, request *pb.GetListRequest) (*pb.TasksResponse, error)
}

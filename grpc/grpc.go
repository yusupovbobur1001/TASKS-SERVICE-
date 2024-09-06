package grpc

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	pb "task_service/genproto/task_service"
	"task_service/service"
	"task_service/storage"
)

func SetUpServer(storage storage.IStorage) *grpc.Server {
	grpcServer := grpc.NewServer()

	pb.RegisterTaskServiceServer(grpcServer, service.NewTaskService(storage.Tasks()))
	pb.RegisterDetailServiceServer(grpcServer, service.NewDetailService(storage.Details()))

	reflection.Register(grpcServer)
	return grpcServer
}

package service

import (
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
	pb "task_service/genproto/task_service"
	"task_service/storage"
)

type DetailService struct {
	storage storage.IDetailStorage
	pb.UnimplementedDetailServiceServer
}

func NewDetailService(storage storage.IDetailStorage) *DetailService {
	return &DetailService{storage: storage}
}
func (service *DetailService) CreateDetail(ctx context.Context, req *pb.DetailRequest) (*pb.Detail, error) {
	return service.storage.CreateDetail(ctx, req)
}
func (service *DetailService) UpdateDetail(ctx context.Context, req *pb.Detail) (*pb.Detail, error) {
	return service.storage.UpdateDetail(ctx, req)
}
func (service *DetailService) DeleteDetail(ctx context.Context, req *pb.PrimaryKey) (*emptypb.Empty, error) {
	return service.storage.DeleteDetail(ctx, req)
}
func (service *DetailService) GetDetail(ctx context.Context, req *pb.PrimaryKey) (*pb.Detail, error) {
	return service.storage.GetDetail(ctx, req)
}
func (service *DetailService) GetAllDetail(ctx context.Context, req *pb.GetListRequest) (*pb.DetailResponse, error) {
	return service.storage.GetAllDetails(ctx, req)
}

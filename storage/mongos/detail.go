package mongo

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/protobuf/types/known/emptypb"
	"task_service/configs"
	pb "task_service/genproto/task_service"
	"task_service/pkg/logger"
	"task_service/storage"
	"time"
)

type DetailMongo struct {
	Coll *mongo.Collection
	log  logger.ILogger
}

func NewDetailMongo(collection *mongo.Collection, log logger.ILogger) storage.IDetailStorage {
	return &DetailMongo{
		Coll: collection,
		log:  log,
	}
}

type TaskDetail struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	TaskID      string             `bson:"task_id"`
	Description string             `bson:"description"`
	Status      string             `bson:"status"`
	Priority    string             `bson:"priority"`
	CreatedAt   time.Time          `bson:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at"`
}

func (mongodb *DetailMongo) CreateDetail(ctx context.Context, request *pb.DetailRequest) (*pb.Detail, error) {
	var (
		err      error
		timeNow  = time.Now()
		response TaskDetail
	)

	detail := TaskDetail{
		TaskID:      request.GetTaskId(),
		Description: request.GetDescription(),
		Status:      request.GetStatus(),
		Priority:    request.GetPriority(),
		CreatedAt:   timeNow,
	}

	detailId, err := mongodb.Coll.InsertOne(ctx, detail)
	if err != nil {
		mongodb.log.Error("Error inserting detail ", logger.Error(err))
		return nil, err
	}

	err = mongodb.Coll.FindOne(ctx, bson.M{"_id": detailId.InsertedID}).Decode(&response)
	if err != nil {
		mongodb.log.Error("Error  inserted detail ", logger.Error(err))
		return nil, err
	}

	return &pb.Detail{
		Id:          response.ID.Hex(),
		TaskId:      response.TaskID,
		Description: response.Description,
		Status:      response.Status,
		Priority:    response.Priority,
		CreatedAt:   response.CreatedAt.Format(configs.Layout),
		UpdatedAt:   response.UpdatedAt.Format(configs.Layout),
	}, nil
}
func (mongodb *DetailMongo) UpdateDetail(ctx context.Context, request *pb.Detail) (*pb.Detail, error) {
	var (
		err      error
		timeNow  = time.Now()
		response TaskDetail
	)

	objectID, err := primitive.ObjectIDFromHex(request.GetId())
	if err != nil {
		mongodb.log.Error("Invalid ObjectID ", logger.Error(err))
		return nil, err
	}

	filter := bson.M{
		"$set": bson.M{
			"task_id":     request.GetTaskId(),
			"description": request.GetDescription(),
			"status":      request.GetStatus(),
			"priority":    request.GetPriority(),
			"updated_at":  timeNow,
		},
	}

	_, err = mongodb.Coll.UpdateOne(ctx, bson.M{"_id": objectID}, filter)
	if err != nil {
		mongodb.log.Error("Error updating detail ", logger.Error(err))
		return nil, err
	}

	err = mongodb.Coll.FindOne(ctx, bson.M{"_id": objectID}).Decode(&response)
	if err != nil {
		mongodb.log.Error("Error  updated detail ", logger.Error(err))
		return nil, err
	}

	return &pb.Detail{
		Id:          response.ID.Hex(),
		TaskId:      response.TaskID,
		Description: response.Description,
		Status:      response.Status,
		Priority:    response.Priority,
		CreatedAt:   response.CreatedAt.Format(configs.Layout),
		UpdatedAt:   response.UpdatedAt.Format(configs.Layout),
	}, nil
}
func (mongodb *DetailMongo) GetDetail(ctx context.Context, request *pb.PrimaryKey) (*pb.Detail, error) {
	var (
		err      error
		response TaskDetail
	)

	objectID, err := primitive.ObjectIDFromHex(request.GetId())
	if err != nil {
		mongodb.log.Error("Invalid ObjectID ", logger.Error(err))
		return nil, err
	}

	err = mongodb.Coll.FindOne(ctx, bson.M{"_id": objectID}).Decode(&response)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			mongodb.log.Error("Detail not found", logger.Error(err))
			return nil, fmt.Errorf("detail not found")
		}
		mongodb.log.Error("Error details not found", logger.Error(err))
		return nil, err
	}

	return &pb.Detail{
		Id:          response.ID.Hex(),
		TaskId:      response.TaskID,
		Description: response.Description,
		Status:      response.Status,
		Priority:    response.Priority,
		CreatedAt:   response.CreatedAt.Format(configs.Layout),
		UpdatedAt:   response.UpdatedAt.Format(configs.Layout),
	}, nil
}
func (mongodb *DetailMongo) GetAllDetails(ctx context.Context, request *pb.GetListRequest) (*pb.DetailResponse, error) {
	var (
		err     error
		details []*pb.Detail
		offset  int64
		limit   int64
		count   int64
		cursor  *mongo.Cursor
		filter  bson.M
		sort    bson.D
	)

	limit = request.GetLimit()
	offset = (request.GetPage() - 1) * request.GetLimit()

	filter = bson.M{}
	if request.GetSearch() != "" {
		filter["description"] = bson.M{"$regex": request.GetSearch(), "$options": "i"}
	}

	sort = bson.D{{Key: "created_at", Value: -1}}

	cursor, err = mongodb.Coll.Find(ctx, filter, &options.FindOptions{
		Sort:  sort,
		Limit: &limit,
		Skip:  &offset,
	})
	if err != nil {
		mongodb.log.Error(" query that details from ", logger.Error(err))
		return &pb.DetailResponse{}, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var detail TaskDetail
		if err := cursor.Decode(&detail); err != nil {
			mongodb.log.Error("Error that  decode  detail ", logger.Error(err))
			return nil, err
		}
		details = append(details, &pb.Detail{
			Id:          detail.ID.Hex(),
			TaskId:      detail.TaskID,
			Description: detail.Description,
			Status:      detail.Status,
			Priority:    detail.Priority,
			CreatedAt:   detail.CreatedAt.Format(configs.Layout),
			UpdatedAt:   detail.UpdatedAt.Format(configs.Layout),
		})
	}

	count, err = mongodb.Coll.CountDocuments(ctx, filter)
	if err != nil {
		mongodb.log.Error("Error that  counting details ", logger.Error(err))
		return &pb.DetailResponse{}, err
	}

	return &pb.DetailResponse{
		Details: details,
		Count:   count,
	}, nil
}
func (mongodb *DetailMongo) DeleteDetail(ctx context.Context, request *pb.PrimaryKey) (*emptypb.Empty, error) {
	objectID, err := primitive.ObjectIDFromHex(request.GetId())
	if err != nil {
		mongodb.log.Error("Error  ObjectID", logger.Error(err))
		return &emptypb.Empty{}, err
	}

	_, err = mongodb.Coll.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		mongodb.log.Error("Error deleting detail ", logger.Error(err))
		return &emptypb.Empty{}, err
	}

	return &emptypb.Empty{}, nil
}

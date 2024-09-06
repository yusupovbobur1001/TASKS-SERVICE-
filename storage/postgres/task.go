package postgres

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/protobuf/types/known/emptypb"
	"task_service/configs"
	pb "task_service/genproto/task_service"
	"task_service/pkg/helper"
	"task_service/pkg/logger"
	"task_service/storage"
	"time"
)

type TaskRepository struct {
	db  *pgxpool.Pool
	log logger.ILogger
}

func NewTaskRepository(db *pgxpool.Pool, log logger.ILogger) storage.ITaskStorage {
	return &TaskRepository{
		db:  db,
		log: log,
	}
}

func (repo TaskRepository) CreateTask(ctx context.Context, request *pb.TaskRequest) (*pb.Task, error) {
	var (
		query     string
		err       error
		response  pb.Task
		id        = uuid.New()
		timeNow   = time.Now()
		createdAt sql.NullTime
		updatedAt sql.NullTime
	)

	query = `insert into tasks(
                  id,
                  user_id,
                  title,
                  created_at
                  )values ($1,$2,$3,$4) returning 
                   id,
                   user_id,
                   title,
                   created_at,
                   updated_at`
	err = repo.db.QueryRow(ctx, query,
		id,
		request.GetUserId(),
		request.GetTitle(),
		timeNow,
	).Scan(
		&response.Id,
		&response.UserId,
		&response.Title,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		repo.log.Error("this error is create task scan ERROR-~~~~~~~~~~`", logger.Error(err))
		return nil, err
	}
	if createdAt.Valid {
		response.CreatedAt = createdAt.Time.Format(configs.Layout)
	}
	if updatedAt.Valid {
		response.UpdatedAt = updatedAt.Time.Format(configs.Layout)
	}
	return &response, nil
}
func (repo TaskRepository) UpdateTask(ctx context.Context, request *pb.Task) (*pb.Task, error) {
	var (
		err       error
		query     string
		response  pb.Task
		timeNow   = time.Now()
		createdAt sql.NullTime
		updatedAt sql.NullTime
	)
	query = `update tasks  set 
                     user_id=$1,
                     title=$2,
                     updated_at=$3
                     where id=$4 returning
                      id,
                   user_id,
                   title,
                   created_at,
                   updated_at`

	err = repo.db.QueryRow(ctx, query,
		request.GetUserId(),
		request.GetTitle(),
		timeNow,
		request.GetId()).Scan(
		&response.Id,
		&response.UserId,
		&response.Title,
		&createdAt,
		&updatedAt)
	if err != nil {
		repo.log.Error("this error is can be used on scan that ERROR -~~~~~~~~~~~`", logger.Error(err))
		return nil, err
	}
	if createdAt.Valid {
		response.CreatedAt = createdAt.Time.Format(configs.Layout)
	}
	if updatedAt.Valid {
		response.UpdatedAt = updatedAt.Time.Format(configs.Layout)
	}
	return &response, nil
}

func (repo *TaskRepository) GetTask(ctx context.Context, request *pb.PrimaryKey) (*pb.Task, error) {
	var (
		query     string
		err       error
		response  pb.Task
		createdAt sql.NullTime
		updatedAt sql.NullTime
	)
	query = `select  
                   id,
                   user_id,
                   title,
                   created_at,
                   updated_at from tasks where id=$1`

	err = repo.db.QueryRow(ctx, query,
		request.GetId()).Scan(
		&response.Id,
		&response.UserId,
		&response.Title,
		&createdAt,
		&updatedAt)
	if err != nil {
		repo.log.Error("this error that can be happened ERROR-~~~~~~~~~~~", logger.Error(err))
		return nil, err
	}

	if createdAt.Valid {
		response.CreatedAt = createdAt.Time.Format(configs.Layout)
	}
	if updatedAt.Valid {
		response.UpdatedAt = updatedAt.Time.Format(configs.Layout)
	}
	return &response, nil
}

func (repo *TaskRepository) DeleteTask(ctx context.Context, request *pb.PrimaryKey) (*emptypb.Empty, error) {
	var (
		err      error
		query    string
		response emptypb.Empty
	)
	query = `delete from tasks where id=$1`
	result, err := repo.db.Exec(ctx, query, request.GetId())
	if err != nil {
		repo.log.Error("this error can be ERROR-~~~~~~~~~~~`", logger.Error(err))
		return nil, err
	}
	if result.RowsAffected() == 0 {
		repo.log.Error("this error is absolutely that ERROR-~~~~~~~~~ this have been delete before time ")
		return nil, err
	}
	return &response, nil

}
func (repo TaskRepository) GetAllTasks(ctx context.Context, request *pb.GetListRequest) (*pb.TasksResponse, error) {
	var (
		err        error
		tasks      []*pb.Task
		offset     = (request.GetLimit() - 1) * request.GetPage()
		query      string
		count      = int64(0)
		countQuery string
		where      string
		createdAt  sql.NullTime
		updatedAt  sql.NullTime
	)
	countQuery = "select count(*) from tasks"
	if request.GetSearch() != "" {
		where, err := helper.MakeWherePartOfQueryWithSearchFieldOfRequest(request.GetSearch())
		if err != nil {
			repo.log.Error("error while taking values from search field of request in storage layer", logger.Error(err))
			return &pb.TasksResponse{}, err
		}
		countQuery += where
	}

	if err = repo.db.QueryRow(ctx, countQuery).Scan(&count); err != nil {
		repo.log.Error("error while selecting tasks count in storage layer", logger.Error(err))
		return &pb.TasksResponse{}, err
	}
	query = ` select 
           id,
           user_id,
           title,
           created_at,
           updated_at from tasks  `
	query += where
	query += ` order by created_at DESC limit $1 offset $2`

	rows, err := repo.db.Query(ctx, query, request.GetLimit(), offset)
	if err != nil {
		repo.log.Error("error while using Query method to take details in storage layer", logger.Error(err))
		return &pb.TasksResponse{}, err
	}
	for rows.Next() {
		var task pb.Task
		err = rows.Scan(
			&task.Id,
			&task.UserId,
			&task.Title,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			repo.log.Error("this error get_all-~~~~~~~ERROR", logger.Error(err))
			return nil, err
		}
		if createdAt.Valid {
			task.CreatedAt = createdAt.Time.Format(configs.Layout)
		}
		if updatedAt.Valid {
			task.UpdatedAt = updatedAt.Time.Format(configs.Layout)
		}
		tasks = append(tasks, &task)

	}

	return &pb.TasksResponse{
		Tasks: tasks,
		Count: count,
	}, nil
}

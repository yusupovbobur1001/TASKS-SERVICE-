package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"net"
	"os"
	"os/signal"
	"syscall"
	"task_service/configs"
	"task_service/grpc"
	"task_service/pkg/logger"
	"task_service/storage"
	mongo "task_service/storage/mongos"
	"task_service/storage/postgres"
)

func main() {
	cfg := configs.Load()

	var loggerLevel string
	switch cfg.Environment {
	case configs.DebugMode:
		loggerLevel = logger.LevelDebug
		gin.SetMode(gin.DebugMode)
	case configs.TestMode:
		loggerLevel = logger.LevelDebug
		gin.SetMode(gin.TestMode)
	default:
		loggerLevel = logger.LevelInfo
		gin.SetMode(gin.ReleaseMode)
	}

	log := logger.NewLogger(cfg.ServiceName, loggerLevel)
	defer logger.Cleanup(log)

	postgresStore, err := postgres.NewStore(context.Background(), log, cfg)
	if err != nil {
		log.Error("error while connecting to mongo", logger.Error(err))
		return
	}
	mongoStore, err := mongo.NewStore(context.Background(), cfg, log)
	if err != nil {
		log.Error("this error that  ERROR-~~~~~~~~~~~~", logger.Error(err))
		return
	}
	Store := storage.NewIStorageRepo(mongoStore.Details(), postgresStore.Tasks())
	grpcServer := grpc.SetUpServer(Store)

	lis, err := net.Listen("tcp", cfg.TaskServiceGrpcHost+cfg.TaskServiceGrpcPort)
	if err != nil {
		log.Error("error while listening grpc host port", logger.Error(err))
		return
	}

	log.Info("Service is running...", logger.Any("grpc port", cfg.TaskServiceGrpcPort))

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Error("error while serving grpc", logger.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down service...")
	grpcServer.GracefulStop()
}

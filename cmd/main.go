package main

import (
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/asadbekGo/book-shop-catalog/config"
	pb "github.com/asadbekGo/book-shop-catalog/genproto/catalog_service"
	"github.com/asadbekGo/book-shop-catalog/pkg/db"
	"github.com/asadbekGo/book-shop-catalog/pkg/logger"
	"github.com/asadbekGo/book-shop-catalog/service"
	"github.com/asadbekGo/book-shop-catalog/storage"
)

func main() {
	cfg := config.Load()

	log := logger.New(cfg.LogLevel, "catalog-service")
	defer func(l logger.Logger) {
		err := logger.Cleanup(l)
		if err != nil {
			log.Fatal("failed cleanup logger", logger.Error(err))
		}
	}(log)

	log.Info("main: sqlxConfig",
		logger.String("host", cfg.PostgresHost),
		logger.Int("port", cfg.PostgresPort),
		logger.String("database", cfg.PostgresDatabase))

	connDB, err := db.ConnectionToDB(cfg)
	if err != nil {
		log.Fatal("sqlx connection to postgres error", logger.Error(err))
	}

	pgStorage := storage.NewStoragePg(connDB)

	catalogService := service.NewCatalogService(pgStorage, log)

	lis, err := net.Listen("tcp", cfg.RPCPort)
	if err != nil {
		log.Fatal("Error while listening: %v", logger.Error(err))
	}

	s := grpc.NewServer()
	pb.RegisterCatalogServiceServer(s, catalogService)
	reflection.Register(s)
	log.Info("main: server running",
		logger.String("port", cfg.RPCPort))
	if err := s.Serve(lis); err != nil {
		log.Fatal("Error while listening: %v", logger.Error(err))
	}
}

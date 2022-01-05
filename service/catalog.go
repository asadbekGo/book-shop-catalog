package service

import (
	"context"

	"github.com/gofrs/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/asadbekGo/book-shop-catalog/genproto"
	l "github.com/asadbekGo/book-shop-catalog/pkg/logger"
	"github.com/asadbekGo/book-shop-catalog/storage"
)

// CatalogService ...
type CatalogService struct {
	storage storage.IStorage
	logger  l.Logger
}

// NewCatalogService ...
func NewCatalogService(storage storage.IStorage, log l.Logger) *CatalogService {
	return &CatalogService{
		storage: storage,
		logger:  log,
	}
}

func (s *CatalogService) CreateAuthor(ctx context.Context, req *pb.Author) (*pb.Author, error) {
	id, err := uuid.NewV4()
	if err != nil {
		s.logger.Error("failed while generating uuid", l.Error(err))
		return nil, status.Error(codes.Internal, "failed generate uuid")
	}
	req.Id = id.String()

	author, err := s.storage.Catalog().CreateAuthor(*req)
	if err != nil {
		s.logger.Error("failed to create author", l.Error(err))
		return nil, status.Error(codes.Internal, "failed to create author")
	}

	return &author, nil
}

func (s *CatalogService) GetAuthor(ctx context.Context, req *pb.Author) (*pb.Author, error) {
	id, err := uuid.NewV4()
	if err != nil {
		s.logger.Error("failed while generating uuid", l.Error(err))
		return nil, status.Error(codes.Internal, "failed generate uuid")
	}
	req.Id = id.String()

	author, err := s.storage.Catalog().CreateAuthor(*req)
	if err != nil {
		s.logger.Error("failed to create author", l.Error(err))
		return nil, status.Error(codes.Internal, "failed to create author")
	}

	return &author, nil
}

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
	author, err := s.storage.Catalog().CreateAuthor(*req)
	if err != nil {
		s.logger.Error("failed to get author", l.Error(err))
		return nil, status.Error(codes.Internal, "failed to get author")
	}

	return &author, nil
}

func (s *CatalogService) GetAuthors(ctx context.Context, req *pb.ListReq) (*pb.AuthorListResp, error) {
	authors, count, err := s.storage.Catalog().GetAuthors(req.Page, req.Limit)
	if err != nil {
		s.logger.Error("failed to author list", l.Error(err))
		return nil, status.Error(codes.Internal, "failed to author list")
	}

	return &pb.AuthorListResp{
		Authors: authors,
		Count:   count,
	}, nil
}

func (s *CatalogService) UpdateAuthor(ctx context.Context, req *pb.Author) (*pb.Author, error) {
	author, err := s.storage.Catalog().UpdateAuthor(*req)
	if err != nil {
		s.logger.Error("failed to update author", l.Error(err))
		return nil, status.Error(codes.Internal, "failed to update author")
	}

	return &author, nil
}

func (s *CatalogService) DeleteAuthor(ctx context.Context, req *pb.ByIdReq) (*pb.Empty, error) {
	err := s.storage.Catalog().DeleteAuthor(req.Id)
	if err != nil {
		s.logger.Error("failed to delete author", l.Error(err))
		return nil, status.Error(codes.Internal, "failed to delete author")
	}

	return &pb.Empty{}, nil
}

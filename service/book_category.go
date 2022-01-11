package service

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/asadbekGo/book-shop-catalog/genproto/catalog_service"
	l "github.com/asadbekGo/book-shop-catalog/pkg/logger"
)

func (s *CatalogService) CreateBookCategory(ctx context.Context, req *pb.BookCategory) (*pb.BookResp, error) {
	bookCategoryResp, err := s.storage.Catalog().CreateBookCategory(*req)
	if err != nil {
		s.logger.Error("failed to create book category", l.Error(err))
		return nil, status.Error(codes.Internal, "failed to create book category")
	}

	return &bookCategoryResp, nil
}

func (s *CatalogService) DeleteBookCategory(ctx context.Context, req *pb.BookCategory) (*pb.Empty, error) {
	err := s.storage.Catalog().DeleteBookCategory(*req)
	if err != nil {
		s.logger.Error("failed delete book category", l.Error(err))
		return nil, status.Error(codes.Internal, "failed delete book category")
	}

	return &pb.Empty{}, nil
}

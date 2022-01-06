package service

import (
	"context"

	pb "github.com/asadbekGo/book-shop-catalog/genproto/catalog_service"
	l "github.com/asadbekGo/book-shop-catalog/pkg/logger"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *CatalogService) CreateBookCategory(ctx context.Context, req *pb.BookCategory) (*pb.BookCategoryResp, error) {
	bookCategoryResp, err := s.storage.Catalog().CreateBookCategory(*req)
	if err != nil {
		s.logger.Error("failed to create book category", l.Error(err))
		return nil, status.Error(codes.Internal, "failed to create book category")
	}

	return &bookCategoryResp, nil
}

func (s *CatalogService) GetBookCategory(ctx context.Context, req *pb.ByIdReq) (*pb.BookCategoryResp, error) {
	book, err := s.storage.Catalog().GetBookCategory(req.Id)
	if err != nil {
		s.logger.Error("failed to get book category", l.Error(err))
		return nil, status.Error(codes.Internal, "failed to get book category")
	}

	return &book, nil
}

func (s *CatalogService) GetBookCategories(ctx context.Context, req *pb.ListReq) (*pb.BookCategoryListResp, error) {
	bookCategoryList, count, err := s.storage.Catalog().GetBookCategories(req.Page, req.Limit)
	if err != nil {
		s.logger.Error("failed to get book categories list", l.Error(err))
		return nil, status.Error(codes.Internal, "failed to get book categories list")
	}

	return &pb.BookCategoryListResp{
		BookCategory: bookCategoryList,
		Count:        count,
	}, nil
}

func (s *CatalogService) DeleteBookCategory(ctx context.Context, req *pb.BookCategory) (*pb.Empty, error) {
	err := s.storage.Catalog().DeleteBookCategory(*req)
	if err != nil {
		s.logger.Error("failed delete book category", l.Error(err))
		return nil, status.Error(codes.Internal, "failed delete book category")
	}

	return &pb.Empty{}, nil
}

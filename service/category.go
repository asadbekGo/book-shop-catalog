package service

import (
	"context"

	"github.com/gofrs/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/asadbekGo/book-shop-catalog/genproto/catalog_service"
	l "github.com/asadbekGo/book-shop-catalog/pkg/logger"
)

func (s *CatalogService) CreateCategory(ctx context.Context, req *pb.Category) (*pb.Category, error) {
	id, err := uuid.NewV4()
	if err != nil {
		s.logger.Error("failed while generating uuid", l.Error(err))
		return nil, status.Error(codes.Internal, "failed generate uuid")
	}
	req.Id = id.String()

	category, err := s.storage.Catalog().CreateCategory(*req)
	if err != nil {
		s.logger.Error("falied to create category", l.Error(err))
		return nil, status.Error(codes.Internal, "failed to create category")
	}

	return &category, nil
}

func (s *CatalogService) GetCategory(ctx context.Context, req *pb.ByIdReq) (*pb.Category, error) {
	category, err := s.storage.Catalog().GetCategory(req.Id)
	if err != nil {
		s.logger.Error("failed to get category", l.Error(err))
		return nil, status.Error(codes.Internal, "failed to get category")
	}

	return &category, nil
}

func (s *CatalogService) GetCategories(ctx context.Context, req *pb.ListReq) (*pb.CategoryListResp, error) {
	categories, count, err := s.storage.Catalog().GetCategories(req.Page, req.Limit)
	if err != nil {
		s.logger.Error("failed to list category", l.Error(err))
		return nil, status.Error(codes.Internal, "failed to list category")
	}

	return &pb.CategoryListResp{
		Categories: categories,
		Count:      count,
	}, nil
}

func (s *CatalogService) UpdateCategory(ctx context.Context, req *pb.Category) (*pb.Category, error) {
	category, err := s.storage.Catalog().UpdateCategory(*req)
	if err != nil {
		s.logger.Error("failed to update category", l.Error(err))
		return nil, status.Error(codes.Internal, "failed to update category")
	}

	return &category, nil
}

func (s *CatalogService) DeleteCategory(ctx context.Context, req *pb.ByIdReq) (*pb.Empty, error) {
	err := s.storage.Catalog().DeleteCategory(req.Id)
	if err != nil {
		s.logger.Error("failed to delete category", l.Error(err))
		return nil, status.Error(codes.Internal, "failed to delete category")
	}

	return &pb.Empty{}, nil
}

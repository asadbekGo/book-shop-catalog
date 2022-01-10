package service

import (
	"context"

	"github.com/gofrs/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/asadbekGo/book-shop-catalog/genproto/catalog_service"
	l "github.com/asadbekGo/book-shop-catalog/pkg/logger"
)

func (s *CatalogService) CreateBook(ctx context.Context, req *pb.Book) (*pb.BookResp, error) {
	id, err := uuid.NewV4()
	if err != nil {
		s.logger.Error("failed while generating uuid", l.Error(err))
		return nil, status.Error(codes.Internal, "failed generate uuid")
	}
	req.Id = id.String()

	book, err := s.storage.Catalog().CreateBook(*req)
	if err != nil {
		s.logger.Error("failed to create book", l.Error(err))
		return nil, status.Error(codes.Internal, "failed to create book")
	}

	return &book, nil
}

func (s *CatalogService) GetBook(ctx context.Context, req *pb.ByIdReq) (*pb.BookResp, error) {
	book, err := s.storage.Catalog().GetBook(req.Id)
	if err != nil {
		s.logger.Error("failed to get book", l.Error(err))
		return nil, status.Error(codes.Internal, "failed to get book")
	}

	return &book, nil
}

func (s *CatalogService) GetBooks(ctx context.Context, req *pb.BookListReq) (*pb.BookListResp, error) {
	books, count, err := s.storage.Catalog().GetBooks(req.Page, req.Limit, req.Filters)
	if err != nil {
		s.logger.Error("failed to get books list", l.Error(err))
		return nil, status.Error(codes.Internal, "failed to get books list")
	}

	return &pb.BookListResp{
		Books: books,
		Count: count,
	}, nil
}

func (s *CatalogService) UpdateBook(ctx context.Context, req *pb.Book) (*pb.BookResp, error) {
	book, err := s.storage.Catalog().UpdateBook(*req)
	if err != nil {
		s.logger.Error("failed to update book", l.Error(err))
		return nil, status.Error(codes.Internal, "failed to update list")
	}

	return &book, nil
}

func (s *CatalogService) DeleteBook(ctx context.Context, req *pb.ByIdReq) (*pb.Empty, error) {
	err := s.storage.Catalog().DeleteBook(req.Id)
	if err != nil {
		s.logger.Error("failed delete book", l.Error(err))
		return nil, status.Error(codes.Internal, "failed delete book")
	}

	return &pb.Empty{}, nil
}

package repo

import (
	pb "github.com/asadbekGo/book-shop-catalog/genproto/catalog_service"
)

// CatalogStorageI ...
type CatalogStorageI interface {
	CreateAuthor(pb.Author) (pb.Author, error)
	GetAuthor(id string) (pb.Author, error)
	GetAuthors(page, limit int64) ([]*pb.Author, int64, error)
	UpdateAuthor(pb.Author) (pb.Author, error)
	DeleteAuthor(id string) error

	CreateCategory(pb.Category) (pb.Category, error)
	GetCategory(id string) (pb.Category, error)
	GetCategories(page, limit int64) ([]*pb.Category, int64, error)
	UpdateCategory(pb.Category) (pb.Category, error)
	DeleteCategory(id string) error

	CreateBook(pb.Book) (pb.BookResp, error)
	GetBook(id string) (pb.BookResp, error)
	GetBooks(page, limit int64, filters map[string]string) ([]*pb.BookResp, int64, error)
	UpdateBook(pb.Book) (pb.BookResp, error)
	DeleteBook(id string) error
}

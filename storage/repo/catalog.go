package repo

import (
	pb "github.com/asadbekGo/book-shop-catalog/genproto"
)

// CatalogStorageI ...
type CatalogStorageI interface {
	CreateAuthor(pb.Author) (pb.Author, error)
	GetAuthor(pb.ByIdReq) (pb.Author, error)
	GetAuthors(pb.ListReq) (pb.AuthorListResp, error)
	UpdateAuthor(pb.Author) (pb.Author, error)
	DeleteAuthor(pb.ByIdReq) (pb.Empty, error)

	CreateCategory(pb.Category) (pb.Category, error)
	GetCategory(pb.ByIdReq) (pb.Category, error)
	GetCategories(pb.ListReq) (pb.CategoryListResp, error)
	UpdateCategory(pb.Category) (pb.Category, error)
	DeleteCategory(pb.ByIdReq) (pb.Empty, error)

	CreateBook(pb.NewBook) (pb.Book, error)
	GetBook(pb.ByIdReq) (pb.Book, error)
	GetBooks(pb.ListReq) (pb.BookListResp, error)
	UpdateBook(pb.NewBook) (pb.Book, error)
	DeleteBook(pb.ByIdReq) (pb.Empty, error)
}

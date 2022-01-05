package postgres

import (
	"database/sql"

	
	"github.com/jmoiron/sqlx"
	
	pb "github.com/asadbekGo/book-shop-catalog/genproto"
)

type catalogRepo struct {
	db *sqlx.DB
}

// NewCatalogRepo ...
func NewCatalogRepo(db *sqlx.DB) *catalogRepo {
	return &catalogRepo{db: db}
}

func (r *catalogRepo) CreateBook(book pb.NewBook) (pb.Book, error) {
	var id string
	err := r.db.QueryRow(`
        INSERT INTO book(book_id, name, author_id, parent_uuid)
        VALUES ($1,$2,$3,$4) returning book_id`, book.Id, book.Name, book.AuthorId, book.CategoryId.Scan(&id)
	if err != nil {
		return pb.Book{}, err
	}
	book, err = r.GetBook(id)
	if err != nil {
		return pb.Book{}, err
	}

	return book, nil
}
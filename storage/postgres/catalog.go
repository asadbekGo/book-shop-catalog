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

	for _, i range := book.Category {
		err := r.db.QueryRow(`
			INSERT INTO book(book_id, name, author_id, parent_uuid)
			VALUES ($1,$2,$3,$4) returning book_id`, book.Id, book.Name, book.AuthorId, i).Scan(&id)
		if err != nil {
			return pb.Book{}, err
		}
	}
	book, err = r.GetBook(id)
	if err != nil {
		return pb.Book{}, err
	}

	return book, nil
}


func (r *catalogRepo) GetBook(id string) (pb.Book, error) {
	var (
		book pb.Book
		authorId string	
		categoryId []string
	)

	err := r.db.QueryRow(`
        SELECT book_id, name, author_id, array_arg(parent_uuid) FROM book
        WHERE id=$1 GROUP BY book_id`, id).Scan(&book.Id, &book.Name, &authorId, &categoryId)
	if err != nil {
		return pb.Task{}, err
	}
	err = r.db.QueryRow(
		`SELECT author_id, name FROM author WHERE author_id = $1`, authorId).Scan(&book.Author.Id, &book.Author.Name)

	if err != nil {
		return pb.Book{}, err
	}
	
	for , i := range categoryId {
		var category pb.Category
		err = r.db.QueryRow(
			`SELECT category_id, name FROM category WHERE category_id = $1`, i).Scan(&category.Id, &category.Name)
	
		if err != nil {
			return pb.Book{}, err
		}
		book.Category = append(book.Category, category)
	}
		
	return book, nil
}


func (r *catalogRepo) GetBooks(page, limit int64) ([]*pb.Book, int64, err){

	offset := (page - 1) * limit
	rows, err := r.db.Queryx(
		`SELECT book_id, name, author_id, array_arg(parent_uuid) FROM book GROUP BY book_id OFFSET $1  LIMIT $2`,
		offset, limit)
	if err != nil {
		return nil, 0, err
	}
	if err = rows.Err(); err != nil {
		return nil, 0, err
	}
	defer rows.Close() // nolint:errcheck


	var (
		books []*pb.Book
		count int64
	)

	for rows.Next() {
		var (
			book pb.Book
			authorId string
			categoryId string
		)
		err = rows.Scan(&book.Id, &book.Name, &authorId, &categoryId)
		if err != nil {
			return nil, 0, err
		}

		err = r.db.QueryRow(
			`SELECT author_id, name FROM author WHERE author_id = $1`, authorId).Scan(&book.Author.Id, &book.Author.Name)
	
		if err != nil {
			return nil, 0, err
		}
		
		for , i := range categoryId {
			var category pb.Category
			err = r.db.QueryRow(
				`SELECT category_id, name FROM category WHERE category_id = $1`, i).Scan(&category.Id, &category.Name)
		
			if err != nil {
				return nil, 0, err
			}
			book.Category = append(book.Category, category)

		}

		books = append(books, book)
	}

		return books, count, err
}




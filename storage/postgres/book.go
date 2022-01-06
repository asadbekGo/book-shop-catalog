package postgres

import (
	"database/sql"

	pb "github.com/asadbekGo/book-shop-catalog/genproto/catalog_service"
)

func (r *catalogRepo) CreateBook(book pb.Book) (pb.Book, error) {
	var id string
	err := r.db.QueryRow(`
		INSERT INTO book(book_id, name, author_id, updated_at) 
		VALUES ($1, $2, $3, current_timestamp) RETURNING book_id`,
		book.Id,
		book.Name,
		book.AuthorId,
	).Scan(&id)
	if err != nil {
		return pb.Book{}, err
	}

	book, err = r.GetBook(id)

	if err != nil {
		return pb.Book{}, err
	}

	return book, nil
}

func (r *catalogRepo) GetBook(id string) (pb.Book, error) {
	var book pb.Book
	err := r.db.QueryRow(`
		SELECT book_id, name, author_id, created_at, updated_at
		FROM book WHERE book_id=$1 AND deleted_at IS NULL`, id).Scan(
		&book.Id,
		&book.Name,
		&book.AuthorId,
		&book.CreatedAt,
		&book.UpdatedAt,
	)
	if err != nil {
		return pb.Book{}, err
	}

	return book, nil
}

func (r *catalogRepo) GetBooks(page, limit int64, filters map[string]string) ([]*pb.Book, int64, error) {
	offset := (page - 1) * limit

	rows, err := r.db.Query(`
		SELECT book_id, name, author_id, created_at, updated_at FROM book WHERE deleted_at IS NULL LIMIT $1 OFFSET $2`, limit, offset)
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
		var book pb.Book
		err = rows.Scan(
			&book.Id,
			&book.Name,
			&book.AuthorId,
			&book.CreatedAt,
			&book.UpdatedAt,
		)

		if err != nil {
			return nil, 0, err
		}

		books = append(books, &book)
	}

	err = r.db.QueryRow(`SELECT count(*) FROM book WHERE deleted_ad IS NULL`).Scan(&count)
	if err != nil {
		return nil, 0, err
	}

	return books, count, nil
}

func (r *catalogRepo) UpdateBook(book pb.Book) (pb.Book, error) {
	result, err := r.db.Exec(`
		UPDATE book SET name=$1, updated_at=current_timestamp
		WHERE book_id=$2 AND deleted_at IS NULL`,
		book.Name, book.Id)
	if err != nil {
		return pb.Book{}, err
	}
	if i, _ := result.RowsAffected(); i == 0 {
		return pb.Book{}, sql.ErrNoRows
	}

	book, err = r.GetBook(book.Id)
	if err != nil {
		return pb.Book{}, err
	}

	return book, nil
}

func (r *catalogRepo) DeleteBook(id string) error {
	result, err := r.db.Exec(`
	UPDATE book SET deleted_at=current_timestamp WHERE book_id=$1 AND deleted_at IS NULL`, id)
	if err != nil {
		return err
	}
	if i, _ := result.RowsAffected(); i == 0 {
		return sql.ErrNoRows
	}

	return nil
}

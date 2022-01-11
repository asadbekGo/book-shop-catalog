package postgres

import (
	"database/sql"

	pb "github.com/asadbekGo/book-shop-catalog/genproto/catalog_service"
)

func (r *catalogRepo) CreateBookCategory(bookCategory pb.BookCategory) (pb.BookResp, error) {
	var id string
	err := r.db.QueryRow(`
		INSERT INTO book_category(book_id, category_id) 
		VALUES ($1, $2) RETURNING book_id`,
		bookCategory.BookId,
		bookCategory.CategoryId,
	).Scan(&id)
	if err != nil {
		return pb.BookResp{}, err
	}

	book, err := r.GetBook(id)
	if err != nil {
		return pb.BookResp{}, err
	}

	return book, nil
}

func (r *catalogRepo) GetBookCategory(id string) ([]*pb.Category, error) {
	rows, err := r.db.Query(`
			SELECT
				c.category_id,
				c.name,
				c.parent_uuid,
				bc.created_at
			FROM
				book_category as bc
			JOIN
				book as b using(book_id)
			JOIN
				category as c using(category_id)
			WHERE b.book_id = $1 AND bc.deleted_at IS NULL`, id)
	if err != nil {
		return nil, err
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	defer rows.Close() // nolint:errcheck

	var categories []*pb.Category

	for rows.Next() {
		var parent_uuid sql.NullString
		var category pb.Category
		err = rows.Scan(
			&category.Id,
			&category.Name,
			&parent_uuid,
			&category.CreatedAt,
		)

		if err != nil {
			return nil, err
		}
		if parent_uuid.Valid {
			category.ParentUUID = parent_uuid.String
		}
		categories = append(categories, &category)
	}

	return categories, nil
}

func (r *catalogRepo) DeleteBookCategory(bookCategory pb.BookCategory) error {
	result, err := r.db.Exec(`
		UPDATE book_category SET deleted_at=current_timestamp WHERE book_id=$1 AND category_id=$2 AND deleted_at IS NULL`,
		bookCategory.BookId, bookCategory.CategoryId)
	if err != nil {
		return err
	}
	if i, _ := result.RowsAffected(); i == 0 {
		return sql.ErrNoRows
	}

	return nil
}

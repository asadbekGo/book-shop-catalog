package postgres

import (
	"database/sql"

	pb "github.com/asadbekGo/book-shop-catalog/genproto/catalog_service"
)

func CreateBookCategory(r *catalogRepo, bookCategory pb.BookCategory) error {
	var id string
	err := r.db.QueryRow(`
		INSERT INTO book_category(book_id, category_id) 
		VALUES ($1, $2) RETURNING book_id`,
		bookCategory.BookId,
		bookCategory.CategoryId,
	).Scan(&id)
	if err != nil {
		return err
	}

	return nil
}

func GetBookCategory(r *catalogRepo, id string) ([]*pb.Category, error) {
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
			WHERE b.book_id = $1`, id)
	if err != nil {
		return nil, err
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	defer rows.Close() // nolint:errcheck

	var categories []*pb.Category

	for rows.Next() {
		var category pb.Category
		err = rows.Scan(
			&category.Id,
			&category.Name,
			&category.ParentUUID,
			&category.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		categories = append(categories, &category)
	}

	return categories, nil
}

func DeleteBookCategory(r *catalogRepo, id string) error {
	result, err := r.db.Exec(`DELETE FROM book_category WHERE book_id=$1`,
		id)
	if err != nil {
		return err
	}
	if i, _ := result.RowsAffected(); i == 0 {
		return sql.ErrNoRows
	}

	return nil
}

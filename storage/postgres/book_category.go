package postgres

import (
	"database/sql"

	pb "github.com/asadbekGo/book-shop-catalog/genproto/catalog_service"
)

func (r *catalogRepo) CreateBookCategory(bookCategory pb.BookCategory) (pb.BookCategoryResp, error) {
	var id string
	err := r.db.QueryRow(`
		INSERT INTO book_category(book_id, category_id) 
		VALUES ($1, $2) RETURNING book_id`,
		bookCategory.BookId,
		bookCategory.CategoryId,
	).Scan(&id)

	if err != nil {
		return pb.BookCategoryResp{}, err
	}

	bookCategoryResp, err := r.GetBookCategory(id)

	if err != nil {
		return pb.BookCategoryResp{}, err
	}

	return bookCategoryResp, nil
}

func (r *catalogRepo) GetBookCategory(id string) (pb.BookCategoryResp, error) {
	var book pb.Book
	err := r.db.QueryRow(`
		SELECT book_id, name, author_id, created_at
		FROM book WHERE book_id=$1 AND deleted_at IS NULL`, id).Scan(
		&book.Id,
		&book.Name,
		&book.AuthorId,
		&book.CreatedAt,
	)

	if err != nil {
		return pb.BookCategoryResp{}, err
	}

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
			WHERE b.book_id = $1 AND bc.deleted_at IS NULL
		`, id)
	if err != nil {
		return pb.BookCategoryResp{}, err
	}
	if err = rows.Err(); err != nil {
		return pb.BookCategoryResp{}, err
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
			return pb.BookCategoryResp{}, err
		}

		categories = append(categories, &category)
	}

	var BookCategoryResp pb.BookCategoryResp

	BookCategoryResp.Book = &book
	BookCategoryResp.Category = categories

	return BookCategoryResp, nil
}

func (r *catalogRepo) GetBookCategories(page, limit int64) ([]*pb.BookCategoryResp, int64, error) {
	offset := (page - 1) * limit

	rows, err := r.db.Query(`
		SELECT book_id FROM book_category WHERE deleted_at IS NULL GROUP BY book_id LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	if err = rows.Err(); err != nil {
		return nil, 0, err
	}
	defer rows.Close() // nolint:errcheck

	var (
		bookCategoryList []*pb.BookCategoryResp
		count            int64
	)

	for rows.Next() {
		var id string
		err = rows.Scan(
			&id,
		)

		if err != nil {
			return nil, 0, err
		}

		bookCategory, err := r.GetBookCategory(id)
		if err != nil {
			return nil, 0, err
		}

		bookCategoryList = append(bookCategoryList, &bookCategory)
	}

	err = r.db.QueryRow(`SELECT count(book_id) FROM book_category WHERE deleted_at IS NULL GROUP BY book_id`).Scan(&count)
	if err != nil {
		return nil, 0, err
	}

	return bookCategoryList, count, nil
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

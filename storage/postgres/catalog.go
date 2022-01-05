package postgres

import (
	"database/sql"

	"github.com/jmoiron/sqlx"

	pb "github.com/asadbekGo/book-shop-catalog/genproto/catalog_service"
)

type catalogRepo struct {
	db *sqlx.DB
}

// NewCatalogRepo ...
func NewCatalogRepo(db *sqlx.DB) *catalogRepo {
	return &catalogRepo{db: db}
}

func (r *catalogRepo) CreateCategory(category pb.Category) (pb.Category, error) {
	var id string
	err := r.db.QueryRow(`
	INSERT INTO category(category_id, name, updated_at)
		VALUES ($1, $2, current_timestamp) returning category_id`,
		category.Id,
		category.Name,
	).Scan(&id)
	if err != nil {
		return pb.Category{}, err
	}

	category, err = r.GetCategory(id)

	if err != nil {
		return pb.Category{}, nil
	}

	return category, nil
}

func (r *catalogRepo) GetCategory(id string) (pb.Category, error) {
	var category pb.Category

	err := r.db.QueryRow(`
	SELECT category_id, name, created_at, updated_at FROM category WHERE category_id=$1 and deleted_at is null`, id).Scan(
		&category.Id,
		&category.Name,
		&category.CreatedAt,
		&category.UpdatedAt,
	)
	if err != nil {
		return pb.Category{}, err
	}

	return category, nil
}

func (r *catalogRepo) GetCategories(page, limit int64) ([]*pb.Category, int64, error) {
	offset := (page - 1) * limit

	rows, err := r.db.Queryx(`
	SELECT category_id, name, created_at, updated_at FROM category WHERE deleted_at is null LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	if err = rows.Err(); err != nil {
		return nil, 0, err
	}
	defer rows.Close() // nolint:errcheck

	var (
		categories []*pb.Category
		count      int64
	)

	for rows.Next() {
		var category pb.Category
		err = rows.Scan(
			&category.Id,
			&category.Name,
			&category.CreatedAt,
			&category.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}

		categories = append(categories, &category)
	}

	err = r.db.QueryRow(`SELECT count(*) FROM category`).Scan(&count)
	if err != nil {
		return nil, 0, err
	}

	return categories, count, nil
}

func (r *catalogRepo) UpdateCategory(category pb.Category) (pb.Category, error) {
	result, err := r.db.Exec(`
		UPDATE Category SET name=$1, updated_at=current_timestamp
		WHERE category_id=$2 and deleted_at is null`,
		category.Name,
		category.Id,
	)
	if err != nil {
		return pb.Category{}, err
	}
	if i, _ := result.RowsAffected(); i == 0 {
		return pb.Category{}, sql.ErrNoRows
	}

	category, err = r.GetCategory(category.Id)
	if err != nil {
		return pb.Category{}, err
	}

	return category, nil
}

func (r *catalogRepo) DeleteCategory(id string) error {
	result, err := r.db.Exec(`
		UPDATE category SET deleted_at=current_timestamp WHERE category_id=$1`, id)
	if err != nil {
		return err
	}
	if i, _ := result.RowsAffected(); i == 0 {
		return sql.ErrNoRows
	}

	return nil
}

package postgres

import (
	"database/sql"
	"time"

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

func (r *catalogRepo) CreateCategory(category pb.Category) (pb.Category, error) {
	var id string
	err := r.db.QueryRow(`
		INSERT INTO category(category_id, name status)
		VALUES ($1, $2, $3) returning category_id`,
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
		SELECT category_id, name FROM category 
		WHERE category_id=$1`, id).Scan(
		&category.Id,
		&category.Name,
			)
	if err != nil {
		return pb.Category{}, err
	}

	return category, nil
}

func (r *catalogRepo) GetCategories(page, limit int64) ([]*pb.Category, int64, error) {
	offset := (page - 1) * limit

	rows, err := r.db.Queryx(`
		SELECT category_id, name FROM category
		LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	if err = rows.Err(); err != nil {
		return nil, 0, err
	}
	defer rows.Close() // nolint:errcheck

	var (
		categories []*pb.Category
		count int64
	)

	for rows.Next() {
		var category pb.Category
		err = rows.Scan(
			&category.Id,
			&category.Name,
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
		UPDATE Category SET name=$1
		WHERE id=$2`,
		category.Id, 
		category.Name,
	)
	if err != nil {
		return pb.Category{}, err
	}
	if i, _ := result.RowsAffected(); i == 0 {
		return pb.Category{}, sql.ErrNoRows
	}

	category, err = r.Get(category.Id)
	if err != nil {
		return pb.Category{}, err
	}

	return category, nil
}

func (r *catalogRepo) Delete(id string) error {
	result, err := r.db.Exec(`
		delete * from category where category_id = $1`, id)
	if err != nil {
		return err
	}
	if i, _ := result.RowsAffected(); i == 0 {
		return sql.ErrNoRows
	}

	return nil
}
package postgres

import (
	"database/sql"

	pb "github.com/asadbekGo/book-shop-catalog/genproto/catalog_service"
)

func (r *catalogRepo) CreateCategory(category pb.Category) (pb.Category, error) {
	var id string
	var parentUUID interface{}

	if category.ParentUUID == "" {
		parentUUID = nil
	} else {
		parentUUID = category.ParentUUID
	}

	err := r.db.QueryRow(`
	INSERT INTO category(category_id, name, parent_uuid, updated_at)
		VALUES ($1, $2, $3, current_timestamp) RETURNING category_id`,
		category.Id,
		category.Name,
		parentUUID,
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
	var parentUUID sql.NullString

	err := r.db.QueryRow(`
	SELECT category_id, name, parent_uuid, created_at, updated_at FROM category WHERE category_id=$1 and deleted_at IS NULL`, id).Scan(
		&category.Id,
		&category.Name,
		&parentUUID,
		&category.CreatedAt,
		&category.UpdatedAt,
	)

	if parentUUID.Valid {
		category.ParentUUID = parentUUID.String
	} else {
		category.ParentUUID = ""
	}

	if err != nil {
		return pb.Category{}, err
	}

	return category, nil
}

func (r *catalogRepo) GetCategories(page, limit int64) ([]*pb.Category, int64, error) {
	offset := (page - 1) * limit

	rows, err := r.db.Queryx(`
	SELECT category_id, name, parent_uuid, created_at, updated_at FROM category WHERE deleted_at IS NULL LIMIT $1 OFFSET $2`, limit, offset)
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
		var parentUUID sql.NullString
		var category pb.Category
		err = rows.Scan(
			&category.Id,
			&category.Name,
			&parentUUID,
			&category.CreatedAt,
			&category.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}

		if parentUUID.Valid {
			category.ParentUUID = parentUUID.String
		} else {
			category.ParentUUID = ""
		}

		categories = append(categories, &category)
	}

	err = r.db.QueryRow(`SELECT count(*) FROM category WHERE deleted_at IS NULL`).Scan(&count)
	if err != nil {
		return nil, 0, err
	}

	return categories, count, nil
}

func (r *catalogRepo) UpdateCategory(category pb.Category) (pb.Category, error) {
	var parentUUID interface{}

	row, err := r.GetCategory(category.Id)
	if err != nil {
		return pb.Category{}, nil
	}

	if category.ParentUUID == "" {
		if row.ParentUUID == "" {
			parentUUID = nil
		} else {
			parentUUID = row.ParentUUID
		}
	} else {
		parentUUID = category.ParentUUID
	}
	result, err := r.db.Exec(`
		UPDATE category SET name=$1, parent_uuid=$2, updated_at=current_timestamp
		WHERE category_id=$3 AND deleted_at IS NULL`,
		category.Name,
		parentUUID,
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
		UPDATE category SET deleted_at=current_timestamp WHERE category_id=$1 AND deleted_at IS NULL`, id)
	if err != nil {
		return err
	}
	if i, _ := result.RowsAffected(); i == 0 {
		return sql.ErrNoRows
	}

	return nil
}

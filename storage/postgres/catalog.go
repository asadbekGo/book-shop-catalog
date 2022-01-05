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

func (r *catalogRepo) CreateAuthor(author pb.Author) (pb.Author, error) {
	var id string
	err := r.db.QueryRow(`
		INSERT INTO author(author_id, name, updated_at)
		VALUES ($1, $2, current_timestamp) returning author_id`,
		author.Id,
		author.Name,
	).Scan(&id)
	if err != nil {
		return pb.Author{}, err
	}

	author, err = r.GetAuthor(id)

	if err != nil {
		return pb.Author{}, nil
	}

	return author, nil
}

func (r *catalogRepo) GetAuthor(id string) (pb.Author, error) {
	var author pb.Author
	err := r.db.QueryRow(`
		SELECT author_id, name, created_at, updated_at FROM author WHERE author_id=$1 and deleted_at is null`, id).Scan(
		&author.Id,
		&author.Name,
		&author.CreatedAt,
		&author.UpdatedAt,
	)

	if err != nil {
		return pb.Author{}, err
	}

	return author, nil
}

func (r *catalogRepo) GetAuthors(page, limit int64) ([]*pb.Author, int64, error) {
	offset := (page - 1) * limit

	rows, err := r.db.Query(`
	SELECT author_id, name, created_at, updated_at FROM author WHERE deleted_at is null LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	if err = rows.Err(); err != nil {
		return nil, 0, err
	}
	defer rows.Close() // nolint:errcheck

	var (
		authors []*pb.Author
		count   int64
	)

	for rows.Next() {
		var author pb.Author
		err = rows.Scan(
			&author.Id,
			&author.Name,
			&author.CreatedAt,
			&author.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}

		authors = append(authors, &author)
	}

	err = r.db.QueryRow(`SELECT count(*) FROM author WHERE deleted_at is null`).Scan(&count)
	if err != nil {
		return nil, 0, err
	}

	return authors, count, nil
}

func (r *catalogRepo) UpdateAuthor(author pb.Author) (pb.Author, error) {
	result, err := r.db.Exec(`
		UPDATE author SET name=$1, updated_at=current_timestamp
		WHERE author_id=$2 and deleted_at is null`,
		author.Name, author.Id)
	if err != nil {
		return pb.Author{}, err
	}
	if i, _ := result.RowsAffected(); i == 0 {
		return pb.Author{}, sql.ErrNoRows
	}

	author, err = r.GetAuthor(author.Id)
	if err != nil {
		return pb.Author{}, err
	}

	return author, nil
}

func (r *catalogRepo) DeleteAuthor(id string) error {
	result, err := r.db.Exec(`
	UPDATE author SET deleted_at=current_timestamp WHERE author_id=$1`, id)
	if err != nil {
		return err
	}
	if i, _ := result.RowsAffected(); i == 0 {
		return sql.ErrNoRows
	}

	return nil
}

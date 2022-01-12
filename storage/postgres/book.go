package postgres

import (
	"database/sql"
	"strings"

	"github.com/huandu/go-sqlbuilder"

	pb "github.com/asadbekGo/book-shop-catalog/genproto/catalog_service"
	"github.com/asadbekGo/book-shop-catalog/pkg/utils"
)

func (r *catalogRepo) CreateBook(book pb.Book) (pb.BookResp, error) {
	var id string
	err := r.db.QueryRow(`
		INSERT INTO book(book_id, name, author_id, updated_at) 
		VALUES ($1, $2, $3, current_timestamp) RETURNING book_id`,
		book.Id,
		book.Name,
		book.AuthorId,
	).Scan(&id)
	if err != nil {
		return pb.BookResp{}, err
	}

	for _, category := range book.Categories {
		var bookCategory pb.BookCategory
		bookCategory.BookId = id
		bookCategory.CategoryId = category
		err = CreateBookCategory(r, bookCategory)
		if err != nil {
			return pb.BookResp{}, err
		}
	}

	bookResp, err := r.GetBook(id)
	if err != nil {
		return pb.BookResp{}, err
	}

	return bookResp, nil
}

func (r *catalogRepo) GetBook(id string) (pb.BookResp, error) {
	var book pb.BookResp
	var authorId string
	err := r.db.QueryRow(`
		SELECT book_id, name, author_id, created_at, updated_at
		FROM book WHERE book_id=$1 AND deleted_at IS NULL`, id).Scan(
		&book.Id,
		&book.Name,
		&authorId,
		&book.CreatedAt,
		&book.UpdatedAt,
	)
	if err != nil {
		return pb.BookResp{}, err
	}
	author, err := r.GetAuthor(authorId)
	if err != nil {
		return pb.BookResp{}, err
	}
	book.Author = &author

	categories, err := GetBookCategory(r, book.Id)
	if err != nil {
		return pb.BookResp{}, err
	}
	book.Category = categories

	return book, nil
}

func (r *catalogRepo) GetBooks(page, limit int64, filters map[string]string) ([]*pb.BookResp, int64, error) {
	offset := (page - 1) * limit
	sb := sqlbuilder.NewSelectBuilder()

	sb.Select("b.book_id", "b.name", "b.author_id", "b.created_at", "b.updated_at")
	sb.From("book b")
	sb.Where("b.deleted_at IS NULL")

	if val, ok := filters["category"]; ok && val != "" {
		args := utils.StringSliceToInterfaceSlice(utils.ParseFilter(val))
		sb.JoinWithOption("LEFT", "book_category bc", "b.book_id=bc.book_id")

		argsId, err := ArgsToCategoryId(r, args)
		if err != nil {
			return nil, 0, err
		}
		sb.Where(sb.In("bc.category_id", argsId...))
	}

	if val, ok := filters["author"]; ok && val != "" {
		var id string
		sb.JoinWithOption("LEFT", "author a", "b.author_id=a.author_id")
		value := strings.ToLower(val)
		err := r.db.QueryRow(`SELECT a.author_id FROM author as a
			WHERE (SELECT lower(name) FROM author WHERE author_id = a.author_id)=$1`, &value).Scan(&id)
		if err != nil {
			return nil, 0, err
		}
		sb.Where(sb.Equal("b.author_id", id))
	}

	sb.Where("(SELECT count(*) FROM author WHERE deleted_at IS NULL AND author_id=b.author_id) <> 0")
	sb.GroupBy("b.book_id", "b.name")
	sb.Limit(int(limit))
	sb.Offset(int(offset))

	query, args := sb.BuildWithFlavor(sqlbuilder.PostgreSQL)
	rows, err := r.db.Queryx(query, args...)
	if err != nil {
		return nil, 0, err
	}
	if err = rows.Err(); err != nil {
		return nil, 0, err
	}
	defer rows.Close() // nolint:errcheck
	var (
		books []*pb.BookResp
		count int64
	)

	for rows.Next() {
		var (
			book     pb.BookResp
			authorId string
			author   pb.Author
		)
		err = rows.Scan(
			&book.Id,
			&book.Name,
			&authorId,
			&book.CreatedAt,
			&book.UpdatedAt,
		)

		if err != nil {
			return nil, 0, err
		}

		author, err = r.GetAuthor(authorId)
		if err != nil {
			return nil, 0, err
		}

		book.Author = &author
		books = append(books, &book)
	}

	for _, book := range books {
		var categories []*pb.Category
		categories, err = GetBookCategory(r, book.Id)
		if err != nil {
			return nil, 0, err
		}
		book.Category = categories
	}

	sbc := sqlbuilder.NewSelectBuilder()
	sbc.Select("ROW_NUMBER() over (order by b.book_id) as r_number")
	sbc.From("book b")
	sbc.Where("b.deleted_at IS NULL")

	if val, ok := filters["category"]; ok && val != "" {
		var argsCategory, argsId []interface{}
		argsCategory = utils.StringSliceToInterfaceSlice(utils.ParseFilter(val))
		sbc.JoinWithOption("LEFT", "book_category bc", "b.book_id=bc.book_id")
		argsId, err = ArgsToCategoryId(r, argsCategory)
		if err != nil {
			return nil, 0, err
		}
		sbc.Where((sbc.In("bc.category_id", argsId...)))
	}

	if val, ok := filters["author"]; ok && val != "" {
		var id string
		value := strings.ToLower(val)
		err = r.db.QueryRow(`SELECT a.author_id FROM author as a
		WHERE (SELECT lower(name) FROM author WHERE author_id = a.author_id)=$1`, &value).Scan(&id)
		if err != nil {
			return nil, 0, err
		}
		sbc.Where(sbc.Equal("author_id", id))
	}

	sbc.GroupBy("b.book_id")
	sbc.OrderBy("r_number desc")
	sbc.Limit(1)
	query, args = sbc.BuildWithFlavor(sqlbuilder.PostgreSQL)
	err = r.db.QueryRow(query, args...).Scan(&count)
	if err != nil {
		return nil, 0, err
	}

	return books, count, nil
}

func (r *catalogRepo) UpdateBook(book pb.Book) (pb.BookResp, error) {
	result, err := r.db.Exec(`
		UPDATE book SET name=$1, author_id=$2, updated_at=current_timestamp
		WHERE book_id=$3 AND deleted_at IS NULL AND 
		(SELECT count(*) FROM author WHERE deleted_at IS NULL AND author_id=$4) <> 0`,
		book.Name, book.AuthorId, book.Id, book.AuthorId)
	if err != nil {
		return pb.BookResp{}, err
	}
	if i, _ := result.RowsAffected(); i == 0 {
		return pb.BookResp{}, sql.ErrNoRows
	}

	err = DeleteBookCategory(r, book.Id)
	if err != nil {
		return pb.BookResp{}, err
	}

	for _, category := range book.Categories {
		var bookCategory pb.BookCategory
		bookCategory.BookId = book.Id
		bookCategory.CategoryId = category

		err = CreateBookCategory(r, bookCategory)
		if err != nil {
			return pb.BookResp{}, err
		}
	}

	bookResp, err := r.GetBook(book.Id)
	if err != nil {
		return pb.BookResp{}, err
	}

	return bookResp, nil
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

// ArgsToCategoryId ...
func ArgsToCategoryId(r *catalogRepo, args []interface{}) ([]interface{}, error) {
	sbCategory := sqlbuilder.NewSelectBuilder()
	sbCategory.Select("category_id").From("category").Where(sbCategory.In("name", args...))
	query, argsCategory := sbCategory.BuildWithFlavor(sqlbuilder.PostgreSQL)

	rows, err := r.db.Query(query, argsCategory...)
	if err != nil {
		return nil, err
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	defer rows.Close() // nolint:errcheck

	var values string
	var i int = 0
	for rows.Next() {

		if i != 0 {
			values += ","
		}
		i++
		var value string

		err = rows.Scan(
			&value,
		)
		if err != nil {
			return nil, err
		}

		values += value
	}

	argsId := utils.StringSliceToInterfaceSlice(utils.ParseFilter(values))

	return argsId, nil
}

package main

import (
	"fmt"
	"strings"

	"github.com/huandu/go-sqlbuilder"
)

func main() {
	sb := sqlbuilder.NewSelectBuilder()

	filters := make(map[string]string)

	filters["categories"] = "1,2"
	filters["author"] = "1,2"

	sb.Select("b.book_id", "b.name")
	sb.From("book b")
	sb.Where("b.deleted_at IS NULL")

	if val, ok := filters["categories"]; ok && val != "" {
		args := StringSliceToInterfaceSlice(ParseFilter(val))
		sb.JoinWithOption("LEFT", "book_category bc", "b.book_id=bc.book_id")
		sb.Where((sb.In("bc.category_id", args...)))
	}

	if val, ok := filters["author"]; ok && val != "" {
		sb.Where(sb.Equal("author_id", val))
	}

	sb.GroupBy("b.book_id", "b.name")
	sb.Limit(10)
	sb.Offset(8)

	query, args := sb.BuildWithFlavor(sqlbuilder.PostgreSQL)

	fmt.Println(query)
	fmt.Println(args)

	// SELECT
	// b.book_id,
	// b.name
	// FROM book b
	// LEFT JOIN book_category bc ON b.book_id=bc.book_id
	// WHERE bc.category_id IN ('39af40c6-7e8a-4495-b16e-a1c58da4cba3', '3bee22f5-96d5-4d6b-bd16-89454cb0cd18', '9272dc0a-f27e-4fb9-8448-100851c1b374')
	// AND author_id = 'f9ad2e1f-7511-40d4-a560-a0a6de712671'
	// GROUP by b.book_id, b.name
	// LIMIT 10 OFFSET 0;

}

func ParseFilter(s string) []string {
	return strings.Split(s, ",")
}

func StringSliceToInterfaceSlice(ss []string) []interface{} {
	is := make([]interface{}, 0, len(ss))
	for _, s := range ss {
		is = append(is, s)
	}

	return is
}

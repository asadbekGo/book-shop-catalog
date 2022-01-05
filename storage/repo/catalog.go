package repo

import (
	pb "github.com/asadbekGo/book-shop-catalog/genproto"
)

// CatalogStorageI ...
type CatalogStorageI interface {
	CreateAuthor(pb.Author) (pb.Author, error)
}

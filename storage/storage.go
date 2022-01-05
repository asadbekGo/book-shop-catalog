package storage

import (
	"github.com/jmoiron/sqlx"

	"github.com/asadbekGo/book-shop-catalog/storage/postgres"
	"github.com/asadbekGo/book-shop-catalog/storage/repo"
)

// IStorage ...
type IStorage interface {
	Catalog() repo.CatalogStorageI
}

type storagePg struct {
	db          *sqlx.DB
	catalogRepo repo.CatalogStorageI
}

//  NewStoragePg ...
func NewStoragePg(db *sqlx.DB) *storagePg {
	return &storagePg{
		db:          db,
		catalogRepo: postgres.NewCatalogRepo(db),
	}
}

func (s storagePg) Catalog() repo.CatalogStorageI {
	return s.catalogRepo
}

package repository

import (
	"azmi17/go-rest-api/model/entity"
	"context"
	"database/sql"
)

type CategoryRepository interface {
	// in params (context|Db transaction|data:category) return value => entity.Category
	Save(ctx context.Context, tx *sql.Tx, category entity.Category) entity.Category
	Update(ctx context.Context, tx *sql.Tx, category entity.Category) entity.Category
	Delete(ctx context.Context, tx *sql.Tx, category entity.Category)
	FindById(ctx context.Context, tx *sql.Tx, categoryId int) (entity.Category, error)
	FindAll(ctx context.Context, tx *sql.Tx) []entity.Category
	/*
	 ^
	 FindAll() return []entity.Category
	 dibuat slice karena akan return semua data (artinya banyak object)
	*/
}

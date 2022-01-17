package repository

import (
	"azmi17/go-rest-api/helper"
	"azmi17/go-rest-api/model/entity"
	"context"
	"database/sql"
	"errors"
)

type CategoryRepositoryImpl struct {
	// masih belum faham, kenapa struct ini empty (?)
}

// Buat objek CategoryRepositoryImpl (constructror)
func NewCategoryRepository() CategoryRepository {
	return &CategoryRepositoryImpl{}
}

/*
 DIBAWAH ADALAH SEMUA FUNCTION KONTRAK DARI REPOSITORY/CATEGORY
*/
func (repository *CategoryRepositoryImpl) Save(ctx context.Context, tx *sql.Tx, category entity.Category) entity.Category {

	// SQL Statement
	SQL := "INSERT INTO categories(name) VALUES (?)" // tidak passing Id, karena set auto-increment

	// Exec
	result, err := tx.ExecContext(ctx, SQL, category.Name)
	helper.PanicIfError(err)

	// Get LastInsertId
	id, err := result.LastInsertId()
	helper.PanicIfError(err)
	category.Id = int(id)

	// Return data
	return category

}

func (repository *CategoryRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, category entity.Category) entity.Category {

	// SQL Statement
	SQL := "UPDATE categories SET name = ? WHERE id = ?"

	// Exec
	_, err := tx.ExecContext(ctx, SQL, category.Name, category.Id)
	helper.PanicIfError(err)

	// Return data
	return category
}

func (repository *CategoryRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, category entity.Category) {
	// SQL Statement
	SQL := "DELETE FROM categories WHERE id = ?"

	// Exec
	_, err := tx.ExecContext(ctx, SQL, category.Id)
	helper.PanicIfError(err)
}

func (repository *CategoryRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, categoryId int) (entity.Category, error) {

	// SQL Statement
	SQL := "SELECT id, name FROM categories where ID = ?"

	// Query
	category := entity.Category{} // buat empty object categories
	rows, err := tx.QueryContext(ctx, SQL, categoryId)
	helper.PanicIfError(err)
	defer rows.Close()

	// Get Data setelah Query
	if rows.Next() {
		rows.Scan(&category.Id, &category.Name)
		return category, nil
	} else {
		return category, errors.New("category is not found")
	}
}

func (repository *CategoryRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx) []entity.Category {

	// SQL Statement
	SQL := "SELECT id, name FROM categories"

	// Query
	rows, err := tx.QueryContext(ctx, SQL)
	helper.PanicIfError(err)
	defer rows.Close()

	// Get All Data Result
	var categories []entity.Category
	for rows.Next() {
		category := entity.Category{} // buat empty object categories
		rows.Scan(&category.Id, &category.Name)
		helper.PanicIfError(err)
		categories = append(categories, category)
	}

	return categories
}

package service

import (
	"azmi17/go-rest-api/exception"
	"azmi17/go-rest-api/helper"
	"azmi17/go-rest-api/model/entity"
	"azmi17/go-rest-api/model/web"
	"azmi17/go-rest-api/repository"
	"context"
	"database/sql"

	"github.com/go-playground/validator"
)

type CategoryServiceImpl struct {
	CategoryRepository repository.CategoryRepository // embedd repository karena manipulasi datanya menggunakan repository (repository Ini type-nya interface{} jd tidak butuh di set sebagai pointer)
	DB                 *sql.DB                       // set sebagai pointer karena struct type
	Validate           *validator.Validate
}

// Buat objek NewCategoryService (constructror)
func NewCategoryService(categoryRepository repository.CategoryRepository, db *sql.DB, validate *validator.Validate) CategoryService {
	return &CategoryServiceImpl{
		CategoryRepository: categoryRepository,
		DB:                 db,
		Validate:           validate,
	}
}

//IMPLEMENTASI KONTRAK DARI CATEGORY SERVICE
func (service *CategoryServiceImpl) Create(ctx context.Context, request web.CategoryCreateRequest) web.CategoryResponse {

	// Buat validasi (menggunakan third-party:validator)
	err := service.Validate.Struct(request)
	helper.PanicIfError(err)

	// Begin Transaction Statement
	tx, err := service.DB.Begin()
	helper.PanicIfError(err)

	// Handle jika terjadi error maka rollback (dibuat dengan helper)
	defer helper.CommitOrRollBack(tx)

	// Request
	category := entity.Category{
		Name: request.Name,
	}
	category = service.CategoryRepository.Save(ctx, tx, category) // memanggil service, dan service memanggil repository

	// Convert dari Category ke CategoryResponse
	return helper.ToCategoryResponse(category)
}

func (service *CategoryServiceImpl) Update(ctx context.Context, request web.CategoryUpdateRequest) web.CategoryResponse {

	// Buat validasi (menggunakan third-party:validator)
	err := service.Validate.Struct(request)
	helper.PanicIfError(err)

	// Begin Transaction Statement
	tx, err := service.DB.Begin()
	helper.PanicIfError(err)

	// Handle jika terjadi error maka rollback (dibuat dengan helper)
	defer helper.CommitOrRollBack(tx)

	// Validasi (apakah data ada/tidak)
	category, err := service.CategoryRepository.FindById(ctx, tx, request.Id)

	// throw exception jika data tidak ada
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	// Request setelah validasi
	category.Name = request.Name

	category = service.CategoryRepository.Update(ctx, tx, category) // memanggil service, dan service memanggil repository

	// Convert dari Category ke CategoryResponse
	return helper.ToCategoryResponse(category)
}

func (service *CategoryServiceImpl) Delete(ctx context.Context, categoryId int) {
	// Begin Transaction Statement
	tx, err := service.DB.Begin()
	helper.PanicIfError(err)

	// Handle jika terjadi error maka rollback (dibuat dengan helper)
	defer helper.CommitOrRollBack(tx)

	// Validasi (apakah data ada/tidak)
	category, err := service.CategoryRepository.FindById(ctx, tx, categoryId)

	// throw exception jika data tidak ada
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	// Proses Delete
	service.CategoryRepository.Delete(ctx, tx, category)

}

func (service *CategoryServiceImpl) FindById(ctx context.Context, categoryId int) web.CategoryResponse {
	// Begin Transaction Statement
	tx, err := service.DB.Begin()
	helper.PanicIfError(err)

	// Handle jika terjadi error maka rollback (dibuat dengan helper)
	defer helper.CommitOrRollBack(tx)

	// Get FindById
	category, err := service.CategoryRepository.FindById(ctx, tx, categoryId)

	// throw exception jika data tidak ada
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	// Convert dari Category ke CategoryResponse
	return helper.ToCategoryResponse(category)
}

func (service *CategoryServiceImpl) FindAll(ctx context.Context) []web.CategoryResponse {
	// Begin Transaction Statement
	tx, err := service.DB.Begin()
	helper.PanicIfError(err)

	// Handle jika terjadi error maka rollback (dibuat dengan helper)
	defer helper.CommitOrRollBack(tx)

	categories := service.CategoryRepository.FindAll(ctx, tx)

	// return semua data cageories
	return helper.ToCategoryResponses(categories)
}

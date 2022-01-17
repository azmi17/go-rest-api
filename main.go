package main

import (
	"azmi17/go-rest-api/app"
	"azmi17/go-rest-api/controller"
	"azmi17/go-rest-api/helper"
	"azmi17/go-rest-api/middleware"
	"azmi17/go-rest-api/repository"
	"azmi17/go-rest-api/service"
	"net/http"

	_ "github.com/go-sql-driver/mysql"

	"github.com/go-playground/validator"
)

func main() {

	// Inisiasi objek db
	db := app.NewDB()

	// Inisiasi validate
	validate := validator.New()

	// Inisiasi masing-masing layer
	categoryRepository := repository.NewCategoryRepository()
	categoryService := service.NewCategoryService(categoryRepository, db, validate)
	categoryController := controller.NewCategoryController(categoryService)

	// Inisiasi httpRouter & errHandler
	router := app.NewRouter(categoryController)

	// Buat Server menggunakan pkg http
	server := http.Server{
		Addr: "localhost:1717",
		// Handler: router,
		Handler: middleware.NewAuthMiddleware(router), // menggunakan middleware utk kebutuhan Auth pada setiap Action Request
	}

	err := server.ListenAndServe()
	helper.PanicIfError(err)

}

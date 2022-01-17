package app

import (
	"azmi17/go-rest-api/controller"
	"azmi17/go-rest-api/exception"

	"github.com/julienschmidt/httprouter"
)

func NewRouter(categoryController controller.CategoryController) *httprouter.Router {
	router := httprouter.New()

	// Mapping Endpoint sesuai path API SPEC
	router.GET("/api/categories", categoryController.FindAll)
	router.GET("/api/categories/:categoryId", categoryController.FindById)
	router.POST("/api/categories", categoryController.Create)
	router.PUT("/api/categories/:categoryId", categoryController.Update)
	router.DELETE("/api/categories/:categoryId", categoryController.Delete)

	// Handle jika terjadi error
	router.PanicHandler = exception.ErrorHandler

	return router
}

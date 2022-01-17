package helper

import (
	"azmi17/go-rest-api/model/entity"
	"azmi17/go-rest-api/model/web"
)

// byId
func ToCategoryResponse(category entity.Category) web.CategoryResponse {

	return web.CategoryResponse{
		Id:   category.Id,
		Name: category.Name,
	}
}

// getAll
func ToCategoryResponses(categories []entity.Category) []web.CategoryResponse {

	var CategoryResponses []web.CategoryResponse
	for _, category := range categories {
		CategoryResponses = append(CategoryResponses, ToCategoryResponse(category))
	}

	return CategoryResponses
}

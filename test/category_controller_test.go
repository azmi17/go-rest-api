package test

import (
	"azmi17/go-rest-api/app"
	"azmi17/go-rest-api/controller"
	"azmi17/go-rest-api/helper"
	"azmi17/go-rest-api/middleware"
	"azmi17/go-rest-api/model/entity"
	"azmi17/go-rest-api/repository"
	"azmi17/go-rest-api/service"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/go-playground/validator"
	"github.com/stretchr/testify/assert"
)

/*
 - Ini adalah Pengujian Integration Test,
   melakukan pengujian terhadap API Endpoint yang sudah dibuat.
*/

// DUMMY DB: go_rest_api_test
func setupTestDB() *sql.DB {
	// Db Config
	db, err := sql.Open("mysql", "root@tcp(localhost:3317)/go_rest_api_test")
	helper.PanicIfError(err)

	// Db Pooling
	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(20)
	db.SetConnMaxLifetime(60 * time.Minute)
	db.SetConnMaxIdleTime(10 * time.Minute)

	return db
}

// SETUP ROUTER
func setupRouter(db *sql.DB) http.Handler {

	// Inisiasi validate
	validate := validator.New()

	// Inisiasi masing-masing layer
	categoryRepository := repository.NewCategoryRepository()
	categoryService := service.NewCategoryService(categoryRepository, db, validate)
	categoryController := controller.NewCategoryController(categoryService)

	// Inisiasi httpRouter & errHandler
	router := app.NewRouter(categoryController)

	// middleware
	return middleware.NewAuthMiddleware(router)
}

// TRUNCATE FUNCTION
func truncateCategory(db *sql.DB) {
	db.Exec("TRUNCATE categories")
}

// SUCCESS SCENARIO (CREATE CATEGORY)
func TestCreateCategorySuccess(t *testing.T) {

	// Db init
	db := setupTestDB()

	// sebelum test, akan truncate dulu baru insert
	truncateCategory(db)

	// Router init
	router := setupRouter(db)

	// Buat data
	requestBody := strings.NewReader(`{"name": "Gadget"}`)

	// Buat request dengan HTTP Test
	request := httptest.NewRequest(http.MethodPost, "http://localhost:1717/api/categories", requestBody)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-API-Key", "p0l1moRphY5m")

	// Buat Recorder
	recorder := httptest.NewRecorder()

	// ServeHTTP
	router.ServeHTTP(recorder, request)

	// Skenario Test: Expected 200:status-code after request
	response := recorder.Result()
	assert.Equal(t, 200, response.StatusCode)

	// read body with IO
	body, _ := io.ReadAll(response.Body)

	// Decode JSON
	var responBody map[string]interface{}
	json.Unmarshal(body, &responBody)
	// fmt.Println(responBody) // => map[code:200 data:map[id:3 name:Computer] status:OK]

	// Skenario Test, expected: 200, Gadget, OK
	assert.Equal(t, 200, int(responBody["code"].(float64)))
	assert.Equal(t, "Gadget", responBody["data"].(map[string]interface{})["name"])
	assert.Equal(t, "OK", responBody["status"])

}

// FAIL SCENARIO (CREATE CATEGORY)
func TestCreateCategoryFail(t *testing.T) {
	// Db init
	db := setupTestDB()

	// sebelum test, akan truncate dulu baru insert
	truncateCategory(db)

	// Router init
	router := setupRouter(db)

	// Buat data
	requestBody := strings.NewReader(`{"name": ""}`)

	// Buat request dengan HTTP Test
	request := httptest.NewRequest(http.MethodPost, "http://localhost:1717/api/categories", requestBody)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-API-Key", "p0l1moRphY5m")

	// Buat Recorder
	recorder := httptest.NewRecorder()

	// ServeHTTP
	router.ServeHTTP(recorder, request)

	// Skenario Test: Expected 200:status-code after request
	response := recorder.Result()
	assert.Equal(t, 400, response.StatusCode)

	// read body with IO
	body, _ := io.ReadAll(response.Body)

	// Decode JSON
	var responBody map[string]interface{}
	json.Unmarshal(body, &responBody)
	// fmt.Println(responBody) // => map[code:200 data:map[id:3 name:Computer] status:OK]

	// Skenario Test, expected: 400, BAD REQUEST
	assert.Equal(t, 400, int(responBody["code"].(float64)))
	assert.Equal(t, "BAD REQUEST", responBody["status"])
}

// SUCCESS SCENARIO (UPDATE CATEGORY)
func TestUpdateCategorySuccess(t *testing.T) {

	// Db init
	db := setupTestDB()

	// sebelum test, akan truncate dulu baru insert
	truncateCategory(db)

	// Buat data dulu sebelum update, memanfaatkan repository
	tx, _ := db.Begin()
	categoryRepository := repository.NewCategoryRepository()
	category := categoryRepository.Save(context.Background(), tx, entity.Category{
		Name: "Computer",
	})
	tx.Commit()

	// Router init
	router := setupRouter(db)

	// Buat data
	requestBody := strings.NewReader(`{"name": "Gadget"}`)

	// Buat request dengan HTTP Test
	request := httptest.NewRequest(http.MethodPut, "http://localhost:1717/api/categories/"+strconv.Itoa(category.Id), requestBody)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-API-Key", "p0l1moRphY5m")

	// Buat Recorder
	recorder := httptest.NewRecorder()

	// ServeHTTP
	router.ServeHTTP(recorder, request)

	// Skenario Test: Expected 200:status-code after request
	response := recorder.Result()
	assert.Equal(t, 200, response.StatusCode)

	// read body with IO
	body, _ := io.ReadAll(response.Body)

	// Decode JSON
	var responBody map[string]interface{}
	json.Unmarshal(body, &responBody)
	// fmt.Println(responBody) // => map[code:200 data:map[id:3 name:Gadget] status:OK]

	// Skenario Test, expected: 200, Data: categoryId categoryName, OK
	assert.Equal(t, 200, int(responBody["code"].(float64)))
	assert.Equal(t, category.Id, int(responBody["data"].(map[string]interface{})["id"].(float64)))
	assert.Equal(t, "Gadget", responBody["data"].(map[string]interface{})["name"])
	assert.Equal(t, "OK", responBody["status"])

}

// FAIL SCENARIO (UPDATE CATEGORY)
func TestUpdateCategoryFail(t *testing.T) {

	// Db init
	db := setupTestDB()

	// sebelum test, akan truncate dulu baru insert
	truncateCategory(db)

	// Buat data dulu sebelum update, memanfaatkan repository
	tx, _ := db.Begin()
	categoryRepository := repository.NewCategoryRepository()
	category := categoryRepository.Save(context.Background(), tx, entity.Category{
		Name: "Computer",
	})
	tx.Commit()

	// Router init
	router := setupRouter(db)

	// Buat data
	requestBody := strings.NewReader(`{"name": ""}`)

	// Buat request dengan HTTP Test
	request := httptest.NewRequest(http.MethodPut, "http://localhost:1717/api/categories/"+strconv.Itoa(category.Id), requestBody)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-API-Key", "p0l1moRphY5m")

	// Buat Recorder
	recorder := httptest.NewRecorder()

	// ServeHTTP
	router.ServeHTTP(recorder, request)

	// Skenario Test: Expected 400:status-code after request
	response := recorder.Result()
	assert.Equal(t, 400, response.StatusCode)

	// read body with IO
	body, _ := io.ReadAll(response.Body)

	// Decode JSON
	var responBody map[string]interface{}
	json.Unmarshal(body, &responBody)
	// fmt.Println(responBody)

	// Skenario Test, expected: 400, BAD REQUEST
	assert.Equal(t, 400, int(responBody["code"].(float64)))
	assert.Equal(t, "BAD REQUEST", responBody["status"])

}

// SUCCESS SCENARIO (GET CATEGORY)
func TestGetCategorySuccess(t *testing.T) {

	// Db init
	db := setupTestDB()

	// sebelum test, akan truncate dulu baru insert
	truncateCategory(db)

	// Buat data dulu sebelum get, memanfaatkan repository
	tx, _ := db.Begin()
	categoryRepository := repository.NewCategoryRepository()
	category := categoryRepository.Save(context.Background(), tx, entity.Category{
		Name: "Fashion",
	})
	tx.Commit()

	// Router init
	router := setupRouter(db)

	// Buat request dengan HTTP Test
	request := httptest.NewRequest(http.MethodGet, "http://localhost:1717/api/categories/"+strconv.Itoa(category.Id), nil)
	request.Header.Add("X-API-Key", "p0l1moRphY5m")

	// Buat Recorder
	recorder := httptest.NewRecorder()

	// ServeHTTP
	router.ServeHTTP(recorder, request)

	// Skenario Test: Expected 200:status-code after request
	response := recorder.Result()
	assert.Equal(t, 200, response.StatusCode)

	// read body with IO
	body, _ := io.ReadAll(response.Body)

	// Decode JSON
	var responBody map[string]interface{}
	json.Unmarshal(body, &responBody)
	fmt.Println(responBody)

	// Skenario Test, expected: 200, OK, category-Name
	assert.Equal(t, 200, int(responBody["code"].(float64)))
	assert.Equal(t, category.Id, int(responBody["data"].(map[string]interface{})["id"].(float64)))
	assert.Equal(t, category.Name, responBody["data"].(map[string]interface{})["name"])

	assert.Equal(t, "OK", responBody["status"])
}

// FAIL SCENARIO (GET CATEGORY)
func TestGetCategoryFail(t *testing.T) {

	// Db init
	db := setupTestDB()

	// sebelum test, akan truncate dulu baru insert
	truncateCategory(db)

	// Router init
	router := setupRouter(db)

	// Buat request dengan HTTP Test
	request := httptest.NewRequest(http.MethodGet, "http://localhost:1717/api/categories/531", nil)
	request.Header.Add("X-API-Key", "p0l1moRphY5m")

	// Buat Recorder
	recorder := httptest.NewRecorder()

	// ServeHTTP
	router.ServeHTTP(recorder, request)

	// Skenario Test: Expected 404:status-code after request
	response := recorder.Result()
	assert.Equal(t, 404, response.StatusCode)

	// read body with IO
	body, _ := io.ReadAll(response.Body)

	// Decode JSON
	var responBody map[string]interface{}
	json.Unmarshal(body, &responBody)
	fmt.Println(responBody)

	// Skenario Test, expected: 404, NOT FOUND
	assert.Equal(t, 404, int(responBody["code"].(float64)))
	assert.Equal(t, "NOT FOUND", responBody["status"])
}

// SUCCESS SCENARIO (DELETE CATEGORY)
func TestDeleteCategorySuccess(t *testing.T) {

	// Db init
	db := setupTestDB()

	// sebelum test, akan truncate dulu baru insert
	truncateCategory(db)

	// Buat data dulu sebelum update, memanfaatkan repository
	tx, _ := db.Begin()
	categoryRepository := repository.NewCategoryRepository()
	category := categoryRepository.Save(context.Background(), tx, entity.Category{
		Name: "Gadget",
	})
	tx.Commit()

	// Router init
	router := setupRouter(db)

	// Buat request dengan HTTP Test
	request := httptest.NewRequest(http.MethodDelete, "http://localhost:1717/api/categories/"+strconv.Itoa(category.Id), nil)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-API-Key", "p0l1moRphY5m")

	// Buat Recorder
	recorder := httptest.NewRecorder()

	// ServeHTTP
	router.ServeHTTP(recorder, request)

	// Skenario Test: Expected 200:status-code after request
	response := recorder.Result()
	assert.Equal(t, 200, response.StatusCode)

	// read body with IO
	body, _ := io.ReadAll(response.Body)

	// Decode JSON
	var responBody map[string]interface{}
	json.Unmarshal(body, &responBody)
	// fmt.Println(responBody)

	// Skenario Test, expected: 200, OK
	assert.Equal(t, 200, int(responBody["code"].(float64)))
	assert.Equal(t, "OK", responBody["status"])
}

// FAIL SCENARIO (DELETE CATEGORY)
func TestDeleteCategoryFail(t *testing.T) {

	// Db init
	db := setupTestDB()

	// sebelum test, akan truncate dulu baru insert
	truncateCategory(db)

	// Router init
	router := setupRouter(db)

	// Buat request dengan HTTP Test
	request := httptest.NewRequest(http.MethodDelete, "http://localhost:1717/api/categories/124", nil)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-API-Key", "p0l1moRphY5m")

	// Buat Recorder
	recorder := httptest.NewRecorder()

	// ServeHTTP
	router.ServeHTTP(recorder, request)

	// Skenario Test: Expected 404:status-code after request
	response := recorder.Result()
	assert.Equal(t, 404, response.StatusCode)

	// read body with IO
	body, _ := io.ReadAll(response.Body)

	// Decode JSON
	var responBody map[string]interface{}
	json.Unmarshal(body, &responBody)
	// fmt.Println(responBody)

	// Skenario Test, expected: 404, NOT FOUND
	assert.Equal(t, 404, int(responBody["code"].(float64)))
	assert.Equal(t, "NOT FOUND", responBody["status"])
}

// SUCCESS SCENARIO (GET CATEGORIES)
func TestListCategorySuccess(t *testing.T) {

	// Db init
	db := setupTestDB()

	// sebelum test, akan truncate dulu baru insert
	truncateCategory(db)

	// Buat data dulu sebelum update, memanfaatkan repository
	tx, _ := db.Begin()
	categoryRepository := repository.NewCategoryRepository()
	category1 := categoryRepository.Save(context.Background(), tx, entity.Category{
		Name: "Computer",
	})
	category2 := categoryRepository.Save(context.Background(), tx, entity.Category{
		Name: "Fashion",
	})
	tx.Commit()

	// Router init
	router := setupRouter(db)

	// Buat request dengan HTTP Test
	request := httptest.NewRequest(http.MethodGet, "http://localhost:1717/api/categories", nil)
	request.Header.Add("X-API-Key", "p0l1moRphY5m")

	// Buat Recorder
	recorder := httptest.NewRecorder()

	// ServeHTTP
	router.ServeHTTP(recorder, request)

	// Skenario Test: Expected 200:status-code after request
	response := recorder.Result()
	assert.Equal(t, 200, response.StatusCode)

	// read body with IO
	body, _ := io.ReadAll(response.Body)

	// Decode JSON
	var responBody map[string]interface{}
	json.Unmarshal(body, &responBody)

	// Skenario Test, expected: 200, OK
	assert.Equal(t, 200, int(responBody["code"].(float64)))
	assert.Equal(t, "OK", responBody["status"])

	// debug
	fmt.Println(responBody)

	//get list categories (bentuknya slice)
	var categories = responBody["data"].([]interface{})

	categoryResponse1 := categories[0].(map[string]interface{})
	categoryResponse2 := categories[1].(map[string]interface{})

	assert.Equal(t, category1.Id, int(categoryResponse1["id"].(float64)))
	assert.Equal(t, category1.Name, categoryResponse1["name"])

	assert.Equal(t, category2.Id, int(categoryResponse2["id"].(float64)))
	assert.Equal(t, category2.Name, categoryResponse2["name"])

}

// FAIL AUTHORIZATION (SECURITY ENDPOINT)
func TestUnauthorized(t *testing.T) {

	// Db init
	db := setupTestDB()

	// sebelum test, akan truncate dulu baru insert
	truncateCategory(db)

	// Router init
	router := setupRouter(db)

	// Buat request dengan HTTP Test
	request := httptest.NewRequest(http.MethodGet, "http://localhost:1717/api/categories", nil)
	request.Header.Add("X-API-Key", "XXX")

	// Buat Recorder
	recorder := httptest.NewRecorder()

	// ServeHTTP
	router.ServeHTTP(recorder, request)

	// Skenario Test: Expected 401:status-code after request
	response := recorder.Result()
	assert.Equal(t, 401, response.StatusCode)

	// read body with IO
	body, _ := io.ReadAll(response.Body)

	// Decode JSON
	var responBody map[string]interface{}
	json.Unmarshal(body, &responBody)

	// Skenario Test, expected: 401, OK
	assert.Equal(t, 401, int(responBody["code"].(float64)))
	assert.Equal(t, "UNAUTHORIZED", responBody["status"])
}

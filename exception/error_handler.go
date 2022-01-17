package exception

import (
	"azmi17/go-rest-api/helper"
	"azmi17/go-rest-api/model/web"
	"net/http"

	"github.com/go-playground/validator"
)

func ErrorHandler(writer http.ResponseWriter, request *http.Request, err interface{}) {

	// Call notFounderror
	if notFoundError(writer, request, err) {
		return
	}

	// Call validationErrors
	if validationErrors(writer, request, err) {
		return
	}

	// Call internalServerError
	internalServerError(writer, request, err)
}

// VALIDATION ERROR HANDLER
func validationErrors(writer http.ResponseWriter, request *http.Request, err interface{}) bool {

	exception, ok := err.(validator.ValidationErrors)
	if ok {

		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusBadRequest)

		webResponse := web.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "BAD REQUEST",
			Data:   exception.Error(),
		}

		helper.WriteToResponseBody(writer, webResponse)

		return true
	} else {
		return false
	}
}

// NOT FOUND ERROR HANDLER
func notFoundError(writer http.ResponseWriter, request *http.Request, err interface{}) bool {

	exception, ok := err.(NotFoundError)

	if ok { // jika masuk kesini artinya data bisa di konversi dan return TRUE
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusNotFound)

		webResponse := web.WebResponse{
			Code:   http.StatusNotFound,
			Status: "NOT FOUND",
			Data:   exception.Error,
		}

		helper.WriteToResponseBody(writer, webResponse)
		return true

	} else { // jika masuk kesini tidak ada error, dan langsung return FALSE
		return false
	}
}

// INTERNAL SERVER ERROR HANDLER
func internalServerError(writer http.ResponseWriter, request *http.Request, err interface{}) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusInternalServerError)

	webResponse := web.WebResponse{
		Code:   http.StatusInternalServerError,
		Status: "INTERNAL SERVER ERROR",
		Data:   err,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

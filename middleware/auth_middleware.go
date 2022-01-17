package middleware

import (
	"azmi17/go-rest-api/helper"
	"azmi17/go-rest-api/model/web"
	"net/http"
)

type AuthMiddleware struct {
	Handler http.Handler
}

// Buat objek AuthMiddleware (constructror)
func NewAuthMiddleware(handler http.Handler) *AuthMiddleware {
	return &AuthMiddleware{Handler: handler}
}

func (middleware *AuthMiddleware) ServeHTTP(writer http.ResponseWriter, request *http.Request) {

	authKey := "p0l1moRphY5m"
	if authKey == request.Header.Get("X-API-Key") {

		// jika auth valid maka request akan diterukan ke handler selanjutnya => OK
		middleware.Handler.ServeHTTP(writer, request)
	} else {

		// jika auth tidak valid maka set unauthorized => ERROR
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusUnauthorized)

		webResponse := web.WebResponse{
			Code:   http.StatusUnauthorized,
			Status: "UNAUTHORIZED",
		}

		helper.WriteToResponseBody(writer, webResponse)
	}
}

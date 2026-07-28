package response

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// Response represents the response structure
type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// Success returns with status code 200
func Success(w http.ResponseWriter, data any) {
	httpx.WriteJson(w, http.StatusOK, Response{
		Code:    http.StatusOK,
		Message: "Success",
		Data:    data,
	})
}

// Error returns with status code 200
func Error(w http.ResponseWriter, code int, message string) {
	httpx.WriteJson(w, http.StatusOK, Response{
		Code:    code,
		Message: message,
	})
}

// BadRequest returns with status code 400
func BadRequest(w http.ResponseWriter, message string) {
	Error(w, http.StatusBadRequest, message)
}

func BadRequestWithCode(w http.ResponseWriter, code int, message string) {
	httpx.WriteJson(w, code, Response{
		Code:    code,
		Message: message,
	})
}

// Unauthorized returns with status code 401
func Unauthorized(w http.ResponseWriter, message string) {
	Error(w, http.StatusUnauthorized, message)
}

// Forbidden returns with status code 403
func Forbidden(w http.ResponseWriter, message string) {
	Error(w, http.StatusForbidden, message)
}

// NotFound returns with status code 404
func NotFound(w http.ResponseWriter, message string) {
	Error(w, http.StatusNotFound, message)
}

// InternalServerError returns with status code 500
func InternalServerError(w http.ResponseWriter, message string) {
	Error(w, http.StatusInternalServerError, message)
}

package middleware

import (
	"log"
	"net/http"
	"runtime/debug"

	"zeebee-backend/common/response"
)

type RecoverMiddleware struct{}

func NewRecoverMiddleware() *RecoverMiddleware {
	return &RecoverMiddleware{}
}

func (m *RecoverMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				// 记录错误日志
				log.Printf("Panic recovered: %v\n%s", err, debug.Stack())

				// 返回500错误
				response.InternalServerError(w, "Internal server error")
			}
		}()

		next(w, r)
	}
}

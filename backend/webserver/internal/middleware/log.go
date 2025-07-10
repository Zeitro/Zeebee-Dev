package middleware

import (
	"log"
	"net/http"
	"time"
)

type LogMiddleware struct{}

func NewLogMiddleware() *LogMiddleware {
	return &LogMiddleware{}
}

func (m *LogMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// 记录请求开始
		log.Printf("Request started: %s %s", r.Method, r.URL.Path)

		// 调用下一个处理器
		next(w, r)

		// 记录请求结束
		duration := time.Since(start)
		log.Printf("Request completed: %s %s - %v", r.Method, r.URL.Path, duration)
	}
}

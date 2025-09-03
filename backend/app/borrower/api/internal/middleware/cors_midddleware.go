package middleware

import "net/http"

func CorsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*") // all origins are allowed
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS PATCH")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Accept")

		// 处理预检请求
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)

	}
}

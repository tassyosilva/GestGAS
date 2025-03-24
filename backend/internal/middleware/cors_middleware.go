package middleware

import (
	"net/http"
)

// CorsMiddleware configura o middleware CORS para permitir requisições cross-origin
func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Obtenha a origem da requisição
		origin := r.Header.Get("Origin")
		
		// LISTA DE ORIGENS PERMITIDAS
		allowedOrigins := map[string]bool{
			"http://localhost:3000": true,
		}
		
		// Verifique se a origem está na lista de permitidas
		if allowedOrigins[origin] {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}
		
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, PATCH")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Authorization")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
package middleware

import (
	"context"
	"database/sql"
	"net/http"
	"strings"
)

// UsuarioKey é a chave utilizada para armazenar o ID do usuário no contexto
type UsuarioKey string

// AuthMiddleware verifica se a requisição possui um token válido
func AuthMiddleware(db *sql.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Obter o token do cabeçalho Authorization
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Token de autenticação não fornecido", http.StatusUnauthorized)
				return
			}

			// Verificar formato do cabeçalho (Bearer TOKEN)
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, "Formato de token inválido", http.StatusUnauthorized)
				return
			}

			token := parts[1]

			// Em um sistema real, validaríamos o token com uma lógica mais robusta
			// Por simplicidade, apenas verificamos se o token existe em nossa "tabela de tokens"
			var userID int
			err := db.QueryRow("SELECT usuario_id FROM tokens WHERE token = $1 AND expiracao > NOW()", token).Scan(&userID)
			if err != nil {
				if err == sql.ErrNoRows {
					http.Error(w, "Token inválido ou expirado", http.StatusUnauthorized)
					return
				}
				http.Error(w, "Erro ao validar token", http.StatusInternalServerError)
				return
			}

			// Adicionar o ID do usuário ao contexto da requisição
			ctx := context.WithValue(r.Context(), UsuarioKey("usuarioID"), userID)
			
			// Chamar o próximo handler com o contexto atualizado
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// ObterUsuarioID extrai o ID do usuário do contexto da requisição
func ObterUsuarioID(r *http.Request) (int, bool) {
	userID, ok := r.Context().Value(UsuarioKey("usuarioID")).(int)
	return userID, ok
}

// VerificarPerfil verifica se o usuário possui um determinado perfil
func VerificarPerfil(db *sql.DB, userID int, perfilRequerido string) bool {
	var perfil string
	err := db.QueryRow("SELECT perfil FROM usuarios WHERE id = $1", userID).Scan(&perfil)
	if err != nil {
		return false
	}
	
	// Se o perfil requerido for "admin", apenas admin pode acessar
	// Se for "gerente", admin e gerente podem acessar
	// Se for "atendente", admin, gerente e atendente podem acessar
	// Se for "entregador", qualquer um pode acessar
	
	switch perfilRequerido {
	case "admin":
		return perfil == "admin"
	case "gerente":
		return perfil == "admin" || perfil == "gerente"
	case "atendente":
		return perfil == "admin" || perfil == "gerente" || perfil == "atendente"
	case "entregador":
		return true // Todos os perfis podem acessar
	default:
		return false
	}
}
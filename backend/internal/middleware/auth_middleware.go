package middleware

import (
	"context"
	"database/sql"
	"net/http"
	"strings"

	"github.com/tassyosilva/GestGAS/internal/auth"
)

// UsuarioKey é a chave utilizada para armazenar o ID do usuário no contexto
type UsuarioKey string

// PerfilKey é a chave utilizada para armazenar o perfil do usuário no contexto
type PerfilKey string

// AuthMiddleware verifica se a requisição possui um token JWT válido
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
			
			tokenString := parts[1]
			
			// Validar o token JWT
			claims, err := auth.ValidarToken(tokenString)
			if err != nil {
				http.Error(w, "Token inválido ou expirado: "+err.Error(), http.StatusUnauthorized)
				return
			}
			
			// Verificar se o usuário ainda existe e está ativo
			var exists bool
			err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM usuarios WHERE id = $1)", claims.UserID).Scan(&exists)
			if err != nil || !exists {
				http.Error(w, "Usuário não encontrado ou inativo", http.StatusUnauthorized)
				return
			}
			
			// Adicionar o ID do usuário e perfil ao contexto da requisição
			ctx := context.WithValue(r.Context(), UsuarioKey("usuarioID"), claims.UserID)
			ctx = context.WithValue(ctx, PerfilKey("perfil"), claims.Perfil)
			
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

// ObterPerfilUsuario extrai o perfil do usuário do contexto da requisição
func ObterPerfilUsuario(r *http.Request) (string, bool) {
	perfil, ok := r.Context().Value(PerfilKey("perfil")).(string)
	return perfil, ok
}

// VerificarPerfil verifica se o usuário possui um determinado perfil
func VerificarPerfil(perfil string, perfilRequerido string) bool {
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
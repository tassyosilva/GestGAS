package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/tassyosilva/GestGAS/internal/auth"
)

// LoginHandler processa requisições de login
func LoginHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Configurar cabeçalhos CORS para esta resposta específica
		w.Header().Set("Content-Type", "application/json")
		
		// Verificar se é uma requisição OPTIONS (preflight)
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		
		// Verificar se é uma requisição POST
		if r.Method != http.MethodPost {
			http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
			return
		}
		
		// Decodificar o corpo da requisição
		var req struct {
			Login string `json:"login"`
			Senha string `json:"senha"`
		}
		
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&req); err != nil {
			http.Error(w, "Erro ao decodificar requisição", http.StatusBadRequest)
			return
		}
		
		// Verificar se os campos obrigatórios estão presentes
		if req.Login == "" || req.Senha == "" {
			http.Error(w, "Login e senha são obrigatórios", http.StatusBadRequest)
			return
		}
		
		// Buscar usuário pelo login
		var usuario struct {
			ID     int    `json:"id"`
			Nome   string `json:"nome"`
			Login  string `json:"login"`
			Senha  string `json:"senha"`
			Perfil string `json:"perfil"`
		}
		
		err := db.QueryRow(
			"SELECT id, nome, login, senha, perfil FROM usuarios WHERE login = $1",
			req.Login,
		).Scan(&usuario.ID, &usuario.Nome, &usuario.Login, &usuario.Senha, &usuario.Perfil)
		
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Credenciais inválidas", http.StatusUnauthorized)
				return
			}
			http.Error(w, "Erro ao buscar usuário", http.StatusInternalServerError)
			return
		}
		
		// Verificar senha usando bcrypt
		if !auth.VerificarSenha(req.Senha, usuario.Senha) {
			http.Error(w, "Credenciais inválidas", http.StatusUnauthorized)
			return
		}
		
		// Gerar token JWT incluindo o ID e o perfil do usuário
		token, err := auth.GerarToken(usuario.ID, usuario.Perfil)
		if err != nil {
			http.Error(w, "Erro ao gerar token", http.StatusInternalServerError)
			return
		}
		
		// Criar resposta
		resp := struct {
			ID     int       `json:"id"`
			Nome   string    `json:"nome"`
			Login  string    `json:"login"`
			Perfil string    `json:"perfil"`
			Token  string    `json:"token"`
			Expira time.Time `json:"expira"`
		}{
			ID:     usuario.ID,
			Nome:   usuario.Nome,
			Login:  usuario.Login,
			Perfil: usuario.Perfil,
			Token:  token,
			Expira: time.Now().Add(24 * time.Hour), // Token expira em 24 horas
		}
		
		// Enviar resposta
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	}
}
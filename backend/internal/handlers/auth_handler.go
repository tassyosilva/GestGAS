package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/tassyosilva/GestGAS/internal/auth"
	"github.com/tassyosilva/GestGAS/internal/models"
)

// LoginHandler processa requisições de login
func LoginHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Verificar se é uma requisição POST
		if r.Method != http.MethodPost {
			http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
			return
		}

		// Decodificar o corpo da requisição
		var req models.LoginRequest
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
		var usuario models.Usuario
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

		// Verificar senha
		if !auth.VerificarSenha(req.Senha, usuario.Senha) {
			http.Error(w, "Credenciais inválidas", http.StatusUnauthorized)
			return
		}

		// Gerar token
		token, err := auth.GerarToken(usuario.ID)
		if err != nil {
			http.Error(w, "Erro ao gerar token", http.StatusInternalServerError)
			return
		}
		
		// Armazenar token no banco de dados (expira em 24 horas)
		_, err = db.Exec(
			"INSERT INTO tokens (usuario_id, token, expiracao) VALUES ($1, $2, NOW() + INTERVAL '24 hours')",
			usuario.ID, token,
		)
		if err != nil {
			http.Error(w, "Erro ao armazenar token", http.StatusInternalServerError)
			return
		}

		// Criar resposta
		resp := models.LoginResponse{
			ID:     usuario.ID,
			Nome:   usuario.Nome,
			Login:  usuario.Login,
			Perfil: usuario.Perfil,
			Token:  token,
		}

		// Configurar cabeçalhos e enviar resposta
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	}
}
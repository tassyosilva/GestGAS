package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/tassyosilva/GestGAS/internal/auth"
	"github.com/tassyosilva/GestGAS/internal/middleware"
	"github.com/tassyosilva/GestGAS/internal/models"
)

// ListarUsuariosHandler retorna a lista de todos os usuários
func ListarUsuariosHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Verificar se é um método GET
		if r.Method != http.MethodGet {
			http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
			return
		}

		// Verificar permissões (apenas admin e gerente podem listar todos os usuários)
		perfil, ok := middleware.ObterPerfilUsuario(r)
		if !ok || !middleware.VerificarPerfil(perfil, "gerente") {
			http.Error(w, "Acesso negado", http.StatusForbidden)
			return
		}

		// Consultar usuários no banco de dados
		rows, err := db.Query(`
			SELECT id, nome, login, cpf, email, perfil, criado_em, atualizado_em 
			FROM usuarios 
			ORDER BY nome
		`)
		if err != nil {
			http.Error(w, "Erro ao consultar usuários", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		// Criar slice para armazenar os usuários
		var usuarios []models.Usuario

		// Iterar sobre os resultados
		for rows.Next() {
			var usuario models.Usuario
			err := rows.Scan(
				&usuario.ID,
				&usuario.Nome,
				&usuario.Login,
				&usuario.CPF,
				&usuario.Email,
				&usuario.Perfil,
				&usuario.CriadoEm,
				&usuario.AtualizadoEm,
			)
			if err != nil {
				http.Error(w, "Erro ao processar dados do usuário", http.StatusInternalServerError)
				return
			}
			usuarios = append(usuarios, usuario)
		}

		// Verificar erros durante a iteração
		if err = rows.Err(); err != nil {
			http.Error(w, "Erro ao percorrer resultados", http.StatusInternalServerError)
			return
		}

		// Retornar a lista de usuários como JSON
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(usuarios)
	}
}

// ObterUsuarioHandler retorna um usuário específico pelo ID
func ObterUsuarioHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Verificar se é um método GET
		if r.Method != http.MethodGet {
			http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
			return
		}

		// Extrair ID do usuário da URL
		path := strings.Split(r.URL.Path, "/")
		if len(path) < 4 {
			http.Error(w, "ID do usuário não especificado", http.StatusBadRequest)
			return
		}
		
		idStr := path[3]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "ID inválido", http.StatusBadRequest)
			return
		}

		// Verificar permissões
		userID, okID := middleware.ObterUsuarioID(r)
		perfil, okPerfil := middleware.ObterPerfilUsuario(r)
		
		// Apenas admin e gerente podem ver detalhes de qualquer usuário
		// Outros usuários só podem ver seus próprios detalhes
		if !okID || !okPerfil || (userID != id && !middleware.VerificarPerfil(perfil, "gerente")) {
			http.Error(w, "Acesso negado", http.StatusForbidden)
			return
		}

		// Consultar usuário no banco de dados
		var usuario models.Usuario
		err = db.QueryRow(`
			SELECT id, nome, login, cpf, email, perfil, criado_em, atualizado_em 
			FROM usuarios 
			WHERE id = $1
		`, id).Scan(
			&usuario.ID,
			&usuario.Nome,
			&usuario.Login,
			&usuario.CPF,
			&usuario.Email,
			&usuario.Perfil,
			&usuario.CriadoEm,
			&usuario.AtualizadoEm,
		)

		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Usuário não encontrado", http.StatusNotFound)
			} else {
				http.Error(w, "Erro ao buscar usuário", http.StatusInternalServerError)
			}
			return
		}

		// Retornar o usuário como JSON
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(usuario)
	}
}

// CriarUsuarioHandler cadastra um novo usuário
func CriarUsuarioHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Verificar se é um método POST
		if r.Method != http.MethodPost {
			http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
			return
		}

		// Verificar permissões (apenas admin pode criar usuários)
		perfil, ok := middleware.ObterPerfilUsuario(r)
		if !ok || perfil != "admin" {
			http.Error(w, "Apenas administradores podem criar usuários", http.StatusForbidden)
			return
		}

		// Decodificar o corpo da requisição
		var req struct {
			Nome    string `json:"nome"`
			Login   string `json:"login"`
			Senha   string `json:"senha"`
			CPF     string `json:"cpf"`
			Email   string `json:"email"`
			Perfil  string `json:"perfil"`
		}

		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&req); err != nil {
			http.Error(w, "Erro ao decodificar requisição", http.StatusBadRequest)
			return
		}

		// Validar campos obrigatórios
		if req.Nome == "" || req.Login == "" || req.Senha == "" || req.Perfil == "" {
			http.Error(w, "Nome, login, senha e perfil são obrigatórios", http.StatusBadRequest)
			return
		}

		// Validar perfil
		if req.Perfil != "admin" && req.Perfil != "gerente" && req.Perfil != "atendente" && req.Perfil != "entregador" {
			http.Error(w, "Perfil inválido", http.StatusBadRequest)
			return
		}

		// Verificar se login já existe
		var count int
		err := db.QueryRow("SELECT COUNT(*) FROM usuarios WHERE login = $1", req.Login).Scan(&count)
		if err != nil {
			http.Error(w, "Erro ao verificar login", http.StatusInternalServerError)
			return
		}
		if count > 0 {
			http.Error(w, "Login já existe", http.StatusBadRequest)
			return
		}

		// Verificar se CPF já existe (se fornecido)
		if req.CPF != "" {
			err = db.QueryRow("SELECT COUNT(*) FROM usuarios WHERE cpf = $1", req.CPF).Scan(&count)
			if err != nil {
				http.Error(w, "Erro ao verificar CPF", http.StatusInternalServerError)
				return
			}
			if count > 0 {
				http.Error(w, "CPF já cadastrado", http.StatusBadRequest)
				return
			}
		}

		// Verificar se Email já existe (se fornecido)
		if req.Email != "" {
			err = db.QueryRow("SELECT COUNT(*) FROM usuarios WHERE email = $1", req.Email).Scan(&count)
			if err != nil {
				http.Error(w, "Erro ao verificar email", http.StatusInternalServerError)
				return
			}
			if count > 0 {
				http.Error(w, "Email já cadastrado", http.StatusBadRequest)
				return
			}
		}

		// Gerar hash da senha
		senhaHash, err := auth.HashSenha(req.Senha)
		if err != nil {
			http.Error(w, "Erro ao processar senha", http.StatusInternalServerError)
			return
		}

		// Inserir novo usuário no banco de dados
		var usuarioID int
		err = db.QueryRow(`
			INSERT INTO usuarios (nome, login, senha, cpf, email, perfil, criado_em, atualizado_em)
			VALUES ($1, $2, $3, $4, $5, $6, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
			RETURNING id
		`, req.Nome, req.Login, senhaHash, req.CPF, req.Email, req.Perfil).Scan(&usuarioID)

		if err != nil {
			http.Error(w, "Erro ao criar usuário: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Buscar o usuário recém-criado (sem a senha)
		var usuario models.Usuario
		err = db.QueryRow(`
			SELECT id, nome, login, cpf, email, perfil, criado_em, atualizado_em 
			FROM usuarios 
			WHERE id = $1
		`, usuarioID).Scan(
			&usuario.ID,
			&usuario.Nome,
			&usuario.Login,
			&usuario.CPF,
			&usuario.Email,
			&usuario.Perfil,
			&usuario.CriadoEm,
			&usuario.AtualizadoEm,
		)

		if err != nil {
			http.Error(w, "Usuário criado, mas erro ao retornar dados", http.StatusInternalServerError)
			return
		}

		// Retornar o novo usuário como JSON
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(usuario)
	}
}

// AtualizarUsuarioHandler atualiza um usuário existente
func AtualizarUsuarioHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Verificar se é um método PUT ou PATCH
		if r.Method != http.MethodPut && r.Method != http.MethodPatch {
			http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
			return
		}

		// Extrair ID do usuário da URL
		path := strings.Split(r.URL.Path, "/")
		if len(path) < 4 {
			http.Error(w, "ID do usuário não especificado", http.StatusBadRequest)
			return
		}
		
		idStr := path[3]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "ID inválido", http.StatusBadRequest)
			return
		}

		// Verificar permissões
		userID, okID := middleware.ObterUsuarioID(r)
		perfil, okPerfil := middleware.ObterPerfilUsuario(r)
		
		// Apenas admin pode atualizar qualquer usuário
		// Outros usuários só podem atualizar seus próprios dados
		if !okID || !okPerfil || (userID != id && perfil != "admin") {
			http.Error(w, "Acesso negado", http.StatusForbidden)
			return
		}

		// Verificar se o usuário existe
		var usuarioAtual models.Usuario
		err = db.QueryRow("SELECT id, perfil FROM usuarios WHERE id = $1", id).Scan(&usuarioAtual.ID, &usuarioAtual.Perfil)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Usuário não encontrado", http.StatusNotFound)
			} else {
				http.Error(w, "Erro ao verificar usuário", http.StatusInternalServerError)
			}
			return
		}

		// Decodificar o corpo da requisição
		var req struct {
			Nome    string `json:"nome"`
			Login   string `json:"login"`
			Senha   string `json:"senha"`
			CPF     string `json:"cpf"`
			Email   string `json:"email"`
			Perfil  string `json:"perfil"`
		}

		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&req); err != nil {
			http.Error(w, "Erro ao decodificar requisição", http.StatusBadRequest)
			return
		}

		// Validar campos obrigatórios
		if req.Nome == "" {
			http.Error(w, "Nome é obrigatório", http.StatusBadRequest)
			return
		}

		// Apenas admin pode mudar o perfil
		if perfil != "admin" && req.Perfil != "" && req.Perfil != usuarioAtual.Perfil {
			http.Error(w, "Apenas administradores podem alterar o perfil", http.StatusForbidden)
			return
		}

		// Verificar se login já existe (se estiver sendo alterado)
		if req.Login != "" {
			var count int
			err = db.QueryRow("SELECT COUNT(*) FROM usuarios WHERE login = $1 AND id != $2", req.Login, id).Scan(&count)
			if err != nil {
				http.Error(w, "Erro ao verificar login", http.StatusInternalServerError)
				return
			}
			if count > 0 {
				http.Error(w, "Login já existe", http.StatusBadRequest)
				return
			}
		}

		// Verificar se CPF já existe (se estiver sendo alterado)
		if req.CPF != "" {
			var count int
			err = db.QueryRow("SELECT COUNT(*) FROM usuarios WHERE cpf = $1 AND id != $2", req.CPF, id).Scan(&count)
			if err != nil {
				http.Error(w, "Erro ao verificar CPF", http.StatusInternalServerError)
				return
			}
			if count > 0 {
				http.Error(w, "CPF já cadastrado", http.StatusBadRequest)
				return
			}
		}

		// Verificar se Email já existe (se estiver sendo alterado)
		if req.Email != "" {
			var count int
			err = db.QueryRow("SELECT COUNT(*) FROM usuarios WHERE email = $1 AND id != $2", req.Email, id).Scan(&count)
			if err != nil {
				http.Error(w, "Erro ao verificar email", http.StatusInternalServerError)
				return
			}
			if count > 0 {
				http.Error(w, "Email já cadastrado", http.StatusBadRequest)
				return
			}
		}

		// Preparar a query de atualização
		query := `
			UPDATE usuarios SET 
				nome = $1, 
				atualizado_em = CURRENT_TIMESTAMP
		`
		params := []interface{}{req.Nome}
		paramCount := 2

		// Adicionar os campos opcionais na query (se fornecidos)
		if req.Login != "" {
			query += fmt.Sprintf(", login = $%d", paramCount)
			params = append(params, req.Login)
			paramCount++
		}

		if req.Senha != "" {
			// Gerar hash da nova senha
			senhaHash, err := auth.HashSenha(req.Senha)
			if err != nil {
				http.Error(w, "Erro ao processar senha", http.StatusInternalServerError)
				return
			}
			query += fmt.Sprintf(", senha = $%d", paramCount)
			params = append(params, senhaHash)
			paramCount++
		}

		if req.CPF != "" {
			query += fmt.Sprintf(", cpf = $%d", paramCount)
			params = append(params, req.CPF)
			paramCount++
		}

		if req.Email != "" {
			query += fmt.Sprintf(", email = $%d", paramCount)
			params = append(params, req.Email)
			paramCount++
		}

		if req.Perfil != "" && perfil == "admin" {
			query += fmt.Sprintf(", perfil = $%d", paramCount)
			params = append(params, req.Perfil)
			paramCount++
		}

		// Finalizar a query com a condição WHERE
		query += " WHERE id = $" + strconv.Itoa(paramCount)
		params = append(params, id)

		// Executar a atualização
		_, err = db.Exec(query, params...)
		if err != nil {
			http.Error(w, "Erro ao atualizar usuário: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Buscar o usuário atualizado
		var usuario models.Usuario
		err = db.QueryRow(`
			SELECT id, nome, login, cpf, email, perfil, criado_em, atualizado_em 
			FROM usuarios 
			WHERE id = $1
		`, id).Scan(
			&usuario.ID,
			&usuario.Nome,
			&usuario.Login,
			&usuario.CPF,
			&usuario.Email,
			&usuario.Perfil,
			&usuario.CriadoEm,
			&usuario.AtualizadoEm,
		)

		if err != nil {
			http.Error(w, "Erro ao buscar usuário atualizado", http.StatusInternalServerError)
			return
		}

		// Retornar o usuário atualizado como JSON
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(usuario)
	}
}

// ExcluirUsuarioHandler remove um usuário
func ExcluirUsuarioHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Verificar se é um método DELETE
		if r.Method != http.MethodDelete {
			http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
			return
		}

		// Verificar permissões (apenas admin pode excluir usuários)
		perfil, ok := middleware.ObterPerfilUsuario(r)
		if !ok || perfil != "admin" {
			http.Error(w, "Apenas administradores podem excluir usuários", http.StatusForbidden)
			return
		}

		// Extrair ID do usuário da URL
		path := strings.Split(r.URL.Path, "/")
		if len(path) < 4 {
			http.Error(w, "ID do usuário não especificado", http.StatusBadRequest)
			return
		}
		
		idStr := path[3]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "ID inválido", http.StatusBadRequest)
			return
		}

		// Impedir exclusão do próprio usuário
		userID, _ := middleware.ObterUsuarioID(r)
		if userID == id {
			http.Error(w, "Não é possível excluir o próprio usuário", http.StatusBadRequest)
			return
		}

		// Verificar se o usuário existe
		var existe bool
		err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM usuarios WHERE id = $1)", id).Scan(&existe)
		if err != nil {
			http.Error(w, "Erro ao verificar usuário", http.StatusInternalServerError)
			return
		}
		if !existe {
			http.Error(w, "Usuário não encontrado", http.StatusNotFound)
			return
		}

		// Verificar se usuário tem registros dependentes
		var temPedidos bool
		err = db.QueryRow(`
			SELECT EXISTS(
				SELECT 1 FROM pedidos 
				WHERE atendente_id = $1 OR entregador_id = $1
			)
		`, id).Scan(&temPedidos)
		if err != nil {
			http.Error(w, "Erro ao verificar pedidos do usuário", http.StatusInternalServerError)
			return
		}
		if temPedidos {
			http.Error(w, "Não é possível excluir o usuário pois existem pedidos associados", http.StatusBadRequest)
			return
		}

		// Excluir o usuário
		_, err = db.Exec("DELETE FROM usuarios WHERE id = $1", id)
		if err != nil {
			http.Error(w, "Erro ao excluir usuário: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Retornar resposta de sucesso
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"mensagem": "Usuário excluído com sucesso"})
	}
}

// ListarEntregadoresHandler retorna a lista de usuários com perfil entregador
func ListarEntregadoresHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Verificar se é um método GET
		if r.Method != http.MethodGet {
			http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
			return
		}

		// Consultar entregadores no banco de dados
		rows, err := db.Query(`
			SELECT id, nome, login, cpf, email, criado_em, atualizado_em 
			FROM usuarios 
			WHERE perfil = 'entregador'
			ORDER BY nome
		`)
		if err != nil {
			http.Error(w, "Erro ao consultar entregadores", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		// Criar slice para armazenar os entregadores
		var entregadores []models.Usuario

		// Iterar sobre os resultados
		for rows.Next() {
			var entregador models.Usuario
			entregador.Perfil = "entregador" // Já sabemos que é entregador
			
			err := rows.Scan(
				&entregador.ID,
				&entregador.Nome,
				&entregador.Login,
				&entregador.CPF,
				&entregador.Email,
				&entregador.CriadoEm,
				&entregador.AtualizadoEm,
			)
			if err != nil {
				http.Error(w, "Erro ao processar dados do entregador", http.StatusInternalServerError)
				return
			}
			entregadores = append(entregadores, entregador)
		}

		// Verificar erros durante a iteração
		if err = rows.Err(); err != nil {
			http.Error(w, "Erro ao percorrer resultados", http.StatusInternalServerError)
			return
		}

		// Retornar a lista de entregadores como JSON
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(entregadores)
	}
}
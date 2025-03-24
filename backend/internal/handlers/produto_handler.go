package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/tassyosilva/GestGAS/internal/middleware"
	"github.com/tassyosilva/GestGAS/internal/models"
)

// ListarProdutosHandler retorna a lista de produtos
func ListarProdutosHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Verificar se é uma requisição GET
		if r.Method != http.MethodGet {
			http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
			return
		}

		// Verificar permissão do usuário
		_, ok := middleware.ObterUsuarioID(r)
		if !ok {
			http.Error(w, "Usuário não autenticado", http.StatusUnauthorized)
			return
		}

		// Consultar produtos no banco de dados
		rows, err := db.Query("SELECT id, nome, descricao, categoria, preco, criado_em, atualizado_em FROM produtos ORDER BY nome")
		if err != nil {
			http.Error(w, "Erro ao buscar produtos", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		// Processar resultados
		produtos := []models.Produto{}
		for rows.Next() {
			var produto models.Produto
			err := rows.Scan(
				&produto.ID,
				&produto.Nome,
				&produto.Descricao,
				&produto.Categoria,
				&produto.Preco,
				&produto.CriadoEm,
				&produto.AtualizadoEm,
			)
			if err != nil {
				http.Error(w, "Erro ao processar produtos", http.StatusInternalServerError)
				return
			}
			produtos = append(produtos, produto)
		}

		// Verificar erros de iteração
		if err = rows.Err(); err != nil {
			http.Error(w, "Erro ao processar produtos", http.StatusInternalServerError)
			return
		}

		// Retornar a lista de produtos
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(produtos)
	}
}

// ObterProdutoHandler retorna um produto específico
func ObterProdutoHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Verificar se é uma requisição GET
		if r.Method != http.MethodGet {
			http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
			return
		}

		// Extrair ID do produto da URL
		path := r.URL.Path
		segments := strings.Split(path, "/")
		if len(segments) < 4 {
			http.Error(w, "ID do produto não especificado", http.StatusBadRequest)
			return
		}

		produtoID, err := strconv.Atoi(segments[3])
		if err != nil {
			http.Error(w, "ID do produto inválido", http.StatusBadRequest)
			return
		}

		// Buscar produto no banco de dados
		var produto models.Produto
		err = db.QueryRow(
			"SELECT id, nome, descricao, categoria, preco, criado_em, atualizado_em FROM produtos WHERE id = $1",
			produtoID,
		).Scan(
			&produto.ID,
			&produto.Nome,
			&produto.Descricao,
			&produto.Categoria,
			&produto.Preco,
			&produto.CriadoEm,
			&produto.AtualizadoEm,
		)

		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Produto não encontrado", http.StatusNotFound)
				return
			}
			http.Error(w, "Erro ao buscar produto", http.StatusInternalServerError)
			return
		}

		// Retornar o produto
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(produto)
	}
}

// CriarProdutoHandler cria um novo produto
func CriarProdutoHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Verificar se é uma requisição POST
		if r.Method != http.MethodPost {
			http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
			return
		}

		// Verificar permissão do usuário
		_, ok := middleware.ObterUsuarioID(r)
		if !ok {
			http.Error(w, "Usuário não autenticado", http.StatusUnauthorized)
			return
		}

		// Verificar se o usuário tem permissão para criar produtos
		perfil, perfilOk := middleware.ObterPerfilUsuario(r)
		if !perfilOk || !middleware.VerificarPerfil(perfil, "gerente") {
			http.Error(w, "Permissão negada", http.StatusForbidden)
			return
		}

		// Decodificar o corpo da requisição
		var produto models.Produto
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&produto); err != nil {
			http.Error(w, "Erro ao decodificar requisição", http.StatusBadRequest)
			return
		}

		// Validar dados do produto
		if produto.Nome == "" || produto.Categoria == "" || produto.Preco <= 0 {
			http.Error(w, "Dados do produto inválidos", http.StatusBadRequest)
			return
		}

		// Definir timestamps
		now := time.Now()
		produto.CriadoEm = now
		produto.AtualizadoEm = now

		// Iniciar uma transação para garantir que tanto o produto quanto o estoque sejam criados
		tx, err := db.Begin()
		if err != nil {
			http.Error(w, "Erro ao iniciar transação", http.StatusInternalServerError)
			return
		}

		// Inserir produto no banco de dados dentro da transação
		err = tx.QueryRow(
			`INSERT INTO produtos (nome, descricao, categoria, preco, criado_em, atualizado_em) 
             VALUES ($1, $2, $3, $4, $5, $6) 
             RETURNING id`,
			produto.Nome,
			produto.Descricao,
			produto.Categoria,
			produto.Preco,
			produto.CriadoEm,
			produto.AtualizadoEm,
		).Scan(&produto.ID)

		if err != nil {
			tx.Rollback()
			http.Error(w, "Erro ao criar produto", http.StatusInternalServerError)
			return
		}

		// Criar o registro de estoque dentro da mesma transação
		_, err = tx.Exec(
			"INSERT INTO estoque (produto_id, quantidade, alerta_minimo) VALUES ($1, 0, 5)",
			produto.ID,
		)
		if err != nil {
			tx.Rollback()
			http.Error(w, "Erro ao configurar estoque para o produto", http.StatusInternalServerError)
			return
		}

		// Commit da transação
		err = tx.Commit()
		if err != nil {
			tx.Rollback()
			http.Error(w, "Erro ao finalizar criação do produto", http.StatusInternalServerError)
			return
		}

		// Retornar o produto criado
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(produto)
	}
}

// AtualizarProdutoHandler atualiza um produto existente
func AtualizarProdutoHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Verificar se é uma requisição PUT
		if r.Method != http.MethodPut {
			http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
			return
		}

		// Extrair ID do produto da URL
		path := r.URL.Path
		segments := strings.Split(path, "/")
		if len(segments) < 4 {
			http.Error(w, "ID do produto não especificado", http.StatusBadRequest)
			return
		}

		produtoID, err := strconv.Atoi(segments[3])
		if err != nil {
			http.Error(w, "ID do produto inválido", http.StatusBadRequest)
			return
		}

		// Verificar permissão do usuário
		_, ok := middleware.ObterUsuarioID(r)
		if !ok {
			http.Error(w, "Usuário não autenticado", http.StatusUnauthorized)
			return
		}

		// Verificar se o usuário tem permissão para atualizar produtos
		perfil, perfilOk := middleware.ObterPerfilUsuario(r)
		if !perfilOk || !middleware.VerificarPerfil(perfil, "gerente") {
			http.Error(w, "Permissão negada", http.StatusForbidden)
			return
		}

		// Verificar se o produto existe
		var existe bool
		err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM produtos WHERE id = $1)", produtoID).Scan(&existe)
		if err != nil {
			http.Error(w, "Erro ao verificar produto", http.StatusInternalServerError)
			return
		}

		if !existe {
			http.Error(w, "Produto não encontrado", http.StatusNotFound)
			return
		}

		// Decodificar o corpo da requisição
		var produto models.Produto
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&produto); err != nil {
			http.Error(w, "Erro ao decodificar requisição", http.StatusBadRequest)
			return
		}

		// Validar dados do produto
		if produto.Nome == "" || produto.Categoria == "" || produto.Preco <= 0 {
			http.Error(w, "Dados do produto inválidos", http.StatusBadRequest)
			return
		}

		// Atualizar produto no banco de dados
		produto.ID = produtoID
		produto.AtualizadoEm = time.Now()

		_, err = db.Exec(
			`UPDATE produtos 
             SET nome = $1, descricao = $2, categoria = $3, preco = $4, atualizado_em = $5
             WHERE id = $6`,
			produto.Nome,
			produto.Descricao,
			produto.Categoria,
			produto.Preco,
			produto.AtualizadoEm,
			produto.ID,
		)

		if err != nil {
			http.Error(w, "Erro ao atualizar produto", http.StatusInternalServerError)
			return
		}

		// Buscar o produto atualizado
		err = db.QueryRow(
			"SELECT id, nome, descricao, categoria, preco, criado_em, atualizado_em FROM produtos WHERE id = $1",
			produto.ID,
		).Scan(
			&produto.ID,
			&produto.Nome,
			&produto.Descricao,
			&produto.Categoria,
			&produto.Preco,
			&produto.CriadoEm,
			&produto.AtualizadoEm,
		)

		if err != nil {
			http.Error(w, "Produto atualizado, mas erro ao buscar dados atualizados", http.StatusInternalServerError)
			return
		}

		// Retornar o produto atualizado
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(produto)
	}
}

// ExcluirProdutoHandler remove um produto
func ExcluirProdutoHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Verificar se é uma requisição DELETE
		if r.Method != http.MethodDelete {
			http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
			return
		}

		// Extrair ID do produto da URL
		path := r.URL.Path
		segments := strings.Split(path, "/")
		if len(segments) < 4 {
			http.Error(w, "ID do produto não especificado", http.StatusBadRequest)
			return
		}

		produtoID, err := strconv.Atoi(segments[3])
		if err != nil {
			http.Error(w, "ID do produto inválido", http.StatusBadRequest)
			return
		}

		// Verificar permissão do usuário
		_, ok := middleware.ObterUsuarioID(r)
		if !ok {
			http.Error(w, "Usuário não autenticado", http.StatusUnauthorized)
			return
		}

		// Apenas administradores podem excluir produtos
		perfil, perfilOk := middleware.ObterPerfilUsuario(r)
		if !perfilOk || !middleware.VerificarPerfil(perfil, "admin") {
			http.Error(w, "Permissão negada", http.StatusForbidden)
			return
		}

		// Verificar se o produto existe
		var existe bool
		err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM produtos WHERE id = $1)", produtoID).Scan(&existe)
		if err != nil {
			http.Error(w, "Erro ao verificar produto", http.StatusInternalServerError)
			return
		}

		if !existe {
			http.Error(w, "Produto não encontrado", http.StatusNotFound)
			return
		}

		// Em um sistema real, verificaríamos se o produto pode ser excluído
		// (por exemplo, se não há pedidos ou estoque vinculados a ele)
		// Aqui, faremos uma exclusão simples

		// Primeiro excluir registros de estoque relacionados
		_, err = db.Exec("DELETE FROM estoque WHERE produto_id = $1", produtoID)
		if err != nil {
			http.Error(w, "Erro ao excluir estoque do produto", http.StatusInternalServerError)
			return
		}

		// Excluir o produto
		_, err = db.Exec("DELETE FROM produtos WHERE id = $1", produtoID)
		if err != nil {
			http.Error(w, "Erro ao excluir produto", http.StatusInternalServerError)
			return
		}

		// Retornar sucesso
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Produto excluído com sucesso"})
	}
}
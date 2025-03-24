package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/tassyosilva/GestGAS/internal/middleware"
	"github.com/tassyosilva/GestGAS/internal/models"
)

// ListarClientesHandler retorna a lista de clientes com paginação e filtros
func ListarClientesHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Verificar se o usuário está autenticado
		_, ok := middleware.ObterUsuarioID(r)
		if !ok {
			http.Error(w, "Usuário não autenticado", http.StatusUnauthorized)
			return
		}

		// Configurar cabeçalhos
		w.Header().Set("Content-Type", "application/json")

		// Verificar método
		if r.Method != http.MethodGet {
			http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
			return
		}

		// Parâmetros de consulta
		query := r.URL.Query()
		page, _ := strconv.Atoi(query.Get("page"))
		limit, _ := strconv.Atoi(query.Get("limit"))
		nome := query.Get("nome")
		telefone := query.Get("telefone")
		canalOrigem := query.Get("canal_origem")

		if page < 1 {
			page = 1
		}
		if limit < 1 || limit > 100 {
			limit = 20
		}
		offset := (page - 1) * limit

		// Construir consulta SQL
		sqlQuery := `
			SELECT 
				id, nome, telefone, cpf, email, 
				endereco, complemento, bairro, cidade, estado, 
				cep, observacoes, canal_origem, criado_em, atualizado_em
			FROM clientes
			WHERE 1=1
		`
		var params []interface{}
		var whereConditions []string

		// Adicionar filtros se fornecidos
		if nome != "" {
			whereConditions = append(whereConditions, fmt.Sprintf("LOWER(nome) LIKE LOWER($%d)", len(params)+1))
			params = append(params, "%"+nome+"%")
		}
		if telefone != "" {
			whereConditions = append(whereConditions, fmt.Sprintf("telefone LIKE $%d", len(params)+1))
			params = append(params, "%"+telefone+"%")
		}
		if canalOrigem != "" {
			whereConditions = append(whereConditions, fmt.Sprintf("canal_origem = $%d", len(params)+1))
			params = append(params, canalOrigem)
		}

		// Adicionar condições WHERE
		if len(whereConditions) > 0 {
			sqlQuery += " AND " + strings.Join(whereConditions, " AND ")
		}

		// Adicionar ordenação e paginação
		sqlQuery += " ORDER BY nome ASC LIMIT $" + strconv.Itoa(len(params)+1) + " OFFSET $" + strconv.Itoa(len(params)+2)
		params = append(params, limit, offset)

		// Executar consulta
		rows, err := db.Query(sqlQuery, params...)
		if err != nil {
			http.Error(w, "Erro ao buscar clientes: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		// Processar resultados
		var clientes []models.Cliente
		for rows.Next() {
			var c models.Cliente
			var cpf, email, endereco, complemento, bairro, cidade, estado, cep, observacoes sql.NullString
			var canalOrigem sql.NullString
			var criadoEm, atualizadoEm time.Time

			err := rows.Scan(
				&c.ID, &c.Nome, &c.Telefone, &cpf, &email,
				&endereco, &complemento, &bairro, &cidade, &estado,
				&cep, &observacoes, &canalOrigem, &criadoEm, &atualizadoEm,
			)
			if err != nil {
				http.Error(w, "Erro ao processar clientes: "+err.Error(), http.StatusInternalServerError)
				return
			}

			// Converter tipos nulos
			if cpf.Valid {
				c.CPF = cpf.String
			}
			if email.Valid {
				c.Email = email.String
			}
			if endereco.Valid {
				c.Endereco = endereco.String
			}
			if complemento.Valid {
				c.Complemento = complemento.String
			}
			if bairro.Valid {
				c.Bairro = bairro.String
			}
			if cidade.Valid {
				c.Cidade = cidade.String
			}
			if estado.Valid {
				c.Estado = estado.String
			}
			if cep.Valid {
				c.CEP = cep.String
			}
			if observacoes.Valid {
				c.Observacoes = observacoes.String
			}
			if canalOrigem.Valid {
				c.CanalOrigem = models.CanalOrigem(canalOrigem.String)
			}
			
			c.CriadoEm = criadoEm
			c.AtualizadoEm = atualizadoEm

			clientes = append(clientes, c)
		}

		// Contar total de registros para paginação
		var total int
		countQuery := `
			SELECT COUNT(*) FROM clientes WHERE 1=1
		`
		if len(whereConditions) > 0 {
			countQuery += " AND " + strings.Join(whereConditions, " AND ")
		}

		err = db.QueryRow(countQuery, params[:len(params)-2]...).Scan(&total)
		if err != nil {
			http.Error(w, "Erro ao contar clientes: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Montar resposta
		response := struct {
			Clientes []models.Cliente `json:"clientes"`
			Total    int             `json:"total"`
			Page     int             `json:"page"`
			Limit    int             `json:"limit"`
			Pages    int             `json:"pages"`
		}{
			Clientes: clientes,
			Total:    total,
			Page:     page,
			Limit:    limit,
			Pages:    (total + limit - 1) / limit,
		}

		// Retornar resposta
		json.NewEncoder(w).Encode(response)
	}
}

// ObterClienteHandler retorna detalhes de um cliente específico
func ObterClienteHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Verificar se o usuário está autenticado
		_, ok := middleware.ObterUsuarioID(r)
		if !ok {
			http.Error(w, "Usuário não autenticado", http.StatusUnauthorized)
			return
		}

		// Configurar cabeçalhos
		w.Header().Set("Content-Type", "application/json")

		// Verificar método
		if r.Method != http.MethodGet {
			http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
			return
		}

		// Extrair ID do cliente da URL
		parts := strings.Split(r.URL.Path, "/")
		if len(parts) < 4 {
			http.Error(w, "ID do cliente não fornecido", http.StatusBadRequest)
			return
		}
		clienteID, err := strconv.Atoi(parts[len(parts)-1])
		if err != nil {
			http.Error(w, "ID do cliente inválido", http.StatusBadRequest)
			return
		}

		// Buscar cliente
		var response models.ClienteResponse
		var cpf, email, endereco, complemento, bairro, cidade, estado, cep, observacoes sql.NullString
		var canalOrigem sql.NullString
		var criadoEm, atualizadoEm time.Time

		err = db.QueryRow(`
			SELECT 
				id, nome, telefone, cpf, email, 
				endereco, complemento, bairro, cidade, estado, 
				cep, observacoes, canal_origem, criado_em, atualizado_em
			FROM clientes
			WHERE id = $1
		`, clienteID).Scan(
			&response.ID, &response.Nome, &response.Telefone, &cpf, &email,
			&endereco, &complemento, &bairro, &cidade, &estado,
			&cep, &observacoes, &canalOrigem, &criadoEm, &atualizadoEm,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Cliente não encontrado", http.StatusNotFound)
				return
			}
			http.Error(w, "Erro ao buscar cliente: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Converter tipos nulos
		if cpf.Valid {
			response.CPF = cpf.String
		}
		if email.Valid {
			response.Email = email.String
		}
		if endereco.Valid {
			response.Endereco = endereco.String
		}
		if complemento.Valid {
			response.Complemento = complemento.String
		}
		if bairro.Valid {
			response.Bairro = bairro.String
		}
		if cidade.Valid {
			response.Cidade = cidade.String
		}
		if estado.Valid {
			response.Estado = estado.String
		}
		if cep.Valid {
			response.CEP = cep.String
		}
		if observacoes.Valid {
			response.Observacoes = observacoes.String
		}
		if canalOrigem.Valid {
			response.CanalOrigem = models.CanalOrigem(canalOrigem.String)
		}
		
		response.CriadoEm = criadoEm
		response.AtualizadoEm = atualizadoEm

		// Contar total de pedidos do cliente
		err = db.QueryRow("SELECT COUNT(*) FROM pedidos WHERE cliente_id = $1", clienteID).Scan(&response.TotalPedidos)
		if err != nil {
			http.Error(w, "Erro ao contar pedidos do cliente: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Buscar últimos pedidos do cliente 
		rows, err := db.Query(`
			SELECT id, status, forma_pagamento, valor_total, criado_em
			FROM pedidos
			WHERE cliente_id = $1
			ORDER BY criado_em DESC
			LIMIT 5
		`, clienteID)
		if err != nil {
			http.Error(w, "Erro ao buscar pedidos do cliente: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var ultimosPedidos []models.PedidoResumido
		for rows.Next() {
			p := models.PedidoResumido{}
			var dataPedido time.Time
			err := rows.Scan(&p.ID, &p.Status, &p.FormaPagamento, &p.ValorTotal, &dataPedido)
			if err != nil {
				http.Error(w, "Erro ao processar pedidos do cliente: "+err.Error(), http.StatusInternalServerError)
				return
			}
			p.DataPedido = dataPedido
			ultimosPedidos = append(ultimosPedidos, p)
		}
		
		response.UltimosPedidos = ultimosPedidos

		// Retornar resposta
		json.NewEncoder(w).Encode(response)
	}
}

// CriarClienteHandler cria um novo cliente
func CriarClienteHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Verificar se o usuário está autenticado
		_, ok := middleware.ObterUsuarioID(r)
		if !ok {
			http.Error(w, "Usuário não autenticado", http.StatusUnauthorized)
			return
		}

		// Verificar permissões (apenas atendentes ou acima podem criar clientes)
		perfil, perfilOk := middleware.ObterPerfilUsuario(r)
		if !perfilOk || !middleware.VerificarPerfil(perfil, "atendente") {
			http.Error(w, "Sem permissão para criar clientes", http.StatusForbidden)
			return
		}

		// Configurar cabeçalhos
		w.Header().Set("Content-Type", "application/json")

		// Verificar método
		if r.Method != http.MethodPost {
			http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
			return
		}

		// Decodificar requisição
		var req models.NovoClienteRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Erro ao processar requisição: "+err.Error(), http.StatusBadRequest)
			return
		}

		// Validar dados
		if req.Nome == "" {
			http.Error(w, "Nome é obrigatório", http.StatusBadRequest)
			return
		}
		if req.Telefone == "" {
			http.Error(w, "Telefone é obrigatório", http.StatusBadRequest)
			return
		}

		// Verificar se já existe cliente com mesmo CPF ou email (se fornecidos)
		if req.CPF != "" {
			var exists bool
			err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM clientes WHERE cpf = $1)", req.CPF).Scan(&exists)
			if err != nil {
				http.Error(w, "Erro ao verificar CPF: "+err.Error(), http.StatusInternalServerError)
				return
			}
			if exists {
				http.Error(w, "CPF já cadastrado", http.StatusBadRequest)
				return
			}
		}

		if req.Email != "" {
			var exists bool
			err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM clientes WHERE email = $1)", req.Email).Scan(&exists)
			if err != nil {
				http.Error(w, "Erro ao verificar email: "+err.Error(), http.StatusInternalServerError)
				return
			}
			if exists {
				http.Error(w, "Email já cadastrado", http.StatusBadRequest)
				return
			}
		}

		// Inserir cliente
		var clienteID int
		err := db.QueryRow(`
			INSERT INTO clientes (
				nome, telefone, cpf, email,
				endereco, complemento, bairro, cidade, estado,
				cep, observacoes, canal_origem,
				criado_em, atualizado_em
			) VALUES (
				$1, $2, NULLIF($3, ''), NULLIF($4, ''),
				NULLIF($5, ''), NULLIF($6, ''), NULLIF($7, ''), NULLIF($8, ''), NULLIF($9, ''),
				NULLIF($10, ''), NULLIF($11, ''), $12,
				NOW(), NOW()
			) RETURNING id
		`, req.Nome, req.Telefone, req.CPF, req.Email,
		   req.Endereco, req.Complemento, req.Bairro, req.Cidade, req.Estado,
		   req.CEP, req.Observacoes, req.CanalOrigem).Scan(&clienteID)

		if err != nil {
			http.Error(w, "Erro ao criar cliente: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Buscar cliente criado para resposta
		var cliente models.Cliente
		var cpf, email, endereco, complemento, bairro, cidade, estado, cep, observacoes sql.NullString
		var canalOrigem sql.NullString
		var criadoEm, atualizadoEm time.Time

		err = db.QueryRow(`
			SELECT 
				id, nome, telefone, cpf, email, 
				endereco, complemento, bairro, cidade, estado, 
				cep, observacoes, canal_origem, criado_em, atualizado_em
			FROM clientes
			WHERE id = $1
		`, clienteID).Scan(
			&cliente.ID, &cliente.Nome, &cliente.Telefone, &cpf, &email,
			&endereco, &complemento, &bairro, &cidade, &estado,
			&cep, &observacoes, &canalOrigem, &criadoEm, &atualizadoEm,
		)
		if err != nil {
			http.Error(w, "Cliente criado, mas erro ao buscar detalhes: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Converter tipos nulos
		if cpf.Valid {
			cliente.CPF = cpf.String
		}
		if email.Valid {
			cliente.Email = email.String
		}
		if endereco.Valid {
			cliente.Endereco = endereco.String
		}
		if complemento.Valid {
			cliente.Complemento = complemento.String
		}
		if bairro.Valid {
			cliente.Bairro = bairro.String
		}
		if cidade.Valid {
			cliente.Cidade = cidade.String
		}
		if estado.Valid {
			cliente.Estado = estado.String
		}
		if cep.Valid {
			cliente.CEP = cep.String
		}
		if observacoes.Valid {
			cliente.Observacoes = observacoes.String
		}
		if canalOrigem.Valid {
			cliente.CanalOrigem = models.CanalOrigem(canalOrigem.String)
		}
		
		cliente.CriadoEm = criadoEm
		cliente.AtualizadoEm = atualizadoEm

		// Retornar resposta
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(cliente)
	}
}

// AtualizarClienteHandler atualiza um cliente existente
func AtualizarClienteHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Verificar se o usuário está autenticado
		_, ok := middleware.ObterUsuarioID(r)
		if !ok {
			http.Error(w, "Usuário não autenticado", http.StatusUnauthorized)
			return
		}

		// Verificar permissões (apenas atendentes ou acima podem atualizar clientes)
		perfil, perfilOk := middleware.ObterPerfilUsuario(r)
		if !perfilOk || !middleware.VerificarPerfil(perfil, "atendente") {
			http.Error(w, "Sem permissão para atualizar clientes", http.StatusForbidden)
			return
		}

		// Configurar cabeçalhos
		w.Header().Set("Content-Type", "application/json")

		// Verificar método
		if r.Method != http.MethodPut && r.Method != http.MethodPatch {
			http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
			return
		}

		// Extrair ID do cliente da URL
		parts := strings.Split(r.URL.Path, "/")
		if len(parts) < 4 {
			http.Error(w, "ID do cliente não fornecido", http.StatusBadRequest)
			return
		}
		clienteID, err := strconv.Atoi(parts[len(parts)-1])
		if err != nil {
			http.Error(w, "ID do cliente inválido", http.StatusBadRequest)
			return
		}

		// Verificar se o cliente existe
		var clienteExiste bool
		err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM clientes WHERE id = $1)", clienteID).Scan(&clienteExiste)
		if err != nil {
			http.Error(w, "Erro ao verificar cliente: "+err.Error(), http.StatusInternalServerError)
			return
		}
		if !clienteExiste {
			http.Error(w, "Cliente não encontrado", http.StatusNotFound)
			return
		}

		// Verificar se é uma atualização apenas de endereço
		if len(parts) >= 5 && parts[len(parts)-1] == "endereco" {
			// Decodificar requisição de endereço
			var req models.ClienteEnderecoRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, "Erro ao processar requisição: "+err.Error(), http.StatusBadRequest)
				return
			}

			// Atualizar endereço
			_, err = db.Exec(`
				UPDATE clientes SET
					endereco = $1,
					complemento = NULLIF($2, ''),
					bairro = NULLIF($3, ''),
					cidade = NULLIF($4, ''),
					estado = NULLIF($5, ''),
					cep = NULLIF($6, ''),
					atualizado_em = NOW()
				WHERE id = $7
			`, req.Endereco, req.Complemento, req.Bairro, req.Cidade, req.Estado, req.CEP, clienteID)

			if err != nil {
				http.Error(w, "Erro ao atualizar endereço: "+err.Error(), http.StatusInternalServerError)
				return
			}

			// Retornar resposta simples
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{
				"mensagem": "Endereço atualizado com sucesso",
			})
			return
		}

		// Atualização completa do cliente
		var req models.NovoClienteRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Erro ao processar requisição: "+err.Error(), http.StatusBadRequest)
			return
		}

		// Validar dados
		if req.Nome == "" {
			http.Error(w, "Nome é obrigatório", http.StatusBadRequest)
			return
		}
		if req.Telefone == "" {
			http.Error(w, "Telefone é obrigatório", http.StatusBadRequest)
			return
		}

		// Verificar se CPF ou email já existem para outro cliente
		if req.CPF != "" {
			var exists bool
			err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM clientes WHERE cpf = $1 AND id != $2)", req.CPF, clienteID).Scan(&exists)
			if err != nil {
				http.Error(w, "Erro ao verificar CPF: "+err.Error(), http.StatusInternalServerError)
				return
			}
			if exists {
				http.Error(w, "CPF já cadastrado para outro cliente", http.StatusBadRequest)
				return
			}
		}

		if req.Email != "" {
			var exists bool
			err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM clientes WHERE email = $1 AND id != $2)", req.Email, clienteID).Scan(&exists)
			if err != nil {
				http.Error(w, "Erro ao verificar email: "+err.Error(), http.StatusInternalServerError)
				return
			}
			if exists {
				http.Error(w, "Email já cadastrado para outro cliente", http.StatusBadRequest)
				return
			}
		}

		// Atualizar cliente
		_, err = db.Exec(`
			UPDATE clientes SET
				nome = $1,
				telefone = $2,
				cpf = NULLIF($3, ''),
				email = NULLIF($4, ''),
				endereco = NULLIF($5, ''),
				complemento = NULLIF($6, ''),
				bairro = NULLIF($7, ''),
				cidade = NULLIF($8, ''),
				estado = NULLIF($9, ''),
				cep = NULLIF($10, ''),
				observacoes = NULLIF($11, ''),
				canal_origem = $12,
				atualizado_em = NOW()
			WHERE id = $13
		`, req.Nome, req.Telefone, req.CPF, req.Email,
		   req.Endereco, req.Complemento, req.Bairro, req.Cidade, req.Estado,
		   req.CEP, req.Observacoes, req.CanalOrigem, clienteID)

		if err != nil {
			http.Error(w, "Erro ao atualizar cliente: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Buscar cliente atualizado para resposta
		var response models.ClienteResponse
		var cpf, email, endereco, complemento, bairro, cidade, estado, cep, observacoes sql.NullString
		var canalOrigem sql.NullString
		var criadoEm, atualizadoEm time.Time

		err = db.QueryRow(`
			SELECT 
				id, nome, telefone, cpf, email, 
				endereco, complemento, bairro, cidade, estado, 
				cep, observacoes, canal_origem, criado_em, atualizado_em
			FROM clientes
			WHERE id = $1
		`, clienteID).Scan(
			&response.ID, &response.Nome, &response.Telefone, &cpf, &email,
			&endereco, &complemento, &bairro, &cidade, &estado,
			&cep, &observacoes, &canalOrigem, &criadoEm, &atualizadoEm,
		)
		if err != nil {
			http.Error(w, "Cliente atualizado, mas erro ao buscar detalhes: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Converter tipos nulos
		if cpf.Valid {
			response.CPF = cpf.String
		}
		if email.Valid {
			response.Email = email.String
		}
		if endereco.Valid {
			response.Endereco = endereco.String
		}
		if complemento.Valid {
			response.Complemento = complemento.String
		}
		if bairro.Valid {
			response.Bairro = bairro.String
		}
		if cidade.Valid {
			response.Cidade = cidade.String
		}
		if estado.Valid {
			response.Estado = estado.String
		}
		if cep.Valid {
			response.CEP = cep.String
		}
		if observacoes.Valid {
			response.Observacoes = observacoes.String
		}
		if canalOrigem.Valid {
			response.CanalOrigem = models.CanalOrigem(canalOrigem.String)
		}
		
		response.CriadoEm = criadoEm
		response.AtualizadoEm = atualizadoEm

		// Retornar resposta
		json.NewEncoder(w).Encode(response)
	}
}

// ExcluirClienteHandler exclui um cliente
func ExcluirClienteHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Verificar se o usuário está autenticado
		_, ok := middleware.ObterUsuarioID(r)
		if !ok {
			http.Error(w, "Usuário não autenticado", http.StatusUnauthorized)
			return
		}

		// Verificar permissões (apenas gerentes ou admin podem excluir clientes)
		perfil, perfilOk := middleware.ObterPerfilUsuario(r)
		if !perfilOk || !middleware.VerificarPerfil(perfil, "gerente") {
			http.Error(w, "Sem permissão para excluir clientes", http.StatusForbidden)
			return
		}

		// Configurar cabeçalhos
		w.Header().Set("Content-Type", "application/json")

		// Verificar método
		if r.Method != http.MethodDelete {
			http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
			return
		}

		// Extrair ID do cliente da URL
		parts := strings.Split(r.URL.Path, "/")
		if len(parts) < 4 {
			http.Error(w, "ID do cliente não fornecido", http.StatusBadRequest)
			return
		}
		clienteID, err := strconv.Atoi(parts[len(parts)-1])
		if err != nil {
			http.Error(w, "ID do cliente inválido", http.StatusBadRequest)
			return
		}

		// Verificar se o cliente existe
		var clienteExiste bool
		err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM clientes WHERE id = $1)", clienteID).Scan(&clienteExiste)
		if err != nil {
			http.Error(w, "Erro ao verificar cliente: "+err.Error(), http.StatusInternalServerError)
			return
		}
		if !clienteExiste {
			http.Error(w, "Cliente não encontrado", http.StatusNotFound)
			return
		}

		// Verificar se o cliente tem pedidos associados
		var temPedidos bool
		err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM pedidos WHERE cliente_id = $1)", clienteID).Scan(&temPedidos)
		if err != nil {
			http.Error(w, "Erro ao verificar pedidos: "+err.Error(), http.StatusInternalServerError)
			return
		}
		if temPedidos {
			http.Error(w, "Não é possível excluir o cliente pois existem pedidos associados a ele", http.StatusBadRequest)
			return
		}

		// Excluir cliente
		_, err = db.Exec("DELETE FROM clientes WHERE id = $1", clienteID)
		if err != nil {
			http.Error(w, "Erro ao excluir cliente: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Retornar resposta
       w.WriteHeader(http.StatusOK)
       json.NewEncoder(w).Encode(map[string]string{
           "mensagem": "Cliente excluído com sucesso",
       })
   }
}

// BuscarClientePorTelefoneHandler busca um cliente pelo telefone
func BuscarClientePorTelefoneHandler(db *sql.DB) http.HandlerFunc {
   return func(w http.ResponseWriter, r *http.Request) {
       // Verificar se o usuário está autenticado
       _, ok := middleware.ObterUsuarioID(r)
       if !ok {
           http.Error(w, "Usuário não autenticado", http.StatusUnauthorized)
           return
       }

       // Configurar cabeçalhos
       w.Header().Set("Content-Type", "application/json")

       // Verificar método
       if r.Method != http.MethodGet {
           http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
           return
       }

       // Obter telefone da query string
       telefone := r.URL.Query().Get("telefone")
       if telefone == "" {
           http.Error(w, "Telefone é obrigatório", http.StatusBadRequest)
           return
       }

       // Buscar cliente pelo telefone
       var cliente models.Cliente
       var cpf, email, endereco, complemento, bairro, cidade, estado, cep, observacoes sql.NullString
       var canalOrigem sql.NullString
       var criadoEm, atualizadoEm time.Time

       err := db.QueryRow(`
           SELECT 
               id, nome, telefone, cpf, email, 
               endereco, complemento, bairro, cidade, estado, 
               cep, observacoes, canal_origem, criado_em, atualizado_em
           FROM clientes
           WHERE telefone LIKE $1
           LIMIT 1
       `, "%"+telefone+"%").Scan(
           &cliente.ID, &cliente.Nome, &cliente.Telefone, &cpf, &email,
           &endereco, &complemento, &bairro, &cidade, &estado,
           &cep, &observacoes, &canalOrigem, &criadoEm, &atualizadoEm,
       )
       if err != nil {
           if err == sql.ErrNoRows {
               // Cliente não encontrado, retornar array vazio
               w.WriteHeader(http.StatusOK)
               json.NewEncoder(w).Encode([]interface{}{})
               return
           }
           http.Error(w, "Erro ao buscar cliente: "+err.Error(), http.StatusInternalServerError)
           return
       }

       // Converter tipos nulos
       if cpf.Valid {
           cliente.CPF = cpf.String
       }
       if email.Valid {
           cliente.Email = email.String
       }
       if endereco.Valid {
           cliente.Endereco = endereco.String
       }
       if complemento.Valid {
           cliente.Complemento = complemento.String
       }
       if bairro.Valid {
           cliente.Bairro = bairro.String
       }
       if cidade.Valid {
           cliente.Cidade = cidade.String
       }
       if estado.Valid {
           cliente.Estado = estado.String
       }
       if cep.Valid {
           cliente.CEP = cep.String
       }
       if observacoes.Valid {
           cliente.Observacoes = observacoes.String
       }
       if canalOrigem.Valid {
           cliente.CanalOrigem = models.CanalOrigem(canalOrigem.String)
       }
       
       cliente.CriadoEm = criadoEm
       cliente.AtualizadoEm = atualizadoEm

       // Retornar resposta como um array contendo o cliente
       w.WriteHeader(http.StatusOK)
       json.NewEncoder(w).Encode([]models.Cliente{cliente})
   }
}
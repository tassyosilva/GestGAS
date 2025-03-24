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

// ListarPedidosHandler retorna a lista de pedidos com paginação e filtros
func ListarPedidosHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Iniciando ListarPedidosHandler")
		// Verificar se o usuário está autenticado
		_, ok := middleware.ObterUsuarioID(r)
		if !ok {
			fmt.Println("ERRO: Usuário não autenticado")
			http.Error(w, "Usuário não autenticado", http.StatusUnauthorized)
			return
		}
		fmt.Println("Usuário autenticado com sucesso")

		// Configurar cabeçalhos
		w.Header().Set("Content-Type", "application/json")

		// Verificar método
		if r.Method != http.MethodGet {
			fmt.Println("ERRO: Método não permitido:", r.Method)
			http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
			return
		}
		fmt.Println("Método validado:", r.Method)

		// Parâmetros de consulta
		query := r.URL.Query()
		page, _ := strconv.Atoi(query.Get("page"))
		limit, _ := strconv.Atoi(query.Get("limit"))
		status := query.Get("status")
		clienteID := query.Get("cliente_id")
		dataInicio := query.Get("data_inicio")
		dataFim := query.Get("data_fim")
		fmt.Println("Parâmetros de consulta:", "page=", page, "limit=", limit, "status=", status, "cliente_id=", clienteID, "data_inicio=", dataInicio, "data_fim=", dataFim)

		if page < 1 {
			page = 1
		}
		if limit < 1 || limit > 100 {
			limit = 20
		}
		offset := (page - 1) * limit
		fmt.Println("Paginação configurada:", "page=", page, "limit=", limit, "offset=", offset)

		// Construir consulta SQL base
		sqlQuery := `
			SELECT p.id, p.cliente_id, p.atendente_id, p.entregador_id, p.status, 
			       p.forma_pagamento, p.valor_total, p.observacoes, p.endereco_entrega, 
				   p.canal_origem, p.data_entrega, p.criado_em, p.atualizado_em
			FROM pedidos p
			WHERE 1=1
		`
		var params []interface{}
		var whereConditions []string

		// Adicionar filtros se fornecidos
		if status != "" {
			whereConditions = append(whereConditions, "p.status = $"+strconv.Itoa(len(params)+1))
			params = append(params, status)
			fmt.Println("Adicionado filtro de status:", status)
		}
		if clienteID != "" {
			id, err := strconv.Atoi(clienteID)
			if err == nil {
				whereConditions = append(whereConditions, "p.cliente_id = $"+strconv.Itoa(len(params)+1))
				params = append(params, id)
				fmt.Println("Adicionado filtro de cliente_id:", id)
			} else {
				fmt.Println("AVISO: cliente_id não é um número válido:", clienteID)
			}
		}
		if dataInicio != "" {
			whereConditions = append(whereConditions, "p.criado_em >= $"+strconv.Itoa(len(params)+1))
			params = append(params, dataInicio)
			fmt.Println("Adicionado filtro de data_inicio:", dataInicio)
		}
		if dataFim != "" {
			whereConditions = append(whereConditions, "p.criado_em <= $"+strconv.Itoa(len(params)+1))
			params = append(params, dataFim)
			fmt.Println("Adicionado filtro de data_fim:", dataFim)
		}

		// Adicionar condições WHERE
		if len(whereConditions) > 0 {
			sqlQuery += " AND " + strings.Join(whereConditions, " AND ")
		}

		// Adicionar ordenação e paginação
		sqlQuery += " ORDER BY p.criado_em DESC LIMIT $" + strconv.Itoa(len(params)+1) + " OFFSET $" + strconv.Itoa(len(params)+2)
		params = append(params, limit, offset)

		fmt.Println("SQL Query:", sqlQuery)
		fmt.Println("SQL Params:", params)

		// Executar consulta
		fmt.Println("Executando consulta no banco de dados...")
		rows, err := db.Query(sqlQuery, params...)
		fmt.Println("Consulta executada, verificando erros...")
		if err != nil {
			fmt.Println("ERRO ao executar consulta:", err)
			http.Error(w, "Erro ao buscar pedidos: "+err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Println("Consulta executada com sucesso")
		defer rows.Close()

		// Processar resultados
		var pedidos []models.Pedido
		fmt.Println("Processando resultados da consulta...")
		for rows.Next() {
			var p models.Pedido
			var entregadorID sql.NullInt64
			var dataEntrega sql.NullTime
			var canalOrigem sql.NullString

			err := rows.Scan(
				&p.ID, &p.ClienteID, &p.AtendenteID, &entregadorID, &p.Status,
				&p.FormaPagamento, &p.ValorTotal, &p.Observacoes, &p.EnderecoEntrega,
				&canalOrigem, &dataEntrega, &p.CriadoEm, &p.AtualizadoEm,
			)
			if err != nil {
				fmt.Println("ERRO ao processar pedido:", err)
				http.Error(w, "Erro ao processar pedidos: "+err.Error(), http.StatusInternalServerError)
				return
			}

			// Converter tipos nulos
			if entregadorID.Valid {
				entregadorIDInt := int(entregadorID.Int64)
				p.EntregadorID = &entregadorIDInt
			}
			if dataEntrega.Valid {
				p.DataEntrega = &dataEntrega.Time
			}
			if canalOrigem.Valid {
				p.CanalOrigem = models.CanalOrigem(canalOrigem.String)
			}

			// Buscar itens do pedido
			fmt.Println("Buscando itens para o pedido ID:", p.ID)
			itens, err := buscarItensPedido(db, p.ID)
			if err != nil {
				fmt.Println("ERRO ao buscar itens do pedido:", err)
				http.Error(w, "Erro ao buscar itens do pedido: "+err.Error(), http.StatusInternalServerError)
				return
			}
			p.Itens = itens
			fmt.Println("Encontrados", len(itens), "itens para o pedido ID:", p.ID)

			pedidos = append(pedidos, p)
		}
		fmt.Println("Total de", len(pedidos), "pedidos processados")

		// Contar total de registros para paginação
		var total int
		countQuery := `
			SELECT COUNT(*) FROM pedidos p WHERE 1=1
		`
		if len(whereConditions) > 0 {
			countQuery += " AND " + strings.Join(whereConditions, " AND ")
		}
		fmt.Println("Executando consulta para contar total de registros...")
		fmt.Println("Count Query:", countQuery)
		fmt.Println("Count Params:", params[:len(params)-2])
		err = db.QueryRow(countQuery, params[:len(params)-2]...).Scan(&total)
		if err != nil {
			fmt.Println("ERRO ao contar pedidos:", err)
			http.Error(w, "Erro ao contar pedidos: "+err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Println("Total de registros encontrados:", total)

		// Montar resposta
		response := struct {
			Pedidos []models.Pedido `json:"pedidos"`
			Total   int             `json:"total"`
			Page    int             `json:"page"`
			Limit   int             `json:"limit"`
			Pages   int             `json:"pages"`
		}{
			Pedidos: pedidos,
			Total:   total,
			Page:    page,
			Limit:   limit,
			Pages:   (total + limit - 1) / limit,
		}

		// Retornar resposta
		fmt.Println("Retornando resposta com", len(pedidos), "pedidos")
		json.NewEncoder(w).Encode(response)
	}
}

// ObterPedidoHandler retorna detalhes de um pedido específico
func ObterPedidoHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Iniciando ObterPedidoHandler")
		// Verificar se o usuário está autenticado
		_, ok := middleware.ObterUsuarioID(r)
		if !ok {
			fmt.Println("ERRO: Usuário não autenticado")
			http.Error(w, "Usuário não autenticado", http.StatusUnauthorized)
			return
		}
		fmt.Println("Usuário autenticado com sucesso")

		// Configurar cabeçalhos
		w.Header().Set("Content-Type", "application/json")

		// Verificar método
		if r.Method != http.MethodGet {
			fmt.Println("ERRO: Método não permitido:", r.Method)
			http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
			return
		}
		fmt.Println("Método validado:", r.Method)

		// Extrair ID do pedido da URL
		parts := strings.Split(r.URL.Path, "/")
		if len(parts) < 4 {
			fmt.Println("ERRO: URL inválida, partes insuficientes:", parts)
			http.Error(w, "ID do pedido não fornecido", http.StatusBadRequest)
			return
		}
		pedidoID, err := strconv.Atoi(parts[len(parts)-1])
		if err != nil {
			fmt.Println("ERRO: ID do pedido inválido:", parts[len(parts)-1])
			http.Error(w, "ID do pedido inválido", http.StatusBadRequest)
			return
		}
		fmt.Println("ID do pedido extraído da URL:", pedidoID)

		// Buscar pedido
		fmt.Println("Buscando detalhes do pedido ID:", pedidoID)
		pedidoResp, err := buscarPedidoDetalhado(db, pedidoID)
		if err != nil {
			if err == sql.ErrNoRows {
				fmt.Println("ERRO: Pedido não encontrado ID:", pedidoID)
				http.Error(w, "Pedido não encontrado", http.StatusNotFound)
				return
			}
			fmt.Println("ERRO ao buscar pedido:", err)
			http.Error(w, "Erro ao buscar pedido: "+err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Println("Pedido encontrado com sucesso ID:", pedidoID)

		// Retornar resposta
		fmt.Println("Retornando detalhes do pedido ID:", pedidoID)
		json.NewEncoder(w).Encode(pedidoResp)
	}
}

// CriarPedidoHandler cria um novo pedido
func CriarPedidoHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Iniciando CriarPedidoHandler")
		// Verificar se o usuário está autenticado
		userID, ok := middleware.ObterUsuarioID(r)
		if !ok {
			fmt.Println("ERRO: Usuário não autenticado")
			http.Error(w, "Usuário não autenticado", http.StatusUnauthorized)
			return
		}
		fmt.Println("Usuário autenticado com sucesso, ID:", userID)

		// Verificar permissões (apenas atendentes ou admins podem criar pedidos)
		fmt.Println("Verificando permissões do usuário ID:", userID)
		perfil, perfilOk := middleware.ObterPerfilUsuario(r)
		if !perfilOk || !middleware.VerificarPerfil(perfil, "atendente") {
			fmt.Println("ERRO: Usuário sem permissão, ID:", userID)
			http.Error(w, "Sem permissão para criar pedidos", http.StatusForbidden)
			return
		}
		fmt.Println("Usuário tem permissão para criar pedidos")

		// Configurar cabeçalhos
		w.Header().Set("Content-Type", "application/json")

		// Verificar método
		if r.Method != http.MethodPost {
			fmt.Println("ERRO: Método não permitido:", r.Method)
			http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
			return
		}
		fmt.Println("Método validado:", r.Method)

		// Decodificar requisição
		fmt.Println("Decodificando corpo da requisição")
		var req models.NovoPedidoRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			fmt.Println("ERRO ao decodificar requisição:", err)
			http.Error(w, "Erro ao processar requisição: "+err.Error(), http.StatusBadRequest)
			return
		}
		fmt.Println("Requisição decodificada com sucesso, cliente ID:", req.ClienteID)

		// Validar dados
		fmt.Println("Validando dados do pedido")
		if req.ClienteID <= 0 {
			fmt.Println("ERRO: ID do cliente inválido ou não fornecido:", req.ClienteID)
			http.Error(w, "ID do cliente é obrigatório", http.StatusBadRequest)
			return
		}
		if req.EnderecoEntrega == "" {
			fmt.Println("ERRO: Endereço de entrega não fornecido")
			http.Error(w, "Endereço de entrega é obrigatório", http.StatusBadRequest)
			return
		}
		if len(req.Itens) == 0 {
			fmt.Println("ERRO: Pedido sem itens")
			http.Error(w, "Pedido deve conter pelo menos um item", http.StatusBadRequest)
			return
		}
		fmt.Println("Dados do pedido validados com sucesso")

		// Iniciar transação
		fmt.Println("Iniciando transação no banco de dados")
		tx, err := db.Begin()
		if err != nil {
			fmt.Println("ERRO ao iniciar transação:", err)
			http.Error(w, "Erro ao iniciar transação: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer func() {
			if err != nil {
				fmt.Println("Realizando rollback da transação devido a erro")
				tx.Rollback()
				return
			}
		}()

		// Verificar se cliente existe
		fmt.Println("Verificando se cliente existe, ID:", req.ClienteID)
		var clienteExiste bool
		err = tx.QueryRow("SELECT EXISTS(SELECT 1 FROM clientes WHERE id = $1)", req.ClienteID).Scan(&clienteExiste)
		if err != nil {
			fmt.Println("ERRO ao verificar cliente:", err)
			http.Error(w, "Erro ao verificar cliente: "+err.Error(), http.StatusInternalServerError)
			return
		}
		if !clienteExiste {
			fmt.Println("ERRO: Cliente não encontrado, ID:", req.ClienteID)
			http.Error(w, "Cliente não encontrado", http.StatusBadRequest)
			return
		}
		fmt.Println("Cliente encontrado com sucesso")

		// Calcular valor total e preparar itens
		var valorTotal float64
		var itensPedido []models.ItemPedido
		fmt.Println("Processando", len(req.Itens), "itens do pedido")
		for i, item := range req.Itens {
			fmt.Println("Processando item", i+1, "produto ID:", item.ProdutoID)
			// Buscar produto
			var produto struct {
				ID    int
				Nome  string
				Preco float64
			}
			err = tx.QueryRow("SELECT id, nome, preco FROM produtos WHERE id = $1", item.ProdutoID).Scan(
				&produto.ID, &produto.Nome, &produto.Preco,
			)
			if err != nil {
				if err == sql.ErrNoRows {
					fmt.Println("ERRO: Produto não encontrado, ID:", item.ProdutoID)
					http.Error(w, fmt.Sprintf("Produto ID %d não encontrado", item.ProdutoID), http.StatusBadRequest)
					return
				}
				fmt.Println("ERRO ao buscar produto:", err)
				http.Error(w, "Erro ao buscar produto: "+err.Error(), http.StatusInternalServerError)
				return
			}
			fmt.Println("Produto encontrado:", produto.Nome, "preço:", produto.Preco)

			// Verificar estoque
			fmt.Println("Verificando estoque para produto ID:", item.ProdutoID)
			var qtdEstoque int
			err = tx.QueryRow("SELECT quantidade FROM estoque WHERE produto_id = $1", item.ProdutoID).Scan(&qtdEstoque)
			if err != nil {
				fmt.Println("ERRO ao verificar estoque:", err)
				http.Error(w, "Erro ao verificar estoque: "+err.Error(), http.StatusInternalServerError)
				return
			}
			fmt.Println("Estoque disponível:", qtdEstoque, "unidades. Solicitado:", item.Quantidade)
			if qtdEstoque < item.Quantidade {
				fmt.Println("ERRO: Estoque insuficiente para produto", produto.Nome)
				http.Error(w, fmt.Sprintf("Estoque insuficiente para produto %s", produto.Nome), http.StatusBadRequest)
				return
			}

			// Calcular subtotal
			subtotal := float64(item.Quantidade) * produto.Preco
			valorTotal += subtotal
			fmt.Println("Subtotal calculado:", subtotal, "Valor total acumulado:", valorTotal)

			// Adicionar ao slice de itens
			itensPedido = append(itensPedido, models.ItemPedido{
				ProdutoID:     item.ProdutoID,
				NomeProduto:   produto.Nome,
				Quantidade:    item.Quantidade,
				PrecoUnitario: produto.Preco,
				Subtotal:      subtotal,
				RetornaBotija: item.RetornaBotija,
			})
		}

		// Inserir pedido
		fmt.Println("Inserindo pedido no banco de dados")
		var pedidoID int
		err = tx.QueryRow(`
			INSERT INTO pedidos
			(cliente_id, atendente_id, status, forma_pagamento, valor_total, observacoes, 
			endereco_entrega, canal_origem, criado_em, atualizado_em)
			VALUES
			($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW())
			RETURNING id
		`, req.ClienteID, userID, models.StatusNovo, req.FormaPagamento, valorTotal, 
		   req.Observacoes, req.EnderecoEntrega, req.CanalOrigem).Scan(&pedidoID)
		if err != nil {
			fmt.Println("ERRO ao inserir pedido:", err)
			http.Error(w, "Erro ao criar pedido: "+err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Println("Pedido inserido com sucesso, ID:", pedidoID)

		// Inserir itens do pedido
		fmt.Println("Inserindo", len(itensPedido), "itens do pedido")
		for i, item := range itensPedido {
			fmt.Println("Inserindo item", i+1, "produto:", item.NomeProduto)
			_, err = tx.Exec(`
				INSERT INTO itens_pedido
				(pedido_id, produto_id, quantidade, preco_unitario, subtotal, retorna_botija)
				VALUES
				($1, $2, $3, $4, $5, $6)
			`, pedidoID, item.ProdutoID, item.Quantidade, item.PrecoUnitario, item.Subtotal, item.RetornaBotija)
			if err != nil {
				fmt.Println("ERRO ao inserir item do pedido:", err)
				http.Error(w, "Erro ao inserir item do pedido: "+err.Error(), http.StatusInternalServerError)
				return
			}

			// Atualizar estoque
			fmt.Println("Atualizando estoque para produto ID:", item.ProdutoID)
			_, err = tx.Exec(`
				UPDATE estoque
				SET quantidade = quantidade - $1, atualizado_em = NOW()
				WHERE produto_id = $2
			`, item.Quantidade, item.ProdutoID)
			if err != nil {
				fmt.Println("ERRO ao atualizar estoque:", err)
				http.Error(w, "Erro ao atualizar estoque: "+err.Error(), http.StatusInternalServerError)
				return
			}

			// Registrar movimentação de estoque
			fmt.Println("Registrando movimentação de estoque para produto ID:", item.ProdutoID)
			_, err = tx.Exec(`
				INSERT INTO movimentacoes_estoque
				(produto_id, tipo, quantidade, usuario_id, pedido_id, criado_em)
				VALUES
				($1, $2, $3, $4, $5, NOW())
			`, item.ProdutoID, models.MovimentacaoSaida, item.Quantidade, userID, pedidoID)
			if err != nil {
				fmt.Println("ERRO ao registrar movimentação de estoque:", err)
				http.Error(w, "Erro ao registrar movimentação de estoque: "+err.Error(), http.StatusInternalServerError)
				return
			}
		}

		// Commit da transação
		fmt.Println("Realizando commit da transação")
		err = tx.Commit()
		if err != nil {
			fmt.Println("ERRO ao realizar commit da transação:", err)
			http.Error(w, "Erro ao finalizar transação: "+err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Println("Transação finalizada com sucesso")

		// Buscar pedido completo para resposta
		fmt.Println("Buscando detalhes completos do pedido para resposta")
		pedidoResp, err := buscarPedidoDetalhado(db, pedidoID)
		if err != nil {
			fmt.Println("ERRO ao buscar detalhes do pedido:", err)
			http.Error(w, "Pedido criado, mas erro ao buscar detalhes: "+err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Println("Detalhes do pedido obtidos com sucesso")

		// Retornar resposta
		fmt.Println("Retornando resposta com detalhes do pedido ID:", pedidoID)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(pedidoResp)
	}
}

// AtualizarStatusPedidoHandler atualiza o status de um pedido
func AtualizarStatusPedidoHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Iniciando AtualizarStatusPedidoHandler")
		// Verificar se o usuário está autenticado
		userID, ok := middleware.ObterUsuarioID(r)
		if !ok {
			fmt.Println("ERRO: Usuário não autenticado")
			http.Error(w, "Usuário não autenticado", http.StatusUnauthorized)
			return
		}
		fmt.Println("Usuário autenticado com sucesso, ID:", userID)

		// Obter perfil do usuário para verificações de permissão
		perfil, perfilOk := middleware.ObterPerfilUsuario(r)
		if !perfilOk {
			fmt.Println("ERRO: Não foi possível obter o perfil do usuário")
			http.Error(w, "Erro ao verificar permissões", http.StatusInternalServerError)
			return
		}
		fmt.Println("Perfil do usuário:", perfil)

		// Configurar cabeçalhos
		w.Header().Set("Content-Type", "application/json")

		// Verificar método
		if r.Method != http.MethodPut && r.Method != http.MethodPatch {
			fmt.Println("ERRO: Método não permitido:", r.Method)
			http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
			return
		}
		fmt.Println("Método validado:", r.Method)

		// Extrair ID do pedido da URL
		parts := strings.Split(r.URL.Path, "/")
		if len(parts) < 5 || parts[len(parts)-1] != "status" {
			fmt.Println("ERRO: URL inválida, partes:", parts)
			http.Error(w, "URL inválida", http.StatusBadRequest)
			return
		}
		pedidoID, err := strconv.Atoi(parts[len(parts)-2])
		if err != nil {
			fmt.Println("ERRO: ID do pedido inválido:", parts[len(parts)-2])
			http.Error(w, "ID do pedido inválido", http.StatusBadRequest)
			return
		}
		fmt.Println("ID do pedido extraído da URL:", pedidoID)

		// Verificar se o pedido existe
		fmt.Println("Verificando se pedido existe, ID:", pedidoID)
		var pedidoExiste bool
		var statusAtual models.StatusPedido
		err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM pedidos WHERE id = $1), status FROM pedidos WHERE id = $1", pedidoID).Scan(&pedidoExiste, &statusAtual)
		if err != nil {
			fmt.Println("ERRO ao verificar pedido:", err)
			http.Error(w, "Erro ao verificar pedido: "+err.Error(), http.StatusInternalServerError)
			return
		}
		if !pedidoExiste {
			fmt.Println("ERRO: Pedido não encontrado, ID:", pedidoID)
			http.Error(w, "Pedido não encontrado", http.StatusNotFound)
			return
		}
		fmt.Println("Pedido encontrado, status atual:", statusAtual)

		// Decodificar requisição
		fmt.Println("Decodificando corpo da requisição")
		var req models.AtualizarStatusRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			fmt.Println("ERRO ao decodificar requisição:", err)
			http.Error(w, "Erro ao processar requisição: "+err.Error(), http.StatusBadRequest)
			return
		}
		fmt.Println("Requisição decodificada, novo status solicitado:", req.Status)

		// Validar status
		if req.Status == "" {
			fmt.Println("ERRO: Status não fornecido")
			http.Error(w, "Status é obrigatório", http.StatusBadRequest)
			return
		}

		// Validar transição de status
		fmt.Println("Validando transição de status:", statusAtual, "->", req.Status)
		if !validarTransicaoStatus(statusAtual, req.Status) {
			fmt.Println("ERRO: Transição de status inválida:", statusAtual, "->", req.Status)
			http.Error(w, fmt.Sprintf("Transição de status inválida: %s -> %s", statusAtual, req.Status), http.StatusBadRequest)
			return
		}
		fmt.Println("Transição de status validada com sucesso")

		// Verificar permissões baseado no status
		fmt.Println("Verificando permissões para atualização de status")
		switch req.Status {
		case models.StatusCancelado, models.StatusEmPreparo:
			// Atendentes e acima podem cancelar ou preparar pedidos
			if !middleware.VerificarPerfil(perfil, "atendente") {
				fmt.Println("ERRO: Usuário sem permissão para esta atualização, ID:", userID)
				http.Error(w, "Sem permissão para esta atualização", http.StatusForbidden)
				return
			}
		case models.StatusEmEntrega, models.StatusEntregue, models.StatusFinalizado:
			// Verificar se o entregador foi definido
			if req.EntregadorID == nil && req.Status == models.StatusEmEntrega {
				fmt.Println("ERRO: Entregador não definido para status em entrega")
				http.Error(w, "É necessário definir um entregador para iniciar a entrega", http.StatusBadRequest)
				return
			}
		}
		fmt.Println("Permissões validadas com sucesso")

		// Iniciar transação
		fmt.Println("Iniciando transação no banco de dados")
		tx,err := db.Begin()
		if err != nil {
			fmt.Println("ERRO ao iniciar transação:", err)
			http.Error(w, "Erro ao iniciar transação: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer func() {
			if err != nil {
				fmt.Println("Realizando rollback da transação devido a erro")
				tx.Rollback()
				return
			}
		}()

		// Atualizar status do pedido
		var updateQuery string
		var args []interface{}
		fmt.Println("Preparando query de atualização")

		// Preparar query de atualização baseada nos campos fornecidos
		if req.Status == models.StatusEmEntrega && req.EntregadorID != nil {
			// Se estiver iniciando entrega, atualizar entregador também
			updateQuery = `
				UPDATE pedidos 
				SET status = $1, entregador_id = $2, atualizado_em = NOW() 
				WHERE id = $3
			`
			args = []interface{}{req.Status, *req.EntregadorID, pedidoID}
			fmt.Println("Atualizando para status Em Entrega com entregador ID:", *req.EntregadorID)
		} else if req.Status == models.StatusEntregue || req.Status == models.StatusFinalizado {
			// Se estiver finalizando entrega, registrar data de entrega
			dataEntrega := time.Now()
			if req.DataEntrega != nil {
				dataEntrega = *req.DataEntrega
			}
			updateQuery = `
				UPDATE pedidos 
				SET status = $1, data_entrega = $2, atualizado_em = NOW() 
				WHERE id = $3
			`
			args = []interface{}{req.Status, dataEntrega, pedidoID}
			fmt.Println("Atualizando para status", req.Status, "com data de entrega:", dataEntrega)

			// Verificar se existem botijas retornadas pelo cliente
			fmt.Println("Verificando existência de botijas retornadas")
			var temnBotijasRetornadas bool
			err = tx.QueryRow(`
				SELECT EXISTS(
					SELECT 1 FROM itens_pedido 
					WHERE pedido_id = $1 AND retorna_botija = TRUE
				)
			`, pedidoID).Scan(&temnBotijasRetornadas)
			if err != nil {
				fmt.Println("ERRO ao verificar botijas retornadas:", err)
				http.Error(w, "Erro ao verificar botijas retornadas: "+err.Error(), http.StatusInternalServerError)
				return
			}
			fmt.Println("Botijas retornadas:", temnBotijasRetornadas)

			// Se houver botijas retornadas, atualizar estoque
			if temnBotijasRetornadas {
				fmt.Println("Processando retorno de botijas")
				// Buscar todas as botijas retornadas
				rows, err := tx.Query(`
					SELECT ip.produto_id, ip.quantidade 
					FROM itens_pedido ip
					JOIN produtos p ON ip.produto_id = p.id
					WHERE ip.pedido_id = $1 AND ip.retorna_botija = TRUE
					AND p.categoria LIKE 'botija_gas%'
				`, pedidoID)
				if err != nil {
					fmt.Println("ERRO ao buscar botijas retornadas:", err)
					http.Error(w, "Erro ao buscar botijas retornadas: "+err.Error(), http.StatusInternalServerError)
					return
				}
				defer rows.Close()

				for rows.Next() {
					var produtoID, quantidade int
					err = rows.Scan(&produtoID, &quantidade)
					if err != nil {
						fmt.Println("ERRO ao processar botijas retornadas:", err)
						http.Error(w, "Erro ao processar botijas retornadas: "+err.Error(), http.StatusInternalServerError)
						return
					}
					fmt.Println("Botija retornada: produto ID:", produtoID, "quantidade:", quantidade)

					// Atualizar estoque de botijas vazias
					fmt.Println("Atualizando estoque de botijas vazias")
					_, err = tx.Exec(`
						UPDATE estoque 
						SET botijas_vazias = botijas_vazias + $1, atualizado_em = NOW() 
						WHERE produto_id = $2
					`, quantidade, produtoID)
					if err != nil {
						fmt.Println("ERRO ao atualizar estoque de botijas vazias:", err)
						http.Error(w, "Erro ao atualizar estoque de botijas vazias: "+err.Error(), http.StatusInternalServerError)
						return
					}

					// Registrar movimentação de estoque
					fmt.Println("Registrando movimentação de estoque para botijas vazias")
					_, err = tx.Exec(`
						INSERT INTO movimentacoes_estoque
						(produto_id, tipo, quantidade, usuario_id, pedido_id, criado_em)
						VALUES
						($1, $2, $3, $4, $5, NOW())
					`, produtoID, models.MovimentacaoBotijasVazias, quantidade, userID, pedidoID)
					if err != nil {
						fmt.Println("ERRO ao registrar movimentação de botijas vazias:", err)
						http.Error(w, "Erro ao registrar movimentação de botijas vazias: "+err.Error(), http.StatusInternalServerError)
						return
					}
				}
			}
		} else {
			// Atualização padrão somente do status
			updateQuery = `
				UPDATE pedidos 
				SET status = $1, atualizado_em = NOW() 
				WHERE id = $2
			`
			args = []interface{}{req.Status, pedidoID}
			fmt.Println("Atualizando para status:", req.Status)
		}

		// Executar atualização
		fmt.Println("Executando query de atualização:", updateQuery)
		fmt.Println("Parâmetros:", args)
		_, err = tx.Exec(updateQuery, args...)
		if err != nil {
			fmt.Println("ERRO ao atualizar status do pedido:", err)
			http.Error(w, "Erro ao atualizar status do pedido: "+err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Println("Status atualizado com sucesso")

		// Commit da transação
		fmt.Println("Realizando commit da transação")
		err = tx.Commit()
		if err != nil {
			fmt.Println("ERRO ao realizar commit da transação:", err)
			http.Error(w, "Erro ao finalizar transação: "+err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Println("Transação finalizada com sucesso")

		// Buscar pedido atualizado para resposta
		fmt.Println("Buscando detalhes atualizados do pedido para resposta")
		pedidoResp, err := buscarPedidoDetalhado(db, pedidoID)
		if err != nil {
			fmt.Println("ERRO ao buscar detalhes do pedido:", err)
			http.Error(w, "Status atualizado, mas erro ao buscar detalhes: "+err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Println("Detalhes do pedido obtidos com sucesso")

		// Retornar resposta
		fmt.Println("Retornando resposta com detalhes do pedido atualizado")
		json.NewEncoder(w).Encode(pedidoResp)
	}
}

// buscarItensPedido é uma função auxiliar para buscar os itens de um pedido
func buscarItensPedido(db *sql.DB, pedidoID int) ([]models.ItemPedido, error) {
	fmt.Println("Iniciando buscarItensPedido para pedido ID:", pedidoID)
	rows, err := db.Query(`
		SELECT ip.id, ip.pedido_id, ip.produto_id, p.nome, ip.quantidade, 
		ip.preco_unitario, ip.subtotal, ip.retorna_botija
		FROM itens_pedido ip
		JOIN produtos p ON ip.produto_id = p.id
		WHERE ip.pedido_id = $1
	`, pedidoID)
	if err != nil {
		fmt.Println("ERRO ao executar query para buscar itens:", err)
		return nil, err
	}
	defer rows.Close()

	var itens []models.ItemPedido
	for rows.Next() {
		var item models.ItemPedido
		var retornaBotija sql.NullBool

		err := rows.Scan(
			&item.ID, &item.PedidoID, &item.ProdutoID, &item.NomeProduto,
			&item.Quantidade, &item.PrecoUnitario, &item.Subtotal, &retornaBotija,
		)
		if err != nil {
			fmt.Println("ERRO ao processar item do pedido:", err)
			return nil, err
		}

		if retornaBotija.Valid {
			item.RetornaBotija = retornaBotija.Bool
		}

		itens = append(itens, item)
	}
	fmt.Println("Encontrados", len(itens), "itens para o pedido ID:", pedidoID)

	return itens, nil
}

// buscarPedidoDetalhado é uma função auxiliar para buscar um pedido com todos os detalhes necessários
func buscarPedidoDetalhado(db *sql.DB, pedidoID int) (models.PedidoResponse, error) {
	fmt.Println("Iniciando buscarPedidoDetalhado para pedido ID:", pedidoID)
	var resp models.PedidoResponse

	// Buscar dados do pedido com cliente, atendente e entregador
	var query = `
		SELECT 
			p.id, p.cliente_id, c.nome, c.telefone,
			p.atendente_id, a.nome, a.perfil,
			p.entregador_id, e.nome, e.perfil,
			p.status, p.forma_pagamento, p.valor_total,
			p.observacoes, p.endereco_entrega,
			p.canal_origem, p.data_entrega,
			p.criado_em, p.atualizado_em
		FROM pedidos p
		JOIN clientes c ON p.cliente_id = c.id
		JOIN usuarios a ON p.atendente_id = a.id
		LEFT JOIN usuarios e ON p.entregador_id = e.id
		WHERE p.id = $1
	`
	fmt.Println("Executando query principal para buscar pedido detalhado")

	var entregadorID sql.NullInt64
	var entregadorNome, entregadorPerfil sql.NullString
	var dataEntrega sql.NullTime
	var observacoes, canalOrigem sql.NullString

	err := db.QueryRow(query, pedidoID).Scan(
		&resp.ID, &resp.Cliente.ID, &resp.Cliente.Nome, &resp.Cliente.Telefone,
		&resp.Atendente.ID, &resp.Atendente.Nome, &resp.Atendente.Perfil,
		&entregadorID, &entregadorNome, &entregadorPerfil,
		&resp.Status, &resp.FormaPagamento, &resp.ValorTotal,
		&observacoes, &resp.EnderecoEntrega,
		&canalOrigem, &dataEntrega,
		&resp.CriadoEm, &resp.AtualizadoEm,
	)
	if err != nil {
		fmt.Println("ERRO ao buscar dados do pedido:", err)
		return resp, err
	}
	fmt.Println("Dados principais do pedido obtidos com sucesso")

	// Converter tipos nulos
	if entregadorID.Valid {
		entregador := models.UsuarioBasico{
			ID:     int(entregadorID.Int64),
			Nome:   entregadorNome.String,
			Perfil: entregadorPerfil.String,
		}
		resp.Entregador = &entregador
		fmt.Println("Entregador atribuído ao pedido:", entregadorNome.String)
	}
	if dataEntrega.Valid {
		resp.DataEntrega = &dataEntrega.Time
		fmt.Println("Data de entrega registrada:", dataEntrega.Time)
	}
	if observacoes.Valid {
		resp.Observacoes = observacoes.String
	}
	if canalOrigem.Valid {
		resp.CanalOrigem = models.CanalOrigem(canalOrigem.String)
		fmt.Println("Canal de origem:", canalOrigem.String)
	}

	// Buscar itens do pedido
	fmt.Println("Buscando itens do pedido")
	itens, err := buscarItensPedido(db, pedidoID)
	if err != nil {
		fmt.Println("ERRO ao buscar itens do pedido:", err)
		return resp, err
	}
	resp.Itens = itens
	fmt.Println("Itens do pedido obtidos com sucesso, total:", len(itens))

	return resp, nil
}

// validarTransicaoStatus verifica se uma transição de status é válida
func validarTransicaoStatus(atual, nova models.StatusPedido) bool {
	fmt.Println("Validando transição de status:", atual, "->", nova)
	switch atual {
	case models.StatusNovo:
		return nova == models.StatusEmPreparo || nova == models.StatusCancelado
	case models.StatusEmPreparo:
		return nova == models.StatusEmEntrega || nova == models.StatusCancelado
	case models.StatusEmEntrega:
		return nova == models.StatusEntregue || nova == models.StatusCancelado
	case models.StatusEntregue:
		return nova == models.StatusFinalizado
	case models.StatusCancelado, models.StatusFinalizado:
		return false // Estados finais, não podem ser alterados
	default:
		return false
	}
}
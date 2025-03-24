package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/tassyosilva/GestGAS/internal/middleware"
	"github.com/tassyosilva/GestGAS/internal/models"
)

// ListarEstoqueHandler retorna a lista de itens no estoque
func ListarEstoqueHandler(db *sql.DB) http.HandlerFunc {
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
		categoria := query.Get("categoria")
		alertas := query.Get("alertas") == "true"
		botijasVazias := query.Get("botijas_vazias") == "true"

		// Construir consulta SQL
		sqlQuery := `
			SELECT e.id, e.produto_id, p.nome, p.categoria, e.quantidade, 
			       e.botijas_vazias, e.botijas_emprestadas, e.alerta_minimo, e.atualizado_em
			FROM estoque e
			JOIN produtos p ON e.produto_id = p.id
			WHERE 1=1
		`
		var params []interface{}

		// Adicionar filtros
		if categoria != "" {
			sqlQuery += fmt.Sprintf(" AND p.categoria = $%d", len(params)+1)
			params = append(params, categoria)
		}

		if alertas {
			sqlQuery += fmt.Sprintf(" AND e.quantidade <= e.alerta_minimo")
		}

		if botijasVazias {
			sqlQuery += fmt.Sprintf(" AND e.botijas_vazias > 0")
		}

		sqlQuery += " ORDER BY p.categoria, p.nome"

		// Executar consulta
		rows, err := db.Query(sqlQuery, params...)
		if err != nil {
			http.Error(w, "Erro ao buscar estoque: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		// Processar resultados
		var estoque []models.EstoqueResponse
		for rows.Next() {
			var e models.EstoqueResponse
			var botijasVazias, botijasEmprestadas, alertaMinimo sql.NullInt64

			err := rows.Scan(
				&e.ID, &e.ProdutoID, &e.NomeProduto, &e.Categoria, &e.Quantidade,
				&botijasVazias, &botijasEmprestadas, &alertaMinimo, &e.AtualizadoEm,
			)
			if err != nil {
				http.Error(w, "Erro ao processar estoque: "+err.Error(), http.StatusInternalServerError)
				return
			}

			// Converter tipos nulos
			if botijasVazias.Valid {
				e.BotijasVazias = int(botijasVazias.Int64)
			}
			if botijasEmprestadas.Valid {
				e.BotijasEmprestadas = int(botijasEmprestadas.Int64)
			}
			if alertaMinimo.Valid {
				e.AlertaMinimo = int(alertaMinimo.Int64)
			}

			// Determinar status do estoque
			e.Status = "normal"
			if alertaMinimo.Valid && e.Quantidade <= int(alertaMinimo.Int64) {
				if e.Quantidade <= int(alertaMinimo.Int64/2) {
					e.Status = "critico"
				} else {
					e.Status = "baixo"
				}
			}

			estoque = append(estoque, e)
		}

		// Retornar resposta
		json.NewEncoder(w).Encode(estoque)
	}
}

// ObterEstoqueItemHandler retorna detalhes de um item específico do estoque
func ObterEstoqueItemHandler(db *sql.DB) http.HandlerFunc {
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

		// Extrair ID do produto da URL
		parts := strings.Split(r.URL.Path, "/")
		if len(parts) < 4 {
			http.Error(w, "ID do produto não fornecido", http.StatusBadRequest)
			return
		}
		produtoID, err := strconv.Atoi(parts[len(parts)-1])
		if err != nil {
			http.Error(w, "ID do produto inválido", http.StatusBadRequest)
			return
		}

		// Buscar item do estoque
		var e models.EstoqueResponse
		var botijasVazias, botijasEmprestadas, alertaMinimo sql.NullInt64

		err = db.QueryRow(`
			SELECT e.id, e.produto_id, p.nome, p.categoria, e.quantidade, 
			       e.botijas_vazias, e.botijas_emprestadas, e.alerta_minimo, e.atualizado_em
			FROM estoque e
			JOIN produtos p ON e.produto_id = p.id
			WHERE e.produto_id = $1
		`, produtoID).Scan(
			&e.ID, &e.ProdutoID, &e.NomeProduto, &e.Categoria, &e.Quantidade,
			&botijasVazias, &botijasEmprestadas, &alertaMinimo, &e.AtualizadoEm,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Item de estoque não encontrado", http.StatusNotFound)
				return
			}
			http.Error(w, "Erro ao buscar item de estoque: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Converter tipos nulos
		if botijasVazias.Valid {
			e.BotijasVazias = int(botijasVazias.Int64)
		}
		if botijasEmprestadas.Valid {
			e.BotijasEmprestadas = int(botijasEmprestadas.Int64)
		}
		if alertaMinimo.Valid {
			e.AlertaMinimo = int(alertaMinimo.Int64)
		}

		// Determinar status do estoque
		e.Status = "normal"
		if alertaMinimo.Valid && e.Quantidade <= int(alertaMinimo.Int64) {
			if e.Quantidade <= int(alertaMinimo.Int64/2) {
				e.Status = "critico"
			} else {
				e.Status = "baixo"
			}
		}

		// Buscar histórico de movimentações (últimas 10)
		rows, err := db.Query(`
			SELECT m.id, m.tipo, m.quantidade, m.observacoes, 
			       m.criado_em, u.nome as usuario_nome
			FROM movimentacoes_estoque m
			JOIN usuarios u ON m.usuario_id = u.id
			WHERE m.produto_id = $1
			ORDER BY m.criado_em DESC
			LIMIT 10
		`, produtoID)
		if err != nil {
			http.Error(w, "Erro ao buscar movimentações: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		type Movimentacao struct {
			ID          int                      `json:"id"`
			Tipo        models.TipoMovimentacao  `json:"tipo"`
			Quantidade  int                      `json:"quantidade"`
			Observacoes string                   `json:"observacoes,omitempty"`
			CriadoEm    string                   `json:"criado_em"`
			Usuario     string                   `json:"usuario"`
		}

		var movimentacoes []Movimentacao
		for rows.Next() {
			var m Movimentacao
			var observacoes sql.NullString

			err := rows.Scan(
				&m.ID, &m.Tipo, &m.Quantidade, &observacoes,
				&m.CriadoEm, &m.Usuario,
			)
			if err != nil {
				http.Error(w, "Erro ao processar movimentações: "+err.Error(), http.StatusInternalServerError)
				return
			}

			if observacoes.Valid {
				m.Observacoes = observacoes.String
			}

			movimentacoes = append(movimentacoes, m)
		}

		// Montar resposta
		response := struct {
			Estoque       models.EstoqueResponse `json:"estoque"`
			Movimentacoes []Movimentacao        `json:"movimentacoes"`
		}{
			Estoque:       e,
			Movimentacoes: movimentacoes,
		}

		// Retornar resposta
		json.NewEncoder(w).Encode(response)
	}
}

// ListarAlertasEstoqueHandler retorna a lista de produtos com estoque baixo
func ListarAlertasEstoqueHandler(db *sql.DB) http.HandlerFunc {
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

		// Buscar produtos com estoque baixo
		rows, err := db.Query(`
			SELECT p.id, p.nome, e.quantidade, e.alerta_minimo
			FROM estoque e
			JOIN produtos p ON e.produto_id = p.id
			WHERE e.quantidade <= e.alerta_minimo AND e.alerta_minimo > 0
			ORDER BY (e.quantidade::float / e.alerta_minimo) ASC
		`)
		if err != nil {
			http.Error(w, "Erro ao buscar alertas de estoque: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		// Processar resultados
		var alertas []models.EstoqueAlertaResponse
		for rows.Next() {
			var a models.EstoqueAlertaResponse
			err := rows.Scan(&a.ProdutoID, &a.NomeProduto, &a.Quantidade, &a.AlertaMinimo)
			if err != nil {
				http.Error(w, "Erro ao processar alertas de estoque: "+err.Error(), http.StatusInternalServerError)
				return
			}

			// Determinar status
			if a.Quantidade <= a.AlertaMinimo/2 {
				a.Status = "critico"
			} else {
				a.Status = "baixo"
			}

			alertas = append(alertas, a)
		}

		// Retornar resposta
		json.NewEncoder(w).Encode(alertas)
	}
}

// AtualizarEstoqueHandler atualiza a quantidade em estoque de um produto
func AtualizarEstoqueHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Verificar se o usuário está autenticado
		userID, ok := middleware.ObterUsuarioID(r)
		if !ok {
			http.Error(w, "Usuário não autenticado", http.StatusUnauthorized)
			return
		}

		// Verificar permissões (apenas gerentes ou admins podem atualizar estoque manualmente)
		perfil, perfilOk := middleware.ObterPerfilUsuario(r)
		if !perfilOk || !middleware.VerificarPerfil(perfil, "gerente") {
			http.Error(w, "Sem permissão para atualizar estoque", http.StatusForbidden)
			return
		}

		// Configurar cabeçalhos
		w.Header().Set("Content-Type", "application/json")

		// Verificar método
		if r.Method != http.MethodPut && r.Method != http.MethodPatch {
			http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
			return
		}

		// Extrair ID do produto da URL
		parts := strings.Split(r.URL.Path, "/")
		if len(parts) < 4 {
			http.Error(w, "ID do produto não fornecido", http.StatusBadRequest)
			return
		}
		produtoID, err := strconv.Atoi(parts[len(parts)-1])
		if err != nil {
			http.Error(w, "ID do produto inválido", http.StatusBadRequest)
			return
		}

		// Decodificar requisição
		var req models.MovimentacaoEstoqueRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Erro ao processar requisição: "+err.Error(), http.StatusBadRequest)
			return
		}

		// Validar dados
		if req.Quantidade <= 0 {
			http.Error(w, "Quantidade deve ser maior que zero", http.StatusBadRequest)
			return
		}

		// Verificar se o produto existe
		var produtoExiste bool
		var produtoNome string
		err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM produtos WHERE id = $1), nome FROM produtos WHERE id = $1", produtoID).Scan(&produtoExiste, &produtoNome)
		if err != nil {
			http.Error(w, "Erro ao verificar produto: "+err.Error(), http.StatusInternalServerError)
			return
		}
		if !produtoExiste {
			http.Error(w, "Produto não encontrado", http.StatusNotFound)
			return
		}

		// Iniciar transação
		tx, err := db.Begin()
		if err != nil {
			http.Error(w, "Erro ao iniciar transação: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer func() {
			if err != nil {
				tx.Rollback()
				return
			}
		}()

		// Atualizar estoque conforme o tipo de movimentação
		var query string
		switch req.Tipo {
		case models.MovimentacaoEntrada:
			query = `
				UPDATE estoque 
				SET quantidade = quantidade + $1, atualizado_em = NOW() 
				WHERE produto_id = $2
			`
		case models.MovimentacaoSaida:
			// Verificar se há estoque suficiente
			var qtdEstoque int
			err = tx.QueryRow("SELECT quantidade FROM estoque WHERE produto_id = $1", produtoID).Scan(&qtdEstoque)
			if err != nil {
				http.Error(w, "Erro ao verificar estoque: "+err.Error(), http.StatusInternalServerError)
				return
			}
			if qtdEstoque < req.Quantidade {
				http.Error(w, fmt.Sprintf("Estoque insuficiente para o produto %s", produtoNome), http.StatusBadRequest)
				return
			}
			query = `
				UPDATE estoque 
				SET quantidade = quantidade - $1, atualizado_em = NOW() 
				WHERE produto_id = $2
			`
		case models.MovimentacaoAjuste:
			// Ajuste direto na quantidade
			query = `
				UPDATE estoque 
				SET quantidade = $1, atualizado_em = NOW() 
				WHERE produto_id = $2
			`
		case models.MovimentacaoBotijasVazias:
			// Atualizar contagem de botijas vazias
			query = `
				UPDATE estoque 
				SET botijas_vazias = botijas_vazias + $1, atualizado_em = NOW() 
				WHERE produto_id = $2
			`
		case models.MovimentacaoEmprestimo:
			// Verificar se há botijas vazias suficientes
			var qtdBotijasVazias int
			err = tx.QueryRow("SELECT botijas_vazias FROM estoque WHERE produto_id = $1", produtoID).Scan(&qtdBotijasVazias)
			if err != nil {
				http.Error(w, "Erro ao verificar botijas vazias: "+err.Error(), http.StatusInternalServerError)
				return
			}
			if qtdBotijasVazias < req.Quantidade {
				http.Error(w, fmt.Sprintf("Botijas vazias insuficientes para empréstimo do produto %s", produtoNome), http.StatusBadRequest)
				return
			}
			query = `
				UPDATE estoque 
				SET botijas_vazias = botijas_vazias - $1, 
					botijas_emprestadas = botijas_emprestadas + $1, 
					atualizado_em = NOW() 
				WHERE produto_id = $2
			`
		case models.MovimentacaoDevolucaoEmprestimo:
			// Verificar se há botijas emprestadas
			var qtdBotijasEmprestadas int
			err = tx.QueryRow("SELECT botijas_emprestadas FROM estoque WHERE produto_id = $1", produtoID).Scan(&qtdBotijasEmprestadas)
			if err != nil {
				http.Error(w, "Erro ao verificar botijas emprestadas: "+err.Error(), http.StatusInternalServerError)
				return
			}
			if qtdBotijasEmprestadas < req.Quantidade {
				http.Error(w, fmt.Sprintf("Quantidade de botijas emprestadas do produto %s menor que a quantidade informada", produtoNome), http.StatusBadRequest)
				return
			}
			query = `
				UPDATE estoque 
				SET botijas_emprestadas = botijas_emprestadas - $1, 
					quantidade = quantidade + $1, 
					atualizado_em = NOW() 
				WHERE produto_id = $2
			`
		default:
			http.Error(w, "Tipo de movimentação inválido", http.StatusBadRequest)
			return
		}

		// Executar atualização
		_, err = tx.Exec(query, req.Quantidade, produtoID)
		if err != nil {
			http.Error(w, "Erro ao atualizar estoque: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Registrar movimentação
		_, err = tx.Exec(`
			INSERT INTO movimentacoes_estoque
			(produto_id, tipo, quantidade, observacoes, usuario_id, pedido_id, criado_em)
			VALUES
			($1, $2, $3, $4, $5, $6, NOW())
		`, produtoID, req.Tipo, req.Quantidade, req.Observacoes, userID, req.PedidoID)
		if err != nil {
			http.Error(w, "Erro ao registrar movimentação: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Commit da transação
		err = tx.Commit()
		if err != nil {
			http.Error(w, "Erro ao finalizar transação: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Buscar informações atualizadas do estoque
		var e models.EstoqueResponse
		var botijasVazias, botijasEmprestadas, alertaMinimo sql.NullInt64

		err = db.QueryRow(`
			SELECT e.id, e.produto_id, p.nome, p.categoria, e.quantidade, 
			       e.botijas_vazias, e.botijas_emprestadas, e.alerta_minimo, e.atualizado_em
			FROM estoque e
			JOIN produtos p ON e.produto_id = p.id
			WHERE e.produto_id = $1
		`, produtoID).Scan(
			&e.ID, &e.ProdutoID, &e.NomeProduto, &e.Categoria, &e.Quantidade,
			&botijasVazias, &botijasEmprestadas, &alertaMinimo, &e.AtualizadoEm,
		)
		if err != nil {
			http.Error(w, "Estoque atualizado, mas erro ao buscar informações: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Converter tipos nulos
		if botijasVazias.Valid {
			e.BotijasVazias = int(botijasVazias.Int64)
		}
		if botijasEmprestadas.Valid {
			e.BotijasEmprestadas = int(botijasEmprestadas.Int64)
		}
		if alertaMinimo.Valid {
			e.AlertaMinimo = int(alertaMinimo.Int64)
		}

		// Determinar status do estoque
		e.Status = "normal"
		if alertaMinimo.Valid && e.Quantidade <= int(alertaMinimo.Int64) {
			if e.Quantidade <= int(alertaMinimo.Int64/2) {
				e.Status = "critico"
			} else {
				e.Status = "baixo"
			}
		}

		// Retornar resposta
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(e)
	}
}

// AtualizarAlertaMinimoHandler atualiza o alerta mínimo de estoque de um produto
func AtualizarAlertaMinimoHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Verificar se o usuário está autenticado
		_, ok := middleware.ObterUsuarioID(r)
		if !ok {
			http.Error(w, "Usuário não autenticado", http.StatusUnauthorized)
			return
		}

		// Verificar permissões (apenas gerentes ou admins podem atualizar alerta mínimo)
		perfil, perfilOk := middleware.ObterPerfilUsuario(r)
		if !perfilOk || !middleware.VerificarPerfil(perfil, "gerente") {
			http.Error(w, "Sem permissão para atualizar alerta mínimo", http.StatusForbidden)
			return
		}

		// Configurar cabeçalhos
		w.Header().Set("Content-Type", "application/json")

		// Verificar método
		if r.Method != http.MethodPut && r.Method != http.MethodPatch {
			http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
			return
		}

		// Extrair ID do produto da URL
		parts := strings.Split(r.URL.Path, "/")
		if len(parts) < 5 || parts[len(parts)-1] != "alerta" {
			http.Error(w, "URL inválida", http.StatusBadRequest)
			return
		}
		produtoID, err := strconv.Atoi(parts[len(parts)-2])
		if err != nil {
			http.Error(w, "ID do produto inválido", http.StatusBadRequest)
			return
		}

		// Decodificar requisição
		var req struct {
			AlertaMinimo int `json:"alerta_minimo"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Erro ao processar requisição: "+err.Error(), http.StatusBadRequest)
			return
		}

		// Validar dados
		if req.AlertaMinimo < 0 {
			http.Error(w, "Alerta mínimo não pode ser negativo", http.StatusBadRequest)
			return
		}

		// Verificar se o produto existe
		var produtoExiste bool
		err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM produtos WHERE id = $1)", produtoID).Scan(&produtoExiste)
		if err != nil {
			http.Error(w, "Erro ao verificar produto: "+err.Error(), http.StatusInternalServerError)
			return
		}
		if !produtoExiste {
			http.Error(w, "Produto não encontrado", http.StatusNotFound)
			return
		}

		// Atualizar alerta mínimo
		_, err = db.Exec(`
			UPDATE estoque 
			SET alerta_minimo = $1, atualizado_em = NOW() 
			WHERE produto_id = $2
		`, req.AlertaMinimo, produtoID)
		if err != nil {
			http.Error(w, "Erro ao atualizar alerta mínimo: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Buscar informações atualizadas do estoque
		var e models.EstoqueResponse
		var botijasVazias, botijasEmprestadas, alertaMinimo sql.NullInt64

		err = db.QueryRow(`
			SELECT e.id, e.produto_id, p.nome, p.categoria, e.quantidade, 
			       e.botijas_vazias, e.botijas_emprestadas, e.alerta_minimo, e.atualizado_em
			FROM estoque e
			JOIN produtos p ON e.produto_id = p.id
			WHERE e.produto_id = $1
		`, produtoID).Scan(
			&e.ID, &e.ProdutoID, &e.NomeProduto, &e.Categoria, &e.Quantidade,
			&botijasVazias, &botijasEmprestadas, &alertaMinimo, &e.AtualizadoEm,
		)
		if err != nil {
			http.Error(w, "Alerta mínimo atualizado, mas erro ao buscar informações: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Converter tipos nulos
		if botijasVazias.Valid {
			e.BotijasVazias = int(botijasVazias.Int64)
		}
		if botijasEmprestadas.Valid {
			e.BotijasEmprestadas = int(botijasEmprestadas.Int64)
		}
		if alertaMinimo.Valid {
			e.AlertaMinimo = int(alertaMinimo.Int64)
		}

		// Determinar status do estoque
		e.Status = "normal"
		if alertaMinimo.Valid && e.Quantidade <= int(alertaMinimo.Int64) {
			if e.Quantidade <= int(alertaMinimo.Int64/2) {
				e.Status = "critico"
			} else {
				e.Status = "baixo"
			}
		}

		// Retornar resposta
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(e)
	}
}

// EmprestimoBotijasHandler registra empréstimo de botijas vazias ao caminhoneiro
func EmprestimoBotijasHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Verificar se o usuário está autenticado
		userID, ok := middleware.ObterUsuarioID(r)
		if !ok {
			http.Error(w, "Usuário não autenticado", http.StatusUnauthorized)
			return
		}

		// Verificar permissões (apenas atendentes ou acima podem registrar empréstimos)
		perfil, perfilOk := middleware.ObterPerfilUsuario(r)
		if !perfilOk || !middleware.VerificarPerfil(perfil, "atendente") {
			http.Error(w, "Sem permissão para registrar empréstimo", http.StatusForbidden)
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
		var req models.EmprestimoBotijasRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Erro ao processar requisição: "+err.Error(), http.StatusBadRequest)
			return
		}

		// Validar dados
		if req.ProdutoID <= 0 {
			http.Error(w, "ID do produto é obrigatório", http.StatusBadRequest)
			return
		}
		if req.Quantidade <= 0 {
			http.Error(w, "Quantidade deve ser maior que zero", http.StatusBadRequest)
			return
		}

		// Iniciar transação
		tx, err := db.Begin()
		if err != nil {
			http.Error(w, "Erro ao iniciar transação: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer func() {
			if err != nil {
				tx.Rollback()
				return
			}
		}()

		// Verificar se há botijas vazias suficientes
		var qtdBotijasVazias int
		var produtoNome string
		err = tx.QueryRow(`
			SELECT e.botijas_vazias, p.nome
			FROM estoque e
			JOIN produtos p ON e.produto_id = p.id
			WHERE e.produto_id = $1
		`, req.ProdutoID).Scan(&qtdBotijasVazias, &produtoNome)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Produto não encontrado", http.StatusNotFound)
				return
			}
			http.Error(w, "Erro ao verificar botijas vazias: "+err.Error(), http.StatusInternalServerError)
			return
		}

		if qtdBotijasVazias < req.Quantidade {
			http.Error(w, fmt.Sprintf("Botijas vazias insuficientes para o produto %s. Disponível: %d", produtoNome, qtdBotijasVazias), http.StatusBadRequest)
			return
		}

		// Atualizar estoque
		_, err = tx.Exec(`
			UPDATE estoque 
			SET botijas_vazias = botijas_vazias - $1,
				botijas_emprestadas = botijas_emprestadas + $1,
				atualizado_em = NOW()
			WHERE produto_id = $2
		`, req.Quantidade, req.ProdutoID)
		if err != nil {
			http.Error(w, "Erro ao atualizar estoque: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Registrar movimentação
		_, err = tx.Exec(`
			INSERT INTO movimentacoes_estoque
			(produto_id, tipo, quantidade, observacoes, usuario_id, criado_em)
			VALUES
			($1, $2, $3, $4, $5, NOW())
		`, req.ProdutoID, models.MovimentacaoEmprestimo, req.Quantidade, req.Observacoes, userID)
		if err != nil {
			http.Error(w, "Erro ao registrar movimentação: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Commit da transação
		err = tx.Commit()
		if err != nil {
			http.Error(w, "Erro ao finalizar transação: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Buscar informações atualizadas do estoque
		var e models.EstoqueResponse
		var botijasVazias, botijasEmprestadas, alertaMinimo sql.NullInt64

		err = db.QueryRow(`
			SELECT e.id, e.produto_id, p.nome, p.categoria, e.quantidade, 
			       e.botijas_vazias, e.botijas_emprestadas, e.alerta_minimo, e.atualizado_em
			FROM estoque e
			JOIN produtos p ON e.produto_id = p.id
			WHERE e.produto_id = $1
		`, req.ProdutoID).Scan(
			&e.ID, &e.ProdutoID, &e.NomeProduto, &e.Categoria, &e.Quantidade,
			&botijasVazias, &botijasEmprestadas, &alertaMinimo, &e.AtualizadoEm,
		)
		if err != nil {
			http.Error(w, "Empréstimo registrado, mas erro ao buscar informações: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Converter tipos nulos
		if botijasVazias.Valid {
			e.BotijasVazias = int(botijasVazias.Int64)
		}
		if botijasEmprestadas.Valid {
			e.BotijasEmprestadas = int(botijasEmprestadas.Int64)
		}
		if alertaMinimo.Valid {
			e.AlertaMinimo = int(alertaMinimo.Int64)
		}

		// Montar resposta
		response := struct {
			Mensagem string                `json:"mensagem"`
			Estoque  models.EstoqueResponse `json:"estoque"`
		}{
			Mensagem: fmt.Sprintf("Empréstimo de %d botijas vazias de %s registrado com sucesso", req.Quantidade, e.NomeProduto),
			Estoque:  e,
		}

		// Retornar resposta
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}

// DevolucaoBotijasEmprestimoHandler registra devolução de botijas emprestadas (troca por botijas cheias)
func DevolucaoBotijasEmprestimoHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Verificar se o usuário está autenticado
		userID, ok := middleware.ObterUsuarioID(r)
		if !ok {
			http.Error(w, "Usuário não autenticado", http.StatusUnauthorized)
			return
		}

		// Verificar permissões (apenas atendentes ou acima podem registrar devoluções)
		perfil, perfilOk := middleware.ObterPerfilUsuario(r)
		if !perfilOk || !middleware.VerificarPerfil(perfil, "atendente") {
			http.Error(w, "Sem permissão para registrar devolução", http.StatusForbidden)
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
		var req models.DevolucaoBotijasRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Erro ao processar requisição: "+err.Error(), http.StatusBadRequest)
			return
		}

		// Validar dados
		if req.ProdutoID <= 0 {
			http.Error(w, "ID do produto é obrigatório", http.StatusBadRequest)
			return
		}
		if req.Quantidade <= 0 {
			http.Error(w, "Quantidade deve ser maior que zero", http.StatusBadRequest)
			return
		}

		// Iniciar transação
		tx, err := db.Begin()
		if err != nil {
			http.Error(w, "Erro ao iniciar transação: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer func() {
			if err != nil {
				tx.Rollback()
				return
			}
		}()

		// Verificar se há botijas emprestadas suficientes
		var qtdBotijasEmprestadas int
		var produtoNome string
		err = tx.QueryRow(`
			SELECT e.botijas_emprestadas, p.nome
			FROM estoque e
			JOIN produtos p ON e.produto_id = p.id
			WHERE e.produto_id = $1
		`, req.ProdutoID).Scan(&qtdBotijasEmprestadas, &produtoNome)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Produto não encontrado", http.StatusNotFound)
				return
			}
			http.Error(w, "Erro ao verificar botijas emprestadas: "+err.Error(), http.StatusInternalServerError)
			return
		}

		if qtdBotijasEmprestadas < req.Quantidade {
			http.Error(w, fmt.Sprintf("Botijas emprestadas insuficientes para o produto %s. Disponível: %d", produtoNome, qtdBotijasEmprestadas), http.StatusBadRequest)
			return
		}

		// Atualizar estoque
		_, err = tx.Exec(`
			UPDATE estoque 
			SET botijas_emprestadas = botijas_emprestadas - $1,
				quantidade = quantidade + $1,
				atualizado_em = NOW()
			WHERE produto_id = $2
		`, req.Quantidade, req.ProdutoID)
		if err != nil {
			http.Error(w, "Erro ao atualizar estoque: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Registrar movimentação
		_, err = tx.Exec(`
			INSERT INTO movimentacoes_estoque
			(produto_id, tipo, quantidade, observacoes, usuario_id, criado_em)
			VALUES
			($1, $2, $3, $4, $5, NOW())
		`, req.ProdutoID, models.MovimentacaoDevolucaoEmprestimo, req.Quantidade, req.Observacoes, userID)
		if err != nil {
			http.Error(w, "Erro ao registrar movimentação: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Commit da transação
		err = tx.Commit()
		if err != nil {
			http.Error(w, "Erro ao finalizar transação: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Buscar informações atualizadas do estoque
		var e models.EstoqueResponse
		var botijasVazias, botijasEmprestadas, alertaMinimo sql.NullInt64

		err = db.QueryRow(`
			SELECT e.id, e.produto_id, p.nome, p.categoria, e.quantidade, 
			       e.botijas_vazias, e.botijas_emprestadas, e.alerta_minimo, e.atualizado_em
			FROM estoque e
			JOIN produtos p ON e.produto_id = p.id
			WHERE e.produto_id = $1
		`, req.ProdutoID).Scan(
			&e.ID, &e.ProdutoID, &e.NomeProduto, &e.Categoria, &e.Quantidade,
			&botijasVazias, &botijasEmprestadas, &alertaMinimo, &e.AtualizadoEm,
		)
		if err != nil {
			http.Error(w, "Devolução registrada, mas erro ao buscar informações: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Converter tipos nulos
		if botijasVazias.Valid {
			e.BotijasVazias = int(botijasVazias.Int64)
		}
		if botijasEmprestadas.Valid {
			e.BotijasEmprestadas = int(botijasEmprestadas.Int64)
		}
		if alertaMinimo.Valid {
			e.AlertaMinimo = int(alertaMinimo.Int64)
		}

		// Montar resposta
		response := struct {
			Mensagem string                `json:"mensagem"`
			Estoque  models.EstoqueResponse `json:"estoque"`
		}{
			Mensagem: fmt.Sprintf("Devolução de %d botijas emprestadas de %s registrada com sucesso", req.Quantidade, e.NomeProduto),
			Estoque:  e,
		}

		// Retornar resposta
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}
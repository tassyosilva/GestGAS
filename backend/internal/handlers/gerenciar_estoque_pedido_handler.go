package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/tassyosilva/GestGAS/internal/middleware"
	"github.com/tassyosilva/GestGAS/internal/models"
)

// GerenciarEstoquePedidoHandler gerencia o estoque durante o ciclo de vida de um pedido
func GerenciarEstoquePedidoHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Verificar se o usuário está autenticado
		userID, ok := middleware.ObterUsuarioID(r)
		if !ok {
			http.Error(w, "Usuário não autenticado", http.StatusUnauthorized)
			return
		}

		// Configurar cabeçalhos
		w.Header().Set("Content-Type", "application/json")

		// Verificar método
		if r.Method != http.MethodPost {
			http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
			return
		}

		// Estrutura para a requisição
		type EstoquePedidoRequest struct {
			PedidoID          int    `json:"pedido_id"`
			Acao              string `json:"acao"` // "confirmar_entrega", "cancelar", etc.
			MotivoCancelamento string `json:"motivo_cancelamento,omitempty"`
		}

		// Decodificar requisição
		var req EstoquePedidoRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Erro ao processar requisição: "+err.Error(), http.StatusBadRequest)
			return
		}

		// Validar dados
		if req.PedidoID <= 0 {
			http.Error(w, "ID do pedido é obrigatório", http.StatusBadRequest)
			return
		}

		// Validar ações permitidas
		if req.Acao != "confirmar_entrega" && req.Acao != "cancelar" {
			http.Error(w, "Ação inválida", http.StatusBadRequest)
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

		// Buscar status atual do pedido
		var statusAtual string
		err = tx.QueryRow("SELECT status FROM pedidos WHERE id = $1", req.PedidoID).Scan(&statusAtual)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Pedido não encontrado", http.StatusNotFound)
				return
			}
			http.Error(w, "Erro ao buscar pedido: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Verificar se o status atual é compatível com a ação solicitada
		if req.Acao == "confirmar_entrega" && statusAtual != "em_entrega" {
			http.Error(w, "O pedido deve estar em entrega para confirmar a entrega", http.StatusBadRequest)
			return
		}

		// Processar ação
		switch req.Acao {
		case "confirmar_entrega":
			err = confirmarEntregaPedido(tx, req.PedidoID, userID)
		case "cancelar":
			if statusAtual == "entregue" || statusAtual == "finalizado" || statusAtual == "cancelado" {
				http.Error(w, "Não é possível cancelar um pedido entregue, finalizado ou já cancelado", http.StatusBadRequest)
				return
			}
			err = cancelarPedido(tx, req.PedidoID, userID, req.MotivoCancelamento)
		}

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Commit da transação
		err = tx.Commit()
		if err != nil {
			http.Error(w, "Erro ao finalizar transação: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Buscar pedido atualizado para resposta
		pedidoResp, err := buscarPedidoDetalhado(db, req.PedidoID)
		if err != nil {
			http.Error(w, "Ação concluída, mas erro ao buscar detalhes: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Retornar resposta
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(pedidoResp)
	}
}

// confirmarEntregaPedido processa a confirmação de entrega de um pedido
func confirmarEntregaPedido(tx *sql.Tx, pedidoID, userID int) error {
	// 1. Atualizar status do pedido para "entregue"
	_, err := tx.Exec(`
		UPDATE pedidos
		SET status = 'entregue', data_entrega = NOW(), atualizado_em = NOW()
		WHERE id = $1
	`, pedidoID)
	if err != nil {
		return fmt.Errorf("erro ao atualizar status do pedido: %w", err)
	}

	// 2. Processar botijas retornadas (se houver)
	// Buscar todos os itens do pedido que são botijas com retorno
	rows, err := tx.Query(`
		SELECT ip.produto_id, ip.quantidade
		FROM itens_pedido ip
		JOIN produtos p ON ip.produto_id = p.id
		WHERE ip.pedido_id = $1 AND ip.retorna_botija = TRUE
		AND p.categoria LIKE 'botija_gas%'
	`, pedidoID)
	if err != nil {
		return fmt.Errorf("erro ao buscar botijas retornadas: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var produtoID, quantidade int
		err = rows.Scan(&produtoID, &quantidade)
		if err != nil {
			return fmt.Errorf("erro ao processar botijas retornadas: %w", err)
		}

		// Atualizar estoque de botijas vazias
		_, err = tx.Exec(`
			UPDATE estoque
			SET botijas_vazias = botijas_vazias + $1, atualizado_em = NOW()
			WHERE produto_id = $2
		`, quantidade, produtoID)
		if err != nil {
			return fmt.Errorf("erro ao atualizar estoque de botijas vazias: %w", err)
		}

		// Registrar movimentação de estoque
		_, err = tx.Exec(`
			INSERT INTO movimentacoes_estoque
			(produto_id, tipo, quantidade, usuario_id, pedido_id, criado_em)
			VALUES
			($1, $2, $3, $4, $5, NOW())
		`, produtoID, models.MovimentacaoBotijasVazias, quantidade, userID, pedidoID)
		if err != nil {
			return fmt.Errorf("erro ao registrar movimentação de botijas vazias: %w", err)
		}
	}

	return nil
}

// cancelarPedido processa o cancelamento de um pedido
func cancelarPedido(tx *sql.Tx, pedidoID, userID int, motivoCancelamento string) error {
	// 1. Buscar status atual do pedido
	var statusAtual string
	err := tx.QueryRow("SELECT status FROM pedidos WHERE id = $1", pedidoID).Scan(&statusAtual)
	if err != nil {
		return fmt.Errorf("erro ao buscar status do pedido: %w", err)
	}
	if statusAtual == "entregue" || statusAtual == "finalizado" || statusAtual == "cancelado" {
		return fmt.Errorf("não é possível cancelar um pedido entregue, finalizado ou já cancelado")
	}

	// 2. Atualizar status do pedido para "cancelado"
	_, err = tx.Exec(`
		UPDATE pedidos 
		SET status = 'cancelado', motivo_cancelamento = $1, atualizado_em = NOW() 
		WHERE id = $2
	`, motivoCancelamento, pedidoID)
	if err != nil {
		return fmt.Errorf("erro ao atualizar status do pedido: %w", err)
	}

	// 3. Obter itens do pedido
	type itemPedido struct {
		ProdutoID  int
		Quantidade int
	}
	itens := []itemPedido{}
	rows, err := tx.Query(`SELECT produto_id, quantidade FROM itens_pedido WHERE pedido_id = $1`, pedidoID)
	if err != nil {
		return fmt.Errorf("erro ao buscar itens do pedido: %w", err)
	}
	for rows.Next() {
		var item itemPedido
		if err := rows.Scan(&item.ProdutoID, &item.Quantidade); err != nil {
			rows.Close()
			return fmt.Errorf("erro ao ler item do pedido: %w", err)
		}
		itens = append(itens, item)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return fmt.Errorf("erro ao iterar sobre itens do pedido: %w", err)
	}
	// 4. Processar cada item
	for _, item := range itens {
		// Verificar se o item existe no estoque
		var existeNoEstoque bool
		err := tx.QueryRow(`SELECT EXISTS(SELECT 1 FROM estoque WHERE produto_id = $1)`,
			item.ProdutoID).Scan(&existeNoEstoque)
		if err != nil {
			return fmt.Errorf("erro ao verificar existência no estoque para produto %d: %w",
				item.ProdutoID, err)
		}
		// Atualizar ou inserir no estoque
		if existeNoEstoque {
			_, err = tx.Exec(`
				UPDATE estoque
				SET quantidade = quantidade + $1, atualizado_em = NOW()
				WHERE produto_id = $2
			`, item.Quantidade, item.ProdutoID)
			if err != nil {
				return fmt.Errorf("erro ao atualizar estoque para produto %d: %w",
					item.ProdutoID, err)
			}
		} else {
			_, err = tx.Exec(`
				INSERT INTO estoque (produto_id, quantidade, botijas_vazias, quantidade_minima, atualizado_em)
				VALUES ($1, $2, 0, 0, NOW())
			`, item.ProdutoID, item.Quantidade)
			if err != nil {
				return fmt.Errorf("erro ao inserir estoque para produto %d: %w",
					item.ProdutoID, err)
			}
		}
		// Registrar movimentação
		_, err = tx.Exec(`
			INSERT INTO movimentacoes_estoque
			(produto_id, tipo, quantidade, usuario_id, pedido_id, criado_em)
			VALUES
			($1, $2, $3, $4, $5, NOW())
		`, item.ProdutoID, models.MovimentacaoDevolucao, item.Quantidade, userID, pedidoID)
		if err != nil {
			return fmt.Errorf("erro ao registrar movimentação para produto %d: %w",
				item.ProdutoID, err)
		}
	}
	return nil
}

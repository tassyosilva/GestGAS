package handlers

import (
    "database/sql"
    "encoding/json"
    "net/http"

    "github.com/tassyosilva/GestGAS/internal/middleware"
)

// RegistrarRetornoBotijasHandler registra botijas vazias retornadas por um cliente após a entrega
func RegistrarRetornoBotijasHandler(db *sql.DB) http.HandlerFunc {
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

        // Decodificar requisição
        var req struct {
            PedidoID int `json:"pedido_id"`
        }
        
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            http.Error(w, "Erro ao processar requisição: "+err.Error(), http.StatusBadRequest)
            return
        }

        // Validar dados
        if req.PedidoID <= 0 {
            http.Error(w, "ID do pedido é obrigatório", http.StatusBadRequest)
            return
        }

        // Verificar se o pedido está entregue
        var statusPedido string
        err := db.QueryRow("SELECT status FROM pedidos WHERE id = $1", req.PedidoID).Scan(&statusPedido)
        if err != nil {
            if err == sql.ErrNoRows {
                http.Error(w, "Pedido não encontrado", http.StatusNotFound)
                return
            }
            http.Error(w, "Erro ao verificar pedido: "+err.Error(), http.StatusInternalServerError)
            return
        }

        if statusPedido != "entregue" && statusPedido != "finalizado" {
            http.Error(w, "O pedido deve estar entregue ou finalizado para registrar botijas retornadas", http.StatusBadRequest)
            return
        }

        // Buscar os itens do pedido que são botijas a serem retornadas
        rows, err := db.Query(`
            SELECT ip.produto_id, ip.quantidade
            FROM itens_pedido ip
            JOIN produtos p ON ip.produto_id = p.id
            WHERE ip.pedido_id = $1 AND ip.retorna_botija = TRUE
            AND p.categoria LIKE 'botija_gas%'
        `, req.PedidoID)
        if err != nil {
            http.Error(w, "Erro ao buscar itens do pedido: "+err.Error(), http.StatusInternalServerError)
            return
        }
        defer rows.Close()

        // Verificar se há itens para registrar
        var itensRegistrados []map[string]interface{}
        var temItens bool

        for rows.Next() {
            temItens = true
            var produtoID, quantidade int
            err := rows.Scan(&produtoID, &quantidade)
            if err != nil {
                http.Error(w, "Erro ao processar itens do pedido: "+err.Error(), http.StatusInternalServerError)
                return
            }

            // Verificar o valor atual de botijas vazias
            var botijasVazias sql.NullInt64
            err = db.QueryRow(`SELECT botijas_vazias FROM estoque WHERE produto_id = $1`, produtoID).Scan(&botijasVazias)
            if err != nil {
                http.Error(w, "Erro ao verificar estoque: "+err.Error(), http.StatusInternalServerError)
                return
            }

            // Calcular novo valor
            var novoValor int
            if botijasVazias.Valid {
                novoValor = int(botijasVazias.Int64) + quantidade
            } else {
                novoValor = quantidade
            }

            // Atualizar estoque usando uma execução direta sem operação matemática
            _, err = db.Exec(`
                UPDATE estoque
                SET botijas_vazias = $1, atualizado_em = NOW()
                WHERE produto_id = $2
            `, novoValor, produtoID)
            if err != nil {
                http.Error(w, "Erro ao atualizar estoque: "+err.Error(), http.StatusInternalServerError)
                return
            }

            // Registrar movimentação
            _, err = db.Exec(`
                INSERT INTO movimentacoes_estoque
                (produto_id, tipo, quantidade, usuario_id, pedido_id, criado_em)
                VALUES
                ($1, 'botijas_vazias', $2, $3, $4, NOW())
            `, produtoID, quantidade, userID, req.PedidoID)
            if err != nil {
                http.Error(w, "Erro ao registrar movimentação: "+err.Error(), http.StatusInternalServerError)
                return
            }

            // Adicionar item para resposta
            var nomeProduto string
            db.QueryRow("SELECT nome FROM produtos WHERE id = $1", produtoID).Scan(&nomeProduto)
            
            itensRegistrados = append(itensRegistrados, map[string]interface{}{
                "produto_id": produtoID,
                "nome_produto": nomeProduto,
                "quantidade": quantidade,
            })
        }

        if !temItens {
            http.Error(w, "Não há botijas para retornar neste pedido", http.StatusBadRequest)
            return
        }

        // Retornar resposta
        w.WriteHeader(http.StatusOK)
        json.NewEncoder(w).Encode(map[string]interface{}{
            "mensagem": "Botijas vazias registradas com sucesso",
            "pedido_id": req.PedidoID,
            "itens": itensRegistrados,
        })
    }
}
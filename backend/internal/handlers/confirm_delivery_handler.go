package handlers

import (
    "database/sql"
    "encoding/json"
    "net/http"
    "strconv"

    "github.com/tassyosilva/GestGAS/internal/middleware"
)

// ConfirmarEntregaSimples é um handler simplificado para confirmar a entrega de um pedido
func ConfirmarEntregaSimples(db *sql.DB) http.HandlerFunc {
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

        // Buscar status atual do pedido
        var statusAtual string
        err := db.QueryRow("SELECT status FROM pedidos WHERE id = $1", req.PedidoID).Scan(&statusAtual)
        if err != nil {
            if err == sql.ErrNoRows {
                http.Error(w, "Pedido não encontrado", http.StatusNotFound)
                return
            }
            http.Error(w, "Erro ao buscar pedido: "+err.Error(), http.StatusInternalServerError)
            return
        }

        // Verificar se o status atual é compatível
        if statusAtual != "em_entrega" {
            http.Error(w, "O pedido deve estar em entrega para confirmar a entrega", http.StatusBadRequest)
            return
        }

        // Atualizar status do pedido para "entregue"
        _, err = db.Exec(`
            UPDATE pedidos
            SET status = 'entregue', data_entrega = NOW(), atualizado_em = NOW()
            WHERE id = $1
        `, req.PedidoID)
        
        if err != nil {
            http.Error(w, "Erro ao atualizar status do pedido: "+err.Error(), http.StatusInternalServerError)
            return
        }

        // Retornar resposta de sucesso
        w.WriteHeader(http.StatusOK)
        json.NewEncoder(w).Encode(map[string]string{
            "mensagem": "Entrega confirmada com sucesso",
            "pedido_id": strconv.Itoa(req.PedidoID),
        })
    }
}
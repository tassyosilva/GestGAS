package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/tassyosilva/GestGAS/internal/middleware"
)

// FinalizarPedidoHandler manipula a finalização de um pedido sem processamento adicional
func FinalizarPedidoHandler(db *sql.DB) http.HandlerFunc {
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

		// Verificar se o pedido existe e está no status 'entregue'
		var status string
		err := db.QueryRow("SELECT status FROM pedidos WHERE id = $1", req.PedidoID).Scan(&status)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Pedido não encontrado", http.StatusNotFound)
				return
			}
			http.Error(w, "Erro ao verificar pedido: "+err.Error(), http.StatusInternalServerError)
			return
		}

		if status != "entregue" {
			http.Error(w, "Apenas pedidos com status 'entregue' podem ser finalizados", http.StatusBadRequest)
			return
		}

		// Atualizar apenas o status para 'finalizado'
		_, err = db.Exec(`
            UPDATE pedidos
            SET status = 'finalizado', atualizado_em = NOW()
            WHERE id = $1
        `, req.PedidoID)
		
		if err != nil {
			http.Error(w, "Erro ao finalizar pedido: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Retornar resposta de sucesso
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"mensagem": "Pedido finalizado com sucesso",
			"pedido_id": strconv.Itoa(req.PedidoID),
		})
	}
}
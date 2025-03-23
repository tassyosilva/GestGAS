package models

import "time"

// Estoque representa o estoque de um produto no sistema
type Estoque struct {
	ID            int       `json:"id"`
	ProdutoID     int       `json:"produto_id"`
	Produto       *Produto  `json:"produto,omitempty"`
	Quantidade    int       `json:"quantidade"`
	BotijasVazias int       `json:"botijas_vazias,omitempty"` // Apenas para botijas de gás
	BotijasEmprestadas int  `json:"botijas_emprestadas,omitempty"` // Botijas emprestadas a caminhoneiros
	AlertaMinimo  int       `json:"alerta_minimo,omitempty"` // Quantidade mínima antes de alertar para reposição
	CriadoEm      time.Time `json:"criado_em"`
	AtualizadoEm  time.Time `json:"atualizado_em"`
}

// MovimentacaoEstoque representa uma movimentação de entrada ou saída no estoque
type MovimentacaoEstoque struct {
	ID            int       `json:"id"`
	ProdutoID     int       `json:"produto_id"`
	Produto       *Produto  `json:"produto,omitempty"`
	Tipo          string    `json:"tipo"` // entrada, saida, emprestimo, devolucao
	Quantidade    int       `json:"quantidade"`
	Observacoes   string    `json:"observacoes,omitempty"`
	UsuarioID     int       `json:"usuario_id"`
	PedidoID      int       `json:"pedido_id,omitempty"` // Se a movimentação está associada a um pedido
	CriadoEm      time.Time `json:"criado_em"`
}

// Tipos de movimentação de estoque
const (
	MovimentacaoEntrada     = "entrada"     // Entrada de produtos no estoque
	MovimentacaoSaida       = "saida"       // Saída de produtos do estoque
	MovimentacaoEmprestimo  = "emprestimo"  // Empréstimo de botijas vazias para caminhoneiros
	MovimentacaoDevolucao   = "devolucao"   // Devolução de botijas vazias por caminhoneiros
	MovimentacaoAjuste      = "ajuste"      // Ajuste manual de estoque
)

// VendaFiada representa uma venda a prazo (fiado)
type VendaFiada struct {
	ID            int       `json:"id"`
	PedidoID      int       `json:"pedido_id"`
	ClienteID     int       `json:"cliente_id"`
	Cliente       *Cliente  `json:"cliente,omitempty"`
	ValorTotal    float64   `json:"valor_total"`
	DataVencimento time.Time `json:"data_vencimento"`
	Status        string    `json:"status"` // pendente, pago, vencido
	DataPagamento *time.Time `json:"data_pagamento,omitempty"`
	CriadoEm      time.Time `json:"criado_em"`
	AtualizadoEm  time.Time `json:"atualizado_em"`
}

// VendaAntecipada representa uma venda com pagamento antecipado
type VendaAntecipada struct {
	ID            int       `json:"id"`
	ClienteID     int       `json:"cliente_id"`
	Cliente       *Cliente  `json:"cliente,omitempty"`
	ProdutoID     int       `json:"produto_id"`
	Produto       *Produto  `json:"produto,omitempty"`
	Quantidade    int       `json:"quantidade"`
	ValorTotal    float64   `json:"valor_total"`
	FormaPagamento string    `json:"forma_pagamento"`
	DataPagamento time.Time `json:"data_pagamento"`
	DataEntregaPrevista time.Time `json:"data_entrega_prevista"`
	Status        string    `json:"status"` // pago, entregue, cancelado
	PedidoID      int       `json:"pedido_id,omitempty"` // ID do pedido quando for entregue
	CriadoEm      time.Time `json:"criado_em"`
	AtualizadoEm  time.Time `json:"atualizado_em"`
}

// Status de vendas fiadas
const (
	StatusFiadoPendente = "pendente"
	StatusFiadoPago     = "pago"
	StatusFiadoVencido  = "vencido"
)

// Status de vendas antecipadas
const (
	StatusAntecipadoPago     = "pago"
	StatusAntecipadoEntregue = "entregue"
	StatusAntecipadoCancelado = "cancelado"
)
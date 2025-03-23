package models

import "time"

// Pedido representa um pedido no sistema
type Pedido struct {
	ID              int                `json:"id"`
	ClienteID       int                `json:"cliente_id"`
	Cliente         *Cliente           `json:"cliente,omitempty"`
	AtendentID      int                `json:"atendente_id"`
	EntregadorID    int                `json:"entregador_id,omitempty"`
	Status          string             `json:"status"`
	FormaPagamento  string             `json:"forma_pagamento"`
	ValorTotal      float64            `json:"valor_total"`
	Observacoes     string             `json:"observacoes,omitempty"`
	EnderecoEntrega string             `json:"endereco_entrega"`
	DataEntrega     *time.Time         `json:"data_entrega,omitempty"`
	ItensPedido     []ItemPedido       `json:"itens_pedido,omitempty"`
	CriadoEm        time.Time          `json:"criado_em"`
	AtualizadoEm    time.Time          `json:"atualizado_em"`
}

// ItemPedido representa um item dentro de um pedido
type ItemPedido struct {
	ID          int     `json:"id"`
	PedidoID    int     `json:"pedido_id"`
	ProdutoID   int     `json:"produto_id"`
	Produto     *Produto `json:"produto,omitempty"`
	Quantidade  int     `json:"quantidade"`
	PrecoUnitario float64 `json:"preco_unitario"`
	Subtotal    float64 `json:"subtotal"`
}

// Status possíveis para um pedido
const (
	StatusPendente    = "pendente"     // Pedido registrado mas não confirmado
	StatusConfirmado  = "confirmado"   // Pedido confirmado, aguardando entrega
	StatusEmEntrega   = "em_entrega"   // Pedido saiu para entrega
	StatusEntregue    = "entregue"     // Pedido entregue com sucesso
	StatusCancelado   = "cancelado"    // Pedido cancelado
)

// Formas de pagamento
const (
	PagamentoDinheiro     = "dinheiro"
	PagamentoCartaoCredito = "cartao_credito"
	PagamentoCartaoDebito  = "cartao_debito"
	PagamentoPix          = "pix"
	PagamentoFiado        = "fiado"
)
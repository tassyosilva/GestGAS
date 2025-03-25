package models

import (
	"time"
)

// StatusPedido define os possíveis estados de um pedido
type StatusPedido string

const (
	StatusNovo       StatusPedido = "novo"        // Pedido recém-criado
	StatusEmPreparo  StatusPedido = "em_preparo"  // Pedido sendo preparado para entrega
	StatusEmEntrega  StatusPedido = "em_entrega"  // Pedido saiu para entrega
	StatusEntregue   StatusPedido = "entregue"    // Pedido foi entregue ao cliente
	StatusCancelado  StatusPedido = "cancelado"   // Pedido foi cancelado
	StatusFinalizado StatusPedido = "finalizado"  // Pedido foi entregue e pago
)

// FormaPagamento define as formas de pagamento aceitas
type FormaPagamento string

const (
	PagamentoDinheiro    FormaPagamento = "dinheiro"
	PagamentoCartaoDebito FormaPagamento = "cartao_debito"
	PagamentoCartaoCredito FormaPagamento = "cartao_credito"
	PagamentoPix         FormaPagamento = "pix"
	PagamentoFiado       FormaPagamento = "fiado"
)

// CanalOrigem define os canais de venda
type CanalOrigem string

const (
	CanalWhatsApp   CanalOrigem = "whatsapp"
	CanalTelefone   CanalOrigem = "telefone"
	CanalPresencial CanalOrigem = "presencial"
	CanalAplicativo CanalOrigem = "aplicativo"
)

// Pedido representa um pedido de compra no sistema
type Pedido struct {
	ID             int             `json:"id"`
	ClienteID      int             `json:"cliente_id"`
	Cliente        *ClienteBasico  `json:"cliente,omitempty"` // Campo adicionado
	AtendenteID    int             `json:"atendente_id"`
	EntregadorID   *int            `json:"entregador_id,omitempty"` // Pode ser nulo inicialmente
	Status         StatusPedido    `json:"status"`
	FormaPagamento FormaPagamento  `json:"forma_pagamento"`
	ValorTotal     float64         `json:"valor_total"`
	Observacoes    string          `json:"observacoes,omitempty"`
	EnderecoEntrega string         `json:"endereco_entrega"`
	CanalOrigem    CanalOrigem     `json:"canal_origem"`
	DataEntrega    *time.Time      `json:"data_entrega,omitempty"` // Pode ser nulo inicialmente
	Itens          []ItemPedido    `json:"itens,omitempty"`        // Itens do pedido
	CriadoEm       time.Time       `json:"criado_em"`
	AtualizadoEm   time.Time       `json:"atualizado_em"`
}

// ItemPedido representa um item dentro de um pedido
type ItemPedido struct {
	ID           int     `json:"id"`
	PedidoID     int     `json:"pedido_id"`
	ProdutoID    int     `json:"produto_id"`
	NomeProduto  string  `json:"nome_produto,omitempty"` // Para facilitar a exibição
	Quantidade   int     `json:"quantidade"`
	PrecoUnitario float64 `json:"preco_unitario"`
	Subtotal     float64 `json:"subtotal"`
	RetornaBotija bool   `json:"retorna_botija,omitempty"` // Indica se o cliente vai devolver uma botija vazia
}

// NovoPedidoRequest é a estrutura para receber um novo pedido via API
type NovoPedidoRequest struct {
	ClienteID      int             `json:"cliente_id"`
	FormaPagamento FormaPagamento  `json:"forma_pagamento"`
	Observacoes    string          `json:"observacoes,omitempty"`
	EnderecoEntrega string         `json:"endereco_entrega"`
	CanalOrigem    CanalOrigem     `json:"canal_origem"`
	Itens          []ItemPedidoRequest `json:"itens"`
}

// ItemPedidoRequest é a estrutura para receber os itens de um novo pedido
type ItemPedidoRequest struct {
	ProdutoID    int     `json:"produto_id"`
	Quantidade   int     `json:"quantidade"`
	RetornaBotija bool   `json:"retorna_botija,omitempty"`
}

// AtualizarStatusRequest é a estrutura para atualizar o status de um pedido
type AtualizarStatusRequest struct {
	Status       StatusPedido   `json:"status"`
	EntregadorID *int           `json:"entregador_id,omitempty"`
	DataEntrega  *time.Time     `json:"data_entrega,omitempty"`
}

// PedidoResponse é a estrutura de resposta para pedidos
type PedidoResponse struct {
	ID             int             `json:"id"`
	Cliente        ClienteBasico   `json:"cliente"`
	Atendente      UsuarioBasico   `json:"atendente"`
	Entregador     *UsuarioBasico  `json:"entregador,omitempty"`
	Status         StatusPedido    `json:"status"`
	FormaPagamento FormaPagamento  `json:"forma_pagamento"`
	ValorTotal     float64         `json:"valor_total"`
	Observacoes    string          `json:"observacoes,omitempty"`
	EnderecoEntrega string         `json:"endereco_entrega"`
	CanalOrigem    CanalOrigem     `json:"canal_origem"`
	DataEntrega    *time.Time      `json:"data_entrega,omitempty"`
	Itens          []ItemPedido    `json:"itens"`
	CriadoEm       time.Time       `json:"criado_em"`
	AtualizadoEm   time.Time       `json:"atualizado_em"`
}

// ClienteBasico é uma versão simplificada do cliente para inclusão no pedido
type ClienteBasico struct {
	ID       int    `json:"id"`
	Nome     string `json:"nome"`
	Telefone string `json:"telefone"`
}

// UsuarioBasico é uma versão simplificada do usuário para inclusão no pedido
type UsuarioBasico struct {
	ID     int    `json:"id"`
	Nome   string `json:"nome"`
	Perfil string `json:"perfil"`
}
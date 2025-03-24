package models

import (
	"time"
)

// Cliente representa um cliente cadastrado no sistema
type Cliente struct {
	ID           int        `json:"id"`
	Nome         string     `json:"nome"`
	Telefone     string     `json:"telefone"`
	CPF          string     `json:"cpf,omitempty"`
	Email        string     `json:"email,omitempty"`
	Endereco     string     `json:"endereco,omitempty"`
	Complemento  string     `json:"complemento,omitempty"`
	Bairro       string     `json:"bairro,omitempty"`
	Cidade       string     `json:"cidade,omitempty"`
	Estado       string     `json:"estado,omitempty"`
	CEP          string     `json:"cep,omitempty"`
	Observacoes  string     `json:"observacoes,omitempty"`
	CanalOrigem  CanalOrigem `json:"canal_origem,omitempty"`
	CriadoEm     time.Time  `json:"criado_em"`
	AtualizadoEm time.Time  `json:"atualizado_em"`
}

// NovoClienteRequest é a estrutura para receber um novo cliente via API
type NovoClienteRequest struct {
	Nome        string     `json:"nome"`
	Telefone    string     `json:"telefone"`
	CPF         string     `json:"cpf,omitempty"`
	Email       string     `json:"email,omitempty"`
	Endereco    string     `json:"endereco,omitempty"`
	Complemento string     `json:"complemento,omitempty"`
	Bairro      string     `json:"bairro,omitempty"`
	Cidade      string     `json:"cidade,omitempty"`
	Estado      string     `json:"estado,omitempty"`
	CEP         string     `json:"cep,omitempty"`
	Observacoes string     `json:"observacoes,omitempty"`
	CanalOrigem CanalOrigem `json:"canal_origem,omitempty"`
}

// ClienteResponse é a estrutura de resposta para consulta de clientes
type ClienteResponse struct {
	ID           int        `json:"id"`
	Nome         string     `json:"nome"`
	Telefone     string     `json:"telefone"`
	CPF          string     `json:"cpf,omitempty"`
	Email        string     `json:"email,omitempty"`
	Endereco     string     `json:"endereco,omitempty"`
	Complemento  string     `json:"complemento,omitempty"`
	Bairro       string     `json:"bairro,omitempty"`
	Cidade       string     `json:"cidade,omitempty"`
	Estado       string     `json:"estado,omitempty"`
	CEP          string     `json:"cep,omitempty"`
	Observacoes  string     `json:"observacoes,omitempty"`
	CanalOrigem  CanalOrigem `json:"canal_origem,omitempty"`
	UltimosPedidos []PedidoResumido `json:"ultimos_pedidos,omitempty"`
	TotalPedidos int        `json:"total_pedidos"`
	CriadoEm     time.Time  `json:"criado_em"`
	AtualizadoEm time.Time  `json:"atualizado_em"`
}

// PedidoResumido é uma versão simplificada do pedido para inclusão no ClienteResponse
type PedidoResumido struct {
	ID            int          `json:"id"`
	Status        StatusPedido `json:"status"`
	FormaPagamento FormaPagamento `json:"forma_pagamento"`
	ValorTotal    float64      `json:"valor_total"`
	DataPedido    time.Time    `json:"data_pedido"`
}

// ClienteEnderecoRequest é uma estrutura para atualização de endereço
type ClienteEnderecoRequest struct {
	Endereco    string `json:"endereco"`
	Complemento string `json:"complemento,omitempty"`
	Bairro      string `json:"bairro,omitempty"`
	Cidade      string `json:"cidade,omitempty"`
	Estado      string `json:"estado,omitempty"`
	CEP         string `json:"cep,omitempty"`
}
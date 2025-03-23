package models

import "time"

// Cliente representa um cliente no sistema
type Cliente struct {
	ID            int       `json:"id"`
	Nome          string    `json:"nome"`
	Telefone      string    `json:"telefone"`
	CPF           string    `json:"cpf,omitempty"`
	Email         string    `json:"email,omitempty"`
	Endereco      string    `json:"endereco,omitempty"`
	Complemento   string    `json:"complemento,omitempty"`
	Bairro        string    `json:"bairro,omitempty"`
	Cidade        string    `json:"cidade,omitempty"`
	Estado        string    `json:"estado,omitempty"`
	CEP           string    `json:"cep,omitempty"`
	Observacoes   string    `json:"observacoes,omitempty"`
	CanalOrigem   string    `json:"canal_origem,omitempty"` // Whatsapp, ligação, portaria, aplicativo
	CriadoEm      time.Time `json:"criado_em"`
	AtualizadoEm  time.Time `json:"atualizado_em"`
}

// Canais de origem do cliente
const (
	CanalWhatsapp  = "whatsapp"
	CanalLigacao   = "ligacao"
	CanalPortaria  = "portaria"
	CanalAplicativo = "aplicativo"
)
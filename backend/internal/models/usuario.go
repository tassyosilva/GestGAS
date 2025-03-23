package models

import "time"

// Usuario representa um usuário do sistema
type Usuario struct {
	ID          int       `json:"id"`
	Nome        string    `json:"nome"`
	Login       string    `json:"login"`
	Senha       string    `json:"senha,omitempty"` // omitempty para não retornar a senha em respostas JSON
	CPF         string    `json:"cpf,omitempty"`
	Email       string    `json:"email,omitempty"`
	Perfil      string    `json:"perfil"`
	CriadoEm    time.Time `json:"criado_em"`
	AtualizadoEm time.Time `json:"atualizado_em"`
}

// Perfis disponíveis no sistema
const (
	PerfilAdmin     = "admin"     // Administrador com acesso total
	PerfilGerente   = "gerente"   // Gerente com acesso amplo
	PerfilAtendente = "atendente" // Atendente com acesso limitado
	PerfilEntregador = "entregador" // Entregador com acesso muito restrito
)

// LoginRequest é a estrutura para receber requisições de login
type LoginRequest struct {
	Login    string `json:"login"`
	Senha    string `json:"senha"`
}

// LoginResponse é a estrutura para respostas de login bem-sucedido
type LoginResponse struct {
	ID     int    `json:"id"`
	Nome   string `json:"nome"`
	Login  string `json:"login"`
	Perfil string `json:"perfil"`
	Token  string `json:"token"` // Token JWT para autenticação
}
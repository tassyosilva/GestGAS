package models

import (
	"time"
)

// TipoMovimentacao define os tipos de movimentação de estoque
type TipoMovimentacao string

const (
	MovimentacaoEntrada     TipoMovimentacao = "entrada"      // Adição de produtos ao estoque
	MovimentacaoSaida       TipoMovimentacao = "saida"        // Remoção de produtos do estoque (vendas)
	MovimentacaoAjuste      TipoMovimentacao = "ajuste"       // Ajuste manual de estoque
	MovimentacaoDevolucao   TipoMovimentacao = "devolucao"    // Devolução de produto
	MovimentacaoBotijasVazias TipoMovimentacao = "botijas_vazias" // Entrada de botijas vazias
	MovimentacaoEmprestimo  TipoMovimentacao = "emprestimo"   // Empréstimo de botijas ao caminhoneiro
	MovimentacaoDevolucaoEmprestimo TipoMovimentacao = "devolucao_emprestimo" // Devolução de botijas emprestadas
)

// Estoque representa o estado atual do estoque de um produto
type Estoque struct {
	ID                int       `json:"id"`
	ProdutoID         int       `json:"produto_id"`
	NomeProduto       string    `json:"nome_produto,omitempty"` // Para facilitar a exibição
	Quantidade        int       `json:"quantidade"`            // Botijas cheias ou produtos regulares
	BotijasVazias     int       `json:"botijas_vazias,omitempty"`    // Para controle de botijas vazias
	BotijasEmprestadas int       `json:"botijas_emprestadas,omitempty"` // Para controle de botijas emprestadas ao caminhoneiro
	AlertaMinimo      int       `json:"alerta_minimo,omitempty"`
	CriadoEm          time.Time `json:"criado_em"`
	AtualizadoEm      time.Time `json:"atualizado_em"`
}

// MovimentacaoEstoque representa uma movimentação no estoque
type MovimentacaoEstoque struct {
	ID         int              `json:"id"`
	ProdutoID  int              `json:"produto_id"`
	Tipo       TipoMovimentacao `json:"tipo"`
	Quantidade int              `json:"quantidade"`
	Observacoes string          `json:"observacoes,omitempty"`
	UsuarioID  int              `json:"usuario_id"`
	PedidoID   *int             `json:"pedido_id,omitempty"` // Pode ser nulo em ajustes manuais
	CriadoEm   time.Time        `json:"criado_em"`
}

// MovimentacaoEstoqueRequest é a estrutura para receber uma movimentação de estoque via API
type MovimentacaoEstoqueRequest struct {
	ProdutoID   int              `json:"produto_id"`
	Tipo        TipoMovimentacao `json:"tipo"`
	Quantidade  int              `json:"quantidade"`
	Observacoes string           `json:"observacoes,omitempty"`
	PedidoID    *int             `json:"pedido_id,omitempty"`
}

// EstoqueResponse é a estrutura de resposta para consulta de estoque
type EstoqueResponse struct {
	ID                int       `json:"id"`
	ProdutoID         int       `json:"produto_id"`
	NomeProduto       string    `json:"nome_produto"`
	Categoria         string    `json:"categoria"` // Categoria do produto
	Quantidade        int       `json:"quantidade"`
	BotijasVazias     int       `json:"botijas_vazias,omitempty"`
	BotijasEmprestadas int       `json:"botijas_emprestadas,omitempty"`
	AlertaMinimo      int       `json:"alerta_minimo,omitempty"`
	Status            string    `json:"status"` // "normal", "baixo", "critico" baseado no alerta mínimo
	AtualizadoEm      time.Time `json:"atualizado_em"`
}

// EmprestimoBotijasRequest é a estrutura para receber solicitação de empréstimo de botijas
type EmprestimoBotijasRequest struct {
	ProdutoID   int    `json:"produto_id"`
	Quantidade  int    `json:"quantidade"`
	Observacoes string `json:"observacoes,omitempty"`
}

// DevolucaoBotijasRequest é a estrutura para receber devolução de botijas emprestadas
type DevolucaoBotijasRequest struct {
	ProdutoID   int    `json:"produto_id"`
	Quantidade  int    `json:"quantidade"`
	Observacoes string `json:"observacoes,omitempty"`
}

// EstoqueAlertaResponse é a estrutura de resposta para alertas de estoque baixo
type EstoqueAlertaResponse struct {
	ProdutoID    int    `json:"produto_id"`
	NomeProduto  string `json:"nome_produto"`
	Quantidade   int    `json:"quantidade"`
	AlertaMinimo int    `json:"alerta_minimo"`
	Status       string `json:"status"` // "baixo" ou "critico"
}
package models

import "time"

// Produto representa um produto no sistema
type Produto struct {
	ID          int       `json:"id"`
	Nome        string    `json:"nome"`
	Descricao   string    `json:"descricao,omitempty"`
	Categoria   string    `json:"categoria"`
	Preco       float64   `json:"preco"`
	CriadoEm    time.Time `json:"criado_em"`
	AtualizadoEm time.Time `json:"atualizado_em"`
}

// Categorias de produtos
const (
	CategoriaBotijaGas = "botija_gas"  // Botijas de gás
	CategoriaAgua      = "agua"        // Água mineral
	CategoriaAcessorio = "acessorio"   // Acessórios (registros, mangueiras)
	CategoriaOutros    = "outros"      // Outros produtos
)

// TiposBotijaGas contém os tamanhos disponíveis de botijas de gás
var TiposBotijaGas = []string{
	"02kg", "05kg", "08kg", "10kg", "13kg", "20kg", "45kg",
}

// TiposAgua contém os tipos disponíveis de água mineral
var TiposAgua = []string{
	"20L",
}

// TiposAcessorios contém os tipos de acessórios disponíveis
var TiposAcessorios = []string{
	"registro_simples",
	"registro_mangueira_80cm",
	"registro_mangueira_120cm",
}
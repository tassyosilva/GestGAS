package routes

import (
	"database/sql"
	"net/http"
	"strings"

	"github.com/tassyosilva/GestGAS/internal/handlers"
	"github.com/tassyosilva/GestGAS/internal/middleware"
)

// ConfigurarRotas configura todas as rotas da API
func ConfigurarRotas(db *sql.DB) http.Handler {
	mux := http.NewServeMux()

	// Definir rotas básicas
	mux.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "ok"}`))
	})

	// Rota de login (pública)
	mux.HandleFunc("/api/login", handlers.LoginHandler(db))

	// Rotas protegidas - Produtos
	mux.Handle("/api/produtos", middleware.AuthMiddleware(db)(http.HandlerFunc(handlers.ListarProdutosHandler(db))))
	mux.Handle("/api/produtos/", middleware.AuthMiddleware(db)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		segments := strings.Split(path, "/")
		// Verificar se é uma requisição para um produto específico
		if len(segments) >= 4 && segments[3] != "" {
			switch r.Method {
			case http.MethodGet:
				handlers.ObterProdutoHandler(db)(w, r)
			case http.MethodPut:
				handlers.AtualizarProdutoHandler(db)(w, r)
			case http.MethodDelete:
				handlers.ExcluirProdutoHandler(db)(w, r)
			default:
				http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
			}
			return
		}
		// Se chegar aqui, é uma requisição para criar um novo produto
		if r.Method == http.MethodPost {
			handlers.CriarProdutoHandler(db)(w, r)
			return
		}
		// Se não for nenhum dos casos acima, método não permitido
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
	})))

	// Rotas para clientes
	mux.Handle("/api/clientes", middleware.AuthMiddleware(db)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verificar se é método GET para listar ou POST para criar
		if r.Method == http.MethodGet {
			handlers.ListarClientesHandler(db)(w, r)
			return
		} else if r.Method == http.MethodPost {
			handlers.CriarClienteHandler(db)(w, r)
			return
		}
		// Se não for nenhum dos casos acima, método não permitido
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
	})))
	mux.Handle("/api/clientes/buscar", middleware.AuthMiddleware(db)(http.HandlerFunc(handlers.BuscarClientePorTelefoneHandler(db))))
	mux.Handle("/api/clientes/", middleware.AuthMiddleware(db)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		segments := strings.Split(path, "/")
		// Verificar se é uma requisição para um cliente específico
		if len(segments) >= 4 && segments[3] != "" {
			// Verificar se é uma atualização de endereço
			if len(segments) >= 5 && segments[4] == "endereco" {
				if r.Method == http.MethodPut || r.Method == http.MethodPatch {
					handlers.AtualizarClienteHandler(db)(w, r)
					return
				}
				http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
				return
			}
			// Operações padrão
			switch r.Method {
			case http.MethodGet:
				handlers.ObterClienteHandler(db)(w, r)
			case http.MethodPut, http.MethodPatch:
				handlers.AtualizarClienteHandler(db)(w, r)
			case http.MethodDelete:
				handlers.ExcluirClienteHandler(db)(w, r)
			default:
				http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
			}
			return
		}
		// Se não for nenhum dos casos acima, método não permitido
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
	})))

	// Rotas para pedidos
	mux.Handle("/api/pedidos", middleware.AuthMiddleware(db)(http.HandlerFunc(handlers.ListarPedidosHandler(db))))
	mux.Handle("/api/pedidos/", middleware.AuthMiddleware(db)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		segments := strings.Split(path, "/")
		// Rota para criar novo pedido
		if len(segments) == 4 && segments[3] == "" && r.Method == http.MethodPost {
			handlers.CriarPedidoHandler(db)(w, r)
			return
		}
		// Rota para atualizar status do pedido
		if len(segments) == 5 && segments[4] == "status" && (r.Method == http.MethodPut || r.Method == http.MethodPatch) {
			handlers.AtualizarStatusPedidoHandler(db)(w, r)
			return
		}
		// Rota para obter pedido específico
		if len(segments) == 4 && segments[3] != "" && r.Method == http.MethodGet {
			handlers.ObterPedidoHandler(db)(w, r)
			return
		}
		http.Error(w, "Rota não encontrada", http.StatusNotFound)
	})))

	// Rota para gerenciamento de estoque durante o ciclo de vida do pedido
	mux.Handle("/api/pedidos/estoque", middleware.AuthMiddleware(db)(http.HandlerFunc(handlers.GerenciarEstoquePedidoHandler(db))))

	// Adicionar nova rota para confirmar entrega
	mux.Handle("/api/pedidos/confirmar-entrega", middleware.AuthMiddleware(db)(http.HandlerFunc(handlers.ConfirmarEntregaSimples(db))))

	// Adicionar nova rota para registrar retorno de botijas
	mux.Handle("/api/pedidos/registrar-botijas", middleware.AuthMiddleware(db)(http.HandlerFunc(handlers.RegistrarRetornoBotijasHandler(db))))

	// Adicionar nova rota para finalizar pedido
	mux.Handle("/api/pedidos/finalizar", middleware.AuthMiddleware(db)(http.HandlerFunc(handlers.FinalizarPedidoHandler(db))))

	// Rotas para estoque
	mux.Handle("/api/estoque", middleware.AuthMiddleware(db)(http.HandlerFunc(handlers.ListarEstoqueHandler(db))))
	mux.Handle("/api/estoque/alertas", middleware.AuthMiddleware(db)(http.HandlerFunc(handlers.ListarAlertasEstoqueHandler(db))))
	mux.Handle("/api/estoque/botijas/emprestimo", middleware.AuthMiddleware(db)(http.HandlerFunc(handlers.EmprestimoBotijasHandler(db))))
	mux.Handle("/api/estoque/botijas/devolucao", middleware.AuthMiddleware(db)(http.HandlerFunc(handlers.DevolucaoBotijasEmprestimoHandler(db))))
	mux.Handle("/api/estoque/", middleware.AuthMiddleware(db)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		segments := strings.Split(path, "/")
		// Rota para obter item específico do estoque
		if len(segments) == 4 && segments[3] != "" && r.Method == http.MethodGet {
			handlers.ObterEstoqueItemHandler(db)(w, r)
			return
		}
		// Rota para atualizar estoque
		if len(segments) == 4 && segments[3] != "" && (r.Method == http.MethodPut || r.Method == http.MethodPatch) {
			handlers.AtualizarEstoqueHandler(db)(w, r)
			return
		}
		// Rota para atualizar alerta mínimo
		if len(segments) == 5 && segments[4] == "alerta" && (r.Method == http.MethodPut || r.Method == http.MethodPatch) {
			handlers.AtualizarAlertaMinimoHandler(db)(w, r)
			return
		}
		http.Error(w, "Rota não encontrada", http.StatusNotFound)
	})))

	// Rotas para usuários
	mux.Handle("/api/usuarios", middleware.AuthMiddleware(db)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verificar se é método GET para listar ou POST para criar
		if r.Method == http.MethodGet {
			handlers.ListarUsuariosHandler(db)(w, r)
			return
		} else if r.Method == http.MethodPost {
			handlers.CriarUsuarioHandler(db)(w, r)
			return
		}
		// Se não for nenhum dos casos acima, método não permitido
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
	})))
	mux.Handle("/api/usuarios/", middleware.AuthMiddleware(db)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		segments := strings.Split(path, "/")
		// Verificar se é uma requisição para um usuário específico
		if len(segments) >= 4 && segments[3] != "" {
			switch r.Method {
			case http.MethodGet:
				handlers.ObterUsuarioHandler(db)(w, r)
			case http.MethodPut, http.MethodPatch:
				handlers.AtualizarUsuarioHandler(db)(w, r)
			case http.MethodDelete:
				handlers.ExcluirUsuarioHandler(db)(w, r)
			default:
				http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
			}
			return
		}
		// Se não for nenhum dos casos acima, método não permitido
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
	})))

	// Rota específica para listar entregadores
	mux.Handle("/api/entregadores", middleware.AuthMiddleware(db)(http.HandlerFunc(handlers.ListarEntregadoresHandler(db))))

	// Aplicar o middleware CORS a todas as rotas
	return middleware.CorsMiddleware(mux)
}

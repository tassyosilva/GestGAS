package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	_ "os"
	"strings"
	"time"

	_ "github.com/lib/pq"
	"github.com/tassyosilva/GestGAS/internal/handlers"
	"github.com/tassyosilva/GestGAS/internal/middleware"
	"golang.org/x/crypto/bcrypt"
)

// Configuração do banco de dados
const (
	host     = "192.168.1.106"
	port     = 5432
	user     = "postgres"
	password = "123456"
	dbname   = "gestgas"
)

// Middleware CORS para permitir requisições cross-origin
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Obtenha a origem da requisição
		origin := r.Header.Get("Origin")
		
		// LISTA DE ORIGENS PERMITIDAS
		allowedOrigins := map[string]bool{
			"http://localhost:3000": true,
		}
		
		// Verifique se a origem está na lista de permitidas
		if allowedOrigins[origin] {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}
		
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, PATCH")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Authorization")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	// Conectar ao banco de dados PostgreSQL
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
		
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("Não foi possível conectar ao banco de dados: %v", err)
	}
	defer db.Close()

	// Verificar conexão com o banco de dados
	err = db.Ping()
	if err != nil {
		log.Fatalf("Erro ao conectar ao banco de dados: %v", err)
	}
	
	fmt.Println("Conectado ao banco de dados com sucesso!")
	
	// Inicializar o banco de dados (criar tabelas se não existirem)
	if err := inicializarBancoDados(db); err != nil {
		log.Fatalf("Erro ao inicializar o banco de dados: %v", err)
	}
	
	// Configurar o servidor HTTP
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
	
	// Rotas para clientes - MODIFICADO PARA CORRIGIR ERRO 405
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
	
	// NOVAS ROTAS PARA PEDIDOS
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

	// NOVAS ROTAS PARA ESTOQUE
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
	
	// Aplicar o middleware CORS a todas as rotas
	corsHandler := corsMiddleware(mux)
	
	// Configurar o servidor
	server := &http.Server{
		Addr:         ":8080",
		Handler:      corsHandler, // Usar o handler com CORS habilitado
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
	
	// Iniciar o servidor
	fmt.Printf("Servidor rodando na porta %s\n", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Falha ao iniciar o servidor: %v", err)
	}
}

// inicializarBancoDados cria as tabelas necessárias se não existirem
func inicializarBancoDados(db *sql.DB) error {
	// Criar tabela de usuários
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS usuarios (
			id SERIAL PRIMARY KEY,
			nome VARCHAR(100) NOT NULL,
			login VARCHAR(50) UNIQUE NOT NULL,
			senha VARCHAR(255) NOT NULL,
			cpf VARCHAR(14) UNIQUE,
			email VARCHAR(100) UNIQUE,
			perfil VARCHAR(20) NOT NULL,
			criado_em TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			atualizado_em TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("erro ao criar tabela de usuários: %w", err)
	}
	
	// Criar tabela de produtos
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS produtos (
			id SERIAL PRIMARY KEY,
			nome VARCHAR(100) NOT NULL,
			descricao TEXT,
			categoria VARCHAR(50) NOT NULL,
			preco DECIMAL(10, 2) NOT NULL,
			criado_em TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			atualizado_em TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("erro ao criar tabela de produtos: %w", err)
	}
	
	// Criar tabela de clientes
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS clientes (
			id SERIAL PRIMARY KEY,
			nome VARCHAR(100) NOT NULL,
			telefone VARCHAR(20) NOT NULL,
			cpf VARCHAR(14) UNIQUE,
			email VARCHAR(100) UNIQUE,
			endereco VARCHAR(255),
			complemento VARCHAR(100),
			bairro VARCHAR(100),
			cidade VARCHAR(100),
			estado VARCHAR(2),
			cep VARCHAR(10),
			observacoes TEXT,
			canal_origem VARCHAR(20),
			criado_em TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			atualizado_em TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("erro ao criar tabela de clientes: %w", err)
	}
	
	// Criar tabela de pedidos
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS pedidos (
			id SERIAL PRIMARY KEY,
			cliente_id INTEGER NOT NULL REFERENCES clientes(id),
			atendente_id INTEGER NOT NULL REFERENCES usuarios(id),
			entregador_id INTEGER REFERENCES usuarios(id),
			status VARCHAR(20) NOT NULL,
			forma_pagamento VARCHAR(20) NOT NULL,
			valor_total DECIMAL(10, 2) NOT NULL,
			observacoes TEXT,
			endereco_entrega VARCHAR(255) NOT NULL,
			canal_origem VARCHAR(20),
			data_entrega TIMESTAMP WITH TIME ZONE,
			criado_em TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			atualizado_em TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("erro ao criar tabela de pedidos: %w", err)
	}
	
	// Criar tabela de itens de pedido
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS itens_pedido (
			id SERIAL PRIMARY KEY,
			pedido_id INTEGER NOT NULL REFERENCES pedidos(id),
			produto_id INTEGER NOT NULL REFERENCES produtos(id),
			quantidade INTEGER NOT NULL,
			preco_unitario DECIMAL(10, 2) NOT NULL,
			subtotal DECIMAL(10, 2) NOT NULL,
			retorna_botija BOOLEAN DEFAULT FALSE
		)
	`)
	if err != nil {
		return fmt.Errorf("erro ao criar tabela de itens de pedido: %w", err)
	}
	
	// Criar tabela de estoque
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS estoque (
			id SERIAL PRIMARY KEY,
			produto_id INTEGER NOT NULL REFERENCES produtos(id),
			quantidade INTEGER NOT NULL DEFAULT 0,
			botijas_vazias INTEGER DEFAULT 0,
			botijas_emprestadas INTEGER DEFAULT 0,
			alerta_minimo INTEGER,
			criado_em TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			atualizado_em TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("erro ao criar tabela de estoque: %w", err)
	}
	
	// Criar tabela de movimentações de estoque
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS movimentacoes_estoque (
			id SERIAL PRIMARY KEY,
			produto_id INTEGER NOT NULL REFERENCES produtos(id),
			tipo VARCHAR(20) NOT NULL,
			quantidade INTEGER NOT NULL,
			observacoes TEXT,
			usuario_id INTEGER NOT NULL REFERENCES usuarios(id),
			pedido_id INTEGER REFERENCES pedidos(id),
			criado_em TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("erro ao criar tabela de movimentações de estoque: %w", err)
	}
	
	// Criar tabela de vendas fiadas
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS vendas_fiadas (
			id SERIAL PRIMARY KEY,
			pedido_id INTEGER NOT NULL REFERENCES pedidos(id),
			cliente_id INTEGER NOT NULL REFERENCES clientes(id),
			valor_total DECIMAL(10, 2) NOT NULL,
			data_vencimento TIMESTAMP WITH TIME ZONE NOT NULL,
			status VARCHAR(20) NOT NULL,
			data_pagamento TIMESTAMP WITH TIME ZONE,
			criado_em TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			atualizado_em TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("erro ao criar tabela de vendas fiadas: %w", err)
	}
	
	// Criar tabela de vendas antecipadas
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS vendas_antecipadas (
			id SERIAL PRIMARY KEY,
			cliente_id INTEGER NOT NULL REFERENCES clientes(id),
			produto_id INTEGER NOT NULL REFERENCES produtos(id),
			quantidade INTEGER NOT NULL,
			valor_total DECIMAL(10, 2) NOT NULL,
			forma_pagamento VARCHAR(20) NOT NULL,
			data_pagamento TIMESTAMP WITH TIME ZONE NOT NULL,
			data_entrega_prevista TIMESTAMP WITH TIME ZONE NOT NULL,
			status VARCHAR(20) NOT NULL,
			pedido_id INTEGER REFERENCES pedidos(id),
			criado_em TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			atualizado_em TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("erro ao criar tabela de vendas antecipadas: %w", err)
	}

	// Verificar se já existe um usuário administrador
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM usuarios WHERE login = 'admin'").Scan(&count)
	if err != nil {
		return fmt.Errorf("erro ao verificar usuário admin: %w", err)
	}
	
	// Se não existir um admin, criar um
	if count == 0 {
		// Gerar hash da senha do admin
		senhaHash, err := bcrypt.GenerateFromPassword([]byte("admin"), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("erro ao gerar hash de senha: %w", err)
		}
		
		_, err = db.Exec(`
			INSERT INTO usuarios (nome, login, senha, cpf, email, perfil)
			VALUES ('Administrador', 'admin', $1, '000.000.000-00', 'admin@gestgas.com', 'admin')
		`, string(senhaHash))
		
		if err != nil {
			return fmt.Errorf("erro ao criar usuário admin: %w", err)
		}
		fmt.Println("Usuário administrador criado com sucesso!")
	}
	
	// Inserir alguns produtos iniciais se ainda não existirem
	err = db.QueryRow("SELECT COUNT(*) FROM produtos").Scan(&count)
	if err != nil {
		return fmt.Errorf("erro ao verificar produtos: %w", err)
	}
	
	if count == 0 {
		// Inserir botijas de gás
		_, err = db.Exec(`
			INSERT INTO produtos (nome, descricao, categoria, preco) VALUES
			('Botija de Gás 13kg', 'Botija de gás GLP 13kg', 'botija_gas', 115.00),
			('Botija de Gás 08kg', 'Botija de gás GLP 08kg', 'botija_gas', 80.00),
			('Botija de Gás 05kg', 'Botija de gás GLP 05kg', 'botija_gas', 65.00),
			('Botija de Gás 02kg', 'Botija de gás GLP 02kg', 'botija_gas', 40.00),
			('Botija de Gás 20kg', 'Botija de gás GLP 20kg', 'botija_gas', 180.00),
			('Botija de Gás 45kg', 'Botija de gás GLP 45kg', 'botija_gas', 360.00),
			('Água Mineral 20L', 'Galão de água mineral 20 litros', 'agua', 12.00),
			('Registro Simples', 'Registro para botija de gás simples', 'acessorio', 35.00),
			('Registro com Mangueira 80cm', 'Registro completo com mangueira de 80cm', 'acessorio', 45.00),
			('Registro com Mangueira 120cm', 'Registro completo com mangueira de 120cm', 'acessorio', 55.00)
		`)
		if err != nil {
			return fmt.Errorf("erro ao criar produtos iniciais: %w", err)
		}
		
		// Configurar estoque inicial
		_, err = db.Exec(`
			INSERT INTO estoque (produto_id, quantidade, alerta_minimo)
			SELECT id, 20, 5 FROM produtos
		`)
		if err != nil {
			return fmt.Errorf("erro ao configurar estoque inicial: %w", err)
		}
		
		fmt.Println("Produtos e estoque inicial configurados com sucesso!")
	}

	// Verificar se a coluna canal_origem existe na tabela pedidos
	var columnExists bool
	err = db.QueryRow(`
		SELECT EXISTS (
			SELECT 1 
			FROM information_schema.columns 
			WHERE table_name = 'pedidos' AND column_name = 'canal_origem'
		)
	`).Scan(&columnExists)
	if err != nil {
		log.Printf("Erro ao verificar coluna canal_origem: %v", err)
	} else if !columnExists {
		// Adicionar a coluna canal_origem se não existir
		_, err = db.Exec(`ALTER TABLE pedidos ADD COLUMN canal_origem VARCHAR(20)`)
		if err != nil {
			log.Printf("Erro ao adicionar coluna canal_origem: %v", err)
		} else {
			log.Println("Coluna canal_origem adicionada com sucesso à tabela pedidos")
		}
	}
	
	fmt.Println("Banco de dados inicializado com sucesso!")
	return nil
}
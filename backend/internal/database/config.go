package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
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

// Conectar estabelece a conexão com o banco de dados PostgreSQL
func Conectar() (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, fmt.Errorf("não foi possível conectar ao banco de dados: %w", err)
	}

	// Verificar conexão com o banco de dados
	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("erro ao conectar ao banco de dados: %w", err)
	}
	log.Println("Conectado ao banco de dados com sucesso!")
	return db, nil
}

// InicializarBancoDados cria as tabelas necessárias se não existirem
func InicializarBancoDados(db *sql.DB) error {
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

	// Adicionar o campo motivo_cancelamento à verificação de colunas
	var motivo_cancelamento_exists bool
	err = db.QueryRow(`
    SELECT EXISTS (
        SELECT 1 
        FROM information_schema.columns 
        WHERE table_name = 'pedidos' AND column_name = 'motivo_cancelamento'
    )
`).Scan(&motivo_cancelamento_exists)
	if err != nil {
		log.Printf("Erro ao verificar coluna motivo_cancelamento: %v", err)
	} else if !motivo_cancelamento_exists {
		// Adicionar a coluna motivo_cancelamento se não existir
		_, err = db.Exec(`ALTER TABLE pedidos ADD COLUMN motivo_cancelamento TEXT`)
		if err != nil {
			log.Printf("Erro ao adicionar coluna motivo_cancelamento: %v", err)
		} else {
			log.Println("Coluna motivo_cancelamento adicionada com sucesso à tabela pedidos")
		}
	}

	log.Println("Banco de dados inicializado com sucesso!")
	return nil
}

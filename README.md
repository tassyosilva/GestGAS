# GestGAS - Sistema de Gerenciamento para Revendedora de Gás e Água

## Sobre o Projeto

GestGAS é um sistema completo para gerenciamento de revendedoras de gás e água, desenvolvido para substituir o uso de planilhas Excel no controle de estoque, vendas, entregas e finanças.

### Funcionalidades Principais

- Gestão de produtos (botijas de gás, água mineral, acessórios)
- Controle de pedidos e entregas
- Gerenciamento de estoque de botijas (cheias e vazias)
- Controle financeiro (vendas por forma de pagamento)
- Relatórios e análises (ranking de entregadores, canais de venda)
- Gestão de vendas fiadas e antecipadas

## Tecnologias Utilizadas

### Backend
- Golang

### Banco de dados
- PostgreSQL

### Frontend
- React
- TypeScript
- Material UI
- Vite

## Estrutura do Projeto

### Backend

```
gasdelivery/backend/
├── cmd
│   └── api
│       └── main.go           # Ponto de entrada da aplicação
├── go.mod                    # Configuração de dependências Go
├── go.sum                    # Checksums das dependências Go
└── internal
    ├── auth
    │   └── jwt.go            # Funções para trabalhar com tokens JWT
    ├── database
    │   └── config.go         # Configuração de conexão com o banco de dados
    ├── handlers
    │   ├── auth_handler.go   # Handler para autenticação
    │   ├── produto_handler.go # Handler para operações com produtos
    │   ├── pedido_handler.go  # Handler para operações com pedidos
    │   ├── cliente_handler.go # Handler para operações com clientes
    │   ├── estoque_handler.go # Handler para operações com estoque
    │   └── usuario_handler.go # Handler para operações com usuários
    ├── middleware
    │   └── auth_middleware.go # Middleware de autenticação e autorização
    ├── models
    │   ├── user.go           # Modelo de usuários
    │   ├── produto.go        # Modelo de produtos
    │   ├── pedido.go         # Modelo de pedidos
    │   ├── cliente.go        # Modelo de clientes
    │   └── estoque.go        # Modelo de estoque
    └── repository
        ├── usuario_repository.go   # Repositório para operações com usuários
        ├── produto_repository.go   # Repositório para operações com produtos
        ├── pedido_repository.go    # Repositório para operações com pedidos
        ├── cliente_repository.go   # Repositório para operações com clientes
        └── estoque_repository.go   # Repositório para operações com estoque
```

### Frontend

```
gasdelivery/frontend/
├── src
│   ├── App.css               # Estilos globais da aplicação
│   ├── App.tsx               # Componente principal da aplicação
│   ├── assets
│   │   ├── logo.png          # Logo da aplicação
│   │   └── images            # Outras imagens
│   ├── components
│   │   ├── Layout.tsx        # Componente de layout compartilhado
│   │   ├── PrivateRoute.tsx  # Componente para rotas protegidas
│   │   ├── pedidos           # Componentes relacionados a pedidos
│   │   ├── produtos          # Componentes relacionados a produtos
│   │   ├── clientes          # Componentes relacionados a clientes
│   │   └── estoque           # Componentes relacionados a estoque
│   ├── config
│   │   └── api.ts            # Configuração da API
│   ├── index.css             # Estilos globais
│   ├── main.tsx              # Ponto de entrada do React
│   ├── pages
│   │   ├── Login.tsx         # Página de login
│   │   ├── Dashboard.tsx     # Dashboard principal
│   │   ├── Pedidos.tsx       # Página de gestão de pedidos
│   │   ├── Produtos.tsx      # Página de gestão de produtos
│   │   ├── Clientes.tsx      # Página de gestão de clientes
│   │   └── Estoque.tsx       # Página de gestão de estoque
│   ├── services
│   │   ├── authService.ts    # Serviço de autenticação
│   │   ├── pedidoService.ts  # Serviço de pedidos
│   │   ├── produtoService.ts # Serviço de produtos
│   │   ├── clienteService.ts # Serviço de clientes
│   │   └── estoqueService.ts # Serviço de estoque
│   ├── theme
│   │   └── theme.ts          # Configuração de tema
│   └── types
│       ├── user.ts           # Tipos para usuários
│       ├── produto.ts        # Tipos para produtos
│       ├── pedido.ts         # Tipos para pedidos
│       ├── cliente.ts        # Tipos para clientes
│       └── estoque.ts        # Tipos para estoque
```

### Pré-requisitos
- Go 1.23+
- PostgreSQL
- Node.js 18+
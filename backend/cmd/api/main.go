package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/tassyosilva/GestGAS/internal/database"
	"github.com/tassyosilva/GestGAS/internal/routes"
)

func main() {
	// Conectar ao banco de dados
	db, err := database.Conectar()
	if err != nil {
		log.Fatalf("Erro na conexão com o banco de dados: %v", err)
	}
	defer db.Close()
	
	// Inicializar o banco de dados (criar tabelas se não existirem)
	if err := database.InicializarBancoDados(db); err != nil {
		log.Fatalf("Erro ao inicializar o banco de dados: %v", err)
	}
	
	// Configurar rotas
	handler := routes.ConfigurarRotas(db)
	
	// Configurar o servidor
	server := &http.Server{
		Addr:         ":8080",
		Handler:      handler,
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
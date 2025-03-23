package auth

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// HashSenha gera um hash bcrypt da senha fornecida
func HashSenha(senha string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(senha), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// VerificarSenha compara a senha fornecida com um hash armazenado
func VerificarSenha(senha, hashArmazenado string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashArmazenado), []byte(senha))
	return err == nil
}

// GerarToken gera um token para autenticação
func GerarToken(userID int) (string, error) {
	// Gerar bytes aleatórios para o token
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	
	// Codificar como base64
	token := base64.StdEncoding.EncodeToString(b)
	
	// Adicionar ID do usuário e timestamp
	tokenCompleto := fmt.Sprintf("%d:%s:%d", userID, token, time.Now().Unix())
	
	return base64.StdEncoding.EncodeToString([]byte(tokenCompleto)), nil
}
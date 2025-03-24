package auth

import (
	"fmt"
	"time"
	"os"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

// Chave secreta para assinar os JWT
// Em produção, isso deve ser carregado a partir de uma variável de ambiente
var jwtKey = []byte(getEnv("JWT_SECRET", "sua_chave_secreta_super_segura"))

// Claims é a estrutura que vai dentro do token JWT
type Claims struct {
	UserID int    `json:"user_id"`
	Perfil string `json:"perfil"`
	jwt.StandardClaims
}

// Função auxiliar para obter variáveis de ambiente com valor padrão
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

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

// GerarToken gera um token JWT para autenticação
func GerarToken(userID int, perfil string) (string, error) {
	// Define a expiração do token (24 horas)
	expirationTime := time.Now().Add(24 * time.Hour)
	
	// Cria o payload do token (claims)
	claims := &Claims{
		UserID: userID,
		Perfil: perfil,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "gestgas-api",
		},
	}
	
	// Gera o token com o payload e assinatura HMAC SHA256
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	
	// Assina o token com a chave secreta
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	
	return tokenString, nil
}

// ValidarToken valida um token JWT e retorna os claims se válido
func ValidarToken(tokenString string) (*Claims, error) {
	// Parseia o token
	claims := &Claims{}
	
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Valida o método de assinatura
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("método de assinatura inesperado: %v", token.Header["alg"])
		}
		return jwtKey, nil
	})
	
	if err != nil {
		return nil, err
	}
	
	if !token.Valid {
		return nil, fmt.Errorf("token inválido")
	}
	
	return claims, nil
}
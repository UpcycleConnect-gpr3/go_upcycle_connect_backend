package jwt

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"go-upcycle_connect-backend/utils/log"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

var (
	publicKey *rsa.PublicKey
)

func init() {
	publicKeyPEM, err := os.ReadFile("public_key.pem")
	if err != nil {
		log.Fatal(err)
	}
	block, _ := pem.Decode(publicKeyPEM)
	if block == nil {
		log.Info("Failed to decode PEM block containing the public key")
	}

	publicKeyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		log.Fatal(err)
	}

	var ok bool
	publicKey, ok = publicKeyInterface.(*rsa.PublicKey)
	if !ok {
		log.Info("Public key is not an RSA key")
	}
}

func VerifyJWT(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return publicKey, nil
	})
	if err != nil {
		return "", fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", fmt.Errorf("invalid token claims or not valid")
	}

	userId, ok := claims["userId"].(string)
	if !ok {
		return "", fmt.Errorf("userId not found in token")
	}
	return userId, nil
}

func Auth(w http.ResponseWriter, r *http.Request) string {
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		http.Error(w, "Authorization token required", http.StatusUnauthorized)
		return ""
	}
	userId, err := VerifyJWT(tokenString)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return ""
	}

	return userId
}

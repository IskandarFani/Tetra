package services

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"net/mail"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func normalizeEmail(userEmail string) string {

	email := strings.ToLower(strings.TrimSpace(userEmail))

	if _, err := mail.ParseAddress(email); err != nil {
		return ""
	}

	return email
}

func hashPassword(userPassword string) (string, error) {

	hashPass, err := bcrypt.GenerateFromPassword([]byte(userPassword), bcrypt.DefaultCost)

	if err != nil {
		return "", err
	}

	return string(hashPass), err

}

func GenerateJWTToken(userID uint, userEmail string) (string, error) {

	getTokenTime := time.Now()
	secretCode := os.Getenv("JWT_SECRET")

	if secretCode == "" {
		return "", errors.New("secret code is empty")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"userID": userID,
			"email":  userEmail,
			"exp":    getTokenTime.Add(time.Hour * 24).Unix(),
			"iat":    getTokenTime.Unix(),
		})

	return token.SignedString([]byte(secretCode))

}

func (serv *Services) GenerateRefreshToken() (string, string, error) {

	bytes := make([]byte, 64)

	_, err := rand.Read(bytes)
	if err != nil {
		return "", "", err
	}

	refToken := base64.RawURLEncoding.EncodeToString(bytes)

	return refToken, hashRefreshToken(refToken), nil

}

func GetUserIDFromJWTToken(accessToken string) (uint, error) {

	secretCode := os.Getenv("JWT_SECRET")

	if secretCode == "" {
		return 0, errors.New("secret code is empty")
	}

	token, err := jwt.Parse(accessToken, getJWTSecretKey)

	if err != nil || !token.Valid {
		return 0, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok || claims["userID"] == nil {
		return 0, errors.New("invalid token claims")
	}

	userIDFloat, ok := claims["userID"].(float64)

	if !ok {
		return 0, errors.New("invalid userID claim")
	}

	return uint(userIDFloat), nil
}

func getJWTSecretKey(token *jwt.Token) (interface{}, error) {

	if token.Method != jwt.SigningMethodHS256 {
		return nil, errors.New("unexpected signing method")
	}

	secretCode := os.Getenv("JWT_SECRET")

	if secretCode == "" {
		return nil, errors.New("secret code is empty")
	}

	return []byte(secretCode), nil

}

func hashRefreshToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

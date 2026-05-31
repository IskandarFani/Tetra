package services

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type AuthTokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (serv *Services) Register(userEmail string, userPassword string) (AuthTokens, error) {

	var err error
	var hashPass string
	var newUserID uint

	if len([]rune(userPassword)) < 9 {
		return AuthTokens{}, errors.New("password is too short")
	}

	normalEmail := normalizeEmail(userEmail)

	if normalEmail == "" {
		return AuthTokens{}, errors.New("invalid email address")
	}

	exists, err := serv.repo.UserExists(normalEmail)

	if err != nil {
		return AuthTokens{}, err
	}

	if exists {
		return AuthTokens{}, errors.New("user already exists")
	}

	hashPass, err = hashPassword(userPassword)

	if err != nil {
		return AuthTokens{}, err
	}

	newUserID, err = serv.repo.CreateUser(normalEmail, hashPass)

	if err != nil {
		return AuthTokens{}, err
	}

	accessToken, err := GenerateJWTToken(newUserID, normalEmail)

	if err != nil {
		return AuthTokens{}, err
	}

	refreshToken, refreshTokenHash, err := GenerateRefreshToken()

	if err != nil {
		return AuthTokens{}, err
	}

	expiresAt := time.Now().Add(30 * time.Hour)

	err = serv.repo.CreateRefreshTokenHashRecord(newUserID, refreshTokenHash, expiresAt)

	if err != nil {
		return AuthTokens{}, err
	}

	return AuthTokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil

}

func (serv *Services) RefreshToken(refreshTokenString string) (AuthTokens, error) {

	hashToken := HashRefreshToken(refreshTokenString)

	userData, err := serv.repo.UseRefreshTokenByHash(hashToken)

	if err != nil {
		return AuthTokens{}, errors.New("invalid refresh token")
	}

	accessToken, err := GenerateJWTToken(userData.ID, userData.Email)

	if err != nil {
		return AuthTokens{}, err
	}

	refreshToken, refreshTokenHash, err := GenerateRefreshToken()

	if err != nil {
		return AuthTokens{}, err
	}

	expiresAt := time.Now().Add(30 * time.Hour)

	err = serv.repo.CreateRefreshTokenHashRecord(userData.ID, refreshTokenHash, expiresAt)

	if err != nil {
		return AuthTokens{}, err
	}

	return AuthTokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil

}

func (serv *Services) LogIn(userEmail string, userPassword string) (AuthTokens, error) {

	var err error

	normalEmail := normalizeEmail(userEmail)

	if normalEmail == "" {
		return AuthTokens{}, errors.New("invalid email address")
	}

	userData, err := serv.repo.FindUserByEmail(normalEmail)

	if err != nil {
		return AuthTokens{}, errors.New("user does not exist")
	}

	err = bcrypt.CompareHashAndPassword([]byte(userData.PasswordHash), []byte(userPassword))

	if err != nil {
		return AuthTokens{}, errors.New("incorrect password")
	}

	accessToken, err := GenerateJWTToken(userData.ID, normalEmail)

	if err != nil {
		return AuthTokens{}, err
	}

	refreshToken, refreshTokenHash, err := GenerateRefreshToken()

	if err != nil {
		return AuthTokens{}, err
	}

	expiresAt := time.Now().Add(30 * time.Hour)

	err = serv.repo.CreateRefreshTokenHashRecord(userData.ID, refreshTokenHash, expiresAt)

	if err != nil {
		return AuthTokens{}, err
	}

	return AuthTokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil

}

func (serv *Services) ParseAccessToken(accessToken string) (uint, string, time.Time, error) {

	userID, err := GetUserIDFromJWTToken(accessToken)

	if err != nil {
		return 0, "", time.Time{}, err
	}

	email, createdAt, err := serv.repo.FindUserByID(userID)

	if err != nil {
		return 0, "", time.Time{}, err
	}

	return userID, email, createdAt, nil

}

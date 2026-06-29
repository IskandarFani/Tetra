package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

type AuthTokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type UserProfile struct {
	ID        uint      `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
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

	refreshToken, refreshTokenHash, err := serv.GenerateRefreshToken()

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

func (serv *Services) UseRefreshToken(refreshToken string) (uint, string, error) {

	hashToken := hashRefreshToken(refreshToken)

	userID, userEmail, err := serv.useRefreshTokenByHash(hashToken)

	return userID, userEmail, err

}

func (serv *Services) useRefreshTokenByHash(refreshTokenHash string) (uint, string, error) {

	userData, err := serv.repo.UseRefreshTokenByHash(refreshTokenHash)

	if err != nil {
		return 0, "", err
	}

	return userData.UserID, userData.Email, nil

}

func (serv *Services) RefreshToken(usrID uint, usrEmail string) (AuthTokens, error) {

	accessToken, err := GenerateJWTToken(usrID, usrEmail)

	if err != nil {
		return AuthTokens{}, err
	}

	refreshToken, refreshTokenHash, err := serv.GenerateRefreshToken()

	if err != nil {
		return AuthTokens{}, err
	}

	expiresAt := time.Now().Add(30 * time.Hour)

	err = serv.repo.CreateRefreshTokenHashRecord(usrID, refreshTokenHash, expiresAt)

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

	refreshToken, refreshTokenHash, err := serv.GenerateRefreshToken()

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

func (serv *Services) ParseAccessToken(accessToken string) (uint, error) {

	userID, err := GetUserIDFromJWTToken(accessToken)

	if err != nil {
		return 0, err
	}

	return userID, nil

}

func (serv *Services) GetUserProfile(userID uint) (UserProfile, error) {

	ctx := context.Background()
	key := fmt.Sprintf("profile:user:%d", userID)

	cachedProfile, err := serv.redisClient.Get(ctx, key).Result()

	if err == nil {

		var profile UserProfile

		err := json.Unmarshal([]byte(cachedProfile), &profile)

		if err != nil {
			return UserProfile{}, err
		}

		return profile, nil
	}

	if err != redis.Nil {
		return UserProfile{}, err
	}

	email, createdAt, err := serv.repo.FindUserByID(userID)

	if err != nil {
		return UserProfile{}, err
	}

	profile := UserProfile{
		ID:        userID,
		Email:     email,
		CreatedAt: createdAt,
	}

	profileJSON, err := json.Marshal(profile)

	if err != nil {
		return UserProfile{}, err
	}

	err = serv.redisClient.Set(ctx, key, profileJSON, time.Hour).Err()

	if err != nil {
		return UserProfile{}, err
	}

	return profile, nil

}

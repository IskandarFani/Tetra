package repository

import (
	"errors"
	"go-cloud/models"
	"time"

	"gorm.io/gorm"
)

func (repo *Repository) FindUserByID(userID uint) (string, time.Time, error) {

	var user models.User

	err := repo.db.Where("id = ?", userID).First(&user).Error

	if err != nil {
		return "", time.Time{}, err
	}

	return user.Email, user.CreatedAt, nil

}

func (repo *Repository) UserExists(email string) (bool, error) {
	var user models.User

	err := repo.db.Where("email = ?", email).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return true, nil
}

func (repo *Repository) CreateUser(userEmail string, hashPass string) (uint, error) {

	newUser := models.User{
		Email:        userEmail,
		PasswordHash: hashPass}

	err := repo.db.Create(&newUser).Error

	if err != nil {
		return 0, err
	}

	return newUser.ID, nil

}

func (repo *Repository) CreateRefreshTokenHashRecord(userID uint, refreshTokenHash string, expiresAt time.Time) error {
	newToken := models.RefreshToken{
		UserID:    userID,
		TokenHash: refreshTokenHash,
		ExpiresAt: expiresAt,
		RevokedAt: nil,
	}

	return repo.db.Create(&newToken).Error
}

func (repo *Repository) FindUserByEmail(userEmail string) (*models.User, error) {

	var user models.User

	err := repo.db.Where("email = ?", userEmail).First(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, nil

}

func (repo *Repository) UseRefreshTokenByHash(refreshTokenHash string) (*models.User, error) {

	var token models.RefreshToken

	err := repo.db.Where("token_hash = ? AND revoked_at IS NULL AND expires_at > ?", refreshTokenHash, time.Now()).First(&token).Error

	if err != nil {
		return nil, err
	}

	err = repo.db.Where("id = ?", token.UserID).First(&token.User).Error

	if err != nil {
		return nil, err
	}

	err = repo.db.Model(&token).Update("revoked_at", time.Now()).Error

	if err != nil {
		return nil, err
	}

	return &token.User, nil

}

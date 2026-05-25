package repository

import "go-cloud/models"

func (repo *Repository) SubmitTrialAccessRequest(newUserEmail string) (string, error) {

	newTrial := models.TrialRequest{Email: newUserEmail}

	err := repo.db.Create(&newTrial).Error

	if err != nil {
		return "", err
	}

	return newUserEmail, nil

}

package repository

import "go-cloud/models"

func (repoStruct *Repository) AddNewUser(newUserName string) (string, error) {

	newUser := models.User{Name: newUserName}

	err := repoStruct.db.Create(&newUser).Error

	if err != nil {
		return "", err
	}

	return newUserName, nil

}

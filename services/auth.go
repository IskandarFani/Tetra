package services

import (
	"errors"
	"fmt"
)

func (servStruct *Services) AddNewUser(userName string) (string, error) {

	runesName := []rune(userName)
	nameLength := len(runesName)

	if nameLength > 7 {
		return "", errors.New("name length must be 7 characters or less")
	}

	newUserName := fmt.Sprintf("%s%d", userName, nameLength)

	return servStruct.repo.AddNewUser(newUserName)
}

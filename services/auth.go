package services

import (
	"fmt"
)

func (serv *Services) SubmitTrialAccessRequest(userEmail string) (string, error) {

	successMsg := fmt.Sprintf("A request has been sent to %s to confirm your email address", userEmail)

	if _, err := serv.repo.SubmitTrialAccessRequest(userEmail); err != nil {
		return "", err
	}

	return successMsg, nil
}

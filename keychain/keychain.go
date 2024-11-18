package keychain

import (
	"fmt"

	"github.com/zalando/go-keyring"
)

const (
	tokenService = "alfred-workflow-jira"
	tokenName    = "jira-api-token"
)

func SaveToken(tokenValue string) error {
	err := keyring.Set(tokenService, tokenName, tokenValue)
	if err != nil {
		return fmt.Errorf("failed to save token: %w", err)
	}
	return nil
}

func DeleteToken() error {
	err := keyring.Delete(tokenService, tokenName)
	if err != nil {
		return fmt.Errorf("failed to delete token: %w", err)
	}
	return nil
}

func GetToken() (string, error) {
	token, err := keyring.Get(tokenService, tokenName)
	if err != nil {
		return "", fmt.Errorf("failed to get token: %w", err)
	}
	return token, nil
}

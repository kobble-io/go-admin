package auth

import (
	"fmt"
	"github.com/valensto/kobble-go-sdk/utils"
)

var ErrorNames = []string{
	"ID_TOKEN_VERIFICATION_FAILED",
	"ACCESS_TOKEN_VERIFICATION_FAILED",
	"UNAUTHENTICATED",
}

type ErrorName string

func isValidAuthErrorName(name ErrorName) bool {
	for _, n := range ErrorNames {
		if n == string(name) {
			return true
		}
	}
	return false
}

type KobbleAuthError struct {
	*utils.ErrorBase
}

func newKobbleAuthError(name ErrorName, message string, cause error) (*KobbleAuthError, error) {
	if !isValidAuthErrorName(name) {
		return nil, fmt.Errorf("invalid AuthErrorName: %s", name)
	}
	return &KobbleAuthError{
		ErrorBase: utils.NewErrorBase(string(name), message, cause),
	}, nil
}

type IdTokenVerificationError struct {
	*KobbleAuthError
}

func newIdTokenVerificationError(cause error) error {
	err, errBaseErr := newKobbleAuthError("ID_TOKEN_VERIFICATION_FAILED", "ID token verification failed. Are you passing the correct ID token?", cause)
	if errBaseErr != nil {
		return errBaseErr
	}
	return IdTokenVerificationError{
		KobbleAuthError: err,
	}
}

type AccessTokenVerificationError struct {
	*KobbleAuthError
}

func newAccessTokenVerificationError(cause error) error {
	err, errBaseErr := newKobbleAuthError("ACCESS_TOKEN_VERIFICATION_FAILED", "Access token verification failed. Are you passing the correct access token?", cause)
	if errBaseErr != nil {
		return errBaseErr
	}
	return AccessTokenVerificationError{
		KobbleAuthError: err,
	}
}

package auth

import (
	"fmt"
	"github.com/kobble-io/go-admin/utils"
)

var errorNames = []string{
	"ID_TOKEN_VERIFICATION_FAILED",
	"ACCESS_TOKEN_VERIFICATION_FAILED",
	"UNAUTHENTICATED",
}

type errorName string

func isValidAuthErrorName(name errorName) bool {
	for _, n := range errorNames {
		if n == string(name) {
			return true
		}
	}
	return false
}

type kobbleAuthError struct {
	*utils.ErrorBase
}

func newKobbleAuthError(name errorName, message string, cause error) (*kobbleAuthError, error) {
	if !isValidAuthErrorName(name) {
		return nil, fmt.Errorf("invalid AuthErrorName: %s", name)
	}
	return &kobbleAuthError{
		ErrorBase: utils.NewErrorBase(string(name), message, cause),
	}, nil
}

type idTokenVerificationError struct {
	*kobbleAuthError
}

func newIdTokenVerificationError(cause error) error {
	err, errBaseErr := newKobbleAuthError("ID_TOKEN_VERIFICATION_FAILED", "ID token verification failed. Are you passing the correct ID token?", cause)
	if errBaseErr != nil {
		return errBaseErr
	}
	return idTokenVerificationError{
		kobbleAuthError: err,
	}
}

type accessTokenVerificationError struct {
	*kobbleAuthError
}

func newAccessTokenVerificationError(cause error) error {
	err, errBaseErr := newKobbleAuthError("ACCESS_TOKEN_VERIFICATION_FAILED", "Access token verification failed. Are you passing the correct access token?", cause)
	if errBaseErr != nil {
		return errBaseErr
	}
	return accessTokenVerificationError{
		kobbleAuthError: err,
	}
}

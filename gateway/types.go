package gateway

import "github.com/valensto/kobble-go-sdk/utils"

import (
	"crypto/ecdsa"
)

type (
	TokenProductQuota struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Used  int    `json:"used"`
		Limit int    `json:"limit"`
	}

	TokenProduct struct {
		ID     string              `json:"id"`
		Quotas []TokenProductQuota `json:"quotas"`
	}

	TokenPayload struct {
		ProjectID string `json:"project_id"`
		User      struct {
			Email    string         `json:"email"`
			ID       string         `json:"id"`
			Name     *string        `json:"name"`
			Products []TokenProduct `json:"products"`
		} `json:"user"`
	}

	ParseTokenOptions struct {
		VerifyIss       bool `json:"verify_iss,omitempty"`
		VerifyAud       bool `json:"verify_aud,omitempty"`
		VerifyExp       bool `json:"verify_exp,omitempty"`
		VerifySignature bool `json:"verify_signature,omitempty"`
	}

	keyInfo struct {
		Key       *ecdsa.PublicKey `json:"key"`
		ProjectID string           `json:"project_id"`
	}

	Config struct {
		Http *utils.HttpClient
	}
)

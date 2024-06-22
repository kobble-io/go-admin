package auth

import (
	"fmt"
	"github.com/valensto/kobble-go-sdk/utils"
)

import (
	"time"
)

type VerifyAccessTokenResult struct {
	UserID    string                      `json:"user_id"`
	ProjectID string                      `json:"project_id"`
	Claims    RawAccessTokenPayloadClaims `json:"claims"`
}

type IdTokenUser struct {
	ID         string `json:"id"`
	Email      string `json:"email"`
	Name       string `json:"name"`
	PictureURL string `json:"picture_url"`
	IsVerified bool   `json:"is_verified"`
	StripeID   string `json:"stripe_id"`
	UpdatedAt  string `json:"updated_at"`
	CreatedAt  string `json:"created_at"`
}

type VerifyIdTokenResult struct {
	UserID string                  `json:"user_id"`
	User   IdTokenUser             `json:"user"`
	Claims RawIdTokenPayloadClaims `json:"claims"`
}

type Config struct {
	Http    *utils.HttpClient
	BaseURL string
}

type Whoami struct {
	ProjectId   string
	ProjectSlug string
	UserId      string
}

type ProjectCache struct {
	ProjectID string
}

type RawIdTokenPayloadClaims struct {
	Sub        string `json:"sub"`
	ID         string `json:"id"`
	Email      string `json:"email"`
	Name       string `json:"name"`
	PictureURL string `json:"picture_url"`
	IsVerified bool   `json:"is_verified"`
	StripeID   string `json:"stripe_id"`
	UpdatedAt  string `json:"updated_at"`
	CreatedAt  string `json:"created_at"`
	Exp        int64  `json:"exp"`
	Iat        int64  `json:"iat"`
	Nbf        int64  `json:"nbf"`
	Iss        string `json:"iss"`
	Aud        string `json:"aud"`
}

func (r RawIdTokenPayloadClaims) Valid() error {
	now := time.Now().Unix()

	if r.Exp != 0 && now > r.Exp {
		return fmt.Errorf("token is expired")
	}

	if r.Nbf != 0 && now < r.Nbf {
		return fmt.Errorf("token is not valid yet")
	}

	if r.Iat != 0 && now < r.Iat {
		return fmt.Errorf("token used before issued")
	}

	return nil
}

type RawAccessTokenPayloadClaims struct {
	Sub       string `json:"sub"`
	ProjectID string `json:"project_id"`
	Exp       int64  `json:"exp"`
	Iat       int64  `json:"iat"`
	Nbf       int64  `json:"nbf"`
	Iss       string `json:"iss"`
	Aud       string `json:"aud"`
}

func (r RawAccessTokenPayloadClaims) Valid() error {
	now := time.Now().Unix()

	if r.Exp != 0 && now > r.Exp {
		return fmt.Errorf("token is expired")
	}

	if r.Nbf != 0 && now < r.Nbf {
		return fmt.Errorf("token is not valid yet")
	}

	if r.Iat != 0 && now < r.Iat {
		return fmt.Errorf("token used before issued")
	}

	return nil
}

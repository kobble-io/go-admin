package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/valensto/kobble-go-sdk/utils"
)

import (
	"time"
)

type VerifyAccessTokenOptions struct {
	ApplicationID *string `json:"application_id,omitempty"`
	Issuer        *string `json:"issuer,omitempty"`
}

type VerifyIdTokenOptions struct {
	ApplicationID *string `json:"application_id,omitempty"`
}

type VerifyAccessTokenResult struct {
	UserID    string                      `json:"user_id"`
	ProjectID string                      `json:"project_id"`
	Claims    RawAccessTokenPayloadClaims `json:"claims"`
}

type IdTokenUser struct {
	ID         string    `json:"id"`
	Email      *string   `json:"email"`
	Name       *string   `json:"name"`
	PictureURL *string   `json:"picture_url"`
	IsVerified bool      `json:"is_verified"`
	StripeID   *string   `json:"stripe_id"`
	UpdatedAt  time.Time `json:"updated_at"`
	CreatedAt  time.Time `json:"created_at"`
}

type VerifyIdTokenResult struct {
	UserID string                  `json:"user_id"`
	User   IdTokenUser             `json:"user"`
	Claims RawIdTokenPayloadClaims `json:"claims"`
}

type Config struct {
	http    *utils.HttpClient
	baseURL string
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
	Sub        string  `json:"sub"`
	ID         string  `json:"id"`
	Email      *string `json:"email"`
	Name       *string `json:"name"`
	PictureURL *string `json:"picture_url"`
	IsVerified bool    `json:"is_verified"`
	StripeID   *string `json:"stripe_id"`
	UpdatedAt  string  `json:"updated_at"`
	CreatedAt  string  `json:"created_at"`
}

func (r RawIdTokenPayloadClaims) GetExpirationTime() (*jwt.NumericDate, error) {
	return nil, nil
}

func (r RawIdTokenPayloadClaims) GetIssuedAt() (*jwt.NumericDate, error) {
	return nil, nil
}

func (r RawIdTokenPayloadClaims) GetNotBefore() (*jwt.NumericDate, error) {
	return nil, nil
}

func (r RawIdTokenPayloadClaims) GetIssuer() (string, error) {
	return "", nil
}

func (r RawIdTokenPayloadClaims) GetSubject() (string, error) {
	return r.Sub, nil
}

func (r RawIdTokenPayloadClaims) GetAudience() (jwt.ClaimStrings, error) {
	return jwt.ClaimStrings{}, nil
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

func (r RawAccessTokenPayloadClaims) GetExpirationTime() (*jwt.NumericDate, error) {
	return jwt.NewNumericDate(time.Unix(r.Exp, 0)), nil
}

func (r RawAccessTokenPayloadClaims) GetIssuedAt() (*jwt.NumericDate, error) {
	return jwt.NewNumericDate(time.Unix(r.Iat, 0)), nil
}

func (r RawAccessTokenPayloadClaims) GetNotBefore() (*jwt.NumericDate, error) {
	return jwt.NewNumericDate(time.Unix(r.Nbf, 0)), nil
}

func (r RawAccessTokenPayloadClaims) GetIssuer() (string, error) {
	return r.Iss, nil
}

func (r RawAccessTokenPayloadClaims) GetSubject() (string, error) {
	return r.Sub, nil
}

func (r RawAccessTokenPayloadClaims) GetAudience() (jwt.ClaimStrings, error) {
	return jwt.ClaimStrings{r.Aud}, nil
}

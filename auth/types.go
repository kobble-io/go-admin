package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/valensto/kobble-go-sdk/utils"
)

import (
	"time"
)

type VerifyAccessTokenResult struct {
	UserID    string                      `json:"user_id"`
	ProjectID string                      `json:"project_id"`
	Claims    rawAccessTokenPayloadClaims `json:"claims"`
}

type idTokenUser struct {
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
	User   idTokenUser             `json:"user"`
	Claims rawIdTokenPayloadClaims `json:"claims"`
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

type projectCache struct {
	ProjectID string
}

type rawIdTokenPayloadClaims struct {
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

func (r rawIdTokenPayloadClaims) GetExpirationTime() (*jwt.NumericDate, error) {
	return jwt.NewNumericDate(time.Unix(r.Exp, 0)), nil
}

func (r rawIdTokenPayloadClaims) GetIssuedAt() (*jwt.NumericDate, error) {
	return jwt.NewNumericDate(time.Unix(r.Iat, 0)), nil
}

func (r rawIdTokenPayloadClaims) GetNotBefore() (*jwt.NumericDate, error) {
	return jwt.NewNumericDate(time.Unix(r.Nbf, 0)), nil
}

func (r rawIdTokenPayloadClaims) GetIssuer() (string, error) {
	return r.Iss, nil
}

func (r rawIdTokenPayloadClaims) GetSubject() (string, error) {
	return r.Sub, nil
}

func (r rawIdTokenPayloadClaims) GetAudience() (jwt.ClaimStrings, error) {
	return jwt.ClaimStrings{r.Aud}, nil
}

type rawAccessTokenPayloadClaims struct {
	Sub       string `json:"sub"`
	ProjectID string `json:"project_id"`
	Exp       int64  `json:"exp"`
	Iat       int64  `json:"iat"`
	Nbf       int64  `json:"nbf"`
	Iss       string `json:"iss"`
	Aud       string `json:"aud"`
}

func (r rawAccessTokenPayloadClaims) GetExpirationTime() (*jwt.NumericDate, error) {
	return jwt.NewNumericDate(time.Unix(r.Exp, 0)), nil
}

func (r rawAccessTokenPayloadClaims) GetIssuedAt() (*jwt.NumericDate, error) {
	return jwt.NewNumericDate(time.Unix(r.Iat, 0)), nil
}

func (r rawAccessTokenPayloadClaims) GetNotBefore() (*jwt.NumericDate, error) {
	return jwt.NewNumericDate(time.Unix(r.Nbf, 0)), nil
}

func (r rawAccessTokenPayloadClaims) GetIssuer() (string, error) {
	return r.Iss, nil
}

func (r rawAccessTokenPayloadClaims) GetSubject() (string, error) {
	return r.Sub, nil
}

func (r rawAccessTokenPayloadClaims) GetAudience() (jwt.ClaimStrings, error) {
	return jwt.ClaimStrings{r.Aud}, nil
}

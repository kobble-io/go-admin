package auth

import (
	"fmt"
	"github.com/MicahParks/keyfunc"
	"github.com/golang-jwt/jwt/v4"
	"github.com/valensto/kobble-go-sdk/utils"
	"strings"
	"time"
)

type KobbleAuth struct {
	issuer       string
	projectCache *utils.Cache[ProjectCache]
	config       Config
}

func NewKobbleAuth(conf Config) *KobbleAuth {
	defaultTtl := 50 * time.Minute / time.Second
	return &KobbleAuth{
		issuer: "https://kobble.io",
		projectCache: utils.NewCache[ProjectCache](utils.CacheConfig{
			DefaultTtl: &defaultTtl,
		}),
		config: conf,
	}
}

func (auth KobbleAuth) getProjectId() (string, error) {
	cacheProject := auth.projectCache.Get("default")
	if cacheProject != nil {
		return cacheProject.ProjectID, nil
	}

	var whoami Whoami
	err := auth.config.Http.GetJson("/auth/whoami", nil, &whoami)
	if err != nil {
		return "", err
	}

	auth.projectCache.Set("projectId", ProjectCache{ProjectID: whoami.ProjectId}, nil)
	return whoami.ProjectId, nil
}

func (auth KobbleAuth) VerifyAccessToken(token string) (VerifyAccessTokenResult, error) {
	var result VerifyAccessTokenResult
	projectId, err := auth.getProjectId()
	if err != nil {
		return result, fmt.Errorf("get project id err: %w", newAccessTokenVerificationError(err))
	}

	jwksURL := fmt.Sprintf("%s/discovery/p/%s/apps/keys", auth.config.BaseURL, projectId)
	jwks, err := keyfunc.Get(jwksURL, keyfunc.Options{
		RefreshErrorHandler: func(err error) {
			fmt.Printf("There was an error with the jwt.Keyfunc\nError: %s", err.Error())
		},
		RefreshInterval:   time.Hour,
		RefreshRateLimit:  time.Minute * 5,
		RefreshTimeout:    time.Second * 10,
		RefreshUnknownKID: true,
	})
	if err != nil {
		return VerifyAccessTokenResult{}, fmt.Errorf("failed to create JWKS from resource at %s: %v", jwksURL, err)
	}

	var rawClaims RawAccessTokenPayloadClaims
	tk, err := jwt.ParseWithClaims(token, &rawClaims, jwks.Keyfunc)
	if err != nil {
		return VerifyAccessTokenResult{}, fmt.Errorf("failed to parse token: %v", err)
	}

	if claims, ok := tk.Claims.(*RawAccessTokenPayloadClaims); ok && tk.Valid {
		return VerifyAccessTokenResult{
			UserID:    claims.Sub,
			ProjectID: claims.ProjectID,
			Claims:    *claims,
		}, nil
	}

	return VerifyAccessTokenResult{}, fmt.Errorf("invalid token")
}

func parseClaimsDate(date string) (string, error) {
	parseFormat := "Mon Jan 2 2006 15:04:05 GMT-0700"
	parsed, err := time.Parse(parseFormat, strings.Replace(date, " (Coordinated Universal Time)", "", 1))
	if err != nil {
		return "", fmt.Errorf("failed to parse date: %w", err)
	}

	return parsed.Format(time.RFC3339), nil
}

func (auth KobbleAuth) VerifyIdToken(token string) (VerifyIdTokenResult, error) {
	var result VerifyIdTokenResult
	projectId, err := auth.getProjectId()
	if err != nil {
		return result, fmt.Errorf("get project id err: %w", newAccessTokenVerificationError(err))
	}

	jwksURL := fmt.Sprintf("%s/discovery/p/%s/apps/keys", auth.config.BaseURL, projectId)
	jwks, err := keyfunc.Get(jwksURL, keyfunc.Options{
		RefreshErrorHandler: func(err error) {
			fmt.Printf("There was an error with the jwt.Keyfunc\nError: %s", err.Error())
		},
		RefreshInterval:   time.Hour,
		RefreshRateLimit:  time.Minute * 5,
		RefreshTimeout:    time.Second * 10,
		RefreshUnknownKID: true,
	})
	if err != nil {
		return VerifyIdTokenResult{}, fmt.Errorf("failed to create JWKS from resource at %s: %v", jwksURL, err)
	}

	var rawClaims RawIdTokenPayloadClaims
	tk, err := jwt.ParseWithClaims(token, &rawClaims, jwks.Keyfunc)
	if err != nil {
		return VerifyIdTokenResult{}, fmt.Errorf("failed to parse token: %v", err)
	}

	if claims, ok := tk.Claims.(*RawIdTokenPayloadClaims); ok && tk.Valid {
		updatedAt, err := parseClaimsDate(claims.UpdatedAt)
		if err != nil {
			return VerifyIdTokenResult{}, err
		}

		createdAt, err := parseClaimsDate(claims.CreatedAt)
		if err != nil {
			return VerifyIdTokenResult{}, err
		}

		return VerifyIdTokenResult{
			UserID: claims.Sub,
			User: IdTokenUser{
				ID:         claims.ID,
				Email:      claims.Email,
				Name:       claims.Name,
				PictureURL: claims.PictureURL,
				IsVerified: claims.IsVerified,
				StripeID:   claims.StripeID,
				UpdatedAt:  updatedAt,
				CreatedAt:  createdAt,
			},
			Claims: *claims,
		}, nil
	}

	return VerifyIdTokenResult{}, fmt.Errorf("invalid token")
}

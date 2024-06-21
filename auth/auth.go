package auth

import (
	"fmt"
	"github.com/valensto/kobble-go-sdk/utils"
	"time"
)

type KobbleAuth struct {
	issuer       string
	projectCache *utils.Cache[ProjectCache]
	config       Config
}

func NewKobbleAuth(http *utils.HttpClient) *KobbleAuth {
	defaultTtl := 50 * time.Minute / time.Second
	return &KobbleAuth{
		issuer: "https://kobble.io",
		projectCache: utils.NewCache[ProjectCache](utils.CacheConfig{
			DefaultTtl: &defaultTtl,
		}),
		config: Config{http: http},
	}
}

func (auth KobbleAuth) getProjectId() (string, error) {
	cacheProject := auth.projectCache.Get("default")
	if cacheProject != nil {
		return cacheProject.ProjectID, nil
	}

	var whoami Whoami
	err := auth.config.http.GetJson("/auth/whoami", nil, &whoami)
	if err != nil {
		return "", err
	}

	auth.projectCache.Set("projectId", ProjectCache{ProjectID: whoami.ProjectId}, nil)
	return whoami.ProjectId, nil
}

func (auth KobbleAuth) VerifyAccessToken(token string, options VerifyAccessTokenOptions) (VerifyAccessTokenResult, error) {
	var result VerifyAccessTokenResult
	projectId, err := auth.getProjectId()
	if err != nil {
		return result, newAccessTokenVerificationError(err)
	}

	jwks, err := getJWKSet(fmt.Sprintf("%s/discovery/p/%s/apps/keys", auth.config.baseURL, projectId))
	if err != nil {
		return VerifyAccessTokenResult{}, newAccessTokenVerificationError(err)
	}

	res, err := jwtVerify[RawAccessTokenPayloadClaims](token, jwks, options)
	if err != nil {
		return VerifyAccessTokenResult{}, newAccessTokenVerificationError(err)
	}

	return VerifyAccessTokenResult{
		UserID:    res.Sub,
		ProjectID: res.ProjectID,
		Claims:    *res,
	}, nil
}

func (auth KobbleAuth) VerifyIdToken(token string, options VerifyAccessTokenOptions) (VerifyIdTokenResult, error) {
	var result VerifyIdTokenResult
	projectId, err := auth.getProjectId()
	if err != nil {
		return result, newIdTokenVerificationError(err)
	}

	jwks, err := getJWKSet(fmt.Sprintf("%s/discovery/p/%s/apps/keys", auth.config.baseURL, projectId))
	if err != nil {
		return VerifyIdTokenResult{}, newIdTokenVerificationError(err)
	}

	res, err := jwtVerify[RawIdTokenPayloadClaims](token, jwks, options)
	if err != nil {
		return VerifyIdTokenResult{}, newIdTokenVerificationError(err)
	}

	updateAt, err := time.Parse(time.RFC3339, res.UpdatedAt)
	if err != nil {
		return VerifyIdTokenResult{}, newIdTokenVerificationError(err)
	}

	createdAt, err := time.Parse(time.RFC3339, res.CreatedAt)
	if err != nil {
		return VerifyIdTokenResult{}, newIdTokenVerificationError(err)
	}

	return VerifyIdTokenResult{
		UserID: res.Sub,
		User: IdTokenUser{
			ID:         res.ID,
			Email:      res.Email,
			Name:       res.Name,
			PictureURL: res.PictureURL,
			IsVerified: res.IsVerified,
			StripeID:   res.StripeID,
			UpdatedAt:  updateAt,
			CreatedAt:  createdAt,
		},
		Claims: *res,
	}, nil
}

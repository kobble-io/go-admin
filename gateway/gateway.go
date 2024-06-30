package gateway

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"github.com/kobble-io/go-admin/utils"
	"net/http"
	"time"
)

// KobbleGateway is the struct that holds the configuration for the gateway service
type KobbleGateway struct {
	config   Config
	keyCache *utils.Cache[keyInfo]
	issuer   string
}

// NewKobbleGateway creates a new instance of KobbleGateway
func NewKobbleGateway(config Config) *KobbleGateway {
	defaultTtl := 15 * time.Minute / time.Second
	return &KobbleGateway{
		config: config,
		keyCache: utils.NewCache[keyInfo](utils.CacheConfig{
			DefaultTtl: &defaultTtl,
		}),
		issuer: "gateway.kobble.io",
	}
}

func (k *KobbleGateway) fetchKeyInfo() (*keyInfo, error) {
	var result struct {
		Pem       string `json:"pem"`
		ProjectID string `json:"project_id"`
	}
	err := k.config.Http.GetJson("/gateway/getPublicKey", nil, &result, http.StatusOK)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode([]byte(result.Pem))
	if block == nil || block.Type != "PUBLIC KEY" {
		return nil, errors.New("failed to decode PEM block containing public key")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	ecdsaPub, ok := pub.(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.New("not ECDSA public key")
	}

	return &keyInfo{
		Key:       ecdsaPub,
		ProjectID: result.ProjectID,
	}, nil
}

func (k *KobbleGateway) getKeyInfo() (*keyInfo, error) {
	info := k.keyCache.Get("default")
	if info != nil {
		return info, nil
	}

	data, err := k.fetchKeyInfo()
	if err != nil {
		return nil, err
	}

	k.keyCache.Set("default", *data, nil)
	return data, nil
}

// ParseToken verify and parse the payload of a Kobble gateway token.
// This kind of token is sent to your backend on each request that is handled by your Kobble project.
//
// By default, this function will verify that:
//   - The required claims are present
//   - The 'aud' claim matches the current project ID
//   - The 'iss' claim matches the issuer used by Kobble to forge such tokens
//   - The token is not expired
//   - The signature is valid (i.e., that this token has not been tampered with and is intended for your project)
//
// Although it is not recommended, some of these verifications can be skipped by passing special options.
func (k *KobbleGateway) ParseToken(tokenString string, options ParseTokenOptions) (TokenPayload, error) {
	ki, err := k.getKeyInfo()
	if err != nil {
		return TokenPayload{}, err
	}

	payload, err := utils.VerifyJwt(tokenString, ki.Key, utils.VerifyJwtOptions{
		VerifyAud:       options.VerifyAud,
		VerifyExp:       options.VerifyExp,
		VerifySignature: options.VerifySignature,
		VerifyIss:       options.VerifyIss,
		Iss:             k.issuer,
		Audience:        ki.ProjectID,
		RequiredClaims:  []string{"iat", "exp", "iss", "sub", "aud", "user"},
	})
	if err != nil {
		return TokenPayload{}, err
	}

	b, err := json.Marshal(payload)
	if err != nil {
		return TokenPayload{}, err
	}

	var raw TokenPayload
	err = json.Unmarshal(b, &raw)
	if err != nil {
		return TokenPayload{}, err
	}

	return raw, nil
}

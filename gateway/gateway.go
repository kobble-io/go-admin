package gateway

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"github.com/valensto/kobble-go-sdk/utils"
	"time"
)

type KobbleGateway struct {
	config   Config
	keyCache *utils.Cache[KeyInfo]
	issuer   string
}

func NewKobbleGateway(config Config) *KobbleGateway {
	defaultTtl := 15 * time.Minute / time.Second
	return &KobbleGateway{
		config: config,
		keyCache: utils.NewCache[KeyInfo](utils.CacheConfig{
			DefaultTtl: &defaultTtl,
		}),
		issuer: "gateway.kobble.io",
	}
}

func (k *KobbleGateway) fetchKeyInfo() (*KeyInfo, error) {
	var result struct {
		Pem       string `json:"pem"`
		ProjectID string `json:"project_id"`
	}
	err := k.config.Http.GetJson("/gateway/getPublicKey", nil, &result)
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

	return &KeyInfo{
		Key:       ecdsaPub,
		ProjectID: result.ProjectID,
	}, nil
}

func (k *KobbleGateway) getKeyInfo() (*KeyInfo, error) {
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

func (k *KobbleGateway) parseToken(tokenString string, options ParseTokenOptions) (TokenPayload, error) {
	keyInfo, err := k.getKeyInfo()
	if err != nil {
		return TokenPayload{}, err
	}

	payload, err := utils.VerifyJwt(tokenString, keyInfo.Key, utils.VerifyJwtOptions{
		VerifyAud:       options.VerifyAud,
		VerifyExp:       options.VerifyExp,
		VerifySignature: options.VerifySignature,
		VerifyIss:       options.VerifyIss,
		Iss:             k.issuer,
		Audience:        keyInfo.ProjectID,
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

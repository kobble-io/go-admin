package auth

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"io/ioutil"
	"math/big"
	"net/http"
	"sync"
	"time"
)

type jwk struct {
	Kid string `json:"kid"`
	Kty string `json:"kty"`
	Use string `json:"use"`
	N   string `json:"n"`
	E   string `json:"e"`
}

type jwksKeys struct {
	Keys []jwk `json:"keys"`
}

type cachedJWKSet struct {
	jwks      jwksKeys
	expiresAt time.Time
}

var jwksCache = struct {
	sync.RWMutex
	data map[string]cachedJWKSet
}{data: make(map[string]cachedJWKSet)}

func (jwk jwk) parseRSAPublicKeyFromJWK() (*rsa.PublicKey, error) {
	nb, err := base64.RawURLEncoding.DecodeString(jwk.N)
	if err != nil {
		return nil, err
	}
	eb, err := base64.RawURLEncoding.DecodeString(jwk.E)
	if err != nil {
		return nil, err
	}
	if len(eb) < 4 {
		b := make([]byte, 4)
		copy(b[4-len(eb):], eb)
		eb = b
	}
	e := int(big.NewInt(0).SetBytes(eb).Int64())
	if e < 0 {
		return nil, errors.New("invalid public exponent")
	}
	key := &rsa.PublicKey{
		N: big.NewInt(0).SetBytes(nb),
		E: e,
	}
	return key, nil
}

func jwtVerify[T jwt.Claims](tokenString string, jwks jwksKeys, options VerifyAccessTokenOptions) (*T, error) {
	keyFunc := func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		kid, ok := token.Header["kid"].(string)
		if !ok {
			return nil, errors.New("kid header not found")
		}

		for _, key := range jwks.Keys {
			if key.Kid == kid {
				return key.parseRSAPublicKeyFromJWK()
			}
		}

		return nil, errors.New("unable to find appropriate key")
	}

	token, err := jwt.ParseWithClaims(tokenString, &RawAccessTokenPayloadClaims{}, keyFunc, jwt.WithAudience(*options.ApplicationID), jwt.WithIssuer(*options.Issuer))
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(T); ok && token.Valid {
		return &claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

func getJWKSet(url string) (jwksKeys, error) {
	jwksCache.RLock()
	if cached, found := jwksCache.data[url]; found && cached.expiresAt.After(time.Now()) {
		jwksCache.RUnlock()
		return cached.jwks, nil
	}
	jwksCache.RUnlock()

	resp, err := http.Get(url)
	if err != nil {
		return jwksKeys{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return jwksKeys{}, errors.New("failed to fetch jwks")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return jwksKeys{}, err
	}

	var jwks jwksKeys
	if err := json.Unmarshal(body, &jwks); err != nil {
		return jwks, err
	}

	cacheDuration := 15 * time.Minute
	jwksCache.Lock()
	jwksCache.data[url] = cachedJWKSet{
		jwks:      jwks,
		expiresAt: time.Now().Add(cacheDuration),
	}
	jwksCache.Unlock()

	return jwks, nil
}

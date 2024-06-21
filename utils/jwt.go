package utils

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"
)

type VerifyJwtOptions struct {
	VerifyAud       bool
	VerifyExp       bool
	VerifySignature bool
	VerifyIss       bool
	Iss             string
	Audience        string
	RequiredClaims  []string
}

type VerifyJwtError struct {
	Message string
}

func (e *VerifyJwtError) Error() string {
	return e.Message
}

func verifyES256Signature(data, signature []byte, key *ecdsa.PublicKey) bool {
	hash := sha256.Sum256(data)
	return ecdsa.VerifyASN1(key, hash[:], signature)
}

func VerifyJwt(tokenString string, key *ecdsa.PublicKey, options VerifyJwtOptions) (map[string]any, error) {
	ss := splitToken(tokenString)
	if len(ss) != 3 {
		return nil, &VerifyJwtError{Message: fmt.Sprintf("Expected three dot separated segments but found %d. Did you pass a valid JWT?", len(ss))}
	}

	encodedHeader, encodedPayload, encodedSignature := ss[0], ss[1], ss[2]
	payload, err := decodePayload(string(encodedPayload))
	if err != nil {
		return nil, err
	}

	if err := checkRequiredClaims(payload, options.RequiredClaims); err != nil {
		return nil, err
	}

	if options.VerifyExp {
		if err := checkExpiration(payload); err != nil {
			return nil, err
		}
	}

	if options.VerifyAud {
		if err := checkAudience(payload, options.Audience); err != nil {
			return nil, err
		}
	}

	if options.VerifyIss {
		if err := checkIssuer(payload, options.Iss); err != nil {
			return nil, err
		}
	}

	if options.VerifySignature {
		if err := verifySignature(encodedHeader, encodedPayload, encodedSignature, key); err != nil {
			return nil, err
		}
	}

	return payload, nil
}

func splitToken(token string) [][]byte {
	return bytes.SplitN([]byte(token), []byte{'.'}, 3)
}

func decodePayload(encodedPayload string) (map[string]any, error) {
	payloadBytes, err := base64.RawURLEncoding.DecodeString(encodedPayload)
	if err != nil {
		return nil, &VerifyJwtError{Message: "Failed to decode payload"}
	}

	var payload map[string]any
	if err := json.Unmarshal(payloadBytes, &payload); err != nil {
		return nil, &VerifyJwtError{Message: "Failed to unmarshal payload"}
	}

	return payload, nil
}

func checkRequiredClaims(payload map[string]any, requiredClaims []string) error {
	if requiredClaims != nil {
		missing := []string{}
		for _, claim := range requiredClaims {
			if _, ok := payload[claim]; !ok {
				missing = append(missing, claim)
			}
		}

		if len(missing) > 0 {
			return &VerifyJwtError{Message: fmt.Sprintf("The following required claims are missing from the token payload: %s", missing)}
		}
	}
	return nil
}

func checkExpiration(payload map[string]any) error {
	now := time.Now().Unix()
	exp, ok := payload["exp"].(float64)
	if !ok {
		return &VerifyJwtError{Message: "Missing or invalid exp claim"}
	}

	if now > int64(exp) {
		return &VerifyJwtError{Message: fmt.Sprintf("This token expired on %s which is BEFORE the current datetime %s", time.Unix(int64(exp), 0), time.Now())}
	}

	return nil
}

func checkAudience(payload map[string]any, audience string) error {
	if payload["aud"] != audience {
		return &VerifyJwtError{Message: fmt.Sprintf("Invalid audience: expected %s but got %s", audience, payload["aud"])}
	}
	return nil
}

func checkIssuer(payload map[string]any, issuer string) error {
	if payload["iss"] != issuer {
		return &VerifyJwtError{Message: fmt.Sprintf("Invalid issuer: expected %s but got %s", issuer, payload["iss"])}
	}
	return nil
}

func verifySignature(encodedHeader, encodedPayload, encodedSignature []byte, key *ecdsa.PublicKey) error {
	data := fmt.Sprintf("%s.%s", encodedHeader, encodedPayload)
	signature, err := base64.RawURLEncoding.DecodeString(string(encodedSignature))
	if err != nil {
		return &VerifyJwtError{Message: "Failed to decode signature"}
	}

	if !verifyES256Signature([]byte(data), signature, key) {
		return &VerifyJwtError{Message: "Invalid signature"}
	}

	return nil
}

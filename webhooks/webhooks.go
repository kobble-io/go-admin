package webhooks

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
)

type KobbleWebhooks struct{}

func NewKobbleWebhooks() *KobbleWebhooks {
	return &KobbleWebhooks{}
}

func (k *KobbleWebhooks) serializeBody(body any) ([]byte, error) {
	switch v := body.(type) {
	case []byte:
		return v, nil
	case string:
		return []byte(v), nil
	case map[string]any, []any:
		return json.Marshal(v)
	default:
		return json.Marshal(v)
	}
}

func (k *KobbleWebhooks) constructEvent(body any, signature string, secret string) (WebhookEvent, error) {
	serializedBody, err := k.serializeBody(body)
	if err != nil {
		return WebhookEvent{}, err
	}

	constructedSignature := k.createHmacSignature(serializedBody, secret)
	if signature != constructedSignature {
		return WebhookEvent{}, errors.New("signature verification failed. Did you pass the correct secret?")
	}

	var event WebhookEvent
	err = json.Unmarshal(serializedBody, &event)
	if err != nil {
		return WebhookEvent{}, err
	}

	return event, nil
}

func (k *KobbleWebhooks) createHmacSignature(body []byte, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write(body)
	return hex.EncodeToString(h.Sum(nil))
}

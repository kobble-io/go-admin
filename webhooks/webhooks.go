package webhooks

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
)

// KobbleWebhooks is a struct that provides methods to work with webhooks.
type KobbleWebhooks struct{}

// NewKobbleWebhooks creates a new instance of the KobbleWebhooks struct.
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

// Construct a webhook event payload and verify its integrity.
// The `body` parameter is eventually serialized to a `Buffer` object in order to compute
// the signature.
// A `Buffer` can therefore be passed directly, but other types are accepted as well:
//
//   - A Go struct is serialized using `json.Marshal`.
//   - For any other type, native string conversion is attempted. The result is assumed to be UTF-8 encoded.
//
// The expected `signature` is the one sent in the webhook header `Kobble-Signature`.
// The `secret` is the one associated with the webhook expected to receive the event.
// The fully typesafe payload is returned if the signature is valid.
func (k *KobbleWebhooks) constructEvent(body any, signature string, secret string) (WebhookEvent, error) {
	serializedBody, err := k.serializeBody(body)
	if err != nil {
		return WebhookEvent{}, err
	}

	constructedSignature := k.createHmacSignature(serializedBody, secret)
	if signature != constructedSignature {
		return WebhookEvent{}, newWebhookConstructEventError("Signature verification failed. Did you pass the correct secret?")
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

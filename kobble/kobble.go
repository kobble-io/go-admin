package kobble

import (
	"github.com/kobble-io/go-admin/auth"
	"github.com/kobble-io/go-admin/gateway"
	"github.com/kobble-io/go-admin/users"
	"github.com/kobble-io/go-admin/utils"
	"github.com/kobble-io/go-admin/webhooks"
)

const DefaultBaseUrl = "https://sdk.kobble.io"

// Kobble is the main struct for the Kobble SDK.
//
//   - http is the HTTP client used to make requests to the Kobble API.
//   - Gateway is the gateway service that allows you to interact with the gateway service.
//   - Users is the users service that allows you to interact with the users service.
//   - Webhooks is the webhooks service that allows you to interact with the webhooks service.
//   - Auth is the auth service that allows you to interact with the auth service.
type Kobble struct {
	http     *utils.HttpClient
	Gateway  *gateway.KobbleGateway
	Users    *users.KobbleUsers
	Webhooks *webhooks.KobbleWebhooks
	Auth     *auth.KobbleAuth
}

// New is the constructor for the Kobble SDK.
//
// secret is the secret key that you can generate from the Kobble dashboard.
// options is the configuration for the SDK.
//
// Returns a new instance of the Kobble SDK.
func New(secret string, options Options) *Kobble {
	baseURL := DefaultBaseUrl
	if options.BaseApiUrl != nil {
		baseURL = *options.BaseApiUrl
	}
	http := utils.NewHttpClient(utils.HttpClientConfig{
		BaseURL: baseURL,
		Secret:  secret,
	})
	return &Kobble{
		http:     http,
		Gateway:  gateway.NewKobbleGateway(gateway.Config{Http: http}),
		Users:    users.NewKobbleUsers(users.Config{Http: http}),
		Webhooks: webhooks.NewKobbleWebhooks(),
		Auth: auth.NewKobbleAuth(auth.Config{
			Http:    http,
			BaseURL: baseURL,
		}),
	}
}

// Whoami get the project and the user associated with the SDK secret used to authenticate.
//
// The user ID is the one of the user that created the secret.
func (k Kobble) Whoami() (auth.Whoami, error) {
	var whoami auth.Whoami
	err := k.http.GetJson("/auth/whoami", nil, &whoami)
	if err != nil {
		return auth.Whoami{}, err
	}

	return whoami, nil
}

// Ping the Kobble SDK API to ensure that it is reachable and that you are properly authenticated.
//
// Returns true if successful.
func (k Kobble) Ping() bool {
	err := k.http.GetJson("/ping", nil, nil)
	return err == nil
}

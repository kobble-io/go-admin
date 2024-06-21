package kobble

import (
	"github.com/valensto/kobble-go-sdk/auth"
	"github.com/valensto/kobble-go-sdk/gateway"
	"github.com/valensto/kobble-go-sdk/users"
	"github.com/valensto/kobble-go-sdk/utils"
	"github.com/valensto/kobble-go-sdk/webhooks"
)

const DefaultBaseUrl = "https://sdk.kobble.io"

type Kobble struct {
	http     *utils.HttpClient
	Gateway  *gateway.KobbleGateway
	Users    *users.KobbleUsers
	Webhooks *webhooks.KobbleWebhooks
	Auth     *auth.KobbleAuth
}

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
		Auth:     auth.NewKobbleAuth(http),
	}
}

func (k Kobble) Whoami() (auth.Whoami, error) {
	var whoami auth.Whoami
	err := k.http.GetJson("/auth/whoami", nil, &whoami)
	if err != nil {
		return auth.Whoami{}, err
	}

	return whoami, nil
}

func (k Kobble) Ping() bool {
	err := k.http.GetJson("/ping", nil, nil)
	return err == nil
}

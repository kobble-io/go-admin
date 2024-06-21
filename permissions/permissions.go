package permissions

import "github.com/valensto/kobble-go-sdk/utils"

type KobblePermission struct {
	config permissionConfig
}

func NewPermission(http *utils.HttpClient) *KobblePermission {
	return &KobblePermission{
		config: permissionConfig{
			http: http,
		},
	}
}

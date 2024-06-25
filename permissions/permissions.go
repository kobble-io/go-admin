package permissions

import "github.com/valensto/kobble-go-sdk/utils"

// KobblePermission is the struct that holds the configuration for the permission service
type KobblePermission struct {
	config permissionConfig
}

// NewPermission creates a new instance of KobblePermission
func NewPermission(http *utils.HttpClient) *KobblePermission {
	return &KobblePermission{
		config: permissionConfig{
			http: http,
		},
	}
}

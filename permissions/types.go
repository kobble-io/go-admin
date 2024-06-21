package permissions

import "github.com/valensto/kobble-go-sdk/utils"

type (
	permissionConfig struct {
		http *utils.HttpClient
	}

	Permission struct {
		Name string
	}
)

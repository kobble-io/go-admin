package permissions

import "github.com/kobble-io/go-admin/utils"

type (
	permissionConfig struct {
		http *utils.HttpClient
	}

	Permission struct {
		Name string
	}
)

package quotas

import "github.com/valensto/kobble-go-sdk/utils"

type Quotas struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Limit int    `json:"limit"`
}

type Config struct {
	http *utils.HttpClient
}

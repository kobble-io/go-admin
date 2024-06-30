package users

import (
	"github.com/kobble-io/go-admin/utils"
	"time"
)

type User struct {
	ID         string         `json:"id"`
	Email      string         `json:"email"`
	Name       *string        `json:"name"`
	CreatedAt  time.Time      `json:"created_at"`
	IsVerified bool           `json:"is_verified"`
	Metadata   map[string]any `json:"metadata"`
}

type UserActiveProduct struct {
	Product
}

type QuotaUsage struct {
	Name      string    `json:"name"`
	Usage     int       `json:"usage"`
	ExpiresAt time.Time `json:"expires_at"`
	Remaining *int      `json:"remaining"`
	Limit     *int      `json:"limit"`
}

type Product struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

type Config struct {
	Http *utils.HttpClient
}

type ApiUser struct {
	ID         string         `json:"id"`
	Email      string         `json:"email"`
	Name       *string        `json:"name"`
	CreatedAt  time.Time      `json:"created_at"`
	IsVerified bool           `json:"is_verified"`
	Metadata   map[string]any `json:"metadata"`
}

type ApiQuota struct {
	Name      string    `json:"name"`
	Usage     int       `json:"usage"`
	ExpiresAt time.Time `json:"expires_at"`
	Remaining int       `json:"remaining"`
	Limit     int       `json:"limit"`
}

type ListApiQuotaResponse struct {
	Quotas []ApiQuota `json:"quotas"`
}

type ApiPermission struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ListApiPermissionResponse struct {
	Permissions []ApiPermission `json:"permissions"`
}

type CreateUserPayload struct {
	Email                     string         `json:"email,omitempty"`
	Name                      string         `json:"name,omitempty"`
	PhoneNumber               string         `json:"phone_number,omitempty"`
	Metadata                  map[string]any `json:"metadata,omitempty"`
	MarkEmailAsVerified       bool           `json:"mark_email_as_verified,omitempty"`
	MarkPhoneNumberAsVerified bool           `json:"mark_phone_number_as_verified,omitempty"`
}

type UrlLink struct {
	Url       string    `json:"url"`
	ExpiresAt time.Time `json:"expires_at"`
}

type IsAllowedPayload struct {
	PermissionNames []string `json:"permissionNames"`
	QuotaNames      []string `json:"quotaNames"`
}

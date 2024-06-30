package users

import (
	"github.com/kobble-io/go-admin/common"
	"github.com/kobble-io/go-admin/permissions"
	"github.com/kobble-io/go-admin/utils"
	"strconv"
	"time"
)

// KobbleUsers is the client to interact with the users API on your Kobble instance.
//
// You can use this client to create, retrieve, update, and delete users on your Kobble instance.
// You can also use it to manage user metadata, permissions, and quotas.
type KobbleUsers struct {
	config           Config
	permissionsCache *utils.Cache[[]permissions.Permission]
	quotasCache      *utils.Cache[[]QuotaUsage]
}

// NewKobbleUsers creates a new instance of the KobbleUsers client.
//
// @param config - The configuration for the Kobble instance.
func NewKobbleUsers(config Config) *KobbleUsers {
	defaultTtl := 1 * time.Minute / time.Second
	return &KobbleUsers{
		config: config,
		permissionsCache: utils.NewCache[[]permissions.Permission](utils.CacheConfig{
			DefaultTtl: &defaultTtl,
		}),
		quotasCache: utils.NewCache[[]QuotaUsage](utils.CacheConfig{
			DefaultTtl: &defaultTtl,
		}),
	}
}

func (k KobbleUsers) userCacheKey(userId string) string {
	return "user:" + userId
}

func (k KobbleUsers) transformApiUser(apiUser ApiUser) *User {
	metadata := make(map[string]any)
	if apiUser.Metadata != nil {
		metadata = apiUser.Metadata
	}
	return &User{
		ID:         apiUser.ID,
		Email:      apiUser.Email,
		Name:       apiUser.Name,
		CreatedAt:  apiUser.CreatedAt,
		IsVerified: apiUser.IsVerified,
		Metadata:   metadata,
	}
}

// CreateLoginLink creates a login link for a user.
//
//   - @param userId - The unique identifier for the user to create a login link for.
func (k KobbleUsers) CreateLoginLink(userId string) (UrlLink, error) {
	var result UrlLink
	err := k.config.Http.PostJson("/users/mintLoginLink", map[string]string{
		"userId": userId,
	}, &result)
	return result, err
}

// Create a new user on your Kobble instance manually.
//
// While both email and phoneNumber are optional, at least one of them must be provided.
// If an email is provided, it will be marked as verified by default.
// Note that the phone number should be in E.164 format (e.g. +14155552671). Other formats will be rejected.
func (k KobbleUsers) Create(payload CreateUserPayload) (*User, error) {
	var result ApiUser
	err := k.config.Http.PostJson("/users/create", payload, &result)
	if err != nil {
		return nil, err
	}
	return k.transformApiUser(result), nil
}

type GetUserOptions struct {
	IncludeMetadata bool `json:"include_metadata,omitempty"`
}

// GetById fetches a user by their ID.
//
// You can also include the user's metadata in the response by setting the `IncludeMetadata` option to `true`.
func (k KobbleUsers) GetById(userId string, options *GetUserOptions) (*User, error) {
	includeMetadata := false
	if options != nil {
		includeMetadata = options.IncludeMetadata
	}
	var result ApiUser
	err := k.config.Http.GetJson("/users/findById", map[string]string{
		"userId":          userId,
		"includeMetadata": strconv.FormatBool(includeMetadata),
	}, &result)
	if err != nil {
		return nil, err
	}
	return k.transformApiUser(result), nil
}

// GetByEmail fetches a user by their email.
//
// You can also include the user's metadata in the response by setting the `IncludeMetadata` option to `true`.
func (k KobbleUsers) GetByEmail(email string, options *GetUserOptions) (*User, error) {
	includeMetadata := false
	if options != nil {
		includeMetadata = options.IncludeMetadata
	}
	var result ApiUser
	err := k.config.Http.GetJson("/users/findByEmail", map[string]string{
		"email":           email,
		"includeMetadata": strconv.FormatBool(includeMetadata),
	}, &result)
	if err != nil {
		return nil, err
	}
	return k.transformApiUser(result), nil
}

// GetByPhoneNumber fetches a user by their phone number.
//
// Note that the phone number should be in E.164 format (e.g. +14155552671). Other formats will be rejected.
// You can also include the user's metadata in the response by setting the `IncludeMetadata` option to `true`.
func (k KobbleUsers) GetByPhoneNumber(phoneNumber string, options *GetUserOptions) (*User, error) {
	includeMetadata := false
	if options != nil {
		includeMetadata = options.IncludeMetadata
	}
	var result ApiUser
	err := k.config.Http.GetJson("/users/findByPhoneNumber", map[string]string{
		"phoneNumber":     phoneNumber,
		"includeMetadata": strconv.FormatBool(includeMetadata),
	}, &result)
	if err != nil {
		return nil, err
	}
	return k.transformApiUser(result), nil
}

type ListUsersOptions struct {
	Limit           int  `json:"limit,omitempty"`
	Page            int  `json:"page,omitempty"`
	IncludeMetadata bool `json:"include_metadata,omitempty"`
}

// FindByMetadata fetches users by their metadata.
//
// You can also include the user's metadata in the response by setting the `IncludeMetadata` option to `true`.
func (k KobbleUsers) FindByMetadata(metadata map[string]any, options *ListUsersOptions) (common.Pagination[User], error) {
	page, limit := 1, 50
	if options != nil {
		if options.Page > page {
			page = options.Page
		}

		if options.Limit != limit && options.Limit > 0 {
			limit = options.Limit
		}
	}

	var result common.Pagination[User]
	err := k.config.Http.PostJson("/users/findByMetadata", map[string]any{
		"metadata": metadata,
		"page":     page,
		"limit":    limit,
	}, &result)
	if err != nil {
		return common.Pagination[User]{}, err
	}

	return result, nil
}

// PatchMetadata updates a user's metadata.
func (k KobbleUsers) PatchMetadata(userId string, metadata map[string]any) (map[string]any, error) {
	err := k.config.Http.PostJson("/users/patchMetadata", map[string]any{
		"userId":   userId,
		"metadata": metadata,
	}, nil)
	if err != nil {
		return nil, err
	}

	return metadata, nil
}

// UpdateMetadata replaces a user's metadata.
func (k KobbleUsers) UpdateMetadata(userId string, metadata map[string]any) (map[string]any, error) {
	err := k.config.Http.PostJson("/users/updateMetadata", map[string]any{
		"userId":   userId,
		"metadata": metadata,
	}, nil)
	if err != nil {
		return nil, err
	}

	return metadata, nil
}

// ListAll fetches all users on your Kobble instance.
//
// Options:
//   - Page: The page number to fetch. Defaults to 1.
//   - Limit: The number of users to fetch per page. Defaults to 50.
//   - IncludeMetadata: Whether to include the user's metadata in the response. Defaults to false.
func (k KobbleUsers) ListAll(options *ListUsersOptions) (common.Pagination[User], error) {
	page, limit, includeMetadata := 1, 50, false
	if options != nil {
		if options.Page > page {
			page = options.Page
		}

		if options.Limit != limit && options.Limit > 0 {
			limit = options.Limit
		}

		includeMetadata = options.IncludeMetadata
	}

	var result common.Pagination[User]
	err := k.config.Http.GetJson("/users/list", map[string]string{
		"page":            strconv.Itoa(page),
		"limit":           strconv.Itoa(limit),
		"includeMetadata": strconv.FormatBool(includeMetadata),
	}, &result)
	if err != nil {
		return common.Pagination[User]{}, err
	}

	return result, nil
}

// GetActiveProducts retrieves the active product a given user is assigned to.
//
//   - @param userId - The unique identifier for the user whose active product is being retrieved.
//   - @returns UserActiveProduct or nil - The active product assigned to the user, or nil if the user has no active product.
func (k KobbleUsers) GetActiveProducts(userId string) (*UserActiveProduct, error) {
	var result []UserActiveProduct
	err := k.config.Http.GetJson("/users/listActiveProducts", map[string]string{
		"userId": userId,
	}, &result)
	if err != nil {
		return nil, err
	}

	if len(result) > 0 {
		return &result[0], nil
	}

	return nil, nil
}

func (k KobbleUsers) getCachedUserPerms(userId string) *[]permissions.Permission {
	key := k.userCacheKey(userId)
	perms := k.permissionsCache.Get(key)
	if perms != nil {
		return perms
	}

	return nil
}

func (k KobbleUsers) getCachedUserQuotas(userId string) *[]QuotaUsage {
	key := k.userCacheKey(userId)
	quotas := k.quotasCache.Get(key)
	if quotas != nil {
		return quotas
	}

	return nil
}

type ListQuotasOptions struct {
	NoCache bool
}

// ListQuotas retrieves the list of quota usages for a given user based on the product assigned to them.
//
//   - @param userId - The unique identifier for the user whose quota usage is being retrieved.
//   - @param noCache - Set to true to bypass cache and fetch fresh data. Default is false.
//   - @returns []QuotaUsage - An array of QuotaUsage objects, each representing a quota for the user.
func (k KobbleUsers) ListQuotas(userId string, opts *ListQuotasOptions) ([]QuotaUsage, error) {
	quotas := k.getCachedUserQuotas(userId)
	if opts != nil && !opts.NoCache && quotas != nil {
		return *quotas, nil
	}

	var result ListApiQuotaResponse
	err := k.config.Http.GetJson("/users/listQuotas", map[string]string{
		"userId": userId,
	}, &result)
	if err != nil {
		return nil, err
	}

	var quotasUsages []QuotaUsage
	for _, quota := range result.Quotas {
		quotasUsages = append(quotasUsages, QuotaUsage{
			Name:      quota.Name,
			Usage:     quota.Usage,
			ExpiresAt: quota.ExpiresAt,
			Remaining: &quota.Remaining,
			Limit:     &quota.Limit,
		})
	}

	return quotasUsages, nil
}

type IncrementQuotaOptions struct {
	IncrementBy int
}

// IncrementQuotaUsage asynchronously increments the quota usage for a specific user and quota.
//
//		This function allows incrementing a user's quota usage by a specified amount, which defaults to 1 if not provided.
//
//	 - @param userId - The unique identifier for the user whose quota is being incremented.
//	 - @param quotaName - The name of the quota to increment.
//	 - @param incrementBy - The amount by which to increment the quota usage. Optional and defaults to 1.
func (k KobbleUsers) IncrementQuotaUsage(userId string, quotaName string, opts *IncrementQuotaOptions) error {
	inc := 1
	if opts != nil {
		inc = opts.IncrementBy
	}
	err := k.config.Http.PostJson("/quotas/incrementUsage", map[string]any{
		"userId":      userId,
		"quotaName":   quotaName,
		"incrementBy": inc,
	}, nil)
	if err != nil {
		return err
	}

	return nil
}

type DecrementQuotaOptions struct {
	DecrementBy int
}

// DecrementQuotaUsage asynchronously decrements the quota usage for a specific user and quota.
//
//		This function allows decrementing a user's quota usage by a specified amount, which defaults to 1 if not provided.
//
//	 - @param userId - The unique identifier for the user whose quota is being decremented.
//	 - @param quotaName - The name of the quota to decrement.
//	 - @param decrementBy - The amount by which to decrement the quota usage. Optional and defaults to 1.
func (k KobbleUsers) DecrementQuotaUsage(userId string, quotaName string, opts *DecrementQuotaOptions) error {
	dec := 1
	if opts != nil {
		dec = opts.DecrementBy
	}

	incrementBy := dec
	if dec > 0 {
		incrementBy = -dec
	}

	err := k.config.Http.PostJson("/quotas/incrementUsage", map[string]any{
		"userId":      userId,
		"quotaName":   quotaName,
		"incrementBy": incrementBy,
	}, nil)
	if err != nil {
		return err
	}

	return nil
}

// SetQuotaUsage asynchronously set the quota usage for a given user to a given number.
//
//		Unlike incrementQuotaUsage and decrementQuotaUsage, this will set the usage to the specific number.
//
//	 - @param userId - The unique identifier for the user whose quota is being changed.
//	 - @param quotaName - The name of the quota to change.
//	 - @param usage - The new usage you want to set.
func (k KobbleUsers) SetQuotaUsage(userId string, quotaName string, usage int) error {
	err := k.config.Http.PostJson("/quotas/setUsage", map[string]any{
		"userId":    userId,
		"quotaName": quotaName,
		"usage":     usage,
	}, nil)
	if err != nil {
		return err
	}
	return nil
}

// GetQuotaUsage retrieves the quota usage for a given user based on the product assigned to them.
//
//   - @param userId - The unique identifier for the user whose quota usage is being retrieved.
//   - @param quotaName - The name of the quota to retrieve.
func (k KobbleUsers) GetQuotaUsage(userId string, quotaName string) (*QuotaUsage, error) {
	quotas, err := k.ListQuotas(userId, nil)
	if err != nil {
		return nil, err
	}

	for _, quota := range quotas {
		if quota.Name == quotaName {
			return &quota, nil
		}
	}

	return nil, nil
}

type ListPermissionsOptions struct {
	NoCache bool
}

// ListPermissions retrieves the list of permissions for a given user based on the product assigned to them.
//
//   - @param userId - The unique identifier for the user whose permissions are being retrieved.
//   - @param noCache - Set to true to bypass cache and fetch fresh data. Default is false.
func (k KobbleUsers) ListPermissions(userId string, opts *ListPermissionsOptions) ([]permissions.Permission, error) {
	perms := k.getCachedUserPerms(userId)
	if opts != nil && !opts.NoCache && perms != nil {
		return *perms, nil
	}

	var result []permissions.Permission
	err := k.config.Http.GetJson("/users/listPermissions", map[string]string{
		"userId": userId,
	}, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

type HasRemainingQuotaOptions struct {
	NoCache bool
}

// HasRemainingQuota checks if a user has remaining credit for all specified quota(s).
//
//   - @param userId - The unique identifier for the user whose quotas are being checked.
//   - @param quotaNames - The names of the quotas to check. Can be a single name or an array of names.
//   - @param noCache - Set to true to bypass cache and fetch fresh data. Default is false.
func (k KobbleUsers) HasRemainingQuota(userId string, quotaNames []string, opts *HasRemainingQuotaOptions) (bool, error) {
	var listQuotaOpts *ListQuotasOptions = nil
	if opts != nil {
		listQuotaOpts = &ListQuotasOptions{
			NoCache: opts.NoCache,
		}
	}
	quotas, err := k.ListQuotas(userId, listQuotaOpts)
	if err != nil {
		return false, err
	}

	for _, quotaName := range quotaNames {
		for _, quota := range quotas {
			if quota.Name == quotaName && quota.Remaining != nil && *quota.Remaining > 0 {
				return true, nil
			}
		}
	}

	return false, nil
}

type HasPermissionOptions struct {
	NoCache bool
}

// HasPermission checks if a user has all permissions specified as arguments.
//
//   - @param userId - The unique identifier for the user whose permissions are being checked.
//   - @param permissionNames - The names of the permission(s) to check. Can be a single permission name or an array of names.
//   - @param noCache - Set to true to bypass cache and fetch fresh data. Default is false.
func (k KobbleUsers) HasPermission(userId string, permissionNames []string, opts *HasPermissionOptions) (bool, error) {
	var listPermissionOpts *ListPermissionsOptions = nil
	if opts != nil {
		listPermissionOpts = &ListPermissionsOptions{
			NoCache: opts.NoCache,
		}
	}
	perms, err := k.ListPermissions(userId, listPermissionOpts)
	if err != nil {
		return false, err
	}

	for _, permName := range permissionNames {
		for _, perm := range perms {
			if perm.Name == permName {
				return true, nil
			}
		}
	}

	return false, nil
}

type IsAllowedOptions struct {
	NoCache bool
}

// IsAllowed this function is a helper to check if a user has all permissions and quotas specified in the payload.
//
//			If both permissionNames and quotaNames are provided, the user must have all permissions and quotas to be allowed.
//			If only permissionNames are provided, the user must have all permissions to be allowed.
//			If only quotaNames are provided, the user must have all quotas to be allowed.
//
//	 - @param userId - The unique identifier for the user whose quotas are being checked.
//	 - @param payload - The payload containing the permission and quota names to check.
//	 - @param payload.permissionNames - The names of the permissions to check.
//	 - @param payload.quotaNames - The names of the quotas to check.
//	 - @param noCache - Set to true to bypass cache and fetch fresh data. Default is false.
func (k KobbleUsers) IsAllowed(userId string, payload IsAllowedPayload, opts *IsAllowedOptions) (bool, error) {
	var hasPermissionOpts *HasPermissionOptions = nil
	var hasRemainingQuotaOpts *HasRemainingQuotaOptions = nil
	if opts != nil {
		hasPermissionOpts = &HasPermissionOptions{
			NoCache: opts.NoCache,
		}
		hasRemainingQuotaOpts = &HasRemainingQuotaOptions{
			NoCache: opts.NoCache,
		}
	}

	if len(payload.PermissionNames) > 0 && len(payload.QuotaNames) > 0 {
		hasPermission, err := k.HasPermission(userId, payload.PermissionNames, hasPermissionOpts)
		if err != nil {
			return false, err
		}

		hasQuota, err := k.HasRemainingQuota(userId, payload.QuotaNames, hasRemainingQuotaOpts)
		if err != nil {
			return false, err
		}

		return hasPermission && hasQuota, nil
	}

	if len(payload.PermissionNames) > 0 {
		return k.HasPermission(userId, payload.PermissionNames, hasPermissionOpts)
	}

	if len(payload.QuotaNames) > 0 {
		return k.HasRemainingQuota(userId, payload.QuotaNames, hasRemainingQuotaOpts)
	}

	return false, nil
}

type IsForbiddenOptions struct {
	NoCache bool
}

// IsForbidden this function is a helper to check if a user is forbidden from performing an action.
//
//		It is the opposite of isAllowed.
//
//	 - @param userId - The unique identifier for the user whose quotas are being checked.
//	 - @param payload - The payload containing the permission and quota names to check.
//	 - @param payload.permissionNames - The names of the permissions to check.
//	 - @param payload.quotaNames - The names of the quotas to check.
//	 - @param noCache - Set to true to bypass cache and fetch fresh data. Default is false.
func (k KobbleUsers) IsForbidden(userId string, payload IsAllowedPayload, opts *IsForbiddenOptions) (bool, error) {
	var isAllowedOpts *IsAllowedOptions = nil
	if opts != nil {
		isAllowedOpts = &IsAllowedOptions{
			NoCache: opts.NoCache,
		}
	}
	isAllowed, err := k.IsAllowed(userId, payload, isAllowedOpts)
	if err != nil {
		return false, err
	}

	return !isAllowed, nil
}

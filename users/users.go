package users

import (
	"github.com/valensto/kobble-go-sdk/common"
	"github.com/valensto/kobble-go-sdk/permissions"
	"github.com/valensto/kobble-go-sdk/utils"
	"strconv"
	"time"
)

type KobbleUsers struct {
	config           Config
	permissionsCache *utils.Cache[[]permissions.Permission]
	quotasCache      *utils.Cache[[]QuotaUsage]
}

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

func (k KobbleUsers) CreateLoginLink(userId string) (UrlLink, error) {
	var result UrlLink
	err := k.config.Http.PostJson("/users/mintLoginLink", map[string]string{
		"userId": userId,
	}, &result)
	return result, err
}

func (k KobbleUsers) Create(payload CreateUserPayload) (*User, error) {
	var result ApiUser
	err := k.config.Http.PostJson("/users", payload, &result)
	if err != nil {
		return nil, err
	}
	return k.transformApiUser(result), nil
}

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

func (k KobbleUsers) FindByMetadata(metadata map[string]any, options *ListUsersOptions) (common.Pagination[User], error) {
	page, limit := 1, 50
	if options != nil {
		if options.Page != nil {
			page = *options.Page
		}

		if options.Limit != nil {
			limit = *options.Limit
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

func (k KobbleUsers) ListAll(options *ListUsersOptions) (common.Pagination[User], error) {
	page, limit, includeMetadata := 1, 50, false
	if options != nil {
		if options.Page != nil {
			page = *options.Page
		}

		if options.Limit != nil {
			limit = *options.Limit
		}

		if options.IncludeMetadata != nil {
			includeMetadata = *options.IncludeMetadata
		}
	}

	var result common.Pagination[User]
	err := k.config.Http.GetJson("/users", map[string]string{
		"page":            strconv.Itoa(page),
		"limit":           strconv.Itoa(limit),
		"includeMetadata": strconv.FormatBool(includeMetadata),
	}, &result)
	if err != nil {
		return common.Pagination[User]{}, err
	}

	return result, nil
}

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

func (k KobbleUsers) GetCachedUserPerms(userId string) *[]permissions.Permission {
	key := k.userCacheKey(userId)
	perms := k.permissionsCache.Get(key)
	if perms != nil {
		return perms
	}

	return nil
}

func (k KobbleUsers) GetCachedUserQuotas(userId string) *[]QuotaUsage {
	key := k.userCacheKey(userId)
	quotas := k.quotasCache.Get(key)
	if quotas != nil {
		return quotas
	}

	return nil
}

func (k KobbleUsers) ListQuotas(userId string, noCache *bool) ([]QuotaUsage, error) {
	quotas := k.GetCachedUserQuotas(userId)
	if noCache != nil && !*noCache {
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

func (k KobbleUsers) IncrementQuotaUsage(userId string, quotaName string, incrementBy *int) error {
	inc := 1
	if incrementBy != nil {
		inc = *incrementBy
	}
	err := k.config.Http.PostJson("/users/incrementUsage", map[string]any{
		"userId":      userId,
		"quotaName":   quotaName,
		"incrementBy": inc,
	}, nil)
	if err != nil {
		return err
	}

	return nil
}

func (k KobbleUsers) DecrementQuotaUsage(userId string, quotaName string, decrementBy *int) error {
	dec := 1
	if decrementBy != nil {
		dec = *decrementBy
	}
	err := k.config.Http.PostJson("/users/decrementUsage", map[string]any{
		"userId":      userId,
		"quotaName":   quotaName,
		"decrementBy": dec,
	}, nil)
	if err != nil {
		return err
	}

	return nil
}

func (k KobbleUsers) SetQuotaUsage(userId string, quotaName string, usage int) error {
	err := k.config.Http.PostJson("/users/setUsage", map[string]any{
		"userId":    userId,
		"quotaName": quotaName,
		"usage":     usage,
	}, nil)
	if err != nil {
		return err
	}
	return nil
}

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

func (k KobbleUsers) ListPermissions(userId string, noCache *bool) ([]permissions.Permission, error) {
	perms := k.GetCachedUserPerms(userId)
	if noCache != nil && !*noCache {
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

func (k KobbleUsers) HasRemainingQuota(userId string, quotaNames []string, noCache *bool) (bool, error) {
	quotas, err := k.ListQuotas(userId, noCache)
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

func (k KobbleUsers) HasPermission(userId string, permissionNames []string, noCache *bool) (bool, error) {
	perms, err := k.ListPermissions(userId, noCache)
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

func (k KobbleUsers) IsAllowed(userId string, payload IsAllowedPayload, noCache *bool) (bool, error) {
	if len(payload.PermissionNames) > 0 && len(payload.QuotaNames) > 0 {
		hasPermission, err := k.HasPermission(userId, payload.PermissionNames, noCache)
		if err != nil {
			return false, err
		}

		hasQuota, err := k.HasRemainingQuota(userId, payload.QuotaNames, noCache)
		if err != nil {
			return false, err
		}

		return hasPermission && hasQuota, nil
	}

	if len(payload.PermissionNames) > 0 {
		return k.HasPermission(userId, payload.PermissionNames, noCache)
	}

	if len(payload.QuotaNames) > 0 {
		return k.HasRemainingQuota(userId, payload.QuotaNames, noCache)
	}

	return false, nil
}

func (k KobbleUsers) IsForbidden(userId string, payload IsAllowedPayload, noCache *bool) (bool, error) {
	isAllowed, err := k.IsAllowed(userId, payload, noCache)
	if err != nil {
		return false, err
	}

	return !isAllowed, nil
}

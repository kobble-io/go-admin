package quotas

// KobbleQuotas is the struct that holds the configuration for the KobbleQuotas service
type KobbleQuotas struct {
	config Config
}

// NewKobbleQuotas creates a new instance of KobbleQuotas
func NewKobbleQuotas(config Config) *KobbleQuotas {
	return &KobbleQuotas{
		config: config,
	}
}

package quotas

type KobbleQuotas struct {
	config Config
}

func NewKobbleQuotas(config Config) *KobbleQuotas {
	return &KobbleQuotas{
		config: config,
	}
}

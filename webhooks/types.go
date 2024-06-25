package webhooks

var WebhookSubscriptions = []string{
	"user.created",
	"quota.reached",
	"subscription.created",
	"subscription.updated",
	"subscription.deleted",
	"ping",
}

type WebhookSubscription string

type WebhookSubscriptionData struct {
	Provider struct {
		SubscriptionID string `json:"subscription_id"`
		ProductID      string `json:"product_id"`
		PriceID        string `json:"price_id"`
	} `json:"provider"`
	ProjectID         string  `json:"project_id"`
	ProductID         string  `json:"product_id"`
	PriceID           string  `json:"price_id"`
	UserID            string  `json:"user_id"`
	Email             string  `json:"email"`
	StartDate         *string `json:"start_date,omitempty"`
	EndedAt           *string `json:"ended_at,omitempty"`
	CancelAt          *string `json:"cancel_at,omitempty"`
	CanceledAt        *string `json:"canceled_at,omitempty"`
	CancelAtPeriodEnd bool    `json:"cancel_at_period_end"`
	Status            string  `json:"status"`
	TrialEnd          *string `json:"trial_end,omitempty"`
	TrialStart        *string `json:"trial_start,omitempty"`
}

type WebhookUserCreatedData struct {
	ID         string  `json:"id"`
	Email      string  `json:"email"`
	Name       *string `json:"name,omitempty"`
	CreatedAt  string  `json:"created_at"`
	IsVerified bool    `json:"is_verified"`
}

type WebhookQuotaReachedData struct {
	QuotaID   string `json:"quota_id"`
	QuotaName string `json:"quota_name"`
	Usage     int    `json:"usage"`
	Limit     int    `json:"limit"`
	UserID    string `json:"user_id"`
}

type WebhookUserCreatedEvent struct {
	Type string                 `json:"type"`
	Data WebhookUserCreatedData `json:"data"`
}

type WebhookQuotaReachedEvent struct {
	Type string                  `json:"type"`
	Data WebhookQuotaReachedData `json:"data"`
}

type WebhookSubscriptionCreatedEvent struct {
	Type string                  `json:"type"`
	Data WebhookSubscriptionData `json:"data"`
}

type WebhookSubscriptionUpdatedEvent struct {
	Type string                  `json:"type"`
	Data WebhookSubscriptionData `json:"data"`
}

type WebhookSubscriptionDeletedEvent struct {
	Type string                  `json:"type"`
	Data WebhookSubscriptionData `json:"data"`
}

type WebhookPingData struct {
	WebhookID string `json:"webhook_id"`
}

type WebhookPingEvent struct {
	Type string          `json:"type"`
	Data WebhookPingData `json:"data"`
}

type WebhookEvent struct {
	Type string `json:"type"`
	Data any    `json:"data"`
}

type webhookError struct {
	Message string
}

func (e *webhookError) Error() string {
	return e.Message
}

type webhookConstructEventError struct {
	webhookError
}

func newWebhookConstructEventError(message string) *webhookConstructEventError {
	return &webhookConstructEventError{
		webhookError: webhookError{
			Message: message,
		},
	}
}

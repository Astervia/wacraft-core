package webhook_model

type SearchableWebhookColumn string

var (
	URL         SearchableWebhookColumn = "url"
	Method      SearchableWebhookColumn = "http_method"
	EventColumn SearchableWebhookColumn = "event"
)

func (t SearchableWebhookColumn) IsValid() bool {
	switch t {
	case URL, Method, EventColumn:
		return true
	default:
		return false
	}
}

package campaign_model

type SendMessage string

var (
	Send   SendMessage = "send"
	Cancel SendMessage = "cancel"
	Status SendMessage = "status"
)

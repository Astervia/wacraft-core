package webhook_model

type Event string

var (
	SendWhatsAppMessage    Event = "send_whatsapp_message"
	ReceiveWhatsAppMessage Event = "receive_whatsapp_message"
)

func (t Event) IsValid() bool {
	switch t {
	case SendWhatsAppMessage, ReceiveWhatsAppMessage:
		return true
	default:
		return false
	}
}

package message_model

type JsonMessageKey string

var (
	SenderDataKey   JsonMessageKey = "sender_data"
	ReceiverDataKey JsonMessageKey = "receiver_data"
	ProductDataKey  JsonMessageKey = "product_data"
)

func (t JsonMessageKey) IsValid() bool {
	switch t {
	case SenderDataKey, ReceiverDataKey, ProductDataKey:
		return true
	default:
		return false
	}
}

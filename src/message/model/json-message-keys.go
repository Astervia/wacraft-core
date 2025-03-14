package message_model

type JsonMessageKey string

var (
	SenderDataKey   JsonMessageKey = "sender_data"
	ReceiverDataKey JsonMessageKey = "receiver_data"
	ProductDataKey  JsonMessageKey = "product_data"
)

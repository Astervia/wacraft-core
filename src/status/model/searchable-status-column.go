package status_model

type SearchableStatusColumn string

var (
	SenderDataColumn   SearchableStatusColumn = "sender_data"
	ReceiverDataColumn SearchableStatusColumn = "receiver_data"
	ProductDataColumn  SearchableStatusColumn = "product_data"
)

func (t SearchableStatusColumn) IsValid() bool {
	switch t {
	case SenderDataColumn, ReceiverDataColumn, ProductDataColumn:
		return true
	default:
		return false
	}
}

package billing_model

import "time"

// Payment represents a single payment read live from the payment provider
// (e.g. Stripe). It is never persisted — the billing API fetches it on demand.
type Payment struct {
	ID            string    `json:"id"`
	AmountCents   int64     `json:"amount_cents"`
	Currency      string    `json:"currency"`
	Status        string    `json:"status"`         // requires_action, processing, succeeded, canceled, ...
	PaymentMethod string    `json:"payment_method"` // e.g. "card", "boleto"
	Description   string    `json:"description,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	BoletoURL     string    `json:"boleto_url,omitempty"`     // Hosted boleto voucher page (when awaiting payment)
	BoletoPDFURL  string    `json:"boleto_pdf_url,omitempty"` // Downloadable boleto voucher PDF
	ReceiptURL    string    `json:"receipt_url,omitempty"`    // Receipt for a succeeded charge
}

// PaymentQuery lists payments from the payment provider. Pagination is
// cursor-based (the provider's opaque page token) because payments are read
// live from the provider rather than from the database.
type PaymentQuery struct {
	Limit  int    `json:"limit" query:"limit" validate:"omitempty,min=1,max=100"`
	Cursor string `json:"cursor,omitempty" query:"cursor"`
}

// PaymentListResponse is a single page of payments plus the cursor for the next page.
type PaymentListResponse struct {
	Data       []Payment `json:"data"`
	HasMore    bool      `json:"has_more"`
	NextCursor string    `json:"next_cursor,omitempty"`
}

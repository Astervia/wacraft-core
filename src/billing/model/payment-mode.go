package billing_model

type PaymentMode string

const (
	PaymentModePayment      PaymentMode = "payment"      // One-time payment (default, current behavior)
	PaymentModeSubscription PaymentMode = "subscription"  // Recurring via Stripe Subscriptions
)

func IsValidPaymentMode(m PaymentMode) bool {
	return m == PaymentModePayment || m == PaymentModeSubscription
}

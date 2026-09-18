package payment

import "time"

type CreatePaymentRequest struct {
	RegistrationID string `json:"registration_id" binding:"required,uuid"`
}

type CreatePaymentResponse struct {
	PaymentID      string `json:"payment_id"`
	GatewayOrderID string `json:"gateway_order_id"`
	AmountPaise    int64  `json:"amount_paise"`
	Currency       string `json:"currency"`
	Provider       string `json:"provider"`
	PublicKey      string `json:"public_key,omitempty"`
	IdempotencyKey string `json:"idempotency_key"`
}

type VerifyPaymentAPIRequest struct {
	PaymentID        string `json:"payment_id" binding:"required,uuid"`
	GatewayOrderID   string `json:"gateway_order_id" binding:"required"`
	GatewayPaymentID string `json:"gateway_payment_id" binding:"required"`
	Signature        string `json:"signature" binding:"required"`
}

type PaymentResponse struct {
	ID               string     `json:"id"`
	RegistrationID   string     `json:"registration_id"`
	UserID           string     `json:"user_id"`
	EventID          string     `json:"event_id"`
	Provider         string     `json:"provider"`
	GatewayPaymentID string     `json:"gateway_payment_id,omitempty"`
	GatewayOrderID   string     `json:"gateway_order_id,omitempty"`
	AmountPaise      int64      `json:"amount_paise"`
	Currency         string     `json:"currency"`
	Status           string     `json:"status"`
	VerifiedAt       *time.Time `json:"verified_at,omitempty"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}
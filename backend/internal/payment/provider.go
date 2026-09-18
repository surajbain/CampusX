package payment

import (
	"context"
	"errors"
)

type Provider interface {
	Name() string
	CreateOrder(ctx context.Context, req CreateOrderRequest) (*CreateOrderResponse, error)
	VerifyPayment(ctx context.Context, req VerifyPaymentRequest) error
	FetchPayment(ctx context.Context, gatewayPaymentID string) (*FetchedPayment, error)
}

type CreateOrderRequest struct {
	AmountPaise int64
	Currency    string
	Receipt     string
	Notes       map[string]string
}

type CreateOrderResponse struct {
	GatewayOrderID string
	AmountPaise    int64
	Currency       string
	PublicKey      string
}

type VerifyPaymentRequest struct {
	GatewayOrderID   string
	GatewayPaymentID string
	Signature        string
}

type FetchedPayment struct {
	GatewayPaymentID string
	GatewayOrderID   string
	AmountPaise      int64
	Currency         string
	Status           string
	Method           string
	RawData          map[string]any
}

var (
	ErrInvalidSignature = errors.New("payment: invalid signature")
	ErrPaymentNotFound  = errors.New("payment: not found on gateway")
	ErrProviderFailure  = errors.New("payment: provider request failed")
)
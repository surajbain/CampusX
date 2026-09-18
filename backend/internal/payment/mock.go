package payment

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"time"
)

type MockProvider struct {
	secret string
}

func NewMockProvider(secret string) *MockProvider {
	if secret == "" {
		secret = "mock_secret_dev_only"
	}
	return &MockProvider{secret: secret}
}

func (m *MockProvider) Name() string { return "MOCK" }

func (m *MockProvider) CreateOrder(ctx context.Context, req CreateOrderRequest) (*CreateOrderResponse, error) {
	orderID := "mock_order_" + randHex(8)
	return &CreateOrderResponse{
		GatewayOrderID: orderID,
		AmountPaise:    req.AmountPaise,
		Currency:       req.Currency,
		PublicKey:      "mock_public_key_dev",
	}, nil
}

func (m *MockProvider) VerifyPayment(ctx context.Context, req VerifyPaymentRequest) error {
	if req.GatewayPaymentID == "" || req.Signature == "" {
		return ErrInvalidSignature
	}
	expected := m.SignPayment(req.GatewayPaymentID)
	if !hmac.Equal([]byte(expected), []byte(req.Signature)) {
		return ErrInvalidSignature
	}
	return nil
}

func (m *MockProvider) FetchPayment(ctx context.Context, gatewayPaymentID string) (*FetchedPayment, error) {
	if gatewayPaymentID == "" {
		return nil, ErrPaymentNotFound
	}
	status := "captured"
	if strings.HasSuffix(gatewayPaymentID, "fail") {
		status = "failed"
	}
	return &FetchedPayment{
		GatewayPaymentID: gatewayPaymentID,
		Status:           status,
		Currency:         "INR",
		Method:           "mock",
		RawData: map[string]any{
			"mock":      true,
			"timestamp": time.Now().Unix(),
		},
	}, nil
}

func (m *MockProvider) SignPayment(gatewayPaymentID string) string {
	mac := hmac.New(sha256.New, []byte(m.secret))
	mac.Write([]byte(gatewayPaymentID))
	return hex.EncodeToString(mac.Sum(nil))
}

func randHex(n int) string {
	b := make([]byte, n)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
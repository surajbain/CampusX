package payment

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const razorpayAPIBase = "https://api.razorpay.com/v1"

type RazorpayProvider struct {
	keyID         string
	keySecret     string
	webhookSecret string
	httpClient    *http.Client
}

func NewRazorpayProvider(keyID, keySecret, webhookSecret string) *RazorpayProvider {
	return &RazorpayProvider{
		keyID:         keyID,
		keySecret:     keySecret,
		webhookSecret: webhookSecret,
		httpClient:    &http.Client{Timeout: 15 * time.Second},
	}
}

func (r *RazorpayProvider) Name() string { return "RAZORPAY" }

type razorpayOrderReq struct {
	Amount   int64             `json:"amount"`
	Currency string            `json:"currency"`
	Receipt  string            `json:"receipt"`
	Notes    map[string]string `json:"notes,omitempty"`
}

type razorpayOrderResp struct {
	ID       string `json:"id"`
	Amount   int64  `json:"amount"`
	Currency string `json:"currency"`
	Status   string `json:"status"`
}

func (r *RazorpayProvider) CreateOrder(ctx context.Context, req CreateOrderRequest) (*CreateOrderResponse, error) {
	body, _ := json.Marshal(razorpayOrderReq{
		Amount:   req.AmountPaise,
		Currency: req.Currency,
		Receipt:  req.Receipt,
		Notes:    req.Notes,
	})
	httpReq, err := http.NewRequestWithContext(ctx, "POST", razorpayAPIBase+"/orders", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrProviderFailure, err)
	}
	httpReq.SetBasicAuth(r.keyID, r.keySecret)
	httpReq.Header.Set("Content-Type", "application/json")

	res, err := r.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrProviderFailure, err)
	}
	defer res.Body.Close()
	respBody, _ := io.ReadAll(res.Body)
	if res.StatusCode >= 400 {
		return nil, fmt.Errorf("%w: status %d: %s", ErrProviderFailure, res.StatusCode, string(respBody))
	}
	var order razorpayOrderResp
	if err := json.Unmarshal(respBody, &order); err != nil {
		return nil, fmt.Errorf("%w: decode: %v", ErrProviderFailure, err)
	}
	return &CreateOrderResponse{
		GatewayOrderID: order.ID,
		AmountPaise:    order.Amount,
		Currency:       order.Currency,
		PublicKey:      r.keyID,
	}, nil
}

func (r *RazorpayProvider) VerifyPayment(ctx context.Context, req VerifyPaymentRequest) error {
	payload := req.GatewayOrderID + "|" + req.GatewayPaymentID
	mac := hmac.New(sha256.New, []byte(r.keySecret))
	mac.Write([]byte(payload))
	expected := hex.EncodeToString(mac.Sum(nil))
	if !hmac.Equal([]byte(expected), []byte(req.Signature)) {
		return ErrInvalidSignature
	}
	return nil
}

type razorpayPaymentResp struct {
	ID       string `json:"id"`
	OrderID  string `json:"order_id"`
	Amount   int64  `json:"amount"`
	Currency string `json:"currency"`
	Status   string `json:"status"`
	Method   string `json:"method"`
}

func (r *RazorpayProvider) FetchPayment(ctx context.Context, gatewayPaymentID string) (*FetchedPayment, error) {
	httpReq, err := http.NewRequestWithContext(ctx, "GET", razorpayAPIBase+"/payments/"+gatewayPaymentID, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrProviderFailure, err)
	}
	httpReq.SetBasicAuth(r.keyID, r.keySecret)

	res, err := r.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrProviderFailure, err)
	}
	defer res.Body.Close()
	respBody, _ := io.ReadAll(res.Body)
	if res.StatusCode == 404 {
		return nil, ErrPaymentNotFound
	}
	if res.StatusCode >= 400 {
		return nil, fmt.Errorf("%w: status %d: %s", ErrProviderFailure, res.StatusCode, string(respBody))
	}
	var p razorpayPaymentResp
	if err := json.Unmarshal(respBody, &p); err != nil {
		return nil, fmt.Errorf("%w: decode: %v", ErrProviderFailure, err)
	}
	return &FetchedPayment{
		GatewayPaymentID: p.ID,
		GatewayOrderID:   p.OrderID,
		AmountPaise:      p.Amount,
		Currency:         p.Currency,
		Status:           p.Status,
		Method:           p.Method,
		RawData:          map[string]any{"id": p.ID, "order_id": p.OrderID, "status": p.Status},
	}, nil
}

func (r *RazorpayProvider) VerifyWebhookSignature(body []byte, signature string) error {
	mac := hmac.New(sha256.New, []byte(r.webhookSecret))
	mac.Write(body)
	expected := hex.EncodeToString(mac.Sum(nil))
	if !hmac.Equal([]byte(expected), []byte(signature)) {
		return ErrInvalidSignature
	}
	return nil
}
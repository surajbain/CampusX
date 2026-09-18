package models

import (
	"time"

	"gorm.io/datatypes"
)

type PaymentProvider string

const (
	ProviderMock     PaymentProvider = "MOCK"
	ProviderRazorpay PaymentProvider = "RAZORPAY"
)

type PaymentStatus string

const (
	PayCreated    PaymentStatus = "CREATED"
	PayAuthorized PaymentStatus = "AUTHORIZED"
	PayCaptured   PaymentStatus = "CAPTURED"
	PayFailed     PaymentStatus = "FAILED"
	PayRefunded   PaymentStatus = "REFUNDED"
)

type Payment struct {
	Base
	RegistrationID string          `gorm:"type:uuid;not null;index" json:"registration_id"`
	UserID         string          `gorm:"type:uuid;not null;index" json:"user_id"`
	EventID        string          `gorm:"type:uuid;not null;index" json:"event_id"`
	Provider       PaymentProvider `gorm:"size:30;not null;default:'MOCK'" json:"provider"`

	GatewayPaymentID *string `gorm:"size:120;uniqueIndex" json:"gateway_payment_id,omitempty"`
	GatewayOrderID   *string `gorm:"size:120;index" json:"gateway_order_id,omitempty"`
	IdempotencyKey   string  `gorm:"size:120;not null;uniqueIndex" json:"idempotency_key"`

	AmountPaise int64         `gorm:"not null" json:"amount_paise"`
	Currency    string        `gorm:"size:8;not null;default:'INR'" json:"currency"`
	Status      PaymentStatus `gorm:"size:20;not null;default:'CREATED';index" json:"status"`

	Signature  string         `gorm:"type:text" json:"-"`
	RawPayload datatypes.JSON `gorm:"type:jsonb" json:"-"`

	VerifiedAt *time.Time `json:"verified_at,omitempty"`

	// ---- Relations ----
	Registration *Registration `gorm:"foreignKey:RegistrationID" json:"registration,omitempty"`
}

func (Payment) TableName() string { return "payments" }
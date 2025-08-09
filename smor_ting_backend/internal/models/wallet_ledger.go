package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LedgerType string
type LedgerDirection string
type LedgerStatus string

const (
	LedgerTopup         LedgerType = "topup"
	LedgerPayment       LedgerType = "payment"
	LedgerEscrowHold    LedgerType = "escrow_hold"
	LedgerEscrowRelease LedgerType = "escrow_release"
	LedgerWithdraw      LedgerType = "withdraw"
)

const (
	LedgerCredit LedgerDirection = "credit"
	LedgerDebit  LedgerDirection = "debit"
)

const (
	LedgerPending   LedgerStatus = "pending"
	LedgerCompleted LedgerStatus = "completed"
	LedgerFailed    LedgerStatus = "failed"
)

type WalletLedgerEntry struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID      primitive.ObjectID `bson:"user_id" json:"user_id"`
	Type        LedgerType         `bson:"type" json:"type"`
	Direction   LedgerDirection    `bson:"direction" json:"direction"`
	Amount      float64            `bson:"amount" json:"amount"`
	Currency    string             `bson:"currency" json:"currency"`
	Status      LedgerStatus       `bson:"status" json:"status"`
	IsEscrow    bool               `bson:"is_escrow" json:"is_escrow"`
	Reference   string             `bson:"reference" json:"reference"`
	ProviderRef string             `bson:"provider_ref" json:"provider_ref"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}

type WalletBalances struct {
	Available   float64 `json:"available"`
	PendingHeld float64 `json:"pending"`
	Total       float64 `json:"total"`
	Currency    string  `json:"currency"`
}

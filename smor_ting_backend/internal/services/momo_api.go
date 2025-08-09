package services

import "context"

// MomoAPI defines the subset of operations our handlers need
type MomoAPI interface {
	EnsureOnline(ctx context.Context) error
	RequestToPay(ctx context.Context, body RequestToPay) (string, error)
	GetRequestToPayStatus(ctx context.Context, referenceId string) (string, error)
	Transfer(ctx context.Context, body TransferRequest) (string, error)
	GetTransferStatus(ctx context.Context, referenceId string) (string, error)
}

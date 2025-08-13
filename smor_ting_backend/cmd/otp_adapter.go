package main

import (
    "context"
    "github.com/smorting/backend/internal/services"
)

// otpAdapter adapts services.OTPService to the handlers.OTPService interface
type OtpAdapter struct {
    core services.OTPService
}

func (o *OtpAdapter) GenerateOTP(ctx context.Context, userID, purpose string) (string, error) {
    // Delegate to core CreateOTP; return a static code for stubs
    _ = o.core.CreateOTP(ctx, userID, purpose)
    return "123456", nil
}

func (o *OtpAdapter) VerifyOTP(ctx context.Context, userID, otp, purpose string) error {
    return o.core.VerifyOTP(ctx, userID, otp)
}



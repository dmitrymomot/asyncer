package main

// Predefined task names.
const (
	SendConfirmationEmailTask = "send_confirmation_email"
	CleanUpExpiredOTP         = "clean_up_expired_otp"
)

// Predefined task payloads.
type (
	SendConfirmationEmailPayload struct {
		Email string `json:"email"`
		OTP   string `json:"otp"`
	}
)

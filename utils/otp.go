package utils

import (
	"time"
)
type otpDetails struct{
	code string
	codeexpiry time.Time
}
const OTPValidityDuration = 5 * time.Minute
func VerifyOTP(otpDetails *otpDetails, providedOTP string) bool {
	// Compare the saved OTP with the provided OTP
	if otpDetails.code != providedOTP {
		return false
	}
	// Check if the OTP has expired based on the timestamp
	currentTime := time.Now()
	return currentTime.Before(otpDetails.codeexpiry)
}	
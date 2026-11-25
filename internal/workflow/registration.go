package workflow

import (
	"fmt"
	"simple-webser/internal/prospect"
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

func generateOTPCode(ctx workflow.Context) (string, error) {
	var otpCode string
	err := workflow.SideEffect(ctx, func(ctx workflow.Context) any {
		otp, err := prospect.NewOTP()
		if err != nil {
			return ""
		}
		return otp.Code
	}).Get(&otpCode)
	if err != nil || otpCode == "" {
		return "", err
	}
	return otpCode, nil
}

func RegistrationWorkflow(ctx workflow.Context, input RegistrationInput) error {
	logger := workflow.GetLogger(ctx)

	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
	}
	actCtx := workflow.WithActivityOptions(ctx, ao)

	otpCode, err := generateOTPCode(ctx)
	if err != nil || otpCode == "" {
		return fmt.Errorf("failed to generate OTP: %w", err)
	}

	var a *Activities

	err = workflow.ExecuteActivity(ctx, a.SendOTPEmail, input.Email, otpCode).Get(ctx, nil)

	otpExpiresAt := workflow.Now(ctx).Add(15 * time.Minute)
	attempts := 0
	verified := false
	done := false
	var workflowErr error

	verifyChan := workflow.GetSignalChannel(ctx, SignalVerifyOTP)
	resendChan := workflow.GetSignalChannel(ctx, SignalResendOTP)

	expiryTimer := workflow.NewTimer(ctx, 24*time.Hour)

	for !verified && !done {
		selector := workflow.NewSelector(ctx)
		selector.AddReceive(verifyChan, func(c workflow.ReceiveChannel, more bool) {
			var sig VerifyOTPSignal
			c.Receive(ctx, &sig)

			//If OTP expires, reject
			if workflow.Now(ctx).After(otpExpiresAt) {
				logger.Warn("OTP expired, ignoring verification attempts", "email", input.Email)
				return
			}

			//wrong otp: increment attempts and lock at 5
			if sig.Code != otpCode {
				attempts++
				logger.Warn("Invalid OPT", "email", input.Email, "attempts", attempts)

				if attempts >= 5 {
					workflowErr = temporal.NewNonRetryableApplicationError(
						"account locked due to too many failed attempts",
						"account locked",
						nil,
					)
					done = true
				}
				return
			}

			verified = true

		})
	}
}

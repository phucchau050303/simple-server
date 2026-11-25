package workflow

import (
	"context"
	credentialv1 "simple-webser/proto/credential/v1"

	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/temporal"
)

type Activities struct {
	CredentialClient credentialv1.CredentialServiceClient
}

func NewActivities(credentialClient credentialv1.CredentialServiceClient) *Activities {
	return &Activities{
		CredentialClient: credentialClient,
	}
}

func (a *Activities) SendOTPEmail(ctx context.Context, email string, otp string) error {
	logger := activity.GetLogger(ctx)
	logger.Info("Sending OTP email", "email", email, "otp", otp)
	return nil
}

func (a *Activities) CreateAccount(ctx context.Context, email string, password string) error {
	response, err := a.CredentialClient.SignUp(ctx, &credentialv1.SignUpRequest{Email: email, Password: password})
	if err != nil {
		return err
	}
	if !response.GetIsValid() {
		return temporal.NewNonRetryableApplicationError(response.GetMessage(), "ValidationError", nil)
	}
	return nil
}

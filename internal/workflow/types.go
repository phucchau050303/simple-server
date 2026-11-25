package workflow

const RegistrationTaskQueue = "REGISTRATION_TASK_QUEUE"

const SignalVerifyOtp = "verify-otp"
const SignalResendOtp = "resend-otp"

type RegistrationInput struct {
	Email    string
	Password string
}

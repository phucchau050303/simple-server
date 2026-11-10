package prospect

import "time"

type Status string

const (
	StausPending   Status = "pending"
	StatusVerified Status = "verify"
	StatusActive   Status = "active"
	StatusLocked   Status = "locked"
	StatusExpired  Status = "expired"
)

type Prospect struct {
	ID           string
	Email        string
	PasswordHash string
	Status       Status
	CurrentOTP   *OTP
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

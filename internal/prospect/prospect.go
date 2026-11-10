package prospect

import (
	"fmt"
	"time"
)

type Status string

const (
	StatusPending  Status = "pending"
	StatusVerified Status = "verified"
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

func (p *Prospect) VerifyOTP(inputCode string) (bool, error) {

	if p.Status == StatusLocked {
		return false, fmt.Errorf("account is locked due to too many attempts")
	}

	if p.Status == StatusExpired {
		return false, fmt.Errorf("Registration has expired")
	}

	if p.Status != StatusPending {
		return false, fmt.Errorf("Account is not pending verification")
	}

	if p.CurrentOTP == nil {
		return false, fmt.Errorf("No active OTP found")
	}

	if time.Now().After(p.CurrentOTP.ExpiresAt) {
		return false, fmt.Errorf("OTP has expired, please request a new one")
	}

	if inputCode != p.CurrentOTP.Code {
		p.CurrentOTP.Attempts++
		p.UpdatedAt = time.Now()
		if p.CurrentOTP.Attempts >= 5 {
			p.Status = StatusLocked
			return false, fmt.Errorf("Too many failed attempts, account has been locked")
		}

		return false, fmt.Errorf("please try again, you have %d attempts left", 5-p.CurrentOTP.Attempts)
	}

	p.Status = StatusVerified
	p.UpdatedAt = time.Now()

	return true, nil

}

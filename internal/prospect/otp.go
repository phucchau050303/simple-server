package prospect

import (
	"fmt"
	"time"

	"crypto/rand"

	"math/big"
)

type OTP struct {
	Code      string
	ExpiresAt time.Time
	Attempts  int
	CreatedAt time.Time
}

func NewOTP() (*OTP, error) {
	max := big.NewInt(1000000)

	num, err := rand.Int(rand.Reader, max)

	if err != nil {
		return nil, err
	}

	code := fmt.Sprintf("%06d", num)

	now := time.Now()

	return &OTP{
		Code:      code,
		ExpiresAt: now.Add(15 * time.Minute),
		Attempts:  0,
		CreatedAt: now,
	}, nil
}

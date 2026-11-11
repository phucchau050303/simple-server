package prospect

import (
	"fmt"
	"testing"
	"time"
)

func TestNewOTP(t *testing.T) {
	otp, err := NewOTP()
	if err != nil || otp == nil {
		t.Fatalf("New OTP can't be created: %v \n", err)
	}
	t.Logf("OTP is %+v \n", otp)
}

func TestVerifyOTP_TooManyAttempts(t *testing.T) {
	//Arrange
	now := time.Now()
	otp := &OTP{
		Code:      "748920",
		ExpiresAt: now.Add(15 * time.Minute),
		Attempts:  0,
		CreatedAt: now,
	}
	prospect := &Prospect{
		ID:           "2od",
		Email:        "abc@abc.com",
		PasswordHash: "3456789i",
		Status:       StatusLocked,
		CurrentOTP:   otp,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	//Act
	valid, err := prospect.VerifyOTP("748920")

	//Assert
	if err == nil && valid != true {
		t.Fatalf(`expected "account is locked due to too many attempts"`)
	}
	fmt.Print("return error: ", err.Error(), " \n")
}

func TestVerifyOTP_StatusExpired(t *testing.T) {
	//Arrange
	now := time.Now()
	otp := &OTP{
		Code:      "748920",
		ExpiresAt: now.Add(15 * time.Minute),
		Attempts:  0,
		CreatedAt: now,
	}
	prospect := &Prospect{
		ID:           "2od",
		Email:        "abc@abc.com",
		PasswordHash: "3456789i",
		Status:       StatusExpired,
		CurrentOTP:   otp,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	//Act
	valid, err := prospect.VerifyOTP("748920")

	//Assert
	if err == nil && valid != true {
		t.Fatalf(`expected "Registration has expired"`)
	}
	fmt.Print("return error: ", err.Error(), " \n")
}

func TestVerifyOTP_StatusPending(t *testing.T) {
	//Arrange
	now := time.Now()
	otp := &OTP{
		Code:      "748920",
		ExpiresAt: now.Add(15 * time.Minute),
		Attempts:  0,
		CreatedAt: now,
	}
	prospect := &Prospect{
		ID:           "2od",
		Email:        "abc@abc.com",
		PasswordHash: "3456789i",
		Status:       StatusActive,
		CurrentOTP:   otp,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	//Act
	valid, err := prospect.VerifyOTP("748920")

	//Assert
	if err == nil && valid != true {
		t.Fatalf(`expected "Account is not pending verification"`)
	}
	fmt.Print("return error: ", err.Error(), " \n")
}

func TestVerifyOTP_NoOPT(t *testing.T) {
	//Arrange
	now := time.Now()

	prospect := &Prospect{
		ID:           "2od",
		Email:        "abc@abc.com",
		PasswordHash: "3456789i",
		Status:       StatusPending,
		CurrentOTP:   nil,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	//Act
	valid, err := prospect.VerifyOTP("748920")

	//Assert
	if err == nil && valid != true {
		t.Fatalf(`expected "No active OTP found"`)
	}
	fmt.Print("return error: ", err.Error(), " \n")
}

func TestVerifyOTP_ExpiredOPT(t *testing.T) {
	//Arrange
	now := time.Now()
	otp := &OTP{
		Code:      "748920",
		ExpiresAt: now.Add(-15 * time.Minute),
		Attempts:  0,
		CreatedAt: now,
	}
	prospect := &Prospect{
		ID:           "2od",
		Email:        "abc@abc.com",
		PasswordHash: "3456789i",
		Status:       StatusPending,
		CurrentOTP:   otp,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	//Act
	valid, err := prospect.VerifyOTP("748920")

	//Assert
	if err == nil && valid != true {
		t.Fatalf(`expected "OTP has expired, please request a new one"`)
	}
	fmt.Print("return error: ", err.Error(), " \n")
}

func TestVerifyOTP_IncorrectOPT_Below5Times(t *testing.T) {
	//Arrange
	now := time.Now()
	otp := &OTP{
		Code:      "748920",
		ExpiresAt: now.Add(15 * time.Minute),
		Attempts:  3,
		CreatedAt: now,
	}
	prospect := &Prospect{
		ID:           "2od",
		Email:        "abc@abc.com",
		PasswordHash: "3456789i",
		Status:       StatusPending,
		CurrentOTP:   otp,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	//Act
	valid, err := prospect.VerifyOTP("881245")

	//Assert
	if err == nil && valid != true {
		t.Fatalf(`expected "please try again, you have %d attempts left"`, 5-prospect.CurrentOTP.Attempts)
	}
	fmt.Print("return error: ", err.Error(), " \n")
}

func TestVerifyOTP_IncorrectOPT_MoreThan5Times(t *testing.T) {
	//Arrange
	now := time.Now()
	otp := &OTP{
		Code:      "748920",
		ExpiresAt: now.Add(15 * time.Minute),
		Attempts:  4,
		CreatedAt: now,
	}
	prospect := &Prospect{
		ID:           "2od",
		Email:        "abc@abc.com",
		PasswordHash: "3456789i",
		Status:       StatusPending,
		CurrentOTP:   otp,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	//Act
	valid, err := prospect.VerifyOTP("748921")

	//Assert
	if err == nil && valid != true {
		t.Fatalf(`expected "Too many failed attempts, account has been locked"`)
	}
	fmt.Print("return error: ", err.Error(), " \n")
}

func TestVerifyOTP_CorrectOTP_NoError(t *testing.T) {
	//Arrange
	now := time.Now()
	otp := &OTP{
		Code:      "748920",
		ExpiresAt: now.Add(15 * time.Minute),
		Attempts:  4,
		CreatedAt: now,
	}
	prospect := &Prospect{
		ID:           "2od",
		Email:        "abc@abc.com",
		PasswordHash: "3456789i",
		Status:       StatusPending,
		CurrentOTP:   otp,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	//Act
	valid, err := prospect.VerifyOTP("748920")

	//Assert
	if err != nil || valid == false {
		t.Fatalf(`OTP can't be verified, error created: %v`, err)
	}
}

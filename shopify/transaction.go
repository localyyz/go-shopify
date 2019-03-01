// Copyright 2019 The go-shopify AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package shopify

import (
	"errors"
	"fmt"
	"time"
)

type Transaction struct {
	ID      int64  `json:"id"`
	Amount  string `json:"amount"`
	OrderID int64  `json:"order_id"`

	// Status and error codes
	ErrorCode string            `json:"error_code"`
	Status    TransactionStatus `json:"status"`
	Message   string            `json:"message"`

	Test     bool   `json:"test"`
	Currency string `json:"currency"`

	CreatedAt *time.Time `json:"created_at"`
}

type TransactionStatus uint32
type TransactionError error

const (
	_ TransactionStatus = iota
	TransactionStatusPending
	TransactionStatusSuccess
	TransactionStatusFailure
	TransactionStatusError
)

var (
	transactionStatuses = []string{
		"-",
		"pending",
		"success",
		"failure",
		"error",
	}

	ErrIncorrectNumber   = errors.New("incorrect_number")
	ErrInvalidNumber     = errors.New("invalid_number")
	ErrInvalidExpiryDate = errors.New("invalid_expiry_date")
	ErrInvalidCvc        = errors.New("invalid_expiry_date")
	ErrExpiredCard       = errors.New("expired_card")
	ErrIncorrectCvc      = errors.New("incorrect_cvc")
	ErrIncorrectZip      = errors.New("incorrect_zip")
	ErrIncorrectAddress  = errors.New("incorrect_address")
	ErrCardDeclined      = errors.New("card_declined")
	ErrProcessingError   = errors.New("processing_error")
	ErrCallIssuer        = errors.New("call_issuer")
	ErrPickUpCard        = errors.New("pick_up_card")
)

// String returns the string value of the status.
func (s TransactionStatus) String() string {
	return transactionStatuses[s]
}

// MarshalText satisfies TextMarshaler
func (s TransactionStatus) MarshalText() ([]byte, error) {
	return []byte(s.String()), nil
}

// UnmarshalText satisfies TextUnmarshaler
func (s *TransactionStatus) UnmarshalText(text []byte) error {
	enum := string(text)
	for i := 0; i < len(transactionStatuses); i++ {
		if enum == transactionStatuses[i] {
			*s = TransactionStatus(i)
			return nil
		}
	}
	return fmt.Errorf("unknown transaction status %s", enum)
}

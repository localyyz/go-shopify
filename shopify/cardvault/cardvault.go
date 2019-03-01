// Copyright 2019 The go-shopify AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cardvault

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
)

type CreditCard struct {
	Number            string `json:"number"`
	FirstName         string `json:"first_name"`
	LastName          string `json:"last_name"`
	Month             string `json:"month"`
	Year              string `json:"year"`
	VerificationValue string `json:"verification_value"`
}

type Payment struct {
	Amount      string      `json:"amount"`
	UniqueToken string      `json:"unique_token"`
	CreditCard  *CreditCard `json:"credit_card"`
}

type PaymentRequest struct {
	Payment *Payment `json:"payment"`
}

const cardVaultURL = "https://elb.deposit.shopifycs.com/sessions"

func AddCard(ctx context.Context, vaultRequest *PaymentRequest) (string, *http.Response, error) {
	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(vaultRequest)
	if err != nil {
		return "", nil, err
	}

	req, err := http.NewRequest("POST", cardVaultURL, buf)
	if err != nil {
		return "", nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", nil, err
	}

	var cardVaultResponse struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(res.Body).Decode(&cardVaultResponse); err != nil {
		return "", res, err
	}

	return cardVaultResponse.ID, res, nil
}

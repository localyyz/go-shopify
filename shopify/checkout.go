// Copyright 2019 The go-shopify AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package shopify

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type CheckoutService service

type Checkout struct {
	LineItems  []*LineItem `json:"line_items,omitempty"`
	Email      string      `json:"email,omitempty"`
	Token      string      `json:"token,omitempty"`
	Name       string      `json:"name,omitempty"`
	CustomerID int64       `json:"customer_id,omitempty"`

	// Totals
	SubtotalPrice string `json:"subtotal_price,omitempty"`
	TotalTax      string `json:"total_tax,omitempty"`
	TotalPrice    string `json:"total_price,omitempty"`
	PaymentDue    string `json:"payment_due,omitempty"`
	Currency      string `json:"currency,omitempty"`

	// Order
	OrderID        int64  `json:"order_id,omitempty"`
	OrderStatusURL string `json:"order_status_url,omitempty"`

	// ShopifyPaymentAccountID is used to use stripe as a payment token provider
	ShopifyPaymentAccountID string `json:"shopify_payments_account_id,omitempty"`
	// Use payment url with the other direct payment providers to generate a token
	PaymentURL       string `json:"payment_url,omitempty"`
	WebURL           string `json:"web_url,omitempty"`
	WebProcessingURL string `json:"web_processing_url,omitempty"`
	// Don't omit empty, need empty to remove
	DiscountCode string `json:"discount_code,omitempty"`

	AppliedDiscount *AppliedDiscount `json:"applied_discount,omitempty"`
	ShippingAddress *CustomerAddress `json:"shipping_address,omitempty"`
	BillingAddress  *CustomerAddress `json:"billing_address,omitempty"`

	TaxesIncluded bool          `json:"taxes_included"`
	ShippingLine  *ShippingLine `json:"shipping_line,omitempty"`
	TaxLines      []*TaxLine    `json:"tax_lines,omitempty"`
	Payments      []*Payment    `json:"payments,omitempty"`

	CompletedAt *time.Time `json:"completed_at,omitempty"`
}

type TaxLine struct {
	Title string `json:"title"`
	// Dollar tax amount
	Price string `json:"price,omitempty"`
	// Percentage tax amount
	Rate float64 `json:"rate,omitempty"`
}

type ShippingLine struct {
	Handle        string      `json:"handle,omitempty"`
	Price         string      `json:"price,omitempty"`
	Title         string      `json:"title,omitempty"`
	DeliveryRange []time.Time `json:"delivery_range,omitempty"`
}

type CheckoutShipping struct {
	ID            string        `json:"id"`
	Price         string        `json:"price"`
	Title         string        `json:"title"`
	Checkout      *ShippingRate `json:"checkout"`
	PhoneRequired bool          `json:"phone_required"`
	DeliveryRange []time.Time   `json:"delivery_range"`
	Handle        string        `json:"handle"`
}

type ShippingRate struct {
	SubtotalPrice string `json:"subtotal_price"`
	TotalTax      string `json:"total_tax"`
	TotalPrice    string `json:"total_price"`
	PaymentDue    string `json:"payment_due"`
}

type BillingAddress struct{}

type LineItem struct {
	VariantID int64 `json:"variant_id"`
	Quantity  int64 `json:"quantity"`
}

type CheckoutRequest struct {
	Checkout *Checkout `json:"checkout"`
}

// custom marshal json
func (c *CheckoutRequest) MarshalJSON() ([]byte, error) {
	if c != nil && c.Checkout != nil {
		cc := *c
		cc.Checkout.AppliedDiscount = nil
		return json.Marshal(cc)
	}
	return nil, nil
}

type ShippingRateRequest struct {
	CheckoutShipping []*CheckoutShipping `json:"shipping_rates"`
}

type Payment struct {
	ID        int64  `json:"id,omitempty"`
	Amount    string `json:"amount"`
	SessionID string `json:"session_id,omitempty"`
	// clientside idempotency token
	UniqueToken string `json:"unique_token"`

	PaymentProcessingErrorMessage string         `json:"payment_processing_error_message,omitempty"`
	PaymentToken                  *PaymentToken  `json:"payment_token,omitempty"`
	RequestDetails                *RequestDetail `json:"request_details,omitempty"`

	// transaction
	Transaction *Transaction `json:"transaction,omitempty"`
}

type RequestDetail struct {
	IPAddress      string `json:"ip_address,omitempty"`
	AcceptLanguage string `json:"accept_language,omitempty"`
	UserAgent      string `json:"user_agent,omitempty"`
}

type PaymentToken struct {
	// Stripe token
	PaymentData string `json:"payment_data"`
	// stripe_vault_token
	Type string `json:"type"`
}

type PaymentRequest struct {
	Payment *Payment `json:"payment"`
}

const StripeVaultToken = `stripe_vault_token`

func (c *CheckoutService) Get(ctx context.Context, token string) (*Checkout, *http.Response, error) {
	req, err := c.client.NewRequest("GET", fmt.Sprintf("/admin/checkouts/%s.json", token), nil)
	if err != nil {
		return nil, nil, err
	}

	checkoutWrapper := new(CheckoutRequest)
	resp, err := c.client.Do(ctx, req, checkoutWrapper)
	if err != nil {
		return nil, resp, err
	}

	return checkoutWrapper.Checkout, resp, nil
}

// helper function get or create based on token
func (c *CheckoutService) CreateOrUpdate(ctx context.Context, checkout *Checkout) (*Checkout, *http.Response, error) {
	if len(checkout.Token) == 0 {
		return c.Create(ctx, checkout)
	}
	return c.Update(ctx, checkout)
}

func (c *CheckoutService) Create(ctx context.Context, checkout *Checkout) (*Checkout, *http.Response, error) {
	req, err := c.client.NewRequest(
		"POST",
		"/admin/checkouts.json",
		&CheckoutRequest{checkout},
	)
	if err != nil {
		return nil, nil, err
	}

	checkoutWrapper := &CheckoutRequest{Checkout: checkout}
	resp, err := c.client.Do(ctx, req, checkoutWrapper)
	if err != nil {
		return nil, resp, err
	}

	return checkoutWrapper.Checkout, resp, nil
}

func (c *CheckoutService) Update(ctx context.Context, checkout *Checkout) (*Checkout, *http.Response, error) {
	// need to re-wrap the incoming checkout with a request
	// and pull out data that we may want to update
	req, err := c.client.NewRequest(
		"PUT",
		fmt.Sprintf("/admin/checkouts/%s.json", checkout.Token),
		&CheckoutRequest{checkout},
	)
	if err != nil {
		return nil, nil, err
	}

	checkoutWrapper := &CheckoutRequest{Checkout: checkout}
	resp, err := c.client.Do(ctx, req, checkoutWrapper)
	if err != nil {
		return nil, resp, err
	}

	return checkoutWrapper.Checkout, resp, nil
}

func (c *CheckoutService) ListShippingRates(ctx context.Context, token string) ([]*CheckoutShipping, *http.Response, error) {
	var (
		resp                *http.Response
		pollWait            = "0"
		pollURL             = fmt.Sprintf("/admin/checkouts/%s/shipping_rates.json", token)
		pollStatus          = http.StatusAccepted
		shippingRateWrapper = new(ShippingRateRequest)
	)

	for {
		if pollStatus != http.StatusAccepted {
			break
		}

		req, err := c.client.NewRequest("GET", pollURL, nil)
		if err != nil {
			return nil, nil, err
		}

		wait, _ := strconv.Atoi(pollWait)
		// TODO: make a proper poller
		time.Sleep(time.Duration(wait) * time.Second)

		resp, err = c.client.Do(ctx, req, shippingRateWrapper)
		if err != nil {
			return nil, resp, err
		}

		// check Location and Retry-After for url and delay
		pollURL = resp.Header.Get("Location")
		pollWait = resp.Header.Get("Retry-After")
		pollStatus = resp.StatusCode
	}

	return shippingRateWrapper.CheckoutShipping, resp, nil
}

func (c *CheckoutService) Payment(ctx context.Context, token string, payment *Payment) (*Payment, *http.Response, error) {
	var (
		resp           *http.Response
		paymentWrapper = &PaymentRequest{payment}
	)

	req, err := c.client.NewRequest("POST", fmt.Sprintf("/admin/checkouts/%s/payments.json", token), paymentWrapper)
	if err != nil {
		return nil, nil, err
	}

	resp, err = c.client.Do(ctx, req, paymentWrapper)
	if err != nil {
		return nil, resp, err
	}

	var (
		pollURL    = resp.Header.Get("Location")
		pollWait   = resp.Header.Get("Retry-After")
		pollStatus = resp.StatusCode
	)

	// poll and wait
	for {
		if pollStatus != http.StatusAccepted {
			break
		}

		req, err = c.client.NewRequest("GET", pollURL, nil)
		if err != nil {
			return nil, nil, err
		}

		wait, _ := strconv.Atoi(pollWait)
		// TODO: make a proper poller
		time.Sleep(time.Duration(wait) * time.Second)

		resp, err = c.client.Do(ctx, req, paymentWrapper)
		if err != nil {
			return nil, resp, err
		}

		// check Location and Retry-After for url and delay
		pollURL = resp.Header.Get("Location")
		pollWait = resp.Header.Get("Retry-After")
		pollStatus = resp.StatusCode
	}

	return paymentWrapper.Payment, resp, nil
}

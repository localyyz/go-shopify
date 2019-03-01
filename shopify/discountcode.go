// Copyright 2019 The go-shopify AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package shopify

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type DiscountCodeService service

type DiscountCode struct {
	ID          int64      `json:"id"`
	PriceRuleID int64      `json:"price_rule_id"`
	Code        string     `json:"code"`
	UsageCount  int32      `json:"usage_count"`
	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
}

type AppliedDiscount struct {
	Amount              string `json:"amount"`
	Title               string `json:"title"`
	Description         string `json:"description"`
	Value               string `json:"value"`
	ValueType           string `json:"value_type"`
	Applicable          bool   `json:"applicable"`
	NonApplicableReason string `json:"non_applicable_reason,omitempty"`
}

type DiscountCodeRequest struct {
	*DiscountCode `json:"discount_code"`
}

func (s *DiscountCodeService) Create(ctx context.Context, discountCode *DiscountCode) (*DiscountCode, *http.Response, error) {
	req, err := s.client.NewRequest(
		"POST",
		fmt.Sprintf("/admin/price_rules/%d/discount_codes.json", discountCode.PriceRuleID),
		&DiscountCodeRequest{discountCode},
	)
	if err != nil {
		return nil, nil, err
	}
	wrapper := &DiscountCodeRequest{discountCode}
	resp, err := s.client.Do(ctx, req, &wrapper)
	if err != nil {
		return nil, resp, err
	}

	return wrapper.DiscountCode, resp, nil
}

func (s *DiscountCodeService) Delete(ctx context.Context, discountCode *DiscountCode) (*DiscountCode, *http.Response, error) {
	req, err := s.client.NewRequest(
		"DELETE",
		fmt.Sprintf("/admin/price_rules/%d/discount_codes/%d.json", discountCode.PriceRuleID, discountCode.ID),
		nil,
	)
	if err != nil {
		return nil, nil, err
	}
	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return nil, resp, err
	}

	return nil, resp, nil
}

func (s *DiscountCodeService) Get(ctx context.Context, discountCode *DiscountCode) (*DiscountCode, *http.Response, error) {
	req, err := s.client.NewRequest(
		"GET",
		fmt.Sprintf("/admin/price_rules/%d/discount_codes/%d.json", discountCode.PriceRuleID, discountCode.ID),
		nil,
	)
	if err != nil {
		return nil, nil, err
	}

	wrapper := &DiscountCodeRequest{discountCode}
	resp, err := s.client.Do(ctx, req, &wrapper)
	if err != nil {
		return nil, resp, err
	}

	return wrapper.DiscountCode, resp, nil
}

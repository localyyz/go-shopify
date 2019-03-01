// Copyright 2019 The go-shopify AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package shopify

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type PriceRuleService service

type PriceRule struct {
	ID                int64                    `json:"id,omitempty"`
	Title             string                   `json:"title,omitempty"`
	ValueType         PriceRuleValueType       `json:"value_type,omitempty"`
	Value             string                   `json:"value,omitempty"`
	CustomerSelection string                   `json:"customer_selection,omitempty"`
	TargetType        PriceRuleTargetType      `json:"target_type,omitempty"`
	TargetSelection   PriceRuleTargetSelection `json:"target_selection,omitempty"`
	AllocationMethod  string                   `json:"allocation_method,omitempty"`
	OncePerCustomer   bool                     `json:"once_per_customer,omitempty"`
	UsageLimit        int                      `json:"usage_limit,omitempty"`

	EntitledProductIds    []int64 `json:"entitled_product_ids,omitempty"`
	EntitledVariantIds    []int64 `json:"entitled_variant_ids,omitempty"`
	EntitledCollectionIds []int64 `json:"entitled_collection_ids,omitempty"`
	EntitledCountryIds    []int64 `json:"entitled_country_ids,omitempty"`

	// Prefreq for BUY X GET Y type deals
	PrerequisiteSavedSearchIds []int64 `json:"prerequisite_saved_search_ids,omitempty"`
	PrerequisiteCustomerIds    []int64 `json:"prerequisite_customer_ids,omitempty"`
	PrerequisiteSubtotalRange  struct {
		Gte string `json:"greater_than_or_equal_to,omitempty"`
	} `json:"prerequisite_subtotal_range,omitempty"`
	PrerequisiteShippingPriceRange struct {
		Lte string `json:"less_than_or_equal_to,omitempty"`
	} `json:"prerequisite_shipping_price_range,omitempty"`
	PrerequisiteQuantityRange struct {
		Gte int `json:"greater_than_or_equal_to,omitempty"`
	} `json:"prerequisite_quantity_range,omitempty"`

	PrerequisiteQuantityRatio struct {
		Quantity         int `json:"prerequisite_quantity,omitempty"`
		EntitledQuantity int `json:"entitled_quantity,omitempty"`
	} `json:"prerequisite_to_entitlement_quantity_ratio,omitempty"`

	PrerequisiteProductIDs    []int64 `json:"prerequisite_product_ids,omitempty"`
	PrerequisiteVariantIDs    []int64 `json:"prerequisite_variant_ids,omitempty"`
	PrerequisiteCollectionIDs []int64 `json:"prerequisite_collection_ids,omitempty"`

	AllocationLimit int `json:"allocation_limit,omitempty"`

	StartsAt  time.Time  `json:"starts_at,omitempty"`
	EndsAt    *time.Time `json:"ends_at,omitempty"`
	CreatedAt time.Time  `json:"created_at,omitempty"`
	UpdatedAt time.Time  `json:"updated_at,omitempty"`
}

type PriceRuleTargetSelection string
type PriceRuleTargetType string
type PriceRuleValueType string

const (
	PriceRuleTargetSelectionAll      PriceRuleTargetSelection = "all"
	PriceRuleTargetSelectionEntitled                          = "entitled"

	PriceRuleTargetTypeLineItem     PriceRuleTargetType = "line_item"     // The price rule applies to the cart's line items
	PriceRuleTargetTypeShippingLine                     = "shipping_line" // The price rule applies to the cart's shipping lines

	PriceRuleValueTypeFixedAmount PriceRuleValueType = "fixed_amount"
	PriceRuleValueTypePercentage                     = "percentage"
)

type PriceRuleParam struct {
	Limit int `json:"limit"`
	Page  int `json:"page"`
	// Show rule starting AFTER date
	StartsAtMin *time.Time `json:"starts_at_min"`
	// Show rule starting BEFORE date
	StartsAtMax *time.Time `json:"starts_at_max"`
	// Show rule ending AFTER date
	EndsAtMin *time.Time `json:"ends_at_min"`
	// Show rule ending BEFORE date
	EndsAtMax *time.Time `json:"ends_at_max"`
	//Show rule created AFTER date
	CreatedAtMin *time.Time `json:"created_at_min"`
	SinceID      int64      `json:"since_id"`
	TimesUsed    int        `json:"times_used"`
}

func (p *PriceRuleParam) EncodeQuery() string {
	if p == nil {
		return ""
	}
	// for now just allow handle
	// TODO: support all params
	v := url.Values{}

	if p.Limit > 0 {
		v.Add("limit", fmt.Sprintf("%d", p.Limit))
	}
	if p.Page > 0 {
		v.Add("page", fmt.Sprintf("%d", p.Page))
	}
	if p.CreatedAtMin != nil {
		v.Add("created_at_min", p.CreatedAtMin.Format(timeFormat))
	}
	if p.EndsAtMin != nil {
		v.Add("ends_at_min", p.EndsAtMin.Format(timeFormat))
	}
	if p.EndsAtMax != nil {
		v.Add("ends_at_max", p.EndsAtMax.Format(timeFormat))
	}
	if p.StartsAtMin != nil {
		v.Add("starts_at_min", p.StartsAtMin.Format(timeFormat))
	}
	if p.StartsAtMax != nil {
		v.Add("starts_at_max", p.StartsAtMax.Format(timeFormat))
	}
	return v.Encode()
}

func (p *PriceRuleService) List(ctx context.Context, params *PriceRuleParam) ([]*PriceRule, *http.Response, error) {
	req, err := p.client.NewRequest("GET", "/admin/price_rules.json", nil)
	if err != nil {
		return nil, nil, err
	}
	req.URL.RawQuery = params.EncodeQuery()

	var priceRuleWrapper struct {
		PriceRules []*PriceRule `json:"price_rules"`
	}
	resp, err := p.client.Do(ctx, req, &priceRuleWrapper)
	if err != nil {
		return nil, resp, err
	}

	return priceRuleWrapper.PriceRules, resp, nil
}

func (p *PriceRuleService) Get(ctx context.Context, ID int64) (*PriceRule, *http.Response, error) {
	req, err := p.client.NewRequest("GET", fmt.Sprintf("/admin/price_rules/%d.json", ID), nil)
	if err != nil {
		return nil, nil, err
	}

	var priceRuleWrapper struct {
		PriceRule *PriceRule `json:"price_rule"`
	}
	resp, err := p.client.Do(ctx, req, &priceRuleWrapper)
	if err != nil {
		return nil, resp, err
	}

	return priceRuleWrapper.PriceRule, resp, nil
}

func (p *PriceRuleService) CreatePriceRule(ctx context.Context, rule *PriceRule) (*PriceRule, *http.Response, error) {

	priceRuleWrapper := struct {
		PriceRule *PriceRule `json:"price_rule"`
	}{
		PriceRule: rule,
	}

	req, err := p.client.NewRequest("POST", "/admin/price_rules.json", priceRuleWrapper)
	if err != nil {
		return nil, nil, err
	}

	resp, err := p.client.Do(ctx, req, &priceRuleWrapper)
	if err != nil {
		return nil, resp, err
	}

	return priceRuleWrapper.PriceRule, resp, nil
}

func (p *PriceRuleService) ListDiscountCodes(ctx context.Context, ID int64) ([]*DiscountCode, *http.Response, error) {
	req, err := p.client.NewRequest("GET", fmt.Sprintf("/admin/price_rules/%d/discount_codes.json", ID), nil)
	if err != nil {
		return nil, nil, err
	}

	var discountCodeWrapper struct {
		DiscountCodes []*DiscountCode `json:"discount_codes"`
	}
	resp, err := p.client.Do(ctx, req, &discountCodeWrapper)
	if err != nil {
		return nil, resp, err
	}

	return discountCodeWrapper.DiscountCodes, resp, nil
}

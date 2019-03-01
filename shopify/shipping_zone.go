// Copyright 2019 The go-shopify AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package shopify

import (
	"context"
	"net/http"
)

type ShippingZoneService service

type ShippingZoneCountry struct {
	ID        int64                   `json:"id,omitempty"`
	Name      string                  `json:"name,omitempty"`
	Tax       float64                 `json:"tax,omitempty"`
	Code      string                  `json:"code,omitempty"`
	TaxName   string                  `json:"tax_name,omitempty"`
	Provinces []*ShippingZoneProvince `json:"provinces,omitempty"`
}

type ShippingZoneProvince struct {
	ID             int64   `json:"id"`
	CountryID      int64   `json:"country_id"`
	Name           string  `json:"name"`
	Code           string  `json:"code"`
	Tax            float64 `json:"tax"`
	TaxName        string  `json:"tax_name"`
	TaxType        string  `json:"tax_type"`
	ShippingZoneID int64   `json:"shipping_zone_id"`
	TaxPercentage  float64 `json:"tax_percentage"`
}

type ShippingZoneRate struct {
	ID             int64  `json:"id"`
	Name           string `json:"name"`
	Price          string `json:"price"`
	ShippingZoneID int64  `json:"shipping_zone_id"`

	// weight based rates
	WeightLow  float64 `json:"weight_low"`
	WeightHigh float64 `json:"weight_high"`

	// price based rates
	MinOrderSubtotal string `json:"min_order_subtotal"`
	MaxOrderSubtotal string `json:"max_order_subtotal"`

	// carrier based rates
	FlatModifier     string  `json:"flat_modifier"`
	PercentModifier  float64 `json:"percent_modifier"`
	CarrierServiceID int64   `json:"carrier_service_id"`
}

type ShippingZone struct {
	ID        int64                  `json:"id"`
	Name      string                 `json:"name"`
	Countries []*ShippingZoneCountry `json:"countries"`

	WeightBasedShippingRates     []*ShippingZoneRate `json:"weight_based_shipping_rates"`
	PriceBasedShippingRates      []*ShippingZoneRate `json:"price_based_shipping_rates"`
	CarrierShippingRateProviders []*ShippingZoneRate `json:"carrier_shipping_rate_providers"`
}

func (s *ShippingZoneService) List(ctx context.Context) ([]*ShippingZone, *http.Response, error) {
	req, err := s.client.NewRequest("GET", "/admin/shipping_zones.json", nil)
	if err != nil {
		return nil, nil, err
	}

	wrapper := struct {
		Zones []*ShippingZone `json:"shipping_zones"`
	}{}
	resp, err := s.client.Do(ctx, req, &wrapper)
	if err != nil {
		return nil, resp, err
	}

	return wrapper.Zones, resp, nil
}

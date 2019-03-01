// Copyright 2019 The go-shopify AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package shopify

import (
	"context"
	"net/http"
)

type ShopService service

type Shop struct {
	ID                              int         `json:"id"`
	Name                            string      `json:"name"`
	Email                           string      `json:"email"`
	Domain                          string      `json:"domain"`
	CreatedAt                       string      `json:"created_at"`
	Province                        string      `json:"province"`
	Country                         string      `json:"country"`
	Address1                        string      `json:"address1"`
	Zip                             string      `json:"zip"`
	City                            string      `json:"city"`
	Source                          interface{} `json:"source"`
	Phone                           string      `json:"phone"`
	UpdatedAt                       string      `json:"updated_at"`
	CustomerEmail                   string      `json:"customer_email"`
	Latitude                        float64     `json:"latitude"`
	Longitude                       float64     `json:"longitude"`
	PrimaryLocale                   string      `json:"primary_locale"`
	Address2                        string      `json:"address2"`
	CountryCode                     string      `json:"country_code"`
	CountryName                     string      `json:"country_name"`
	Currency                        string      `json:"currency"`
	Timezone                        string      `json:"timezone"`
	IanaTimezone                    string      `json:"iana_timezone"`
	ShopOwner                       string      `json:"shop_owner"`
	MoneyFormat                     string      `json:"money_format"`
	MoneyWithCurrencyFormat         string      `json:"money_with_currency_format"`
	WeightUnit                      string      `json:"weight_unit"`
	ProvinceCode                    string      `json:"province_code"`
	TaxesIncluded                   interface{} `json:"taxes_included"`
	TaxShipping                     interface{} `json:"tax_shipping"`
	CountyTaxes                     bool        `json:"county_taxes"`
	PlanDisplayName                 string      `json:"plan_display_name"`
	PlanName                        string      `json:"plan_name"`
	HasDiscounts                    bool        `json:"has_discounts"`
	HasGiftCards                    bool        `json:"has_gift_cards"`
	MyshopifyDomain                 string      `json:"myshopify_domain"`
	GoogleAppsDomain                interface{} `json:"google_apps_domain"`
	GoogleAppsLoginEnabled          interface{} `json:"google_apps_login_enabled"`
	MoneyInEmailsFormat             string      `json:"money_in_emails_format"`
	MoneyWithCurrencyInEmailsFormat string      `json:"money_with_currency_in_emails_format"`
	EligibleForPayments             bool        `json:"eligible_for_payments"`
	RequiresExtraPaymentsAgreement  bool        `json:"requires_extra_payments_agreement"`
	PasswordEnabled                 bool        `json:"password_enabled"`
	HasStorefront                   bool        `json:"has_storefront"`
	EligibleForCardReaderGiveaway   bool        `json:"eligible_for_card_reader_giveaway"`
	Finances                        bool        `json:"finances"`
	SetupRequired                   bool        `json:"setup_required"`
	ForceSsl                        bool        `json:"force_ssl"`
}

func (s *ShopService) Get(ctx context.Context) (*Shop, *http.Response, error) {
	req, err := s.client.NewRequest("GET", "/admin/shop.json", nil)
	if err != nil {
		return nil, nil, err
	}

	var shopWrapper struct {
		Shop *Shop `json:"shop"`
	}
	resp, err := s.client.Do(ctx, req, &shopWrapper)
	if err != nil {
		return nil, resp, err
	}

	return shopWrapper.Shop, resp, nil
}

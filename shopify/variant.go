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
)

type VariantService service

type Variant struct {
	ID                   int64   `json:"id"`
	ProductID            int64   `json:"product_id"`
	Title                string  `json:"title"`
	Price                string  `json:"price"`
	CompareAtPrice       string  `json:"compare_at_price"`
	Sku                  string  `json:"sku"`
	Position             int     `json:"position"`
	InventoryPolicy      string  `json:"inventory_policy"`
	FulfillmentService   string  `json:"fulfillment_service"`
	InventoryManagement  string  `json:"inventory_management"`
	Option1              string  `json:"option1"`
	Option2              string  `json:"option2"`
	Option3              string  `json:"option3"`
	CreatedAt            string  `json:"created_at"`
	UpdatedAt            string  `json:"updated_at"`
	Taxable              bool    `json:"taxable"`
	Barcode              string  `json:"barcode"`
	Grams                int     `json:"grams"`
	ImageID              int64   `json:"image_id"`
	InventoryQuantity    int     `json:"inventory_quantity"`
	Weight               float64 `json:"weight"`
	WeightUnit           string  `json:"weight_unit"`
	InventoryItemID      int     `json:"inventory_item_id"`
	OldInventoryQuantity int     `json:"old_inventory_quantity"`
	RequiresShipping     bool    `json:"requires_shipping"`
	AdminGraphqlAPIID    string  `json:"admin_graphql_api_id"`
}

type VariantParam struct {
	Limit int
	Page  int
}

func (param *VariantParam) EncodeQuery() string {
	if param == nil {
		return ""
	}
	// for now just allow handle
	// TODO: support all params
	v := url.Values{}
	v.Add("page", fmt.Sprintf("%d", param.Page))
	if param.Limit > 0 {
		v.Add("limit", fmt.Sprintf("%d", param.Limit))
	}
	return v.Encode()
}

func (p *VariantService) Get(ctx context.Context, params *VariantParam) ([]*Variant, *http.Response, error) {
	req, err := p.client.NewRequest("GET", "/admin/variants.json", nil)
	if err != nil {
		return nil, nil, err
	}
	// encode param to query
	req.URL.RawQuery = params.EncodeQuery()

	var wrapper struct {
		Variants []*Variant `json:"variants"`
	}
	resp, err := p.client.Do(ctx, req, &wrapper)
	if err != nil {
		return nil, resp, err
	}

	return wrapper.Variants, resp, nil
}

// fetch one product by the given product id
func (p *VariantService) GetVariant(ctx context.Context, ID int64) (*Variant, *http.Response, error) {
	req, err := p.client.NewRequest("GET", fmt.Sprintf("/admin/variants/%d.json", ID), nil)
	if err != nil {
		return nil, nil, err
	}

	var wrapper struct {
		Variant *Variant `json:"variant"`
	}
	resp, err := p.client.Do(ctx, req, &wrapper)
	if err != nil {
		return nil, resp, err
	}

	return wrapper.Variant, resp, nil
}

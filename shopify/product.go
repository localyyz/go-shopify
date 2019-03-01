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

type ProductService service

type Product struct {
	ProductID      int64       `json:"product_id"`
	Title          string      `json:"title"`
	BodyHTML       string      `json:"body_html"`
	Vendor         string      `json:"vendor"`
	ProductType    string      `json:"product_type"`
	Handle         string      `json:"handle"`
	TemplateSuffix interface{} `json:"template_suffix"`
	PublishedScope string      `json:"published_scope"`
	Tags           string      `json:"tags"`
	Available      bool        `json:"available"`

	Variants []*ProductVariant `json:"variants"`
	Options  []*ProductOption  `json:"options"`
	Images   []*ProductImage   `json:"images"`
	Image    *ProductImage     `json:"image"`

	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	PublishedAt time.Time `json:"published_at"`
}

type ProductVariant struct {
	ID                   int64           `json:"id"`
	ProductID            int64           `json:"product_id"`
	Title                string          `json:"title"`
	Price                string          `json:"price"`
	Sku                  string          `json:"sku"`
	Position             int             `json:"position"`
	Grams                int             `json:"grams"`
	InventoryPolicy      string          `json:"inventory_policy"`
	FulfillmentService   string          `json:"fulfillment_service"`
	InventoryManagement  string          `json:"inventory_management"`
	Option1              string          `json:"option1"`
	Option2              string          `json:"option2"`
	Option3              string          `json:"option3"`
	OptionValues         []VariantOption `json:"option_values"`
	Taxable              bool            `json:"taxable"`
	Barcode              string          `json:"barcode"`
	ImageID              interface{}     `json:"image_id"`
	CompareAtPrice       string          `json:"compare_at_price"`
	Available            bool            `json:"available"`
	InventoryQuantity    int             `json:"inventory_quantity"`
	Weight               float64         `json:"weight"`
	WeightUnit           string          `json:"weight_unit"`
	OldInventoryQuantity int             `json:"old_inventory_quantity"`
	RequiresShipping     bool            `json:"requires_shipping"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type VariantOption struct {
	OptionID int64  `json:"option_id"`
	Name     string `json:"name"`
	Value    string `json:"value"`
}

type ProductOption struct {
	ID        int64    `json:"id"`
	ProductID int64    `json:"product_id"`
	Name      string   `json:"name"`
	Position  int      `json:"position"`
	Values    []string `json:"values"`
}

type ProductImage struct {
	ID         int64   `json:"id"`
	ProductID  int64   `json:"product_id"`
	Position   int     `json:"position"`
	CreatedAt  string  `json:"created_at"`
	UpdatedAt  string  `json:"updated_at"`
	Width      int64   `json:"width"`
	Height     int64   `json:"height"`
	Src        string  `json:"src"`
	VariantIds []int64 `json:"variant_ids"`
}

func (p *ProductService) List(ctx context.Context) ([]*Product, *http.Response, error) {
	req, err := p.client.NewRequest("GET", "/admin/products.json", nil)
	if err != nil {
		return nil, nil, err
	}

	var productWrapper struct {
		Products []*Product `json:"products"`
	}
	resp, err := p.client.Do(ctx, req, &productWrapper)
	if err != nil {
		return nil, resp, err
	}

	return productWrapper.Products, resp, nil
}

func (p *ProductService) GetVariant(ctx context.Context, variantID int64) (*ProductVariant, *http.Response, error) {
	req, err := p.client.NewRequest("GET", fmt.Sprintf("/admin/variants/%d.json", variantID), nil)
	if err != nil {
		return nil, nil, err
	}

	var variantWrapper struct {
		Variant *ProductVariant `json:"variant"`
	}
	resp, err := p.client.Do(ctx, req, &variantWrapper)
	if err != nil {
		return nil, resp, err
	}

	return variantWrapper.Variant, resp, nil
}

func (p *ProductService) GetImages(ctx context.Context, productID int64) ([]*ProductImage, *http.Response, error) {
	req, err := p.client.NewRequest("GET", fmt.Sprintf("/admin/products/%d/images.json", productID), nil)
	if err != nil {
		return nil, nil, err
	}

	var imagesWrapper struct {
		Images []*ProductImage `json:"images"`
	}
	resp, err := p.client.Do(ctx, req, &imagesWrapper)
	if err != nil {
		return nil, resp, err
	}

	return imagesWrapper.Images, resp, nil
}

func (p *ProductService) GetStock(ctx context.Context, variantID int64) (int, *http.Response, error) {
	v, resp, err := p.GetVariant(ctx, variantID)
	if err != nil {
		return 0, resp, err
	}

	return v.InventoryQuantity, resp, err
}

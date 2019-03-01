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
	"strings"
	"time"
)

type ProductListService service

type ProductList Product

type ProductListParam struct {
	ProductIDs   []int64
	CollectionID int64
	Handle       string
	Limit        int
	Page         int
	UpdatedAtMin time.Time
}

const timeFormat = "2006-01-02T15:04:05-07:00"

func (p *ProductListParam) EncodeQuery() string {
	if p == nil {
		return ""
	}
	// for now just allow handle
	// TODO: support all params
	v := url.Values{}
	v.Add("handle", p.Handle)
	v.Add("page", fmt.Sprintf("%d", p.Page))
	if len(p.ProductIDs) > 0 {
		s := make([]string, len(p.ProductIDs))
		for i, id := range p.ProductIDs {
			s[i] = fmt.Sprintf("%d", id)
		}
		v.Add("product_ids", strings.Join(s, ","))
	}
	if p.Limit > 0 {
		v.Add("limit", fmt.Sprintf("%d", p.Limit))
	}
	return v.Encode()
}

func (p *ProductListService) Get(ctx context.Context, params *ProductListParam) ([]*ProductList, *http.Response, error) {
	req, err := p.client.NewRequest("GET", "/admin/product_listings.json", nil)
	if err != nil {
		return nil, nil, err
	}
	// encode param to query
	req.URL.RawQuery = params.EncodeQuery()

	var productListWrapper struct {
		ProductListings []*ProductList `json:"product_listings"`
	}
	resp, err := p.client.Do(ctx, req, &productListWrapper)
	if err != nil {
		return nil, resp, err
	}

	return productListWrapper.ProductListings, resp, nil
}

// fetch one product by the given product id
func (p *ProductListService) GetProduct(ctx context.Context, ID int64) (*ProductList, *http.Response, error) {
	req, err := p.client.NewRequest("GET", fmt.Sprintf("/admin/product_listings/%d.json", ID), nil)
	if err != nil {
		return nil, nil, err
	}

	var productListWrapper struct {
		ProductListing *ProductList `json:"product_listing"`
	}
	resp, err := p.client.Do(ctx, req, &productListWrapper)
	if err != nil {
		return nil, resp, err
	}

	return productListWrapper.ProductListing, resp, nil
}

func (p *ProductListService) Count(ctx context.Context) (int, *http.Response, error) {
	req, err := p.client.NewRequest("GET", "/admin/product_listings/count.json", nil)
	if err != nil {
		return 0, nil, err
	}

	var productCount struct {
		Count int `json:"count"`
	}
	resp, err := p.client.Do(ctx, req, &productCount)
	if err != nil {
		return 0, resp, err
	}

	return productCount.Count, resp, nil
}

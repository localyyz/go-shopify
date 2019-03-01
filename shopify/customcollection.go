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

// api reference: https://help.shopify.com/api/reference/customcollection

type CustomCollectionService service

type CustomCollection struct {
	ID     int64  `json:"id"`
	Handle string `json:"handle"`
	Title  string `json:"title"`

	BodyHTML       string                `json:"body_html"`
	SortOrder      string                `json:"sort_order"`
	PublishedScope string                `json:"published_scope"`
	Image          CustomCollectionImage `json:"image"`

	UpdatedAt   time.Time `json:"updated_at"`
	PublishedAt time.Time `json:"published_at"`
}

type CustomCollectionImage struct {
	Src       string    `json:"src"`
	Width     int       `json:"width"`
	Height    int       `json:"height"`
	CreatedAt time.Time `json:"created_at"`
}

type CustomCollectionParam struct {
	ProductID int64
}

func (p *CustomCollectionParam) EncodeQuery() string {
	if p == nil {
		return ""
	}
	// for now just allow handle
	// TODO: support all params
	v := url.Values{}
	v.Add("product_id", fmt.Sprintf("%d", p.ProductID))
	return v.Encode()
}

func (c *CustomCollectionService) Get(ctx context.Context, params *CustomCollectionParam) ([]*CustomCollection, *http.Response, error) {
	req, err := c.client.NewRequest("GET", "/admin/custom_collections.json", nil)
	if err != nil {
		return nil, nil, err
	}
	// encode param to query
	req.URL.RawQuery = params.EncodeQuery()

	var wrapper struct {
		CustomCollections []*CustomCollection `json:"custom_collections"`
	}
	resp, err := c.client.Do(ctx, req, &wrapper)
	if err != nil {
		return nil, resp, err
	}
	return wrapper.CustomCollections, resp, nil
}

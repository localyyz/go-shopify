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

type CollectionListService service

type CollectionList struct {
	ID       int64  `json:"collection_id"`
	BodyHTML string `json:"body_html"`
	//DefaultProductImage []*CollectionListImage `json:"default_product_image,omitempty"`
	Handle      string                  `json:"handle"`
	Image       *CollectionListImage    `json:"image"`
	Title       string                  `json:"title"`
	SortOrder   CollectionListSortOrder `json:"sort_order"`
	UpdatedAt   time.Time               `json:"updated_at"`
	PublishedAt time.Time               `json:"published_at"`
}

type CollectionListSortOrder uint32

type CollectionListImage struct {
	CreatedAt time.Time `json:"created_at"`
	Src       string    `json:"src"`
}

type CollectionListParam struct {
	Limit int
	Page  int
}

const (
	_ CollectionListSortOrder = iota
	CollectionListSortOrderAlphaAsc
	CollectionListSortOrderAlphaDesc
	CollectionListSortOrderBestSelling
	CollectionListSortOrderCreated
	CollectionListSortOrderCreatedDesc
	CollectionListSortOrderManual
	CollectionListSortOrderPriceAsc
	CollectionListSortOrderPriceDesc
)

var collectionListSortOrders = []string{
	"alpha-asc",    // Alphabetically, in ascending order (A - Z).
	"alpha-desc",   // Alphabetically, in descending order (Z - A).
	"best-selling", // By best-selling products.
	"created",      // By date created, in ascending order (oldest - newest).
	"created-desc", // By date created, in descending order (newest - oldest).
	"manual",       // Order created by the shop owner.
	"price-asc",    // By price, in ascending order (lowest - highest).
	"price-desc",   // By price, in descending order (highest - lowest).
}

func (p *CollectionListParam) EncodeQuery() string {
	if p == nil {
		return ""
	}
	v := url.Values{}
	v.Add("page", fmt.Sprintf("%d", p.Page))
	v.Add("limit", fmt.Sprintf("%d", p.Limit))
	return v.Encode()
}

// list all collections
func (p *CollectionListService) List(ctx context.Context, params *CollectionListParam) ([]*CollectionList, *http.Response, error) {
	req, err := p.client.NewRequest("GET", "/admin/collection_listings.json", nil)
	if err != nil {
		return nil, nil, err
	}
	// encode param to query
	req.URL.RawQuery = params.EncodeQuery()

	var collectionListWrapper struct {
		CollectionListings []*CollectionList `json:"collection_listings"`
	}
	resp, err := p.client.Do(ctx, req, &collectionListWrapper)
	if err != nil {
		return nil, resp, err
	}

	return collectionListWrapper.CollectionListings, resp, nil
}

// fetch one product by the given collection id
func (p *CollectionListService) Get(ctx context.Context, ID int64) (*CollectionList, *http.Response, error) {
	req, err := p.client.NewRequest("GET", fmt.Sprintf("/admin/collection_listings/%d.json", ID), nil)
	if err != nil {
		return nil, nil, err
	}
	var collectionListWrapper struct {
		CollectionListing *CollectionList `json:"collection_listing"`
	}
	resp, err := p.client.Do(ctx, req, &collectionListWrapper)
	if err != nil {
		return nil, resp, err
	}

	return collectionListWrapper.CollectionListing, resp, nil
}

// fetch collection product ids
func (p *CollectionListService) ListProductIDs(ctx context.Context, ID int64) ([]int64, *http.Response, error) {
	req, err := p.client.NewRequest("GET", fmt.Sprintf("/admin/collection_listings/%d/product_ids.json", ID), nil)
	if err != nil {
		return nil, nil, err
	}
	var wrapper struct {
		ProductIDs []int64 `json:"product_ids"`
	}
	resp, err := p.client.Do(ctx, req, &wrapper)
	if err != nil {
		return nil, resp, err
	}

	return wrapper.ProductIDs, resp, nil
}

// String returns the string value of the status.
func (s CollectionListSortOrder) String() string {
	return collectionListSortOrders[s]
}

// MarshalText satisfies TextMarshaler
func (s CollectionListSortOrder) MarshalText() ([]byte, error) {
	return []byte(s.String()), nil
}

// UnmarshalText satisfies TextUnmarshaler
func (s *CollectionListSortOrder) UnmarshalText(text []byte) error {
	enum := string(text)
	for i := 0; i < len(collectionListSortOrders); i++ {
		if enum == collectionListSortOrders[i] {
			*s = CollectionListSortOrder(i)
			return nil
		}
	}
	return fmt.Errorf("unknown collection list sort order %s", enum)
}

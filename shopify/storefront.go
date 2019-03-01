// Copyright 2019 The go-shopify AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package shopify

import (
	"context"
	"net/http"
)

type StorefrontService service

type Storefront struct {
	ID          int64  `json:"id,omitempty"`
	AccessToken string `json:"access_token,omitempty"`
	AccessScope string `json:"access_scope,omitempty"`
	CreatedAt   string `json:"created_at,omitempty"`
	Title       string `json:"title,omitempty"`
}

type StorefrontWrapper struct {
	Storefront *Storefront `json:"storefront_access_token"`
}

func (f *StorefrontService) Create(ctx context.Context, title string) (*Storefront, *http.Response, error) {
	w := struct {
		Title string `json:"title"`
	}{"Localyyz"}

	req, err := f.client.NewRequest("POST", "/admin/storefront_access_tokens.json", w)
	if err != nil {
		return nil, nil, err
	}

	ss := new(Storefront)
	resp, err := f.client.Do(ctx, req, ss)
	if err != nil {
		return nil, resp, err
	}

	return ss, resp, nil
}

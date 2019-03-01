// Copyright 2019 The go-shopify AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package shopify

import (
	"context"
	"net/http"
	"time"
)

type PolicyService service

type Policy struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	URL       string    `json:"url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

const PolicyRefund = "Refund Policy"

func (s *PolicyService) List(ctx context.Context) ([]*Policy, *http.Response, error) {
	req, err := s.client.NewRequest("GET", "/admin/policies.json", nil)
	if err != nil {
		return nil, nil, err
	}

	wrapper := struct {
		Policies []*Policy `json:"policies"`
	}{}
	resp, err := s.client.Do(ctx, req, &wrapper)
	if err != nil {
		return nil, resp, err
	}

	return wrapper.Policies, resp, nil
}

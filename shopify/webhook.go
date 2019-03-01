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

type WebhookService service

// Webhook Topics for subscription
//
// Shopify API docs: https://help.shopify.com/api/reference/webhook#create
type Topic string

const (
	TopicUnknown        Topic = ""
	TopicAppUninstalled Topic = "app/uninstalled"

	TopicCartsCreate Topic = "carts/create"
	TopicCartsUpdate Topic = "carts/update"

	TopicCheckoutsCreate Topic = "checkouts/create"
	TopicCheckoutsDelete Topic = "checkouts/delete"
	TopicCheckoutsUpdate Topic = "checkouts/update"

	TopicCollectionListingsAdd    Topic = "collection_listings/add"
	TopicCollectionListingsRemove Topic = "collection_listings/remove"
	TopicCollectionListingsUpdate Topic = "collection_listings/update"

	TopicCollectionsCreate Topic = "collections/create"
	TopicCollectionsDelete Topic = "collections/delete"
	TopicCollectionsUpdate Topic = "collections/update"

	TopicCustomerGroupsCreate Topic = "customer_groups/create"
	TopicCustomerGroupsDelete Topic = "customer_groups/delete"
	TopicCustomerGroupsUpdate Topic = "customer_groups/update"

	TopicCustomersCreate  Topic = "customers/create"
	TopicCustomersDelete  Topic = "customers/delete"
	TopicCustomersDisable Topic = "customers/disable"
	TopicCustomersEnable  Topic = "customers/enable"
	TopicCustomersUpdate  Topic = "customers/update"

	TopicDraftOrdersCreate Topic = "draft_orders/create"
	TopicDraftOrdersDelete Topic = "draft_orders/delete"
	TopicDraftOrdersUpdate Topic = "draft_orders/update"

	TopicFulfillmentEventsCreate Topic = "fulfillment_events/create"
	TopicFulfillmentEventsDelete Topic = "fulfillment_events/delete"
	TopicFulfillmentsCreate      Topic = "fulfillments/create"
	TopicFulfillmentsUpdate      Topic = "fulfillments/update"

	TopicOrderTransactionsCreate  Topic = "order_transactions/create"
	TopicOrdersCancelled          Topic = "orders/cancelled"
	TopicOrdersCreate             Topic = "orders/create"
	TopicOrdersDelete             Topic = "orders/delete"
	TopicOrdersFulfilled          Topic = "orders/fulfilled"
	TopicOrdersPaid               Topic = "orders/paid"
	TopicOrdersPartiallyFulfilled Topic = "orders/partially_fulfilled"
	TopicOrdersUpdated            Topic = "orders/updated"

	TopicProductListingsAdd    Topic = "product_listings/add"
	TopicProductListingsRemove Topic = "product_listings/remove"
	TopicProductListingsUpdate Topic = "product_listings/update"

	TopicProductsCreate Topic = "products/create"
	TopicProductsDelete Topic = "products/delete"
	TopicProductsUpdate Topic = "products/update"

	TopicRefundsCreate Topic = "refunds/create"

	TopicShopUpdate Topic = "shop/update"

	TopicThemesCreate  Topic = "themes/create"
	TopicThemesDelete  Topic = "themes/delete"
	TopicThemesPublish Topic = "themes/publish"
	TopicThemesUpdate  Topic = "themes/update"
)

const (
	WebhookHeaderHmac       = "X-Shopify-Hmac-Sha256"
	WebhookHeaderShopDomain = "X-Shopify-Shop-Domain"
	WebhookHeaderTopic      = "X-Shopify-Topic"
)

type Webhook struct {
	ID                  int       `json:"id,omitempty"`
	Address             string    `json:"address"`
	Topic               Topic     `json:"topic"`
	Format              string    `json:"format"`
	CreatedAt           time.Time `json:"created_at,omitempty"`
	UpdatedAt           time.Time `json:"updated_at,omitempty"`
	Fields              []string  `json:"fields,omitempty"`
	MetafieldNamespaces []string  `json:"metafield_namespaces,omitempty"`
}

// WebhookRequest represents a request to create/edit a webhook.
// It is just an alias of Webhook struct
type WebhookRequest struct {
	Webhook *Webhook `json:"webhook"`
}

func (s *WebhookService) List(ctx context.Context) ([]*Webhook, *http.Response, error) {
	req, err := s.client.NewRequest("GET", "/admin/webhooks.json", nil)
	if err != nil {
		return nil, nil, err
	}
	var webhooksWrapper struct {
		Webhooks []*Webhook `json:"webhooks"`
	}
	resp, err := s.client.Do(ctx, req, &webhooksWrapper)
	if err != nil {
		return nil, resp, err
	}
	return webhooksWrapper.Webhooks, resp, nil
}

func (s *WebhookService) Create(ctx context.Context, webhook *WebhookRequest) (*Webhook, *http.Response, error) {
	req, err := s.client.NewRequest("POST", "/admin/webhooks.json", webhook)
	if err != nil {
		return nil, nil, err
	}

	ww := new(WebhookRequest)
	resp, err := s.client.Do(ctx, req, ww)
	if err != nil {
		return nil, resp, err
	}
	return ww.Webhook, resp, nil
}

func (s WebhookService) Delete(ctx context.Context, ID int64) (*http.Response, error) {
	req, err := s.client.NewRequest("DELETE", fmt.Sprintf("/admin/webhooks/%d.json", ID), nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// Copyright 2019 The go-shopify AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package shopify

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
)

// Shopify errors usually have the form:
// {
//   "errors": {
//     "title": [
//       "something is wrong"
//     ]
//   }
// }
//

// TODO: https://github.com/Jeffail/gabs (????)

type ShopifyErrorer interface {
	Type() string
}

type lineItemErrorValue struct {
	Message string `json:"message"`
	Options struct {
		Remaining int `json:"remaining"`
	} `json:"options"`
	Code string `json:"code"`
}
type lineItemErrorField map[string][]lineItemErrorValue

type LineItemError struct {
	ShopifyErrorer

	Field    string
	Message  string
	Code     string
	Position string

	//Quantity []struct {
	//Message string `json:"message"`
	//Options struct {
	//Remaining int `json:"remaining"`
	//} `json:"options"`
	//Code string `json:"code"`
	//} `json:"quantity"`
}

type DiscountCodeError struct {
	Reason interface{} `json:"reason"`
	ShopifyErrorer
}

type AddressError struct {
	ShopifyErrorer

	Key     string `json:"key"`
	Field   string `json:"field"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

type EmailError struct {
	ShopifyErrorer
	Message string `json:"message"`
}

type ErrorResponse struct {
	Errors interface{} `json:"errors"`
}

var (
	// TODO: make this an unmarshall type
	ErrNotEnoughInStock = `not_enough_in_stock`
)

func (e *LineItemError) Error() string {
	return fmt.Sprintf("%s at pos(%s): %s %s", e.Type(), e.Position, e.Field, e.Message)
}

func (e *LineItemError) Type() string {
	return `line_items`
}

func (e *DiscountCodeError) Error() string {
	return fmt.Sprintf("%+v", e.Reason)
}

func (e *DiscountCodeError) Type() string {
	return `discount_code`
}

func (e *AddressError) Error() string {
	return fmt.Sprintf("%s: %s %s", e.Type(), e.Field, e.Message)
}

func (e *AddressError) Type() string {
	return e.Key
}

func (e *EmailError) Error() string {
	return fmt.Sprintf("email %s", e.Message)
}

func (e *EmailError) Type() string {
	return `email`
}

func (r *ErrorResponse) Error() string {
	if e, ok := r.Errors.(map[string]interface{}); ok {
		for k, v := range e {
			// value here can be a slice
			return fmt.Sprintf("%s: %+v", k, v)
		}
	}
	if e, ok := r.Errors.(string); ok {
		return e
	}
	return "unknown, unparsed error"
}

func toAddressError(key, field string, listError []interface{}) *AddressError {
	for _, ee := range listError {
		// NOTE: parse the first error found
		if ex, _ := ee.(map[string]interface{}); ex != nil {
			code, _ := ex["code"].(string)
			message, _ := ex["message"].(string)
			return &AddressError{
				Key:     key,
				Field:   field,
				Code:    code,
				Message: message,
			}
		}
	}
	return nil
}

// CheckResponse checks the API response for errors, and returns them if
// present. A response is considered an error if it has a status code outside
// the 200 range or equal to 202 Accepted.
// API error responses are expected to have either no response
// body, or a JSON response body that maps to ErrorResponse. Any other
// response body will be silently ignored.
func CheckResponse(r *http.Response) error {
	if r.StatusCode == http.StatusAccepted {
		return nil
	}
	if r.StatusCode == http.StatusForbidden {
		return errors.New("forbidden")
	}
	if c := r.StatusCode; 200 <= c && c <= 299 {
		return nil
	}
	errorResponse := &ErrorResponse{}
	data, err := ioutil.ReadAll(r.Body)
	if err == nil && data != nil {
		json.Unmarshal(data, errorResponse)
	}
	return findFirstError(errorResponse)
}

func findFirstError(r *ErrorResponse) error {
	rr, ok := r.Errors.(map[string]interface{})
	if !ok {
		return r
	}

	// find the first error, and return
	for k, v := range rr {
		if k == "email" {
			return &EmailError{
				Message: "is invalid",
			}
		}
		if vv, ok := v.(map[string]interface{}); ok {
			switch k {
			//TODO: shipping_line: map[id:[map[code:expired message:has expired options:map[]]]]
			case "line_items":
				for pos, vvv := range vv {
					b, _ := json.Marshal(vvv)

					var e lineItemErrorField
					json.Unmarshal(b, &e)

					ee := &LineItemError{
						Position: pos,
					}
					if e != nil {
						for ek, ev := range e {
							ee.Field = ek
							ee.Message = ev[0].Message
							ee.Code = ev[0].Code
							break
						}
					}

					return ee
				}

			case "checkout":
				for kk, vvv := range vv {
					switch kk {
					case "discount_code":
						// list of fields
						if e, _ := vvv.([]interface{}); e != nil {
							for _, ee := range e {
								if ex, _ := ee.(map[string]interface{}); ex != nil {
									for _, reason := range ex {
										return &DiscountCodeError{Reason: reason}
									}
								}
							}
						}
					}
				}
			case "shipping_address", "billing_address":
				for kk, vvv := range vv {
					if e, ok := vvv.([]interface{}); ok && e != nil {
						return toAddressError(k, kk, e)
					}
				}
			}
		} else if vv, ok := v.([]interface{}); ok {
			switch k {
			case "discount_code":
				for _, vvv := range vv {
					vvvv := vvv.(map[string]interface{})
					return &DiscountCodeError{
						Reason: vvvv["message"],
					}
				}

			}
		} else {
			// TODO: concrete error type..
			// I'm not keen on using another package like errorx to wrap error
			// cause here. maybe need to expand the error handling a lot more
			// for things to make sense.
			return fmt.Errorf("unknown shopify error key %s, %+v, %v", k, v, reflect.TypeOf(v))
		}
	}

	return r
}

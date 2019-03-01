# go-shopify

Go library for accessing the [Shopify REST API](https://help.shopify.com/en/api/reference) - (GoDocs coming soon)

go-shopify requires Go version 1.9 or greater.

## Usage ##

```go
import "github.com/localyyz/go-shopify" // v1
```

Construct a new Shopify client, then use the various services on the client to
access different parts of the Shopify API. For example:

```go
client := shopify.NewClient(nil, shopify.ShopURL("https://x.myshopify.com/admin"), shopify.Token("xxx-yyy-zzz"))

// get the store metadata
store, _, err := client.Shop.Get(context.Background())
```

Some API methods have optional parameters that can be passed. For example:

```go
client := shopify.NewClient(nil, shopify.ShopURL("https://x.myshopify.com/admin"), shopify.Token("xxx-yyy-zzz"))

// create a new checkout
checkout := &shopify.Checkout{Email: "paul@somebuyer.com"}
checkout, _, err := client.Checkout.Create(context.Background(), checkout)
```

The services of a client divide the API into logical chunks and correspond to
the structure of the Shopify API documentation at
https://help.shopify.com/en/api/reference.

NOTE: Using the [context](https://godoc.org/context) package, one can easily
pass cancelation signals and deadlines to various services of the client for
handling a request. In case there is no context available, then `context.Background()`
can be used as a starting point.

## Roadmap ##

This library is being initially developed for production use at
Localyyz, it's now open sourced with the community to as the library have been
fairly stable for the past year or so. For future work / improvements:

- tests coverage
- example usage
- documentation and GoDoc generation
- complete REST api converage

## Contributing ##

Feel free to fork and submit changes upstream. Please follow the pattern
established in this package. If you have questions feel free to submit issues,
for inspiration and guidelines please refer to the [google/go-github][] library.

## License ##

This library is distributed under the BSD-style license found in the [LICENSE](./LICENSE)
file.

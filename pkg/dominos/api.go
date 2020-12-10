package dominos

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/cirocosta/pizza-controller/pkg/dominos/internal/api"
)

const (
	CanadaURL       = "https://order.dominos.ca"
	UnitedStatesURL = "https://order.dominos.com"

	PathPlaceOrder   = "/power/place-order"
	PathPriceOrder   = "/power/price-order"
	PathStoreLocator = "/power/store-locator"
	PathStoreMenu    = "/power/store/%s/menu"
)

type Client struct {
	host   *url.URL
	client *http.Client
}

func NewClient(host string, debug bool) (*Client, error) {
	h, err := url.Parse(host)
	if err != nil {
		return nil, fmt.Errorf("url parse '%s': %w", host, err)
	}

	var transport http.RoundTripper = &DumpTransport{r: http.DefaultTransport}
	if !debug {
		transport = http.DefaultTransport
	}

	httpClient := &http.Client{
		Timeout:   15 * time.Second,
		Transport: transport,
	}

	return &Client{
		host:   h,
		client: httpClient,
	}, nil
}

func (c *Client) PlaceOrder(ctx context.Context, order Order) (string, error) {
	url := *c.host
	url.Path = PathPlaceOrder

	msg := c.orderMessage(order)
	buf := &bytes.Buffer{}
	if err := json.NewEncoder(buf).Encode(&msg); err != nil {
		return "", fmt.Errorf("encode order: %w", err)
	}

	resp, err := c.client.Post(url.String(), "application/json", buf)
	if err != nil {
		return "", fmt.Errorf("post %s: %w", url, err)
	}
	defer resp.Body.Close()

	body := api.PlaceOrderResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return "", fmt.Errorf("decode response: %w", err)
	}

	if body.Status == -1 {
		return "", fmt.Errorf("status -1: %s", body.Order.StatusItems.String())
	}

	return body.Order.OrderID, nil
}

func (c *Client) PriceOrder(ctx context.Context, order Order) (string, error) {
	url := *c.host
	url.Path = PathPriceOrder

	order.CreditCard = CreditCard{}

	msg := c.orderMessage(order)
	buf := &bytes.Buffer{}
	if err := json.NewEncoder(buf).Encode(&msg); err != nil {
		return "", fmt.Errorf("encode order: %w", err)
	}

	resp, err := c.client.Post(url.String(), "application/json", buf)
	if err != nil {
		return "", fmt.Errorf("post %s: %w", url, err)
	}
	defer resp.Body.Close()

	body := api.PriceResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return "", fmt.Errorf("decode response: %w", err)
	}

	if body.Status == -1 {
		return "", fmt.Errorf("status -1: %s", body.Order.CorrectiveAction.Code)
	}

	return fmt.Sprintf("%f", body.Order.Amounts.Customer), nil
}

func (c *Client) StoreMenu(ctx context.Context, storeID string) ([]*Product, error) {
	url := *c.host
	url.Path = fmt.Sprintf(PathStoreMenu, storeID)

	v := url.Query()
	v.Set("lang", "en")
	v.Set("structured", "true")

	url.RawQuery = v.Encode()

	resp, err := c.client.Get(url.String())
	if err != nil {
		return nil, fmt.Errorf("get '%s': %w", url, err)
	}

	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("get status code: %d", resp.StatusCode)
	}

	body := api.MenuResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	res := []*Product{}
	for _, product := range body.Preconfigured {
		res = append(res, &Product{
			ID:          product.Code,
			Description: product.Description,
			Name:        product.Name,
			Size:        product.Size,
		})
	}

	return res, nil
}

func (c *Client) StoresNearby(ctx context.Context, addr Address, service Service) ([]*Store, error) {
	url := *c.host
	url.Path = PathStoreLocator

	v := url.Query()
	v.Set("s", fmt.Sprintf("%s %s", addr.StreetNumber, addr.StreetName))
	v.Set("c", fmt.Sprintf("%s, %s %s", addr.City, addr.State, addr.Zip))
	v.Set("type", string(service))

	url.RawQuery = v.Encode()

	resp, err := c.client.Get(url.String())
	if err != nil {
		return nil, fmt.Errorf("get '%s': %w", url, err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("get status code: %d", resp.StatusCode)
	}

	defer resp.Body.Close()

	body := api.StoreLocatorResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	stores := []*Store{}
	for _, store := range body.Stores {
		if !store.IsOpen {
			continue
		}

		if service == ServiceCarryout && !store.ServiceIsOpen.Carryout {
			continue
		}

		if service == ServiceDelivery && !store.ServiceIsOpen.Delivery {
			continue
		}

		stores = append(stores, &Store{
			ID:      store.StoreID,
			Phone:   store.Phone,
			Address: store.AddressDescription,
		})
	}

	return stores, nil
}

func (c *Client) orderMessage(order Order) api.OrderMessage {
	msg := api.OrderMessage{
		Order: api.Order{
			Address: &api.StreetAddr{
				Street:       order.Address.StreetNumber + " " + order.Address.StreetName,
				StreetName:   order.Address.StreetName,
				StreetNumber: order.Address.StreetNumber,
				City:         order.Address.City,
				State:        order.Address.State,
				Zipcode:      order.Address.Zip,
				AddrType:     api.AddressTypeHouse,
			},
			LanguageCode:  "en",
			StoreID:       order.StoreID,
			ServiceMethod: string(order.Service),
			Payments:      []*api.OrderPayment{},
			Products:      []*api.OrderProduct{},
		},
	}

	if order.PersonalInformation.FirstName != "" {
		msg.Order.FirstName = order.PersonalInformation.FirstName
		msg.Order.LastName = order.PersonalInformation.LastName
		msg.Order.Email = order.PersonalInformation.Email
		msg.Order.Phone = order.PersonalInformation.Phone
	}

	if order.CreditCard.Number != "" {
		msg.Order.Payments = append(msg.Order.Payments, &api.OrderPayment{
			Type:     "DoorCredit",
			CardType: string(order.CreditCard.Type),
			Amount:   order.Amount,
			// Number:       order.CreditCard.Number,
			// Expiration:   order.CreditCard.Expiration,
			// SecurityCode: order.CreditCard.SecurityCode,
			// PostalCode:   order.CreditCard.PostalCode,
		})
	}

	for idx, product := range order.Products {
		msg.Order.Products = append(msg.Order.Products, &api.OrderProduct{
			ID:  idx,
			Qty: 1,
			ItemCommon: api.ItemCommon{
				Code: product.ID,
			},
		})
	}

	return msg
}

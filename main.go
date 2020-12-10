package main

import (
	"context"
	"fmt"

	"github.com/cirocosta/pizza-controller/pkg/dominos"
)

func run() error {
	ctx := context.Background()

	client, err := dominos.NewClient(dominos.CanadaURL, true)
	if err != nil {
		return fmt.Errorf("new client: %w", err)
	}

	var (
		addr = dominos.Address{
			StreetNumber: "90",
			StreetName:   "Queens Wharf Rd",
			City:         "Toronto",
			State:        "ON",
			Zip:          "M5V0J4",
		}
		// pi = dominos.PersonalInformation{
		// 	FirstName: "",
		// 	LastName:  "",
		// 	Email:     "",
		// 	Phone:     "",
		// }
		// cc = dominos.CreditCard{
		// 	Type:         dominos.CreditCardTypeVisa,
		// 	Number:       "",
		// 	Expiration:   "",
		// 	SecurityCode: "",
		// 	PostalCode:   "",
		// }
		products = []dominos.Product{
			{
				ID: "10SCREEN",
			},
		}
		service = dominos.ServiceCarryout
	)

	stores, err := client.StoresNearby(ctx, addr, service)
	if err != nil {
		return fmt.Errorf("stores nearby: %w", err)
	}

	storeID := stores[0].ID
	items, err := client.StoreMenu(ctx, storeID)
	if err != nil {
		return fmt.Errorf("store '%s' menu: %w", storeID, err)
	}

	for _, item := range items {
		fmt.Println(item.ID, "\t", item.Name)
	}

	order := dominos.Order{
		StoreID:  storeID,
		Address:  addr,
		Products: products,
		Service:  service,
	}

	price, err := client.PriceOrder(ctx, order)
	if err != nil {
		return fmt.Errorf("submit order: %w", err)
	}

	fmt.Println(price)

	return nil
}

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

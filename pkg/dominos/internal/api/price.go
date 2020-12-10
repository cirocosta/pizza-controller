package api

import "strings"

type PriceResponse struct {
	Order struct {
		Amounts struct {
			Customer float64 `json:"Customer"`
		} `json:"Amounts"`
		CorrectiveAction struct {
			Action string `json:"Action"`
			Code   string `json:"Code"`
			Detail string `json:"Detail"`
		} `json:"CorrectiveAction"`
	} `json:"Order"`
	Status int `json:"Status"`
}

type PlaceOrderResponse struct {
	Order struct {
		OrderID string `json:"OrderID"`
		Amounts struct {
			Customer float64 `json:"Customer"`
		} `json:"Amounts"`
		EstimatedWaitMinutes string                        `json:"EstimatedWaitMinutes"`
		StatusItems          PlaceOrderResponseStatusItems `json:"StatusItems"`
	} `json:"Order"`
	Status int `json:"Status"`
}

type PlaceOrderResponseStatusItems []struct {
	Code      string `json:"Code"`
	PulseText string `json:"PulseText"`
}

func (p PlaceOrderResponseStatusItems) String() string {
	res := []string{}
	for _, item := range p {
		txt := item.Code
		if item.PulseText != "" {
			txt += " " + item.PulseText
		}

		res = append(res, txt)
	}

	return strings.Join(res, ",")
}

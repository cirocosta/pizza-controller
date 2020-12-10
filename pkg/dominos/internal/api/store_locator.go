package api

type StoreLocatorResponse struct {
	Granularity string `json:"Granularity"`
	Address     struct {
		Street       string `json:"Street"`
		StreetNumber string `json:"StreetNumber"`
		StreetName   string `json:"StreetName"`
		UnitType     string `json:"UnitType"`
		UnitNumber   string `json:"UnitNumber"`
		City         string `json:"City"`
		Region       string `json:"Region"`
		PostalCode   string `json:"PostalCode"`
	} `json:"Address"`
	AlternativeAddress []interface{} `json:"AlternativeAddress"`
	Stores             []struct {
		StoreID                 string      `json:"StoreID"`
		IsDeliveryStore         bool        `json:"IsDeliveryStore"`
		MinDistance             interface{} `json:"MinDistance"`
		MaxDistance             interface{} `json:"MaxDistance"`
		Phone                   string      `json:"Phone"`
		AddressDescription      string      `json:"AddressDescription"`
		HolidaysDescription     string      `json:"HolidaysDescription"`
		HoursDescription        string      `json:"HoursDescription"`
		ServiceHoursDescription struct {
			Carryout        string `json:"Carryout"`
			Delivery        string `json:"Delivery"`
			DriveUpCarryout string `json:"DriveUpCarryout"`
		} `json:"ServiceHoursDescription"`
		IsOnlineCapable      bool   `json:"IsOnlineCapable"`
		IsOnlineNow          bool   `json:"IsOnlineNow"`
		IsNEONow             bool   `json:"IsNEONow"`
		IsSpanish            bool   `json:"IsSpanish"`
		LocationInfo         string `json:"LocationInfo"`
		LanguageLocationInfo struct {
			En string `json:"en"`
		} `json:"LanguageLocationInfo"`
		AllowDeliveryOrders               bool `json:"AllowDeliveryOrders"`
		AllowCarryoutOrders               bool `json:"AllowCarryoutOrders"`
		AllowDuc                          bool `json:"AllowDuc"`
		ServiceMethodEstimatedWaitMinutes struct {
			Delivery struct {
				Min int `json:"Min"`
				Max int `json:"Max"`
			} `json:"Delivery"`
			Carryout struct {
				Min int `json:"Min"`
				Max int `json:"Max"`
			} `json:"Carryout"`
		} `json:"ServiceMethodEstimatedWaitMinutes"`
		StoreCoordinates struct {
			StoreLatitude  interface{} `json:"StoreLatitude"`
			StoreLongitude interface{} `json:"StoreLongitude"`
		} `json:"StoreCoordinates"`
		AllowPickupWindowOrders bool   `json:"AllowPickupWindowOrders"`
		ContactlessDelivery     string `json:"ContactlessDelivery"`
		ContactlessCarryout     string `json:"ContactlessCarryout"`
		IsOpen                  bool   `json:"IsOpen"`
		ServiceIsOpen           struct {
			Carryout        bool `json:"Carryout"`
			Delivery        bool `json:"Delivery"`
			DriveUpCarryout bool `json:"DriveUpCarryout"`
		} `json:"ServiceIsOpen"`
	} `json:"Stores"`
}

package api

type OrderMessage struct {
	Order Order `json:"Order"`
}

type Order struct {
	Address       *StreetAddr            `json:"Address"`
	CustomerID    string                 `json:",omitempty"` // leave empty
	Email         string                 `json:"Email"`
	FirstName     string                 `json:"FirstName"`
	LanguageCode  string                 `json:"LanguageCode"`
	LastName      string                 `json:"LastName"`
	MetaData      map[string]interface{} `json:"metaData"`
	OrderID       string                 `json:"OrderID"`
	OrderName     string                 `json:"-"`
	Payments      []*OrderPayment        `json:"Payments"`
	Phone         string                 `json:"Phone"`
	Products      []*OrderProduct        `json:"Products"`
	ServiceMethod string                 `json:"ServiceMethod"`
	StoreID       string                 `json:"StoreID"`
}

type AddressType string

const (
	AddressTypeHouse     AddressType = "House"
	AddressTypeApartment AddressType = "Apartment"
	AddressTypeBusiness  AddressType = "Business"
	AddressTypeHotel     AddressType = "Hotel"
	AddressTypeOther     AddressType = "Other"
)

// StreetAddr represents a street address
type StreetAddr struct {
	Street       string      `json:"Street"`       // street number followed by street name
	StreetNumber string      `json:"StreetNumber"` // just the street number
	StreetName   string      `json:"StreetName"`   // just the street name
	City         string      `json:"City"`
	State        string      `json:"Region"`
	Zipcode      string      `json:"PostalCode"`
	AddrType     AddressType `json:"Type"`
}

type OrderProduct struct {
	ItemCommon

	Qty                int                    `json:"Qty"`
	ID                 int                    `json:"ID"` // index of the product within an order
	IsNew              bool                   `json:"isNew"`
	NeedsCustomization bool                   `json:"NeedsCustomization"`
	Opts               map[string]interface{} `json:"Options"`
}

// this is the struct that will actually be turning into json an will
// be sent to dominos.
type OrderPayment struct {
	Number       string `json:"Number"`
	Expiration   string `json:"Expiration"`
	SecurityCode string `json:"SecurityCode"`
	Type         string `json:"Type"`
	CardType     string `json:"CardType"`
	PostalCode   string `json:"PostalCode"`

	// These next fields are just for dominos
	Amount         float64 `json:"Amount"`
	CardID         string  `json:"CardID,omitempty"`
	ProviderID     string  `json:"ProviderID"`
	OTP            string  `json:"OTP"`
	GpmPaymentType string  `json:"gpmPaymentType,omitempty"`
}

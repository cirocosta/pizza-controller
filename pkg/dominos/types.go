package dominos

type Service string

const (
	ServiceDelivery Service = "Delivery"
	ServiceCarryout Service = "Carryout"
)

type CreditCardType string

const (
	CreditCardTypeMastercard CreditCardType = "MASTERCARD"
	CreditCardTypeVisa       CreditCardType = "VISA"
	CreditCardTypeAmex       CreditCardType = "AMEX"
)

type CreditCard struct {
	Type         CreditCardType
	Expiration   string
	Number       string
	PostalCode   string
	SecurityCode string
}

type Address struct {
	StreetName   string
	StreetNumber string
	City         string
	State        string
	Zip          string
}

type Store struct {
	ID      string
	Phone   string
	Address string
}

type Product struct {
	ID          string
	Description string
	Name        string
	Size        string
}

type PersonalInformation struct {
	FirstName string
	LastName  string
	Email     string
	Phone     string
}

type Order struct {
	StoreID             string
	PersonalInformation PersonalInformation
	Address             Address
	Products            []Product
	CreditCard          CreditCard
	Service             Service
	Amount              float64
}

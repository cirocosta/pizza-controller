package api

type MenuResponse struct {
	Categorization struct {
		Food          MenuCategory
		Coupons       MenuCategory
		Preconfigured MenuCategory `json:"PreconfiguredProducts"`
	} `json:"Categorization"`
	Products      map[string]*Product
	Variants      map[string]*Variant
	Toppings      map[string]map[string]Topping
	Preconfigured map[string]*PreConfiguredProduct `json:"PreconfiguredProducts"`
	Sides         map[string]map[string]struct {
		ItemCommon
		Description string
	}
}

// Variant is a structure that represents a base component of the Dominos menu.
// It will be a subset of a Product (see Product).
type Variant struct {
	ItemCommon

	Price       string
	ProductCode string
	Prepared    bool
}

// MenuCategory is a category on the dominos menu.
type MenuCategory struct {
	Name        string
	Code        string
	Description string
	Categories  []MenuCategory
	Products    []string
}

type ItemCommon struct {
	Code  string
	Name  string
	Tags  map[string]interface{}
	Local bool
}

// PreConfiguredProduct is pre-configured product.
type PreConfiguredProduct struct {
	ItemCommon
	Description string `json:"Description"`
	Opts        string `json:"Options"`
	Size        string `json:"Size"`
}

// Topping is a simple struct that represents a topping on the menu.
//
// Note: this struct does not rempresent a topping that is added to an Item
// and sent to dominos.
type Topping struct {
	ItemCommon

	Description  string
	Availability []interface{}
}

// Product is the structure representing a dominos product. The Product struct
// is meant to instaniated with json data and should be treated as such.
//
// Product is not a the most basic component of the Dominos menu; this is where
// the Variant structure comes in. The Product structure can be seen as more of
// a category that houses a list of Variants. Products are still able to be ordered,
// however.
type Product struct {
	ItemCommon

	Variants          []string
	Description       string
	AvailableToppings string
	AvailableSides    string
	DefaultToppings   string
	DefaultSides      string
	ProductType       string
}

package providerA

type Request struct {
	Weight              Weight               `json:"weight"`
	Shipper             Party                `json:"shipper"`
	Consignee           Party                `json:"consignee"`
	Dimensions          Dimensions           `json:"dimensions"`
	Account             Account              `json:"account"`
	ProductCode         string               `json:"productCode"`
	ServiceType         string               `json:"serviceType"`
	PrintType           string               `json:"printType"`
	IsInsured           bool                 `json:"isInsured"`
	CustomsDeclarations []CustomsDeclaration `json:"customsDeclarations"`
	DeclaredValue       DeclaredValue        `json:"declaredValue"`
	NumberOfPieces      int                  `json:"numberOfPieces"`
	ReferenceNumber1    string               `json:"referenceNumber1"`
	ReferenceNumber2    string               `json:"referenceNumber2"`
	ReferenceNumber3    string               `json:"referenceNumber3"`
	ReferenceNumber4    string               `json:"referenceNumber4"`
	SpecialNotes        string               `json:"specialNotes"`
	Remarks             string               `json:"remarks"`
	DeliveryType        string               `json:"deliveryType"`
	ContentType         string               `json:"contentType"`
	IsCOD               bool                 `json:"isCod"`
}

type Weight struct {
	Value float64 `json:"value"`
	Unit  string  `json:"unit"`
}

type Party struct {
	Contact      Contact `json:"contact"`
	Address      Address `json:"address"`
	ReferenceNo1 string  `json:"referenceNo1"`
	ReferenceNo2 string  `json:"referenceNo2"`
}

type Contact struct {
	Name         string `json:"name"`
	MobileNumber string `json:"mobileNumber"`
	PhoneNumber  string `json:"phoneNumber"`
	EmailAddress string `json:"emailAddress"`
	CompanyName  string `json:"companyName"`
}

type Address struct {
	Line1       string `json:"line1"`
	City        string `json:"city"`
	CountryCode string `json:"countryCode"`
	ZipCode     string `json:"zipCode"`
}

type Dimensions struct {
	Length float64 `json:"length"`
	Height float64 `json:"height"`
	Width  float64 `json:"width"`
	Unit   string  `json:"unit"`
}

type Account struct {
	Number int `json:"number"`
}

type CustomsDeclaration struct {
	Reference       string     `json:"reference"`
	Description     string     `json:"description"`
	CountryOfOrigin string     `json:"countryOfOrigin"`
	Weight          float64    `json:"weight"`
	Dimensions      Dimensions `json:"dimensions"`
	Quantity        int        `json:"quantity"`
	HSCode          string     `json:"hsCode"`
	Value           float64    `json:"value"`
}

type DeclaredValue struct {
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
}

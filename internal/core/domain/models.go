package domain

import "time"

type GenericShippingRequest struct {
	Weight              WeightInfo             `json:"weight"`
	Shipper             Party                  `json:"shipper"`
	Consignee           Party                  `json:"consignee"`
	Dimensions          Dimensions             `json:"dimensions"`
	Account             AccountInfo            `json:"account"`
	ProductCode         string                 `json:"productCode"`
	ServiceType         string                 `json:"serviceType"`
	IsInsured           bool                   `json:"isInsured"`
	CustomsDeclarations []CustomsDeclaration   `json:"customsDeclarations"`
	DeclaredValue       DeclaredValue          `json:"declaredValue"`
	NumberOfPieces      int                    `json:"numberOfPieces"`
	ReferenceNumbers    []string               `json:"referenceNumbers"`
	SpecialNotes        string                 `json:"specialNotes"`
	Remarks             string                 `json:"remarks"`
	DeliveryType        string                 `json:"deliveryType"`
	ContentType         string                 `json:"contentType"`
	IsCOD               bool                   `json:"isCod"`
	CODAmount           float64                `json:"codAmount"`
	Packages            []Package              `json:"packages"`
}

type WeightInfo struct {
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
	Line2       string `json:"line2"`
	City        string `json:"city"`
	State       string `json:"state"`
	CountryCode string `json:"countryCode"`
	ZipCode     string `json:"zipCode"`
}

type Dimensions struct {
	Length float64 `json:"length"`
	Height float64 `json:"height"`
	Width  float64 `json:"width"`
	Unit   string  `json:"unit"`
}

type AccountInfo struct {
	Number   string `json:"number"`
	Username string `json:"username"`
	Password string `json:"password"`
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

type Package struct {
	Width    float64 `json:"width"`
	Height   float64 `json:"height"`
	Length   float64 `json:"length"`
	Weight   float64 `json:"weight"`
	Pieces   int     `json:"pieces"`
	Value    float64 `json:"value"`
}

type ShipmentResponse struct {
	Provider    string                 `json:"provider"`
	Success     bool                   `json:"success"`
	TrackingID  string                 `json:"trackingId,omitempty"`
	AWB         string                 `json:"awb,omitempty"`
	Message     string                 `json:"message,omitempty"`
	RawResponse map[string]interface{} `json:"rawResponse,omitempty"`
}

type ShipmentRecord struct {
	ID                   string    `json:"id" db:"id"`
	Provider             string    `json:"provider" db:"provider"`
	GenericPayload       []byte    `json:"genericPayload" db:"generic_payload"`
	TransformedPayload   []byte    `json:"transformedPayload" db:"transformed_payload"`
	ProviderResponse     []byte    `json:"providerResponse" db:"provider_response"`
	Success              bool      `json:"success" db:"success"`
	CreatedAt            time.Time `json:"createdAt" db:"created_at"`
}

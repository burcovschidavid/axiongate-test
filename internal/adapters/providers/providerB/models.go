package providerB

type Request struct {
	Origin                           string                  `json:"Origin"`
	Destination                      string                  `json:"Destination"`
	ProductType                      string                  `json:"ProductType"`
	ServiceType                      string                  `json:"ServiceType"`
	CODAmount                        string                  `json:"CODAmount"`
	CODCurrency                      string                  `json:"CODCurrency"`
	SpecialInstruction               string                  `json:"SpecialInstruction"`
	Shipper                          string                  `json:"Shipper"`
	ShipperCPerson                   string                  `json:"ShipperCPErson"`
	ShipperAddress1                  string                  `json:"ShipperAddress1"`
	ShipperAddress2                  string                  `json:"ShipperAddress2"`
	ShipperCity                      string                  `json:"ShipperCity"`
	ShipperEmail                     string                  `json:"ShipperEmail"`
	ShipperPhone                     string                  `json:"ShipperPhone"`
	ShipperMobile                    string                  `json:"ShipperMobile"`
	ShipperRefNo                     string                  `json:"ShipperRefNo"`
	Consignee                        string                  `json:"Consignee"`
	ConsigneeCPerson                 string                  `json:"ConsigneeCPerson"`
	ConsigneeAddress1                string                  `json:"ConsigneeAddress1"`
	ConsigneeAddress2                string                  `json:"ConsigneeAddress2"`
	ConsigneeCity                    string                  `json:"ConsigneeCity"`
	ConsigneePhone                   string                  `json:"ConsigneePhone"`
	ConsigneeMob                     string                  `json:"ConsigneeMob"`
	ConsigneeEmail                   string                  `json:"ConsigneeEmail"`
	ConsigneeState                   string                  `json:"ConsigneeState"`
	ConsigneeZipCode                 string                  `json:"ConsigneeZipCode"`
	ConsigneeID                      string                  `json:"ConsigneeID"`
	ConsigneeIDType                  string                  `json:"ConsigneeIDType"`
	ValueOfShipment                  float64                 `json:"ValueOfShipment"`
	ValueCurrency                    string                  `json:"ValueCurrency"`
	GoodsDescription                 string                  `json:"GoodsDescription"`
	NumberofPieces                   int                     `json:"NumberofPeices"`
	Weight                           float64                 `json:"Weight"`
	PackageRequest                   []PackageRequest        `json:"PackageRequest"`
	ExportItemDeclarationRequest     []ExportItemDeclaration `json:"ExportItemDeclarationRequest"`
	UserName                         string                  `json:"UserName"`
	Password                         string                  `json:"Password"`
	AccountNo                        string                  `json:"AccountNo"`
}

type PackageRequest struct {
	DimWidth      float64 `json:"DimWidth"`
	DimHeight     float64 `json:"DimHeight"`
	DimLength     float64 `json:"DimLength"`
	DimWeight     float64 `json:"DimWeight"`
	NoofPieces    int     `json:"NoofPeices"`
	ShipmentValue float64 `json:"ShipmentValue"`
}

type ExportItemDeclaration struct {
	HSNCODE         string  `json:"HSNCODE"`
	ItemDesc        string  `json:"ItemDesc"`
	DimWeight       float64 `json:"DimWeight"`
	NoofPieces      int     `json:"NoofPeices"`
	ShipmentValue   float64 `json:"ShipmentValue"`
	CountryofOrigin string  `json:"CountryofOrigin"`
}

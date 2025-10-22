package providerB

import (
	"fmt"
	"shipping-api/internal/core/domain"
	"strings"
)

func MapToProviderB(req *domain.GenericShippingRequest) *Request {
	result := &Request{
		ProductType:        req.ProductCode,
		ServiceType:        mapServiceType(req.IsCOD),
		CODAmount:          fmt.Sprintf("%.2f", req.CODAmount),
		CODCurrency:        req.DeclaredValue.Currency,
		SpecialInstruction: req.SpecialNotes,
		Shipper:            req.Shipper.Contact.CompanyName,
		ShipperCPerson:     req.Shipper.Contact.Name,
		ShipperAddress1:    req.Shipper.Address.Line1,
		ShipperAddress2:    req.Shipper.Address.Line2,
		ShipperCity:        req.Shipper.Address.City,
		ShipperEmail:       req.Shipper.Contact.EmailAddress,
		ShipperPhone:       req.Shipper.Contact.PhoneNumber,
		ShipperMobile:      req.Shipper.Contact.MobileNumber,
		ShipperRefNo:       req.Shipper.ReferenceNo1,
		Consignee:          req.Consignee.Contact.CompanyName,
		ConsigneeCPerson:   req.Consignee.Contact.Name,
		ConsigneeAddress1:  req.Consignee.Address.Line1,
		ConsigneeAddress2:  req.Consignee.Address.Line2,
		ConsigneeCity:      req.Consignee.Address.City,
		ConsigneePhone:     req.Consignee.Contact.PhoneNumber,
		ConsigneeMob:       req.Consignee.Contact.MobileNumber,
		ConsigneeEmail:     req.Consignee.Contact.EmailAddress,
		ConsigneeState:     req.Consignee.Address.State,
		ConsigneeZipCode:   req.Consignee.Address.ZipCode,
		ConsigneeID:        "",
		ConsigneeIDType:    "",
		ValueOfShipment:    req.DeclaredValue.Amount,
		ValueCurrency:      req.DeclaredValue.Currency,
		GoodsDescription:   buildGoodsDescription(req.CustomsDeclarations),
		NumberofPieces:     req.NumberOfPieces,
		Weight:             req.Weight.Value / 1000,
		UserName:           req.Account.Username,
		Password:           req.Account.Password,
		AccountNo:          req.Account.Number,
	}

	if req.Shipper.Address.Line1 != "" {
		result.Origin = extractCityCode(req.Shipper.Address.City)
	}
	if req.Consignee.Address.Line1 != "" {
		result.Destination = extractCityCode(req.Consignee.Address.City)
	}

	if len(req.Packages) > 0 {
		result.PackageRequest = make([]PackageRequest, len(req.Packages))
		for i, pkg := range req.Packages {
			result.PackageRequest[i] = PackageRequest{
				DimWidth:      pkg.Width,
				DimHeight:     pkg.Height,
				DimLength:     pkg.Length,
				DimWeight:     pkg.Weight,
				NoofPieces:    pkg.Pieces,
				ShipmentValue: pkg.Value,
			}
		}
	}

	if len(req.CustomsDeclarations) > 0 {
		result.ExportItemDeclarationRequest = make([]ExportItemDeclaration, len(req.CustomsDeclarations))
		for i, cd := range req.CustomsDeclarations {
			result.ExportItemDeclarationRequest[i] = ExportItemDeclaration{
				HSNCODE:         cd.HSCode,
				ItemDesc:        cd.Description,
				DimWeight:       cd.Weight / 1000,
				NoofPieces:      cd.Quantity,
				ShipmentValue:   cd.Value,
				CountryofOrigin: cd.CountryOfOrigin,
			}
		}
	}

	return result
}

func mapServiceType(isCOD bool) string {
	if isCOD {
		return "COD"
	}
	return "NOR"
}

func buildGoodsDescription(declarations []domain.CustomsDeclaration) string {
	if len(declarations) == 0 {
		return ""
	}

	descriptions := make([]string, len(declarations))
	for i, d := range declarations {
		descriptions[i] = d.Description
	}
	return strings.Join(descriptions, ", ")
}

func extractCityCode(city string) string {
	if len(city) >= 3 {
		return strings.ToUpper(city[:3])
	}
	return strings.ToUpper(city)
}

package providerA

import (
	"shipping-api/internal/core/domain"
	"strconv"
)

func MapToProviderA(req *domain.GenericShippingRequest) *Request {
	result := &Request{
		Weight: Weight{
			Value: req.Weight.Value,
			Unit:  req.Weight.Unit,
		},
		Shipper: Party{
			Contact: Contact{
				Name:         req.Shipper.Contact.Name,
				MobileNumber: req.Shipper.Contact.MobileNumber,
				PhoneNumber:  req.Shipper.Contact.PhoneNumber,
				EmailAddress: req.Shipper.Contact.EmailAddress,
				CompanyName:  req.Shipper.Contact.CompanyName,
			},
			Address: Address{
				Line1:       req.Shipper.Address.Line1,
				City:        req.Shipper.Address.City,
				CountryCode: req.Shipper.Address.CountryCode,
				ZipCode:     req.Shipper.Address.ZipCode,
			},
			ReferenceNo1: req.Shipper.ReferenceNo1,
			ReferenceNo2: req.Shipper.ReferenceNo2,
		},
		Consignee: Party{
			Contact: Contact{
				Name:         req.Consignee.Contact.Name,
				MobileNumber: req.Consignee.Contact.MobileNumber,
				PhoneNumber:  req.Consignee.Contact.PhoneNumber,
				EmailAddress: req.Consignee.Contact.EmailAddress,
				CompanyName:  req.Consignee.Contact.CompanyName,
			},
			Address: Address{
				Line1:       req.Consignee.Address.Line1,
				City:        req.Consignee.Address.City,
				CountryCode: req.Consignee.Address.CountryCode,
				ZipCode:     req.Consignee.Address.ZipCode,
			},
			ReferenceNo1: req.Consignee.ReferenceNo1,
			ReferenceNo2: req.Consignee.ReferenceNo2,
		},
		Dimensions: Dimensions{
			Length: req.Dimensions.Length,
			Height: req.Dimensions.Height,
			Width:  req.Dimensions.Width,
			Unit:   req.Dimensions.Unit,
		},
		ProductCode:    req.ProductCode,
		ServiceType:    req.ServiceType,
		PrintType:      "AWBOnly",
		IsInsured:      req.IsInsured,
		DeclaredValue: DeclaredValue{
			Amount:   req.DeclaredValue.Amount,
			Currency: req.DeclaredValue.Currency,
		},
		NumberOfPieces: req.NumberOfPieces,
		SpecialNotes:   req.SpecialNotes,
		Remarks:        req.Remarks,
		DeliveryType:   req.DeliveryType,
		ContentType:    req.ContentType,
		IsCOD:          req.IsCOD,
	}

	accountNum, _ := strconv.Atoi(req.Account.Number)
	result.Account = Account{Number: accountNum}

	if len(req.ReferenceNumbers) > 0 {
		result.ReferenceNumber1 = req.ReferenceNumbers[0]
	}
	if len(req.ReferenceNumbers) > 1 {
		result.ReferenceNumber2 = req.ReferenceNumbers[1]
	}
	if len(req.ReferenceNumbers) > 2 {
		result.ReferenceNumber3 = req.ReferenceNumbers[2]
	}
	if len(req.ReferenceNumbers) > 3 {
		result.ReferenceNumber4 = req.ReferenceNumbers[3]
	}

	result.CustomsDeclarations = make([]CustomsDeclaration, len(req.CustomsDeclarations))
	for i, cd := range req.CustomsDeclarations {
		result.CustomsDeclarations[i] = CustomsDeclaration{
			Reference:       cd.Reference,
			Description:     cd.Description,
			CountryOfOrigin: cd.CountryOfOrigin,
			Weight:          cd.Weight,
			Dimensions: Dimensions{
				Length: cd.Dimensions.Length,
				Height: cd.Dimensions.Height,
				Width:  cd.Dimensions.Width,
			},
			Quantity: cd.Quantity,
			HSCode:   cd.HSCode,
			Value:    cd.Value,
		}
	}

	return result
}

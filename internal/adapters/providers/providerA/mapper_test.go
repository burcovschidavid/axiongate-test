package providerA

import (
	"shipping-api/internal/core/domain"
	"testing"
)

func TestMapToProviderA(t *testing.T) {
	genericReq := &domain.GenericShippingRequest{
		Weight: domain.WeightInfo{
			Value: 1000,
			Unit:  "Grams",
		},
		Shipper: domain.Party{
			Contact: domain.Contact{
				Name:         "ABC Associates",
				MobileNumber: "0506356566",
				PhoneNumber:  "041234567",
				EmailAddress: "orders@test.com",
				CompanyName:  "ABC Associates",
			},
			Address: domain.Address{
				Line1:       "Umm Rammool",
				City:        "Dubai",
				CountryCode: "AE",
				ZipCode:     "00000",
			},
			ReferenceNo1: "ShipperRef1",
			ReferenceNo2: "ShipperRef2",
		},
		Consignee: domain.Party{
			Contact: domain.Contact{
				Name:         "John Doe",
				MobileNumber: "+919441234567",
				PhoneNumber:  "+919441234567",
				EmailAddress: "test@test.com",
				CompanyName:  "John Doe",
			},
			Address: domain.Address{
				Line1:       "Test Address",
				City:        "Bangalore",
				CountryCode: "IN",
				ZipCode:     "1001",
			},
			ReferenceNo1: "ConsigneeRef1",
			ReferenceNo2: "ConsigneeRef2",
		},
		Dimensions: domain.Dimensions{
			Length: 10,
			Height: 10,
			Width:  10,
			Unit:   "Meter",
		},
		Account: domain.AccountInfo{
			Number: "123",
		},
		ProductCode:    "International",
		ServiceType:    "Express",
		IsInsured:      true,
		NumberOfPieces: 1,
		DeclaredValue: domain.DeclaredValue{
			Amount:   250,
			Currency: "AED",
		},
		ReferenceNumbers: []string{"Ref1", "Ref2", "Ref3", "Ref4"},
		SpecialNotes:     "Handle with care",
		Remarks:          "Fragile",
		DeliveryType:     "DoorToDoor",
		ContentType:      "NonDocument",
		IsCOD:            false,
		CustomsDeclarations: []domain.CustomsDeclaration{
			{
				Reference:       "RefNum",
				Description:     "Medicine",
				CountryOfOrigin: "AE",
				Weight:          1000,
				Dimensions: domain.Dimensions{
					Length: 10,
					Height: 10,
					Width:  10,
				},
				Quantity: 1,
				HSCode:   "3004909",
				Value:    250,
			},
		},
	}

	result := MapToProviderA(genericReq)

	if result.Weight.Value != 1000 {
		t.Errorf("expected weight value 1000, got %f", result.Weight.Value)
	}

	if result.Weight.Unit != "Grams" {
		t.Errorf("expected weight unit Grams, got %s", result.Weight.Unit)
	}

	if result.Shipper.Contact.Name != "ABC Associates" {
		t.Errorf("expected shipper name ABC Associates, got %s", result.Shipper.Contact.Name)
	}

	if result.Shipper.Address.City != "Dubai" {
		t.Errorf("expected shipper city Dubai, got %s", result.Shipper.Address.City)
	}

	if result.Consignee.Contact.Name != "John Doe" {
		t.Errorf("expected consignee name John Doe, got %s", result.Consignee.Contact.Name)
	}

	if result.Consignee.Address.City != "Bangalore" {
		t.Errorf("expected consignee city Bangalore, got %s", result.Consignee.Address.City)
	}

	if result.Dimensions.Length != 10 {
		t.Errorf("expected dimensions length 10, got %f", result.Dimensions.Length)
	}

	if result.Account.Number != 123 {
		t.Errorf("expected account number 123, got %d", result.Account.Number)
	}

	if result.ProductCode != "International" {
		t.Errorf("expected product code International, got %s", result.ProductCode)
	}

	if result.PrintType != "AWBOnly" {
		t.Errorf("expected print type AWBOnly, got %s", result.PrintType)
	}

	if !result.IsInsured {
		t.Error("expected isInsured to be true")
	}

	if result.ReferenceNumber1 != "Ref1" {
		t.Errorf("expected reference number 1 to be Ref1, got %s", result.ReferenceNumber1)
	}

	if result.ReferenceNumber2 != "Ref2" {
		t.Errorf("expected reference number 2 to be Ref2, got %s", result.ReferenceNumber2)
	}

	if result.ReferenceNumber3 != "Ref3" {
		t.Errorf("expected reference number 3 to be Ref3, got %s", result.ReferenceNumber3)
	}

	if result.ReferenceNumber4 != "Ref4" {
		t.Errorf("expected reference number 4 to be Ref4, got %s", result.ReferenceNumber4)
	}

	if len(result.CustomsDeclarations) != 1 {
		t.Errorf("expected 1 customs declaration, got %d", len(result.CustomsDeclarations))
	}

	if result.CustomsDeclarations[0].Description != "Medicine" {
		t.Errorf("expected customs description Medicine, got %s", result.CustomsDeclarations[0].Description)
	}

	if result.CustomsDeclarations[0].HSCode != "3004909" {
		t.Errorf("expected HS code 3004909, got %s", result.CustomsDeclarations[0].HSCode)
	}
}

func TestMapToProviderA_EmptyReferenceNumbers(t *testing.T) {
	genericReq := &domain.GenericShippingRequest{
		Weight: domain.WeightInfo{
			Value: 500,
			Unit:  "Grams",
		},
		Shipper: domain.Party{
			Contact: domain.Contact{
				Name: "Test",
			},
			Address: domain.Address{
				City: "City",
			},
		},
		Consignee: domain.Party{
			Contact: domain.Contact{
				Name: "Test",
			},
			Address: domain.Address{
				City: "City",
			},
		},
		Dimensions: domain.Dimensions{
			Length: 1,
			Height: 1,
			Width:  1,
		},
		Account: domain.AccountInfo{
			Number: "100",
		},
		ReferenceNumbers: []string{},
	}

	result := MapToProviderA(genericReq)

	if result.ReferenceNumber1 != "" {
		t.Errorf("expected empty reference number 1, got %s", result.ReferenceNumber1)
	}
}

func TestMapToProviderA_MultipleCustomsDeclarations(t *testing.T) {
	genericReq := &domain.GenericShippingRequest{
		Weight: domain.WeightInfo{
			Value: 2000,
			Unit:  "Grams",
		},
		Shipper: domain.Party{
			Contact: domain.Contact{Name: "Test"},
			Address: domain.Address{City: "City"},
		},
		Consignee: domain.Party{
			Contact: domain.Contact{Name: "Test"},
			Address: domain.Address{City: "City"},
		},
		Dimensions: domain.Dimensions{Length: 1, Height: 1, Width: 1},
		Account:    domain.AccountInfo{Number: "100"},
		CustomsDeclarations: []domain.CustomsDeclaration{
			{
				Description: "Item 1",
				Weight:      1000,
				Value:       100,
			},
			{
				Description: "Item 2",
				Weight:      1000,
				Value:       150,
			},
		},
	}

	result := MapToProviderA(genericReq)

	if len(result.CustomsDeclarations) != 2 {
		t.Errorf("expected 2 customs declarations, got %d", len(result.CustomsDeclarations))
	}

	if result.CustomsDeclarations[0].Description != "Item 1" {
		t.Errorf("expected first item description Item 1, got %s", result.CustomsDeclarations[0].Description)
	}

	if result.CustomsDeclarations[1].Description != "Item 2" {
		t.Errorf("expected second item description Item 2, got %s", result.CustomsDeclarations[1].Description)
	}
}

package providerB

import (
	"shipping-api/internal/core/domain"
	"testing"
)

func TestMapToProviderB(t *testing.T) {
	genericReq := &domain.GenericShippingRequest{
		Weight: domain.WeightInfo{
			Value: 1000,
			Unit:  "Grams",
		},
		Shipper: domain.Party{
			Contact: domain.Contact{
				Name:         "Sender Name",
				MobileNumber: "502009622",
				PhoneNumber:  "502009622",
				EmailAddress: "sender@email.com",
				CompanyName:  "Test Sender Company",
			},
			Address: domain.Address{
				Line1:       "Address Line 1",
				Line2:       "Address Line 2",
				City:        "Dubai",
				CountryCode: "AE",
			},
			ReferenceNo1: "4565165",
		},
		Consignee: domain.Party{
			Contact: domain.Contact{
				Name:         "Receiver Name",
				MobileNumber: "8008333",
				PhoneNumber:  "8008333",
				EmailAddress: "receiver@email.com",
				CompanyName:  "Test Receiver Company",
			},
			Address: domain.Address{
				Line1:       "Receiver Address Line 1",
				Line2:       "Receiver Address Line 2",
				City:        "AURORA",
				State:       "New York",
				CountryCode: "US",
				ZipCode:     "10001",
			},
		},
		Dimensions: domain.Dimensions{
			Length: 10,
			Height: 10,
			Width:  10,
		},
		Account: domain.AccountInfo{
			Number:   "123",
			Username: "testuser",
			Password: "testpass",
		},
		ProductCode: "XPS",
		IsCOD:       false,
		CODAmount:   0,
		DeclaredValue: domain.DeclaredValue{
			Amount:   1686,
			Currency: "USD",
		},
		NumberOfPieces: 2,
		SpecialNotes:   "Test instruction",
		CustomsDeclarations: []domain.CustomsDeclaration{
			{
				Description:     "Women Shirt",
				Weight:          100,
				Quantity:        10,
				Value:           49,
				CountryOfOrigin: "AE",
				HSCode:          "123456",
			},
			{
				Description:     "Women Pant",
				Weight:          100,
				Quantity:        5,
				Value:           25,
				CountryOfOrigin: "AE",
				HSCode:          "789012",
			},
		},
		Packages: []domain.Package{
			{
				Width:  10,
				Height: 10,
				Length: 10,
				Weight: 0.2,
				Pieces: 1,
				Value:  50,
			},
			{
				Width:  10,
				Height: 15,
				Length: 15,
				Weight: 0.45,
				Pieces: 1,
				Value:  50,
			},
		},
	}

	result := MapToProviderB(genericReq)

	if result.ProductType != "XPS" {
		t.Errorf("expected product type XPS, got %s", result.ProductType)
	}

	if result.ServiceType != "NOR" {
		t.Errorf("expected service type NOR, got %s", result.ServiceType)
	}

	if result.Shipper != "Test Sender Company" {
		t.Errorf("expected shipper Test Sender Company, got %s", result.Shipper)
	}

	if result.ShipperCPerson != "Sender Name" {
		t.Errorf("expected shipper contact Sender Name, got %s", result.ShipperCPerson)
	}

	if result.ShipperAddress1 != "Address Line 1" {
		t.Errorf("expected shipper address1 Address Line 1, got %s", result.ShipperAddress1)
	}

	if result.ShipperCity != "Dubai" {
		t.Errorf("expected shipper city Dubai, got %s", result.ShipperCity)
	}

	if result.Consignee != "Test Receiver Company" {
		t.Errorf("expected consignee Test Receiver Company, got %s", result.Consignee)
	}

	if result.ConsigneeCPerson != "Receiver Name" {
		t.Errorf("expected consignee contact Receiver Name, got %s", result.ConsigneeCPerson)
	}

	if result.ConsigneeState != "New York" {
		t.Errorf("expected consignee state New York, got %s", result.ConsigneeState)
	}

	if result.ConsigneeZipCode != "10001" {
		t.Errorf("expected consignee zip 10001, got %s", result.ConsigneeZipCode)
	}

	if result.ValueOfShipment != 1686 {
		t.Errorf("expected value of shipment 1686, got %f", result.ValueOfShipment)
	}

	if result.ValueCurrency != "USD" {
		t.Errorf("expected value currency USD, got %s", result.ValueCurrency)
	}

	if result.Weight != 1.0 {
		t.Errorf("expected weight 1.0 kg, got %f", result.Weight)
	}

	if result.NumberofPieces != 2 {
		t.Errorf("expected 2 pieces, got %d", result.NumberofPieces)
	}

	if result.UserName != "testuser" {
		t.Errorf("expected username testuser, got %s", result.UserName)
	}

	if result.Password != "testpass" {
		t.Errorf("expected password testpass, got %s", result.Password)
	}

	if result.AccountNo != "123" {
		t.Errorf("expected account number 123, got %s", result.AccountNo)
	}

	if len(result.PackageRequest) != 2 {
		t.Errorf("expected 2 packages, got %d", len(result.PackageRequest))
	}

	if result.PackageRequest[0].DimWidth != 10 {
		t.Errorf("expected first package width 10, got %f", result.PackageRequest[0].DimWidth)
	}

	if result.PackageRequest[0].DimWeight != 0.2 {
		t.Errorf("expected first package weight 0.2, got %f", result.PackageRequest[0].DimWeight)
	}

	if len(result.ExportItemDeclarationRequest) != 2 {
		t.Errorf("expected 2 export items, got %d", len(result.ExportItemDeclarationRequest))
	}

	if result.ExportItemDeclarationRequest[0].ItemDesc != "Women Shirt" {
		t.Errorf("expected first item Women Shirt, got %s", result.ExportItemDeclarationRequest[0].ItemDesc)
	}

	if result.ExportItemDeclarationRequest[0].NoofPieces != 10 {
		t.Errorf("expected first item 10 pieces, got %d", result.ExportItemDeclarationRequest[0].NoofPieces)
	}

	if result.ExportItemDeclarationRequest[1].ItemDesc != "Women Pant" {
		t.Errorf("expected second item Women Pant, got %s", result.ExportItemDeclarationRequest[1].ItemDesc)
	}
}

func TestMapToProviderB_CODService(t *testing.T) {
	genericReq := &domain.GenericShippingRequest{
		Weight: domain.WeightInfo{Value: 1000, Unit: "Grams"},
		Shipper: domain.Party{
			Contact: domain.Contact{Name: "Test"},
			Address: domain.Address{City: "Dubai"},
		},
		Consignee: domain.Party{
			Contact: domain.Contact{Name: "Test"},
			Address: domain.Address{City: "City"},
		},
		Account:   domain.AccountInfo{Number: "123"},
		IsCOD:     true,
		CODAmount: 500,
		DeclaredValue: domain.DeclaredValue{
			Amount:   500,
			Currency: "AED",
		},
	}

	result := MapToProviderB(genericReq)

	if result.ServiceType != "COD" {
		t.Errorf("expected service type COD, got %s", result.ServiceType)
	}

	if result.CODAmount != "500.00" {
		t.Errorf("expected COD amount 500.00, got %s", result.CODAmount)
	}
}

func TestMapToProviderB_GoodsDescription(t *testing.T) {
	genericReq := &domain.GenericShippingRequest{
		Weight: domain.WeightInfo{Value: 1000, Unit: "Grams"},
		Shipper: domain.Party{
			Contact: domain.Contact{Name: "Test"},
			Address: domain.Address{City: "Dubai"},
		},
		Consignee: domain.Party{
			Contact: domain.Contact{Name: "Test"},
			Address: domain.Address{City: "City"},
		},
		Account: domain.AccountInfo{Number: "123"},
		CustomsDeclarations: []domain.CustomsDeclaration{
			{Description: "Item 1"},
			{Description: "Item 2"},
			{Description: "Item 3"},
		},
	}

	result := MapToProviderB(genericReq)

	expected := "Item 1, Item 2, Item 3"
	if result.GoodsDescription != expected {
		t.Errorf("expected goods description %s, got %s", expected, result.GoodsDescription)
	}
}

func TestExtractCityCode(t *testing.T) {
	tests := []struct {
		city     string
		expected string
	}{
		{"Dubai", "DUB"},
		{"London", "LON"},
		{"NY", "NY"},
		{"A", "A"},
		{"", ""},
	}

	for _, tt := range tests {
		result := extractCityCode(tt.city)
		if result != tt.expected {
			t.Errorf("extractCityCode(%s) = %s; expected %s", tt.city, result, tt.expected)
		}
	}
}

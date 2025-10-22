package testutil

import "shipping-api/internal/core/domain"

func CreateSampleShippingRequest() *domain.GenericShippingRequest {
	return &domain.GenericShippingRequest{
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
			Number:   "123",
			Username: "testuser",
			Password: "testpass",
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
		Packages: []domain.Package{
			{
				Width:  10,
				Height: 10,
				Length: 10,
				Weight: 0.5,
				Pieces: 1,
				Value:  250,
			},
		},
	}
}

func CreateMinimalShippingRequest() *domain.GenericShippingRequest {
	return &domain.GenericShippingRequest{
		Weight: domain.WeightInfo{
			Value: 500,
			Unit:  "Grams",
		},
		Shipper: domain.Party{
			Contact: domain.Contact{
				Name:        "Sender",
				EmailAddress: "sender@test.com",
			},
			Address: domain.Address{
				City:        "Dubai",
				CountryCode: "AE",
			},
		},
		Consignee: domain.Party{
			Contact: domain.Contact{
				Name:        "Receiver",
				EmailAddress: "receiver@test.com",
			},
			Address: domain.Address{
				City:        "London",
				CountryCode: "GB",
			},
		},
		Dimensions: domain.Dimensions{
			Length: 5,
			Height: 5,
			Width:  5,
			Unit:   "CM",
		},
		Account: domain.AccountInfo{
			Number: "100",
		},
		NumberOfPieces: 1,
		DeclaredValue: domain.DeclaredValue{
			Amount:   100,
			Currency: "USD",
		},
	}
}

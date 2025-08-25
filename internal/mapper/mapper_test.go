package mapper

import (
	"demoL0/internal/domain"
	"demoL0/internal/dto"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func createSampleDto() *dto.MainOrder {
	return &dto.MainOrder{
		OrderUID:          "order123",
		TrackNumber:       "track123",
		Entry:             "entry1",
		Locale:            "en",
		InternalSignature: "sig",
		CustomerID:        "cust123",
		DeliveryService:   "deliveryX",
		Shardkey:          "shard1",
		SmID:              42,
		DateCreated:       time.Date(2023, 10, 1, 12, 0, 0, 0, time.UTC),
		OofShard:          "shardA",
		Delivery: dto.Delivery{
			Name:    "John Doe",
			Phone:   "123456789",
			Zip:     "12345",
			City:    "CityX",
			Address: "Street 123",
			Region:  "RegionY",
			Email:   "john@example.com",
		},
		Payment: dto.Payment{
			Transaction:  "trans123",
			RequestID:    "req456",
			Currency:     "USD",
			Provider:     "providerA",
			Amount:       1000,
			PaymentDt:    1633084800,
			Bank:         "BankX",
			DeliveryCost: 50,
			GoodsTotal:   950,
			CustomFee:    10,
		},
		Items: []dto.Item{
			{
				ChrtID:      1,
				TrackNumber: "tn1",
				Price:       100,
				RID:         "rid1",
				Name:        "Item1",
				Sale:        0,
				Size:        "M",
				TotalPrice:  100,
				NmID:        123,
				Brand:       "BrandA",
				Status:      1,
			},
		},
	}
}

func createSampleDomain() *domain.MainOrder {
	return &domain.MainOrder{
		OrderUID:          "order123",
		TrackNumber:       "track123",
		Entry:             "entry1",
		Locale:            "en",
		InternalSignature: "sig",
		CustomerID:        "cust123",
		DeliveryService:   "deliveryX",
		Shardkey:          "shard1",
		SmID:              42,
		DateCreated:       time.Date(2023, 10, 1, 12, 0, 0, 0, time.UTC),
		OofShard:          "shardA",
		Delivery: domain.Delivery{
			Name:    "John Doe",
			Phone:   "123456789",
			Zip:     "12345",
			City:    "CityX",
			Address: "Street 123",
			Region:  "RegionY",
			Email:   "john@example.com",
		},
		Payment: domain.Payment{
			Transaction:  "trans123",
			RequestID:    "req456",
			Currency:     "USD",
			Provider:     "providerA",
			Amount:       1000,
			PaymentDt:    1633084800,
			Bank:         "BankX",
			DeliveryCost: 50,
			GoodsTotal:   950,
			CustomFee:    10,
		},
		Items: []domain.Item{
			{
				ID:          1,
				OrderUID:    "order123",
				ChrtID:      1,
				TrackNumber: "tn1",
				Price:       100,
				RID:         "rid1",
				Name:        "Item1",
				Sale:        0,
				Size:        "M",
				TotalPrice:  100,
				NmID:        123,
				Brand:       "BrandA",
				Status:      1,
			},
		},
	}
}

func TestMapDtoToModel(t *testing.T) {
	dtoOrder := createSampleDto()

	domainOrder := MapDtoToModel(dtoOrder)

	assert.NotNil(t, domainOrder)
	assert.Equal(t, dtoOrder.OrderUID, domainOrder.OrderUID)
	assert.Equal(t, dtoOrder.Items[0].NmID, domainOrder.Items[0].NmID)
	assert.Equal(t, dtoOrder.Items[0].Name, domainOrder.Items[0].Name)
}

func TestMapModelToDto(t *testing.T) {
	domainOrder := createSampleDomain()

	dtoOrder := MapModelToDto(domainOrder)

	assert.NotNil(t, dtoOrder)
	assert.Equal(t, domainOrder.OrderUID, dtoOrder.OrderUID)
	assert.Equal(t, domainOrder.Items[0].NmID, dtoOrder.Items[0].NmID)
	assert.Equal(t, domainOrder.Items[0].Name, dtoOrder.Items[0].Name)
}

func TestMapListModelToDto(t *testing.T) {
	domainOrders := []*domain.MainOrder{
		createSampleDomain(),
		createSampleDomain(),
	}

	resultDtos := MapListModelToDto(domainOrders)

	assert.Len(t, resultDtos, 2)
	for _, orderDto := range resultDtos {
		assert.Equal(t, domainOrders[0].OrderUID, orderDto.OrderUID)
	}
}

package order

import (
	"demoL0/internal/domain"
	"reflect"
	"testing"
	"time"
)

func Test_mapDtoToModel(t *testing.T) {
	type args struct {
		dto *MainOrder
	}
	tests := []struct {
		name string
		args args
		want *domain.MainOrder
	}{
		{
			name: "success",
			args: args{
				dto: &MainOrder{
					OrderUID:    "test_order_uid",
					TrackNumber: "test_track_number",
					Entry:       "test_entry",
					Delivery: Delivery{
						Name:    "test_name",
						Phone:   "test_phone",
						Zip:     "test_zip",
						City:    "test_city",
						Address: "test_address",
						Region:  "test_region",
						Email:   "test_email",
					},
					Payment: Payment{
						Transaction:  "test_transaction",
						RequestID:    "test_request_id",
						Currency:     "test_currency",
						Provider:     "test_provider",
						Amount:       1,
						PaymentDt:    1,
						Bank:         "test_bank",
						DeliveryCost: 1,
						GoodsTotal:   1,
						CustomFee:    1,
					},
					Items: []Item{
						{
							ChrtID:      1,
							TrackNumber: "test_track_number",
							Price:       1,
							RID:         "test_rid",
							Name:        "test_name",
							Sale:        1,
							Size:        "test_size",
							TotalPrice:  1,
							NmID:        1,
							Brand:       "test_brand",
							Status:      1,
						},
					},
					Locale:            "test_locale",
					InternalSignature: "test_internal_signature",
					CustomerID:        "test_customer_id",
					DeliveryService:   "test_delivery_service",
					Shardkey:          "test_shardkey",
					SmID:              1,
					DateCreated:       time.Time{},
					OofShard:          "test_oof_shard",
				},
			},
			want: &domain.MainOrder{
				OrderUID:    "test_order_uid",
				TrackNumber: "test_track_number",
				Entry:       "test_entry",
				Delivery: domain.Delivery{
					Name:    "test_name",
					Phone:   "test_phone",
					Zip:     "test_zip",
					City:    "test_city",
					Address: "test_address",
					Region:  "test_region",
					Email:   "test_email",
				},
				Payment: domain.Payment{
					Transaction:  "test_transaction",
					RequestID:    "test_request_id",
					Currency:     "test_currency",
					Provider:     "test_provider",
					Amount:       1,
					PaymentDt:    1,
					Bank:         "test_bank",
					DeliveryCost: 1,
					GoodsTotal:   1,
					CustomFee:    1,
				},
				Items: []domain.Item{
					{
						ChrtID:      1,
						TrackNumber: "test_track_number",
						Price:       1,
						RID:         "test_rid",
						Name:        "test_name",
						Sale:        1,
						Size:        "test_size",
						TotalPrice:  1,
						NmID:        1,
						Brand:       "test_brand",
						Status:      1,
					},
				},
				Locale:            "test_locale",
				InternalSignature: "test_internal_signature",
				CustomerID:        "test_customer_id",
				DeliveryService:   "test_delivery_service",
				Shardkey:          "test_shardkey",
				SmID:              1,
				DateCreated:       time.Time{},
				OofShard:          "test_oof_shard",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mapDtoToModel(tt.args.dto); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("mapDtoToModel() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_mapModelToDto(t *testing.T) {
	type args struct {
		model *domain.MainOrder
	}
	tests := []struct {
		name string
		args args
		want *MainOrder
	}{
		{
			name: "success",
			args: args{
				model: &domain.MainOrder{
					OrderUID:    "test_order_uid",
					TrackNumber: "test_track_number",
					Entry:       "test_entry",
					Delivery: domain.Delivery{
						Name:    "test_name",
						Phone:   "test_phone",
						Zip:     "test_zip",
						City:    "test_city",
						Address: "test_address",
						Region:  "test_region",
						Email:   "test_email",
					},
					Payment: domain.Payment{
						Transaction:  "test_transaction",
						RequestID:    "test_request_id",
						Currency:     "test_currency",
						Provider:     "test_provider",
						Amount:       1,
						PaymentDt:    1,
						Bank:         "test_bank",
						DeliveryCost: 1,
						GoodsTotal:   1,
						CustomFee:    1,
					},
					Items: []domain.Item{
						{
							ChrtID:      1,
							TrackNumber: "test_track_number",
							Price:       1,
							RID:         "test_rid",
							Name:        "test_name",
							Sale:        1,
							Size:        "test_size",
							TotalPrice:  1,
							NmID:        1,
							Brand:       "test_brand",
							Status:      1,
						},
					},
					Locale:            "test_locale",
					InternalSignature: "test_internal_signature",
					CustomerID:        "test_customer_id",
					DeliveryService:   "test_delivery_service",
					Shardkey:          "test_shardkey",
					SmID:              1,
					DateCreated:       time.Time{},
					OofShard:          "test_oof_shard",
				},
			},
			want: &MainOrder{
				OrderUID:    "test_order_uid",
				TrackNumber: "test_track_number",
				Entry:       "test_entry",
				Delivery: Delivery{
					Name:    "test_name",
					Phone:   "test_phone",
					Zip:     "test_zip",
					City:    "test_city",
					Address: "test_address",
					Region:  "test_region",
					Email:   "test_email",
				},
				Payment: Payment{
					Transaction:  "test_transaction",
					RequestID:    "test_request_id",
					Currency:     "test_currency",
					Provider:     "test_provider",
					Amount:       1,
					PaymentDt:    1,
					Bank:         "test_bank",
					DeliveryCost: 1,
					GoodsTotal:   1,
					CustomFee:    1,
				},
				Items: []Item{
					{
						ChrtID:      1,
						TrackNumber: "test_track_number",
						Price:       1,
						RID:         "test_rid",
						Name:        "test_name",
						Sale:        1,
						Size:        "test_size",
						TotalPrice:  1,
						NmID:        1,
						Brand:       "test_brand",
						Status:      1,
					},
				},
				Locale:            "test_locale",
				InternalSignature: "test_internal_signature",
				CustomerID:        "test_customer_id",
				DeliveryService:   "test_delivery_service",
				Shardkey:          "test_shardkey",
				SmID:              1,
				DateCreated:       time.Time{},
				OofShard:          "test_oof_shard",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mapModelToDto(tt.args.model); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("mapModelToDto() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_mapListModelToDto(t *testing.T) {
	type args struct {
		models []*domain.MainOrder
	}
	tests := []struct {
		name string
		args args
		want []*MainOrder
	}{
		{
			name: "success",
			args: args{
				models: []*domain.MainOrder{
					{
						OrderUID:    "test_order_uid",
						TrackNumber: "test_track_number",
						Entry:       "test_entry",
						Delivery: domain.Delivery{
							Name:    "test_name",
							Phone:   "test_phone",
							Zip:     "test_zip",
							City:    "test_city",
							Address: "test_address",
							Region:  "test_region",
							Email:   "test_email",
						},
						Payment: domain.Payment{
							Transaction:  "test_transaction",
							RequestID:    "test_request_id",
							Currency:     "test_currency",
							Provider:     "test_provider",
							Amount:       1,
							PaymentDt:    1,
							Bank:         "test_bank",
							DeliveryCost: 1,
							GoodsTotal:   1,
							CustomFee:    1,
						},
						Items: []domain.Item{
							{
								ChrtID:      1,
								TrackNumber: "test_track_number",
								Price:       1,
								RID:         "test_rid",
								Name:        "test_name",
								Sale:        1,
								Size:        "test_size",
								TotalPrice:  1,
								NmID:        1,
								Brand:       "test_brand",
								Status:      1,
							},
						},
						Locale:            "test_locale",
						InternalSignature: "test_internal_signature",
						CustomerID:        "test_customer_id",
						DeliveryService:   "test_delivery_service",
						Shardkey:          "test_shardkey",
						SmID:              1,
						DateCreated:       time.Time{},
						OofShard:          "test_oof_shard",
					},
				},
			},
			want: []*MainOrder{
				{
					OrderUID:    "test_order_uid",
					TrackNumber: "test_track_number",
					Entry:       "test_entry",
					Delivery: Delivery{
						Name:    "test_name",
						Phone:   "test_phone",
						Zip:     "test_zip",
						City:    "test_city",
						Address: "test_address",
						Region:  "test_region",
						Email:   "test_email",
					},
					Payment: Payment{
						Transaction:  "test_transaction",
						RequestID:    "test_request_id",
						Currency:     "test_currency",
						Provider:     "test_provider",
						Amount:       1,
						PaymentDt:    1,
						Bank:         "test_bank",
						DeliveryCost: 1,
						GoodsTotal:   1,
						CustomFee:    1,
					},
					Items: []Item{
						{
							ChrtID:      1,
							TrackNumber: "test_track_number",
							Price:       1,
							RID:         "test_rid",
							Name:        "test_name",
							Sale:        1,
							Size:        "test_size",
							TotalPrice:  1,
							NmID:        1,
							Brand:       "test_brand",
							Status:      1,
						},
					},
					Locale:            "test_locale",
					InternalSignature: "test_internal_signature",
					CustomerID:        "test_customer_id",
					DeliveryService:   "test_delivery_service",
					Shardkey:          "test_shardkey",
					SmID:              1,
					DateCreated:       time.Time{},
					OofShard:          "test_oof_shard",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mapListModelToDto(tt.args.models); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("mapListModelToDto() = %v, want %v", got, tt.want)
			}
		})
	}
}

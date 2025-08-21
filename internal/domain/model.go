package domain

import "time"

type MainOrder struct {
	OrderUID          string `gorm:"column:order_uid;primaryKey"`
	TrackNumber       string `gorm:"column:track_number"`
	Entry             string `gorm:"column:entry"`
	DeliveryID        uint
	Delivery          Delivery `gorm:"foreignKey:DeliveryID"`
	PaymentID         uint
	Payment           Payment   `gorm:"foreignKey:PaymentID"`
	Items             []Item    `gorm:"foreignKey:OrderUID"`
	Locale            string    `gorm:"column:locale"`
	InternalSignature string    `gorm:"column:internal_signature"`
	CustomerID        string    `gorm:"column:customer_id"`
	DeliveryService   string    `gorm:"column:delivery_service"`
	Shardkey          string    `gorm:"column:shardkey"`
	SmID              int       `gorm:"column:sm_id"`
	DateCreated       time.Time `gorm:"column:date_created"`
	OofShard          string    `gorm:"column:oof_shard"`
}

type Delivery struct {
	ID      uint   `gorm:"primaryKey"`
	Name    string `gorm:"column:name"`
	Phone   string `gorm:"column:phone"`
	Zip     string `gorm:"column:zip"`
	City    string `gorm:"column:city"`
	Address string `gorm:"column:address"`
	Region  string `gorm:"column:region"`
	Email   string `gorm:"column:email"`
}

type Payment struct {
	ID           uint   `gorm:"primaryKey"`
	Transaction  string `gorm:"column:transaction"`
	RequestID    string `gorm:"column:request_id"`
	Currency     string `gorm:"column:currency"`
	Provider     string `gorm:"column:provider"`
	Amount       int    `gorm:"column:amount"`
	PaymentDt    int64  `gorm:"column:payment_dt"`
	Bank         string `gorm:"column:bank"`
	DeliveryCost int    `gorm:"column:delivery_cost"`
	GoodsTotal   int    `gorm:"column:goods_total"`
	CustomFee    int    `gorm:"column:custom_fee"`
}

type Item struct {
	ID          uint   `gorm:"primaryKey"`
	OrderUID    string `gorm:"column:order_uid"`
	ChrtID      int    `gorm:"column:chrt_id"`
	TrackNumber string `gorm:"column:track_number"`
	Price       int    `gorm:"column:price"`
	RID         string `gorm:"column:rid"`
	Name        string `gorm:"column:name"`
	Sale        int    `gorm:"column:sale"`
	Size        string `gorm:"column:size"`
	TotalPrice  int    `gorm:"column:total_price"`
	NmID        int    `gorm:"column:nm_id"`
	Brand       string `gorm:"column:brand"`
	Status      int    `gorm:"column:status"`
}

package modelsNats

import "time"

type Order struct {
	OrderUid          string    `json:"order_uid"`
	TrackNumber       string    `json:"track_number"`
	Entry             string    `json:"entry"`
	Delivery          *Delivery `json:"delivery"`
	Payment           *Payment  `json:"payment"`
	Items             *[]Item   `json:"items"`
	Locale            string    `json:"locale"`
	InternalSignature string    `json:"internal_signature"`
	CustomerId        string    `json:"customer_id"`
	DeliveryService   string    `json:"delivery_service"`
	Shardkey          string    `json:"shardkey"`
	SmId              int       `json:"sm_id"`
	DateCreated       time.Time `json:"date_created"`
	OofShard          string    `json:"oof_shard"`
}

func NewOrder(orderUid string, trackNumber string, entry string, delivery *Delivery, payment *Payment, items *[]Item, locale string, internalSignature string, customerId string, deliveryService string, shardkey string, smId int, dateCreated time.Time, oofShard string) *Order {
	return &Order{OrderUid: orderUid, TrackNumber: trackNumber, Entry: entry, Delivery: delivery, Payment: payment, Items: items, Locale: locale, InternalSignature: internalSignature, CustomerId: customerId, DeliveryService: deliveryService, Shardkey: shardkey, SmId: smId, DateCreated: time.Now(), OofShard: oofShard}
}

type Delivery struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Zip     string `json:"zip"`
	City    string `json:"city"`
	Address string `json:"address"`
	Region  string `json:"region"`
	Email   string `json:"email"`
}
type Payment struct {
	Transaction  string `json:"transaction"`
	RequestID    string `json:"_"`
	Currency     string `json:"currency"`
	Provider     string `json:"provider"`
	Amount       int    `json:"amount"`
	PaymentDT    int    `json:"payment_dt"`
	Bank         string `json:"bank"`
	DeliveryCost int    `json:"delivery_cost"`
	GoodsTotal   int    `json:"goods_total"`
	CustomFee    int    `json:"custom_fee"`
}

type Item struct {
	ChrtId      int    `json:"chrt_id"`
	TrackNumber string `json:"track_number"`
	Price       int    `json:"price"`
	Rid         string `json:"rid"`
	Name        string `json:"name"`
	Sale        int    `json:"sale"`
	Size        string `json:"size"`
	TotalPrice  int    `json:"total_price"`
	NmId        int    `json:"nm_id"`
	Brand       string `json:"brand"`
	Status      int    `json:"status"`
}

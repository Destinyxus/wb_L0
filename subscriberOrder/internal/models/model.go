package models

import (
	"time"
)

type Order struct {
	OrderUid          string    `json:"order_uid" validate:"required"`
	TrackNumber       string    `json:"track_number" validate:"required"`
	Entry             string    `json:"entry" validate:"required"`
	Delivery          *Delivery `json:"delivery" validate:"required"`
	Payment           *Payment  `json:"payment" validate:"required" `
	Items             *[]Item   `json:"items" validate:"required"`
	Locale            string    `json:"locale" validate:"required"`
	InternalSignature string    `json:"internal_signature"`
	CustomerId        string    `json:"customer_id" validate:"required"`
	DeliveryService   string    `json:"delivery_service" validate:"required"`
	Shardkey          string    `json:"shardkey" validate:"required"`
	SmId              int       `json:"sm_id" validate:"required"`
	DateCreated       time.Time `json:"date_created" validate:"required"`
	OofShard          string    `json:"oof_shard" validate:"required"`
}

var OrderID int

func NewOrder(orderUid string, trackNumber string, entry string, delivery *Delivery, payment *Payment, items *[]Item, locale string, internalSignature string, customerId string, deliveryService string, shardkey string, smId int, dateCreated time.Time, oofShard string) *Order {
	return &Order{OrderUid: orderUid, TrackNumber: trackNumber, Entry: entry, Delivery: delivery, Payment: payment, Items: items, Locale: locale, InternalSignature: internalSignature, CustomerId: customerId, DeliveryService: deliveryService, Shardkey: shardkey, SmId: smId, DateCreated: time.Now(), OofShard: oofShard}
}

type Delivery struct {
	Name    string `json:"name" validate:"required"`
	Phone   string `json:"phone" validate:"required"`
	Zip     string `json:"zip" validate:"required"`
	City    string `json:"city" validate:"required"`
	Address string `json:"address" validate:"required"`
	Region  string `json:"region" validate:"required"`
	Email   string `json:"email" validate:"required,email"`
}
type Payment struct {
	Transaction  string `json:"transaction" validate:"required"`
	RequestID    string `json:"_" `
	Currency     string `json:"currency"  validate:"required,currency"`
	Provider     string `json:"provider" validate:"required"`
	Amount       int    `json:"amount"  validate:"required,gt=0"`
	PaymentDT    int    `json:"payment_dt" validate:"required"`
	Bank         string `json:"bank" validate:"required"`
	DeliveryCost int    `json:"delivery_cost" validate:"required,gt=0"`
	GoodsTotal   int    `json:"goods_total" validate:"required,gt=0"`
	CustomFee    int    `json:"custom_fee"`
}

type Item struct {
	ChrtId      int    `json:"chrt_id" validate:"required"`
	TrackNumber string `json:"track_number" validate:"required"`
	Price       int    `json:"price"  validate:"required,gt=0"`
	Rid         string `json:"rid" validate:"required"`
	Name        string `json:"name" validate:"required"`
	Sale        int    `json:"sale" validate:"required"`
	Size        string `json:"size" validate:"required"`
	TotalPrice  int    `json:"total_price" validate:"required"`
	NmId        int    `json:"nm_id" validate:"required"`
	Brand       string `json:"brand" validate:"required"`
	Status      int    `json:"status" validate:"required"`
}

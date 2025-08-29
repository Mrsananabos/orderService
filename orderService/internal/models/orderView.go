package models

import (
	"time"
)

type DeliveryView struct {
	Name    string
	Phone   string
	Zip     string
	City    string
	Address string
	Region  string
	Email   string
}

type PaymentView struct {
	Currency     string
	Provider     string
	Amount       int
	DeliveryCost int
	GoodsTotal   int
}

type ItemView struct {
	Name       string
	TotalPrice int
	Brand      string
}

//easyjson:json
type OrderView struct {
	DeliveryService string
	DateCreated     time.Time
	Delivery        DeliveryView
	Payment         PaymentView
	Items           []ItemView
}

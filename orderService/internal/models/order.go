package models

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"time"
)

type Delivery struct {
	ID      uint   `gorm:"primaryKey"`
	Name    string `json:"name" validate:"required"`
	Phone   string `json:"phone" validate:"required_without=Email,omitempty,e164"`
	Zip     string `json:"zip" validate:"numeric"`
	City    string `json:"city" validate:"required"`
	Address string `json:"address" validate:"required"`
	Region  string `json:"region"`
	Email   string `json:"email" validate:"required_without=Phone,omitempty,email"`
}

func (d *Delivery) TableName() string {
	return "delivery"
}

type Payment struct {
	ID           uint   `gorm:"primaryKey"`
	Transaction  string `json:"transaction" validate:"alphanum"`
	RequestID    string `json:"request_id"`
	Currency     string `json:"currency" validate:"alpha"`
	Provider     string `json:"provider" validate:"required"`
	Amount       int    `json:"amount"  validate:"required"`
	PaymentDt    int    `json:"payment_dt" validate:"required"`
	Bank         string `json:"bank" validate:"required"`
	DeliveryCost int    `json:"delivery_cost" validate:"required"`
	GoodsTotal   int    `json:"goods_total" validate:"required"`
	CustomFee    int    `json:"custom_fee"`
}

func (p *Payment) TableName() string {
	return "payment"
}

type Item struct {
	Id          uint32    `gorm:"primaryKey" json:"-"`
	OrderUid    uuid.UUID `json:"_" gorm:"column:order_uid"`
	ChrtID      int       `json:"chrt_id" validate:"required"`
	TrackNumber string    `json:"track_number" validate:"alphanum"`
	Price       int       `json:"price" validate:"required"`
	RID         string    `json:"rid" gorm:"column:rid" validate:"alphanum"`
	Name        string    `json:"name" validate:"alphanum"`
	Sale        int       `json:"sale"`
	Size        string    `json:"size"`
	TotalPrice  int       `json:"total_price" validate:"required"`
	NmID        int       `json:"nm_id" validate:"required"`
	Brand       string    `json:"brand" validate:"required"`
	Status      int       `json:"status" validate:"required"`
}

func (i *Item) TableName() string {
	return "item"
}

type Order struct {
	Uid               uuid.UUID `gorm:"primaryKey" json:"order_uid" validate:"required"`
	TrackNumber       string    `json:"track_number" validate:"alphanum"`
	Entry             string    `json:"entry" validate:"required"`
	Locale            string    `json:"locale"`
	InternalSignature string    `json:"internal_signature"`
	CustomerID        string    `json:"customer_id" validate:"required"`
	DeliveryService   string    `json:"delivery_service" validate:"required"`
	ShardKey          string    `json:"shardkey" validate:"required"`
	SmID              int       `json:"sm_id" validate:"required"`
	DateCreated       time.Time `json:"date_created" validate:"required"`
	OofShard          string    `json:"oof_shard" validate:"required"`
	DeliveryID        uint      `json:"-" gorm:"column:delivery_id"`
	Delivery          Delivery  `json:"delivery"`
	PaymentID         uint      `json:"-" gorm:"column:payment_id"`
	Payment           Payment   `json:"payment"`
	Items             []Item    `gorm:"foreignKey:OrderUid;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"items"`
}

func (o *Order) TableName() string {
	return "order"
}

func (o *Order) Validate() error {
	validate := validator.New()

	if err := validate.Struct(o); err != nil {
		return err
	}

	for _, item := range o.Items {
		if err := validate.Struct(item); err != nil {
			return err
		}
	}

	return nil
}

func (o *Order) ToOrderView() OrderView {
	viewItems := make([]ItemView, 0, len(o.Items))
	for _, item := range o.Items {
		viewItems = append(viewItems, ItemView{
			Name:       item.Name,
			TotalPrice: item.TotalPrice,
			Brand:      item.Brand,
		})
	}
	return OrderView{
		DeliveryService: o.DeliveryService,
		DateCreated:     o.DateCreated,
		Delivery: DeliveryView{
			Name:    o.Delivery.Name,
			Phone:   o.Delivery.Phone,
			Zip:     o.Delivery.Zip,
			City:    o.Delivery.City,
			Address: o.Delivery.Address,
			Region:  o.Delivery.Region,
			Email:   o.Delivery.Email,
		},
		Payment: PaymentView{
			Currency:     o.Payment.Currency,
			Provider:     o.Payment.Provider,
			Amount:       o.Payment.Amount,
			DeliveryCost: o.Payment.DeliveryCost,
			GoodsTotal:   o.Payment.GoodsTotal,
		},
		Items: viewItems,
	}
}

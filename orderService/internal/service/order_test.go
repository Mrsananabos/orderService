package service

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"log"
	cache "orderService/internal/cache/mocks"
	"orderService/internal/models"
	repo "orderService/internal/repository/mocks"
	"testing"
	"time"
)

var uid uuid.UUID
var dateCreated time.Time
var validDelivery models.Delivery
var validPayment models.Payment
var validItems []models.Item
var validOrder models.Order
var orderView models.OrderView

func init() {
	var err error
	uid, err = uuid.Parse("1e9ad4fb-2615-46f9-9458-20b59253086b")
	if err != nil {
		log.Fatalf("Error while parse uuid: %v", err)
	}

	dateCreatedStr := "2021-11-26T06:22:19Z"
	dateCreated, err = time.Parse(time.RFC3339, dateCreatedStr)
	if err != nil {
		log.Fatalf("Error while parse date: %v", err)
	}

	validDelivery = models.Delivery{
		Name:    "Test Testov",
		Phone:   "+9720000000",
		Zip:     "2639809",
		City:    "Moscow",
		Address: "Prospekt Mira 15",
		Region:  "Moscow",
		Email:   "test@gmail.com",
	}
	validPayment = models.Payment{
		Transaction:  "b563feb7b2b84b6test",
		Currency:     "USD",
		Provider:     "wbpay",
		Amount:       1817,
		PaymentDt:    1637907727,
		Bank:         "alpha",
		DeliveryCost: 1500,
		GoodsTotal:   317,
	}
	validItems = []models.Item{{
		ChrtID:      9934930,
		TrackNumber: "WBILMTESTTRACK",
		Name:        "Mascaras",
		Price:       317,
		TotalPrice:  317,
		RID:         "ab4219087a764ae0btest",
		Brand:       "Vivienne Sabo",
		NmID:        2389212,
		Status:      202,
	},
	}
	validOrder = models.Order{
		Uid:             uid,
		TrackNumber:     "WBILMTESTTRACK",
		Entry:           "WBIL",
		CustomerID:      "100900",
		DeliveryService: "meest",
		ShardKey:        "2",
		SmID:            99,
		DateCreated:     dateCreated,
		OofShard:        "1",
		Delivery:        validDelivery,
		Payment:         validPayment,
		Items:           validItems,
	}
	orderView = models.OrderView{
		DeliveryService: "meest",
		DateCreated:     dateCreated,
		Delivery: models.DeliveryView{
			Name:    "Test Testov",
			Phone:   "+9720000000",
			Zip:     "2639809",
			City:    "Moscow",
			Address: "Prospekt Mira 15",
			Region:  "Moscow",
			Email:   "test@gmail.com",
		},
		Payment: models.PaymentView{
			Currency:     "USD",
			Provider:     "wbpay",
			Amount:       1817,
			DeliveryCost: 1500,
			GoodsTotal:   317,
		},
		Items: []models.ItemView{{
			Name:       "Mascaras",
			TotalPrice: 317,
			Brand:      "Vivienne Sabo",
		},
		}}
}

func TestHandler_GetById(t *testing.T) {
	t.Run("SuccessFromRepo", func(t *testing.T) {
		mockRepo := new(repo.IOrderRepository)
		mockCache := new(cache.ILruCache)

		mockCache.On("Get", uid.String()).Return(models.OrderView{}, false)
		mockRepo.On("GetByUid", uid).Return(validOrder, nil)

		service := NewService(mockRepo, mockCache)

		actualOrder, actualErr := service.GetById(uid)

		assert.Equal(t, actualOrder, orderView)
		assert.Nil(t, actualErr)
		mockCache.AssertCalled(t, "Get", uid.String())
		mockRepo.AssertCalled(t, "GetByUid", uid)
	})

	t.Run("SuccessFromCache", func(t *testing.T) {
		mockRepo := new(repo.IOrderRepository)
		mockCache := new(cache.ILruCache)

		mockCache.On("Get", uid.String()).Return(orderView, true)

		service := NewService(mockRepo, mockCache)

		actualOrder, actualErr := service.GetById(uid)

		assert.Equal(t, actualOrder, orderView)
		assert.Nil(t, actualErr)
		mockCache.AssertCalled(t, "Get", uid.String())
		mockRepo.AssertNotCalled(t, "GetByUid", uid)
	})

	t.Run("NotFoundInRepo", func(t *testing.T) {
		mockRepo := new(repo.IOrderRepository)
		mockCache := new(cache.ILruCache)

		mockCache.On("Get", uid.String()).Return(models.OrderView{}, false)
		mockRepo.On("GetByUid", uid).Return(models.Order{}, fmt.Errorf("record not found"))

		service := NewService(mockRepo, mockCache)

		actualOrder, actualErr := service.GetById(uid)

		assert.Equal(t, actualOrder, models.OrderView{})
		assert.Equal(t, actualErr.Error(), "record not found")
		mockCache.AssertCalled(t, "Get", uid.String())
		mockRepo.AssertCalled(t, "GetByUid", uid)
	})
}

func TestHandler_Create(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockRepo := new(repo.IOrderRepository)
		mockCache := new(cache.ILruCache)

		mockRepo.On("Create", validOrder).Return(nil)
		mockCache.On("Add", uid.String(), orderView).Return(true)

		service := NewService(mockRepo, mockCache)

		actualErr := service.Create(validOrder)

		assert.Nil(t, actualErr)
		mockRepo.AssertCalled(t, "Create", validOrder)
		mockCache.AssertCalled(t, "Add", uid.String(), orderView)

	})

	t.Run("FailedCreateInRepo", func(t *testing.T) {
		mockRepo := new(repo.IOrderRepository)
		mockCache := new(cache.ILruCache)

		mockRepo.On("Create", validOrder).Return(fmt.Errorf("key (%s)=(%s) already exists", "uid", uid.String()))

		service := NewService(mockRepo, mockCache)

		actualErr := service.Create(validOrder)

		assert.Equal(t, actualErr.Error(), fmt.Sprintf("key (%s)=(%s) already exists", "uid", uid.String()))
		mockRepo.AssertCalled(t, "Create", validOrder)
		mockCache.AssertNotCalled(t, "Add")
	})

	tableData := []struct {
		name     string
		order    models.Order
		errorMsg string
	}{
		{
			name: "FailedNotValidOrderDataRequireOrderUid",
			order: models.Order{
				TrackNumber:     "WBILMTESTTRACK",
				Entry:           "WBIL",
				CustomerID:      "100900",
				DeliveryService: "meest",
				ShardKey:        "2", SmID: 99,
				DateCreated: dateCreated,
				OofShard:    "1",
				Delivery:    validDelivery,
				Payment:     validPayment,
				Items:       validItems,
			},
			errorMsg: "Key: 'Order.Uid' Error:Field validation for 'Uid' failed on the 'required' tag",
		},
		{
			name: "FailedNotValidOrderDataRequireAlphanumTrackNumber",
			order: models.Order{
				Uid:             uid,
				Entry:           "WBIL",
				CustomerID:      "100900",
				DeliveryService: "meest",
				ShardKey:        "2", SmID: 99,
				DateCreated: dateCreated,
				OofShard:    "1",
				Delivery:    validDelivery,
				Payment:     validPayment,
				Items:       validItems,
			},
			errorMsg: "Key: 'Order.TrackNumber' Error:Field validation for 'TrackNumber' failed on the 'alphanum' tag",
		},
		{
			name: "FailedNotValidDeliveryDataRequireName",
			order: models.Order{
				Uid:             uid,
				TrackNumber:     "WBILMTESTTRACK",
				Entry:           "WBIL",
				CustomerID:      "100900",
				DeliveryService: "meest",
				ShardKey:        "2", SmID: 99,
				DateCreated: dateCreated,
				OofShard:    "1",
				Delivery: models.Delivery{
					//Name:    "Test Testov",
					Phone:   "+9720000000",
					Zip:     "2639809",
					City:    "Moscow",
					Address: "Prospekt Mira 15",
					Region:  "Moscow",
					Email:   "test@gmail.com",
				},
				Payment: validPayment,
				Items:   validItems,
			},
			errorMsg: "Key: 'Order.Delivery.Name' Error:Field validation for 'Name' failed on the 'required' tag",
		},
		{
			name: "FailedNotValidDeliveryDataPhoneAndEmailAreEmpty",
			order: models.Order{
				Uid:             uid,
				TrackNumber:     "WBILMTESTTRACK",
				Entry:           "WBIL",
				CustomerID:      "100900",
				DeliveryService: "meest",
				ShardKey:        "2", SmID: 99,
				DateCreated: dateCreated,
				OofShard:    "1",
				Delivery: models.Delivery{
					Name:    "Test Testov",
					Zip:     "2639809",
					City:    "Moscow",
					Address: "Prospekt Mira 15",
					Region:  "Moscow",
				},
				Payment: validPayment,
				Items:   validItems,
			},
			errorMsg: "Key: 'Order.Delivery.Phone' Error:Field validation for 'Phone' failed on the 'required_without' tag\nKey: 'Order.Delivery.Email' Error:Field validation for 'Email' failed on the 'required_without' tag",
		},
		{
			name: "FailedNotValidPaymentDataRequireAmount",
			order: models.Order{
				Uid:             uid,
				TrackNumber:     "WBILMTESTTRACK",
				Entry:           "WBIL",
				CustomerID:      "100900",
				DeliveryService: "meest",
				ShardKey:        "2", SmID: 99,
				DateCreated: dateCreated,
				OofShard:    "1",
				Delivery:    validDelivery,
				Payment: models.Payment{
					Transaction:  "b563feb7b2b84b6test",
					Currency:     "USD",
					Provider:     "wbpay",
					PaymentDt:    1637907727,
					Bank:         "alpha",
					DeliveryCost: 1500,
					GoodsTotal:   317,
				},
				Items: validItems,
			},
			errorMsg: "Key: 'Order.Payment.Amount' Error:Field validation for 'Amount' failed on the 'required' tag",
		},
		{
			name: "FailedNotValidPaymentDataRequireAlphanumTransaction",
			order: models.Order{
				Uid:             uid,
				TrackNumber:     "WBILMTESTTRACK",
				Entry:           "WBIL",
				CustomerID:      "100900",
				DeliveryService: "meest",
				ShardKey:        "2", SmID: 99,
				DateCreated: dateCreated,
				OofShard:    "1",
				Delivery:    validDelivery,
				Payment: models.Payment{
					Transaction:  "1A-2B",
					Currency:     "USD",
					Provider:     "wbpay",
					Amount:       1817,
					PaymentDt:    1637907727,
					Bank:         "alpha",
					DeliveryCost: 1500,
					GoodsTotal:   317,
				},
				Items: validItems,
			},
			errorMsg: "Key: 'Order.Payment.Transaction' Error:Field validation for 'Transaction' failed on the 'alphanum' tag",
		},
		{
			name: "FailedNotValidPaymentDataRequireAlphaCurrency",
			order: models.Order{
				Uid:             uid,
				TrackNumber:     "WBILMTESTTRACK",
				Entry:           "WBIL",
				CustomerID:      "100900",
				DeliveryService: "meest",
				ShardKey:        "2", SmID: 99,
				DateCreated: dateCreated,
				OofShard:    "1",
				Delivery:    validDelivery,
				Payment: models.Payment{
					Transaction:  "b563feb7b2b84b6test",
					Currency:     "USDv2",
					Provider:     "wbpay",
					Amount:       1817,
					PaymentDt:    1637907727,
					Bank:         "alpha",
					DeliveryCost: 1500,
					GoodsTotal:   317,
				},
				Items: validItems,
			},
			errorMsg: "Key: 'Order.Payment.Currency' Error:Field validation for 'Currency' failed on the 'alpha' tag",
		},
		{
			name: "FailedNotValidItemDataRequirePrice",
			order: models.Order{
				Uid:             uid,
				TrackNumber:     "WBILMTESTTRACK",
				Entry:           "WBIL",
				CustomerID:      "100900",
				DeliveryService: "meest",
				ShardKey:        "2", SmID: 99,
				DateCreated: dateCreated,
				OofShard:    "1",
				Delivery:    validDelivery,
				Payment:     validPayment,
				Items: []models.Item{
					{
						ChrtID:      9934930,
						TrackNumber: "WBILMTESTTRACK",
						Name:        "Mascaras",
						TotalPrice:  317,
						RID:         "ab4219087a764ae0btest",
						Brand:       "Vivienne Sabo",
						NmID:        2389212,
						Status:      202,
					},
				},
			},
			errorMsg: "Key: 'Item.Price' Error:Field validation for 'Price' failed on the 'required' tag",
		},
		{
			name: "FailedNotValidItemDataRequireAlphanumRID",
			order: models.Order{
				Uid:             uid,
				TrackNumber:     "WBILMTESTTRACK",
				Entry:           "WBIL",
				CustomerID:      "100900",
				DeliveryService: "meest",
				ShardKey:        "2", SmID: 99,
				DateCreated: dateCreated,
				OofShard:    "1",
				Delivery:    validDelivery,
				Payment:     validPayment,
				Items: []models.Item{
					{
						ChrtID:      9934930,
						TrackNumber: "WBILMTESTTRACK",
						Name:        "Mascaras",
						TotalPrice:  317,
						Price:       100,
						RID:         "a1-b2",
						Brand:       "Vivienne Sabo",
						NmID:        2389212,
						Status:      202,
					},
				},
			},
			errorMsg: "Key: 'Item.RID' Error:Field validation for 'RID' failed on the 'alphanum' tag",
		},
	}

	for _, td := range tableData {
		t.Run(td.name, func(t *testing.T) {
			mockRepo := new(repo.IOrderRepository)
			mockCache := new(cache.ILruCache)

			service := NewService(mockRepo, mockCache)

			actualErr := service.Create(td.order)

			assert.Equal(t, actualErr.Error(), td.errorMsg)
			mockRepo.AssertNotCalled(t, "Create")
			mockCache.AssertNotCalled(t, "Add")
		})
	}
}

package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"log"
	"orderService/internal/models"
)

//go:generate mockery --name=IOrderRepository --output=mocks --outpkg=mocks --case=snake --with-expecter
type IOrderRepository interface {
	GetByUid(uuid uuid.UUID) (models.Order, error)
	Create(order models.Order) error
	GetRecentOrders(limit int) ([]models.Order, error)
}

type Repository struct {
	DB *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return Repository{DB: db}
}

func (r Repository) GetByUid(uuid uuid.UUID) (models.Order, error) {
	var order models.Order
	if err := r.DB.Preload("Items").Preload("Delivery").Preload("Payment").Take(&order, "uid = ?", uuid.String()).Error; err != nil {
		log.Printf("Error fetching order: %v\n", err)
		return models.Order{}, err
	}

	return order, nil
}

func (r Repository) GetRecentOrders(limit int) ([]models.Order, error) {
	var orders []models.Order
	if err := r.DB.Preload("Items").Preload("Delivery").Preload("Payment").
		Order("date_created DESC").
		Limit(limit).
		Find(&orders).Error; err != nil {
		log.Printf("Error fetching recent orders: %v\n", err)
		return nil, err
	}

	return orders, nil
}

func (r Repository) Create(order models.Order) error {
	if err := r.DB.Create(&order).Error; err != nil {
		log.Printf("Error create order: %v\n", err)
		return err
	}
	return nil
}

package service

import (
	"github.com/google/uuid"
	"log"
	"orderService/internal/cache"
	"orderService/internal/models"
	"orderService/internal/repository"
)

//go:generate mockery --name=IOrderService --output=mocks --outpkg=mocks --case=snake --with-expecter
type IOrderService interface {
	GetById(uid uuid.UUID) (models.OrderView, error)
	Create(order models.Order) error
	HandleMessage(message []byte) error
}

type OrderService struct {
	repo  repository.IOrderRepository
	cache cache.ILruCache
}

func NewService(r repository.IOrderRepository, c cache.ILruCache) OrderService {
	return OrderService{
		repo:  r,
		cache: c,
	}
}

func (s OrderService) GetById(uid uuid.UUID) (models.OrderView, error) {
	orderInCache, ok := s.cache.Get(uid.String())
	if ok {
		log.Printf("Get from cache by key %s\n", uid.String())
		return orderInCache, nil
	}

	order, err := s.repo.GetByUid(uid)
	if err != nil {
		return models.OrderView{}, err
	}

	return order.ToOrderView(), nil
}

func (s OrderService) Create(order models.Order) error {
	if err := order.Validate(); err != nil {
		return err
	}

	if err := s.repo.Create(order); err != nil {
		return err
	}

	s.cache.Add(order.Uid.String(), order.ToOrderView())
	return nil
}

func (s OrderService) HandleMessage(message []byte) error {
	var order models.Order
	if err := order.UnmarshalJSON(message); err != nil {
		return err
	}

	return s.Create(order)
}

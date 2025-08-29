package rest

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/pressly/goose/v3"
	"log"
	"orderService/configs"
	"orderService/http/rest/handlers"
	"orderService/internal/cache"
	consumer "orderService/internal/kafka"
	"orderService/internal/repository"
	"orderService/internal/service"
	"orderService/pkg/db"
)

type Server struct {
	config   configs.Config
	gin      *gin.Engine
	consumer *consumer.Consumer
	ctx      context.Context
}

func NewServer(ctx context.Context) (*Server, error) {
	cnf, err := configs.NewParsedConfig()
	if err != nil {
		log.Fatal(err.Error())
		return nil, err
	}

	gorm, err := db.Connect(cnf.Database)
	if err != nil {
		log.Fatal(err.Error())
		return nil, err
	}

	//Накатываем миграции
	dbConnect, err := gorm.DB()
	if err != nil {
		log.Fatal(err.Error())
	}
	if err = goose.Up(dbConnect, "migrations"); err != nil {
		log.Fatal(err.Error())
	}

	repo := repository.NewRepository(gorm)
	lruCache := cache.NewCache(cnf.Cache.Size, cnf.Cache.TTL)
	lruCacheLoader := cache.NewLCacheLoader(repo, lruCache)
	orderService := service.NewService(repo, lruCache)

	//Наполнение кеша при инициализации сервера
	if err = lruCacheLoader.LoadCache(lruCache, cnf.Cache.Size); err != nil {
		log.Printf("Error initializing cache: %s\n", err.Error())
	}

	engine := gin.Default()
	handlers.Register(engine, orderService)

	consumer, err := consumer.CreateConsumer(cnf.Kafka, orderService)
	if err != nil {
		log.Fatalf("Error creating kafka consumer %s", err.Error())
	}

	return &Server{
		config:   cnf,
		gin:      engine,
		consumer: consumer,
		ctx:      ctx}, nil
}

func (s *Server) Run() error {
	go s.consumer.Start(s.ctx)
	err := s.gin.Run(s.config.Port)

	if err != nil {
		log.Fatal(err.Error())
		return err
	}

	return nil
}

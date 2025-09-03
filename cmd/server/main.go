package main

import (
	"context"
	"demoL0/internal/adapters/cache/inmemory"
	"demoL0/internal/adapters/kafka"
	"demoL0/internal/adapters/listener"
	"demoL0/internal/adapters/migrations"
	"demoL0/internal/adapters/repository/postgres"
	"demoL0/internal/domain"
	"demoL0/internal/front"
	kafkaPort "demoL0/internal/ports/kafka"
	httpTransport "demoL0/internal/transport/http"
	"demoL0/internal/usecase/order"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
	"gorm.io/gorm"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	log.Println("SERVER STARTING", time.Now().Format(time.RFC3339))

	// environment properties .env
	err := godotenv.Load()
	if err != nil {
		log.Println("Ошибка загрузки файла .env:", err)
	}

	// init database
	db := initDatabase()
	// migrations
	migrations.InitialMigration()
	migrations.RunMigrations()

	// init kafka
	kafka.InitializeTopics()
	producer, err := kafka.NewProducer()
	if err != nil {
		log.Fatalf("Failed to create Kafka producer: %v", err)
	}
	defer func(producer kafkaPort.Producer) {
		err2 := producer.Close()
		if err2 != nil {
			log.Println(err2)
		}
	}(producer)

	// init dependencies
	orderRepo := postgres.NewOrderRepo(db)
	cache := inmemory.NewCache[string, *domain.MainOrder](5*time.Minute, 10*time.Minute)
	cachedRepo := postgres.NewCachedOrderRepository(orderRepo, cache)
	orderService := order.NewService(db, cachedRepo)
	handler := httpTransport.NewHandler(orderService)
	validator := kafka.NewValidator()

	// init kafka listener
	kafkaListener := listener.NewListener(ctx, producer, orderService, validator)

	// start kafka listeners
	wg := &sync.WaitGroup{}
	kafkaListener.Start(wg)

	// init & run server
	router := handler.InitRouter()
	handlers := front.CreatePageHandlers()
	for key, hand := range handlers.Handlers {
		router.HandleFunc(key, hand)
	}

	srv := &http.Server{
		Addr:    ":8081",
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("ListenAndServe error: %v", err)
		}
	}()

	// graceful shutdown
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	<-sigs
	log.Println("Stop signal received")
	cancel()

	// stop kafka listeners
	kafkaListener.Stop()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("server shutdown error: %v", err)
	}

	wg.Wait()
	log.Println("Exit")
}

// initDatabase инициализирует подключение к базе данных
func initDatabase() *gorm.DB {
	return postgres.GetGormInstance()
}

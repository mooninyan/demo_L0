package main

import (
	"context"
	"demoL0/internal/domain"
	"demoL0/internal/front"
	"demoL0/internal/listener"
	"demoL0/internal/mapcache"
	"demoL0/internal/repository"
	"demoL0/internal/utils"
	"demoL0/pkg/database"
	"demoL0/pkg/kfk"
	"errors"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"demoL0/internal/service"
	"demoL0/internal/transport/rest"
	_ "github.com/lib/pq"
	_ "gorm.io/gorm"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	log.Println("SERVER STARTING", time.Now().Format(time.RFC3339))
	// environment properties .env
	err := godotenv.Load()
	if err != nil {
		log.Println("Ошибка загрузки файла .env:", err)
	}
	// migrations
	database.InitialMigration()
	database.RunMigrations()
	// init orm
	gorm := database.GetGormInstance()
	// kafka
	kfk.InitializeTopics()
	kafkaWriter := kfk.CreateKafkaWriter()
	defer kafkaWriter.Close()
	// init deps
	producer := service.NewProducer(kafkaWriter)
	repo := repository.CreateOrderRepo(gorm)
	cache := mapcache.NewCache[string, *domain.MainOrder]()
	cachedRepo := repository.NewCachedRepository(&repo, cache)
	orderService := service.CreateOrderService(gorm, &cachedRepo)
	handler := rest.NewHandler(orderService)
	listenDlq := utils.GetDlqFlagOrDefault()
	consumerGroup := kfk.CreateConsumerGroup()
	liHandler := listener.NewConsumerGroupHandler(false, producer, orderService)

	ghMap := make(map[string]*listener.GroupAndHandler)
	ghMap["main_group"] = &listener.GroupAndHandler{
		Group:   consumerGroup,
		Handler: liHandler,
	}
	if listenDlq {
		var dlqConsumerGroup = kfk.CreateDlqConsumerGroup()
		dlqHandler := listener.NewConsumerGroupHandler(true, producer, orderService)
		ghMap["dlq_group"] = &listener.GroupAndHandler{
			Group:   dlqConsumerGroup,
			Handler: dlqHandler,
		}
	}

	li := listener.NewListener(ctx, ghMap)

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go li.Listen(wg)
	if listenDlq {
		wg.Add(1)
		go li.ListenDlq(wg)
	}
	// close kafka
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

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

	<-sigs
	log.Println("Stop signal received")
	cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("server shutdown error: %v", err)
	}
	wg.Wait()

	log.Println("Exit")
}

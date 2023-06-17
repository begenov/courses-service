package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	delivery "github.com/begenov/courses-service/internal/delivery/http"

	"github.com/begenov/courses-service/internal/config"
	"github.com/begenov/courses-service/internal/repository"
	"github.com/begenov/courses-service/internal/server"
	"github.com/begenov/courses-service/internal/service"
	"github.com/begenov/courses-service/pkg/cache"
	"github.com/begenov/courses-service/pkg/database"
	"github.com/begenov/courses-service/pkg/kafka"
)

const (
	path = "./.env"
)

// @title Courses  Service API
// @version 1.0
// @description API Server for Courses Application

// @host localhost:8000
// @BasePath /api/v1/
func main() {
	cfg, err := config.Init(path)
	if err != nil {
		log.Fatal(err)
	}

	db, err := database.OpenDB(cfg.Database.Driver, cfg.Database.DSN)
	if err != nil {
		log.Fatal(err)
	}

	memCache, err := cache.NewMemoryCache(context.Background(), cfg.Redis)
	if err != nil {
		log.Fatalf("error mem cache init: %v", err)
	}

	producer, err := kafka.NewProducer(cfg.Kafka.Brokers)
	if err != nil {
		log.Fatalf("error creating Kafka producer: %v", err)
	}

	consumer, err := kafka.NewConsumer(cfg.Kafka.Brokers)
	if err != nil {
		log.Fatalf("error creating Kafka consumer: %v", err)
	}

	repos := repository.NewRepository(db)

	service := service.NewService(repos, memCache, cfg.Redis.Ttl, producer, consumer)
	go service.Kafka.Read(context.Background())

	handler := delivery.NewHandler(service)

	srv := server.NewServer(cfg, handler.Init())
	go func() {
		if err := srv.Run(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("error occurred while running http server: %s\n", err.Error())
		}
	}()

	log.Println("server started ", cfg.Server.Port)

	quit := make(chan os.Signal, 1)

	<-quit

	const timeout = 5 * time.Second

	ctx, shutdown := context.WithTimeout(context.Background(), timeout)
	defer shutdown()

	if err := srv.Stop(ctx); err != nil {
		log.Fatalf("failed to stop server %v", err)
	}
}

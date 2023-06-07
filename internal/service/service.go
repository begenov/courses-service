package service

import (
	"context"
	"time"

	"github.com/begenov/courses-service/internal/domain"
	"github.com/begenov/courses-service/internal/repository"
	"github.com/begenov/courses-service/pkg/cache"
	"github.com/begenov/courses-service/pkg/kafka"
	_ "github.com/golang/mock/mockgen/model"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go
type Courses interface {
	Create(ctx context.Context, course domain.Courses) error
	GetByID(ctx context.Context, id int) (domain.Courses, error)
	Update(ctx context.Context, course domain.Courses) error
	Delete(ctx context.Context, id int) error
	GetCoursesByIdStudent(ctx context.Context, studentId string) ([]domain.Courses, error)
}

type Kafka interface {
	SendMessages(topic string, message string) error
	ConsumeMessages(topic string, handler func(message string)) error
	// Read(ctx context.Context)
	Close()
}

type Service struct {
	Courses Courses
	Kafka   Kafka
}

func NewService(repo *repository.Repository, cache cache.Cache, ttl time.Duration, producer *kafka.Producer, concumer *kafka.Consumer) *Service {
	return &Service{
		Courses: NewCoursesService(repo.Courses, cache, ttl),
		Kafka:   NewKafkaSerivce(producer, concumer, repo.Courses),
	}
}

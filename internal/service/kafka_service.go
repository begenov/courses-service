package service

import (
	"log"

	"github.com/begenov/courses-service/internal/repository"
	"github.com/begenov/courses-service/pkg/kafka"
)

type KafkaService struct {
	producer   *kafka.Producer
	consumer   *kafka.Consumer
	repo       repository.Courses
	ResponseCh chan []byte
}

func NewKafkaSerivce(producer *kafka.Producer, consumer *kafka.Consumer, repo repository.Courses) *KafkaService {
	return &KafkaService{
		producer:   producer,
		consumer:   consumer,
		repo:       repo,
		ResponseCh: make(chan []byte),
	}
}

func (s *KafkaService) SendMessages(topic string, message string) error {
	if err := s.producer.SendMessage(topic, message); err != nil {
		return err
	}

	log.Println("Message sent to Kafka:", message)
	return nil

}

func (s *KafkaService) ConsumeMessages(topic string, handler func(message string)) error {
	err := s.consumer.ConsumeMessages(topic, handler)
	if err != nil {
		log.Println("Failed to consume messages from Kafka:", err)
		return err
	}

	return nil
}

func (s *KafkaService) Close() {
	_ = s.consumer.Close()
	_ = s.producer.Close()
}

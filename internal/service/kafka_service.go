package service

import (
	"context"
	"log"

	"github.com/Shopify/sarama"
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

func (s *KafkaService) Read(ctx context.Context) {
	partitions, err := s.consumer.Consumer.Partitions("courses-request")
	if err != nil {
		log.Fatalln("Failed to get partitions:", err)
	}

	for _, partition := range partitions {
		pc, err := s.consumer.Consumer.ConsumePartition("courses-request", partition, sarama.OffsetNewest)
		if err != nil {
			log.Fatalln("Failed to start consumer for partition", partition, ":", err)
		}
		go func(pc sarama.PartitionConsumer) {
			defer pc.Close()

			for message := range pc.Messages() {
				id := ""
				for _, v := range message.Value {
					if v == '"' {
						continue
					}
					id += string(v)
				}
				courses, err := s.repo.GetCoursesByIdStudent(ctx, id)
				if err != nil {
					log.Println(err)
					return
				}
				if err := s.producer.SendMessage("courses-response", courses); err != nil {
					log.Println(err, "send message")
					return
				}

			}
		}(pc)
	}
	<-ctx.Done()

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

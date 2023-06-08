package service

import (
	"bytes"
	"context"
	"log"
	"sync"

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

	var wg sync.WaitGroup
	wg.Add(len(partitions))

	for _, partition := range partitions {
		pc, err := s.consumer.Consumer.ConsumePartition("courses-request", partition, sarama.OffsetNewest)
		if err != nil {
			log.Fatalln("Failed to start consumer for partition", partition, ":", err)
		}

		go func(pc sarama.PartitionConsumer) {
			defer func() {
				pc.Close()
				wg.Done()
			}()

			for message := range pc.Messages() {
				res := getStringWithoutQuotes(message.Value)

				student, err := s.repo.GetCoursesByIdStudent(ctx, res)
				if err != nil {
					log.Println(err)
					return
				}

				if err := s.producer.SendMessage("courses-response", student); err != nil {
					log.Println(err, "send message")
					return
				}
			}
		}(pc)
	}

	wg.Wait()
}

func (s *KafkaService) ConsumeMessages(topic string, handler func(message string)) error {
	return s.consumer.ConsumeMessages(topic, handler)
}

func (s *KafkaService) Close() {
	_ = s.consumer.Close()
	_ = s.producer.Close()
}

func getStringWithoutQuotes(input []byte) string {
	var buffer bytes.Buffer

	for _, v := range input {
		if v == '"' {
			continue
		}
		buffer.WriteByte(v)
	}

	return buffer.String()
}

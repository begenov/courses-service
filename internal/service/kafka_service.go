package service

import (
	"log"

	"github.com/begenov/courses-service/internal/repository"
	"github.com/begenov/courses-service/pkg/kafka"
)

type KafkaService struct {
	producer *kafka.Producer
	consumer *kafka.Consumer
	repo     repository.Courses
}

func NewKafkaSerivce(producer *kafka.Producer, consumer *kafka.Consumer, repo repository.Courses) *KafkaService {
	return &KafkaService{
		producer: producer,
		consumer: consumer,
	}
}

func (s *KafkaService) SendMessages(topic string, message string) error {
	if err := s.producer.SendMessage(topic, message); err != nil {
		return err
	}

	log.Println("Message sent to Kafka:", message)
	return nil

}

/*
	func (s *KafkaService) Read(ctx context.Context) {
		for {
			responseHandler := func(message string) {
				// Добавьте здесь логику обработки ответа от Kafka

				students, err := s.repo.GetCoursesByIdStudent(ctx, message)
				if err != nil {

					log.Println("Failed to send students to course:", err)
					return
				}

				var buf bytes.Buffer
				encoder := gob.NewEncoder(&buf)
				err = encoder.Encode(students)
				if err != nil {
					log.Println("Failed to serialize student:", err)
					return
				}

				m := buf.Bytes()

				if err := s.producer.SendMessage("courses", string(m)); err != nil {
					log.Println("Failed to send message:", err)
					return

				}
			}

			// Потребляем сообщения из Kafka
			err := s.consumer.ConsumeMessages("students", responseHandler)
			if err != nil {
				log.Println("Failed to consume messages from Kafka:", err)
				return
			}
		}
	}
*/
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

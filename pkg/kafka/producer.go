package kafka

import (
	"fmt"
	"log"
	"time"

	"github.com/Shopify/sarama"
)

type Producer struct {
	producer sarama.AsyncProducer
}

func NewProducer(brokers []string) (*Producer, error) {
	config := sarama.NewConfig()

	config.Producer.RequiredAcks = sarama.WaitForLocal       // Only wait for the leader to ack
	config.Producer.Compression = sarama.CompressionSnappy   // Compress messages
	config.Producer.Flush.Frequency = 500 * time.Millisecond // Flush batches every 500ms
	fmt.Println(brokers)
	producer, err := sarama.NewAsyncProducer(brokers, config)
	if err != nil {
		log.Fatalln("Failed to start Sarama producer:", err)
	}

	// We will just log to STDOUT if we're not able to produce messages.
	// Note: messages will only be returned here after all retry attempts are exhausted.
	go func() {
		for err := range producer.Errors() {
			log.Println("Failed to write access log entry:", err)
		}
	}()

	return &Producer{
		producer: producer,
	}, nil
}

func (p *Producer) SendMessage(topic string, data string) error {

	// jsonData, err := json.Marshal(data)
	// if err != nil {
	// 	log.Fatalln("error json marsal", err)
	// 	return err
	// }

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(data),
	}

	p.producer.Input() <- msg
	return nil
}

func (p *Producer) Close() error {
	if p.producer != nil {
		return p.producer.Close()
	}
	return nil
}

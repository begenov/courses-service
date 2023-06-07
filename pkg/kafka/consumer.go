package kafka

import (
	"log"

	"github.com/Shopify/sarama"
)

type Consumer struct {
	consumer sarama.Consumer
	done     chan struct{} // Channel to control the consumption process
}

func NewConsumer(brokers []string) (*Consumer, error) {
	config := sarama.NewConfig()

	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		consumer: consumer,
		done:     make(chan struct{}),
	}, nil
}

func (c *Consumer) ConsumeMessages(topic string, handler func(message string)) error {
	partitions, err := c.consumer.Partitions(topic)
	if err != nil {
		log.Println("Failed to retrieve partitions:", err)
		return err
	}

	if len(partitions) == 0 {
		log.Println("No partitions found for the topic:", topic)
		return nil
	}

	for _, partition := range partitions {
		pc, err := c.consumer.ConsumePartition(topic, partition, sarama.OffsetNewest)
		if err != nil {
			log.Println("Failed to start consumer for partition", partition, ":", err)
			return err
		}
		go func(pc sarama.PartitionConsumer) {
			defer pc.Close()

			for message := range pc.Messages() {
				// log.Println(string(message.Value))
				// Обработка прочитанного сообщения
				handler(string(message.Value))
			}
		}(pc)
	}

	return nil
}

func (c *Consumer) Stop() {
	close(c.done)
}

func (c *Consumer) Close() error {
	c.Stop() // Stop the consumption process before closing the consumer

	if c.consumer != nil {
		return c.consumer.Close()
	}
	return nil
}

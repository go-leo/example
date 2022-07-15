package main

import (
	"context"
	"math/rand"
	"time"

	"github.com/Shopify/sarama"
	"github.com/ThreeDotsLabs/watermill-kafka/v2/pkg/kafka"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/go-leo/leo"
	"github.com/go-leo/leo/global"
	"github.com/go-leo/leo/log"
	"github.com/go-leo/leo/log/zap"
	"github.com/go-leo/leo/runner/task/pubsub"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	ctx := context.Background()
	global.SetLogger(zap.New(zap.LevelAdapt(log.Debug), zap.Console(true), zap.JSON()))
	app := leo.NewApp(
		leo.Name("redissubdemo"),
		leo.Logger(global.Logger()),
		leo.PubSub(&leo.PubSubOptions{
			Jobs: []*pubsub.Job{NewJob()},
		}),
	)
	// 运行app
	if err := app.Run(ctx); err != nil {
		panic(err)
	}
}

func NewJob() *pubsub.Job {
	saramaSubscriberConfig := kafka.DefaultSaramaSubscriberConfig()
	// equivalent of auto.offset.reset: earliest
	saramaSubscriberConfig.Consumer.Offsets.Initial = sarama.OffsetOldest
	subscriber, err := kafka.NewSubscriber(
		kafka.SubscriberConfig{
			Brokers:               []string{":9092"},
			Unmarshaler:           kafka.DefaultMarshaler{},
			OverwriteSaramaConfig: saramaSubscriberConfig,
			ConsumerGroup:         "test_consumer_group",
		},
		pubsub.NewLogger(global.Logger()),
	)
	if err != nil {
		panic(err)
	}
	publisher, err := kafka.NewPublisher(
		kafka.PublisherConfig{
			Brokers:   []string{":9092"},
			Marshaler: kafka.DefaultMarshaler{},
		},
		pubsub.NewLogger(global.Logger()),
	)
	if err != nil {
		panic(err)
	}
	rand.Seed(time.Now().UnixNano())
	return pubsub.NewPubSubJob(
		"kafkademo",
		"hello",
		subscriber,
		"hello-tmp",
		publisher,
		func(msg *message.Message) ([]*message.Message, error) {
			if rand.Int()%2 == 0 {
				global.Logger().Info("ack")
				msg.Ack()
				return []*message.Message{msg}, nil
			} else {
				global.Logger().Info("nack")
				msg.Nack()
				return nil, nil
			}
		},
	)
}

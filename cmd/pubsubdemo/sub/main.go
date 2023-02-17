package main

import (
	"context"
	"math/rand"
	"time"

	"github.com/Shopify/sarama"
	"github.com/ThreeDotsLabs/watermill-kafka/v2/pkg/kafka"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/go-leo/leo/v2"
	"github.com/go-leo/leo/v2/global"
	"github.com/go-leo/leo/v2/log"
	"github.com/go-leo/leo/v2/log/zap"
	"github.com/go-leo/pubsub"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	ctx := context.Background()
	global.SetLogger(zap.New(zap.LevelAdapt(log.Debug), zap.Console(true), zap.JSON()))
	app := leo.NewApp(
		leo.Name("sub"),
		leo.Logger(global.Logger()),
		leo.Runnable(pubsub.New([]*pubsub.Job{NewJob()})),
	)
	// 运行app
	if err := app.Run(ctx); err != nil {
		panic(err)
	}
}

func NewJob() *pubsub.Job {
	saramaSubscriberConfig := kafka.DefaultSaramaSubscriberConfig()
	saramaSubscriberConfig.Consumer.Offsets.Initial = sarama.OffsetNewest
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
	rand.Seed(time.Now().UnixNano())
	return pubsub.NewSubJob(
		"sub",
		"sub",
		subscriber,
		func(msg *message.Message) error {
			if rand.Int()%2 == 0 {
				global.Logger().Info("ack")
				msg.Ack()
			} else {
				global.Logger().Info("nack")
				msg.Nack()
			}
			return nil
		},
	)
}

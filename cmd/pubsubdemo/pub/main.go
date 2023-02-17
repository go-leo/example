package main

import (
	"context"
	"time"

	"github.com/ThreeDotsLabs/watermill-kafka/v2/pkg/kafka"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/go-leo/leo/v2/global"
	"github.com/go-leo/leo/v2/log"
	"github.com/go-leo/leo/v2/log/zap"
	"github.com/go-leo/pubsub"
)

func main() {
	global.SetLogger(zap.New(zap.LevelAdapt(log.Debug), zap.Console(true), zap.JSON()))

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

	newMessage := message.NewMessage("", message.Payload("this from pub"))
	newMessage.SetContext(context.Background())
	for {
		err = publisher.Publish("stream", newMessage)
		if err != nil {
			panic(err)
		}
		<-time.After(time.Second)
	}
}

package main

import (
	"context"
	"math/rand"
	"time"

	"github.com/ThreeDotsLabs/watermill-kafka/v2/pkg/kafka"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/go-leo/leo/global"
	"github.com/go-leo/leo/log"
	"github.com/go-leo/leo/log/zap"
	"github.com/go-leo/leo/runner/task/pubsub"
)

func main() {
	rand.Seed(time.Now().UnixNano())
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
	newMessage := message.NewMessage("", message.Payload("world"))
	newMessage.SetContext(context.Background())
	for {
		err = publisher.Publish("hello", newMessage)
		if err != nil {
			panic(err)
		}
		<-time.After(time.Second)
	}

}

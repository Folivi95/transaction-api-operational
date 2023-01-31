package kafka

import (
	"context"
	"fmt"
	"github.com/saltpay/go-kafka-driver"
	"github.com/saltpay/transaction-api-operational/src-mario-api/internal/adapters/logger/zapctx"
	"github.com/saltpay/transaction-api-operational/src-mario-api/internal/configuration"
	"strings"
)

type MessageListener interface {
	Listen(context.Context, context.CancelFunc, func(context.Context, kafka.Message)) error
}

type TransactionListener struct {
	KafkaConfig configuration.KafkaConfig
}

func NewTransactionListener(config configuration.KafkaConfig) TransactionListener {
	return TransactionListener{
		KafkaConfig: config,
	}
}

func (t *TransactionListener) Listen(ctx context.Context, callback func(ctx context.Context, consumer *kafka.Consumer, message kafka.Message) error) error {
	zapctx.From(ctx).Info("[TransactionListener] Starting to listen for messages")

	consumer, err := kafka.NewConsumer(ctx, kafka.ConsumerConfig{
		Brokers:  strings.Split(t.KafkaConfig.KafkaEndpoint, ","),
		GroupID:  fmt.Sprintf("%s-groupID", t.KafkaConfig.KafkaTransactionTopic),
		Topic:    t.KafkaConfig.KafkaTransactionTopic,
		Username: t.KafkaConfig.KafkaUsername,
		Password: t.KafkaConfig.KafkaPassword,
	})

	if err != nil {
		return err
	}
	go consumer.Listen(ctx, func(ctx context.Context, message kafka.Message) error {
		zapctx.From(ctx).Info("[TransactionListener] Processing message")

		err = callback(ctx, consumer, message)
		if err != nil {
			zapctx.From(ctx).Error(fmt.Sprintf("[TransactionListener] Unable to start listening: %s", err))
			return err
		}
		return nil
	}, kafka.NeverCommit, func(ctx context.Context) bool {
		return false
	})

	return nil
}

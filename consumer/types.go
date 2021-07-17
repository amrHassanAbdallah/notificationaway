package consumer

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
)

type MessageBrokerReadInterface interface {
	ReadAndAutoCommit(ctx context.Context) (*kafka.Message, error)
	FetchAndIgnoreOffsetCommit(ctx context.Context) (*kafka.Message, error)
	CommitOffset(ctx context.Context, msgs []kafka.Message) error
	Close() error
}

type ErrUnrecoverableMsgHandling struct {
	Message string
}

func (e *ErrUnrecoverableMsgHandling) Error() string {
	return "No retry for this err, because:" + e.Message
}

type ErrNotSupportedNotification struct {
	Subject string
	Message string
}

func (e *ErrNotSupportedNotification) Error() string {
	return fmt.Sprintf("not supporterd subject (%v) to match a message with language value (%v)", e.Subject, e.Message)
}

type ErrContextTimeOut struct {
}

func (e *ErrContextTimeOut) Error() string {
	return "context deadline exceeded"
}

type EventMessage struct {
}

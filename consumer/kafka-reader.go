package consumer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/amrHassanAbdallah/notificationaway/logger"
	"github.com/segmentio/kafka-go"
)

type ErrLoggerWrapper struct {
	Logger logger.Logger
}

func (logWrap *ErrLoggerWrapper) Printf(tmp string, args ...interface{}) {
	logWrap.Logger.Warnw(fmt.Sprintf(tmp, args...))
}

type DebLoggerWrapper struct {
	Logger logger.Logger
}

func (logWrap *DebLoggerWrapper) Printf(tmp string, args ...interface{}) {
	logWrap.Logger.Debugw(fmt.Sprintf(tmp, args...))
}

type KafkaHandler struct {
	Reader *kafka.Reader
}

func (k *KafkaHandler) ReadAndAutoCommit(ctx context.Context) (*kafka.Message, error) {
	m, err := k.Reader.ReadMessage(ctx)
	if err != nil {
		return nil, err
	}
	return &m, nil
}
func (k *KafkaHandler) FetchAndIgnoreOffsetCommit(ctx context.Context) (*kafka.Message, error) {
	m, err := k.Reader.FetchMessage(ctx)
	if err != nil {
		return nil, err
	}
	return &m, nil
}
func (k *KafkaHandler) CommitOffset(ctx context.Context, msgs []kafka.Message) error {
	return k.Reader.CommitMessages(ctx, msgs...)
}

func (k *KafkaHandler) Close() error {
	return k.Reader.Close()
}

type WorkerConfig struct {
	MessageBrokerReader MessageBrokerReadInterface
	EventMessageHandler *NotificationEventHandler
	RequestTimeOut      time.Duration
	MaxRetriesOnFailure int
}

func HandleKafkaMessages(ctx context.Context, s *sync.WaitGroup, config WorkerConfig) {
	m := new(kafka.Message)
	var err error
	currRetries := 0
	isMsgHandled := false
outer:
	for {
		select {
		case <-ctx.Done():
			err := config.MessageBrokerReader.Close()
			if err != nil {
				logger.Errorw("failed to close connection with the message broker", "err", err)
			}
			s.Done()
			break outer
		default:
			//to retry on the same message whenever there is a failure in the msg handling
			if currRetries >= config.MaxRetriesOnFailure && err != nil && m != nil {
				logger.Fatalw("failed to process message, due to a failure", "msg-topic", m.Topic, "msg-partition", m.Partition, "msg-offset", m.Offset, "err", err)
			}
			timeN := time.Now()
			timeN2 := time.Now()
			if m == nil || m != nil && len(m.Value) == 0 {
				kctx, cancel := context.WithTimeout(ctx, config.RequestTimeOut)
				m, err = config.MessageBrokerReader.FetchAndIgnoreOffsetCommit(kctx)
				cancel()

				if err != nil {
					if kctx.Err() == nil {
						logger.Errorw("failed to read message", "err", err)
					}
					break
				}
				logger.Debugw("time it took to get message from kafka", "time-it-took-in-s", time.Now().Sub(timeN2).Seconds(), "msg-topic", m.Topic, "msg-partition", m.Partition, "msg-offset", m.Offset)
			} else {
				currRetries++
				logger.Warnw("failed to handle a message", "err", err, "retry", currRetries+1, "max-retries", config.MaxRetriesOnFailure)
			}

			if m != nil && len(m.Value) >= 0 {
				//to not handle message twice when kafka commit fail
				if !isMsgHandled {
					timeN2 = time.Now()
					formattedMessage := map[string]interface{}{}
					decoder := json.NewDecoder(bytes.NewReader(m.Value))
					if err = decoder.Decode(&formattedMessage); err != nil {
						logger.Warnw("failed to cast message", "msg-topic", m.Topic, "msg-partition", m.Partition, "msg-offset", m.Offset, "err", err)
						err = nil
						continue
					}
					for _, msgHeader := range m.Headers {
						formattedMessage[string(msgHeader.Key)] = string(msgHeader.Value)
					}

					timeN2 = time.Now()
					err = config.EventMessageHandler.Handle(ctx, formattedMessage)
					if err != nil {
						switch err.(type) {
						case *ErrUnrecoverableMsgHandling:
							logger.Warnw("skipping failure in handling the message, because :"+err.Error(), "msg-topic", m.Topic, "msg-partition", m.Partition, "msg-offset", m.Offset)
							kctx, cancel := context.WithTimeout(ctx, 10*time.Second)
							err = config.MessageBrokerReader.CommitOffset(kctx, []kafka.Message{*m})
							cancel()
							if err != nil {
								err = fmt.Errorf("failed to commit message, err:" + err.Error())
								continue
							}
							err = nil
							m = nil
							continue
						default:
							//other than the above err type, it's K to retry again.
							logger.Errorf("failed to handle the message", "err", err.Error())
							continue
						}
					}

					isMsgHandled = true
					logger.Debugw("time it took to handle the event message", "time-it-took-in-s", time.Now().Sub(timeN2).Seconds())
				}

				timeN2 = time.Now()
				kctx, cancel := context.WithTimeout(ctx, 10*time.Second)
				err = config.MessageBrokerReader.CommitOffset(kctx, []kafka.Message{*m})
				cancel()
				if err != nil {
					err = fmt.Errorf("failed to commit message, err:" + err.Error())
					continue
				}
				currRetries = 0
				isMsgHandled = false
				m = nil
				logger.Debugw("time it took to commit the kafka offset", "time-it-took-in-s", time.Now().Sub(timeN2).Seconds())
			}
			logger.Debugw("finished processing", "is-job-succeeded", err == nil, "total-time-it-took-in-s", time.Now().Sub(timeN).Seconds())
		}
	}

}

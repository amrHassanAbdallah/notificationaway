package persistence

import (
	"context"
	"github.com/amrHassanAbdallah/notificationaway/logger"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
)

func (m *MongoAdapter) GetMessage(ctx context.Context, Id string) (*Message, error) {
	message := new(Message)

	query, _ := constructQueryNotifiers(QueryMessagesOptions{
		Id: Id,
	})

	err := m.MessagesCollection.FindOne(ctx, query).Decode(message)
	if err != nil {
		if errors.As(err, &mongo.ErrNoDocuments) {
			err = &ErrRecordNotFound{entity: "message", Id: Id}
		} else {
			logger.Error(err)
		}
	}
	return message, err

}

package persistence

import (
	"context"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
)

func (m *MongoAdapter) AddMessage(ctx context.Context, message Message) (*Message, error) {
	message.generateID()
	_, err := m.MessagesCollection.InsertOne(ctx, message)
	if err != nil {
		var e mongo.WriteException
		if errors.As(err, &e) {
			if len(e.WriteErrors) > 0 {
				if e.WriteErrors[0].Code == 11000 {
					return nil, &DuplicateEntityException{entity: "message"}
				}
			}
		} else {
			err = &WriteOperationFailure{e.Error()}
		}
	}
	return &message, err

}

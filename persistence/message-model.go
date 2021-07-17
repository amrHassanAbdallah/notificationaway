package persistence

import (
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

type MessagesDB struct{}

type Notifiers struct {
	MessagesDB `json:"-" bson:"-"`
	Message    `bson:",inline"`
}

func (a MessagesDB) CollectionName() string {
	return "messages"
}
func (a MessagesDB) Indexes() []Index {
	return []Index{
		{
			Name: "name uniqueness",
			Keys: bsonx.Doc{
				{Key: "type", Value: bsonx.Int32(1)},
				{Key: "language", Value: bsonx.Int32(1)},
			},
			Unique:     true,
			Collection: a.CollectionName(),
		},
		{
			Name: "id with created at",
			Keys: bsonx.Doc{
				{Key: "_id", Value: bsonx.Int32(-1)},
				{Key: "created_at", Value: bsonx.Int32(-1)},
			},
			Unique:     false,
			Collection: a.CollectionName(),
		},
	}
}

func (m *Message) generateID() *Message {
	m.Id = uuid.New().String()
	return m

}

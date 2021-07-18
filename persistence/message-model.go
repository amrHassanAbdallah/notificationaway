package persistence

import (
	"fmt"
	"github.com/satori/go.uuid"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"time"
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

func NewMessage(lang, providerType, template, mtype string, templateKeys []string) *Message {
	m := Message{
		Language:     lang,
		ProviderType: providerType,
		Template:     template,
		TemplateKeys: templateKeys,
		Type:         mtype,
		CreatedAt:    time.Now(),
	}
	return m.generateID()
}
func (m *Message) generateID() *Message {
	m.Id = uuid.NewV5(uuid.UUID{}, fmt.Sprintf("%s_%s", m.Type, m.Language)).String()
	return m

}

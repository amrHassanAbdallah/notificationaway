package persistence

import (
	"context"
	"fmt"
	"time"
)

type PersistenceLayerInterface interface {
	AddMessage(ctx context.Context, message Message) (*Message, error)
}

type DuplicateEntityException struct {
	entity string
}

func (nc *DuplicateEntityException) Error() string {
	return fmt.Sprintf("duplicate %s", nc.entity)
}

// NewMessage defines model for NewMessage.
type Message struct {
	Id           string `json:"id" bson:"_id"`
	Language     string `json:"language" bson:"language" validate:"required,oneof=en ar"`
	ProviderType string `json:"provider_type" bson:"provider_type" validate:"required,oneof=sms push"`

	// Message content
	Template     string    `json:"template" bson:"template" validate:"required,gte=5"`
	TemplateKeys *[]string `json:"template_keys,omitempty" bson:"template_keys"`

	// will be used as part of the uniqunes of the message for example type could be greetings, driver arrived,...etc
	Type      string    `json:"type" bson:"type" validate:"required,max=256"` //greetings,
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
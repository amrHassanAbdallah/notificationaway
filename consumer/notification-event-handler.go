package consumer

import (
	"context"
	"github.com/amrHassanAbdallah/notificationaway/persistence"
)

type NotificationEventHandler struct {
	persistenceLayer persistence.ReadLayerInterface
}

func NewNotificationEventHandler(p persistence.ReadLayerInterface) *NotificationEventHandler {
	return &NotificationEventHandler{persistenceLayer: p}
}

func (n *NotificationEventHandler) Handle(ctx context.Context, msg map[string]interface{}) error {
	return nil
}

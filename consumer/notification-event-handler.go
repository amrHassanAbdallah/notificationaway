package consumer

import (
	"context"
	"fmt"
	"github.com/amrHassanAbdallah/notificationaway/persistence"
	"github.com/amrHassanAbdallah/notificationaway/proxy"
)

type NotificationEventHandler struct {
	MessagesService persistence.ReadLayerInterface
}

func NewNotificationEventHandler(p persistence.ReadLayerInterface) *NotificationEventHandler {
	return &NotificationEventHandler{MessagesService: p}
}

func (n *NotificationEventHandler) Handle(ctx context.Context, msg map[string]interface{}) error {
	mtype, k := msg["type"].(string)
	if !k {
		return &ErrUnrecoverableMsgHandling{Message: "Missing mtype field"}
	}
	language, k := msg["lang"].(string)
	if !k {
		return &ErrUnrecoverableMsgHandling{Message: "Missing lang field"}
	}
	to, k := msg["to"].(string)
	if !k {
		return &ErrUnrecoverableMsgHandling{Message: "Missing lang field"}
	}
	m := persistence.NewMessage(language, "", "", mtype, nil)
	message, err := n.MessagesService.GetMessage(ctx, m.Id)
	if err != nil {
		switch err.(type) {
		case *persistence.ErrRecordNotFound:
			return &ErrUnrecoverableMsgHandling{Message: fmt.Sprintf("no message found for type%v and language %v", mtype, language)}
		default:
			return err
		}
	}
	sender, err := proxy.GetSender(message.ProviderType)
	if err != nil {
		return &ErrUnrecoverableMsgHandling{Message: err.Error()}
	}

	return sender.SendNotification(proxy.Notification{
		To:      to,
		Content: message.Template,
	})
}

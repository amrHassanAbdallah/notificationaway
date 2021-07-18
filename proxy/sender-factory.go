package proxy

import "fmt"

func GetSender(senderType string) (NotificationSenderInterface, error) {
	switch senderType {
	case "sms":
		return &SMSNotificationSenderAdapter{}, nil
	case "push":
		return &PushNotificationSenderAdapter{}, nil
	case "webhook":
		return &WebhookNotificationSenderAdapter{}, nil
	default:
		return nil, fmt.Errorf("not supported sender type (%s) supported ones are [sms, push]", senderType)
	}
}

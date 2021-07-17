package proxy

type PushNotificationSenderAdapter struct {
}

func (p *PushNotificationSenderAdapter) SendNotification(n Notification) error {
	return nil
}

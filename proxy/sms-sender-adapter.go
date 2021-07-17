package proxy

type SMSNotificationSenderAdapter struct {
}

func (p *SMSNotificationSenderAdapter) SendNotification(n Notification) error {
	return nil
}

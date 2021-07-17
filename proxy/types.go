package proxy

type Notification struct {
	To      string
	Content string
}

type NotificationSenderInterface interface {
	SendNotification(n Notification) error
}

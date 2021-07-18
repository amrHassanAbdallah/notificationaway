package proxy

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type WebhookNotificationSenderAdapter struct {
}

func (p *WebhookNotificationSenderAdapter) SendNotification(n Notification) error {
	res := map[string]interface{}{
		"message": n.Content,
	}
	bytesRes, err := json.Marshal(res)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", n.To, bytes.NewBuffer(bytesRes))
	client := &http.Client{}
	_, _ = client.Do(req)
	/*if err != nil {//should change it to check if the error related to the network or context timeout maybe for now ignore the error
		return err
	}*/
	return nil
}

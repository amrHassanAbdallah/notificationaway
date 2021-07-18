package main

import (
	"encoding/json"
	"github.com/amrHassanAbdallah/notificationaway/logger"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	var payload map[string]interface{}

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		logger.Fatalf("failed to decode payload, error: %v", err)
	}

	logger.Infow("recieved data", "data", payload)
}
func main() {
	http.HandleFunc("/webhook", handler)
	logger.Fatal(http.ListenAndServe(":5000", nil))
}

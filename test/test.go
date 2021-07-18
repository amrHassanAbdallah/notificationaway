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
	//check the svc some times then fail....
	for i := 0; i < 5; i++ {
		logger.Infow("checking the svc")
		resp, err := http.Get("http://localhost:7981/health")
		if err != nil || resp != nil && resp.StatusCode != 200 {
			code := 0
			if resp != nil {
				code = resp.StatusCode
			}
			logger.Infow("failed to connect to the svc", "status", code, "err", err)
			continue
		}
		logger.Infow("service is up and running on http://localhost:7981")
	}

	http.HandleFunc("/webhook", handler)
	logger.Fatal(http.ListenAndServe(":5000", nil))
}

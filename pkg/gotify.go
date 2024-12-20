package pkg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func SendMessageToGotify(date, message, gotifyURL, gotifyToken string) error {
	url := fmt.Sprintf("%smessage?token=%s", gotifyURL, gotifyToken)

	payload := map[string]interface{}{
		"message":  message,
		"title":    fmt.Sprintf("Training for %s", date),
		"priority": 5,
		"extras": map[string]interface{}{
			"client::display": map[string]interface{}{
				"contentType": "text/markdown",
			},
		},
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	resp, err := http.Post(url, "application/json", bytes.NewReader(payloadBytes))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("error sending gotify message : %s", string(bodyBytes))
	}

	return nil
}

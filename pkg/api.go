package pkg

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Event struct {
	ID          int           `json:"id"`
	StartDate   string        `json:"start_date_local"`
	EndDate     string        `json:"end_date_local"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Type        string        `json:"type"`
	Category    string        `json:"category"`
	Notes       string        `json:"notes"`
	Tags        []string      `json:"tags"`
	Workouts    []interface{} `json:"workouts"`
}

func FetchEventsData(intervalsAPIKey, intervalsUserID, date string) ([]Event, error) {
	var events []Event

	authToken := base64.StdEncoding.EncodeToString([]byte("API_KEY:" + intervalsAPIKey))

	client := &http.Client{}

	eventsURL := fmt.Sprintf("https://intervals.icu/api/v1/athlete/%s/events?oldest=%s&newest=%s", intervalsUserID, date, date)

	req, err := http.NewRequest("GET", eventsURL, nil)
	if err != nil {
		return events, err
	}
	req.Header.Set("Authorization", "Basic "+authToken)

	resp, err := client.Do(req)
	if err != nil {
		return events, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return events, fmt.Errorf("Erreur lors de la récupération des événements : %s", string(bodyBytes))
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return events, err
	}

	err = json.Unmarshal(bodyBytes, &events)
	if err != nil {
		return events, err
	}

	return events, nil
}

func SendMessageToGotify(date, message, gotifyURL, gotifyToken string) error {
	url := fmt.Sprintf("%s/message?token=%s", gotifyURL, gotifyToken)

	payload := map[string]interface{}{
		"message":  message,
		"title":    fmt.Sprintf("Training of the day %s", date),
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
		return fmt.Errorf("Erreur lors de l'envoi du message à Gotify : %s", string(bodyBytes))
	}

	return nil
}

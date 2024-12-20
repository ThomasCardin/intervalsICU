package pkg

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Day struct {
	Events          []Event
	WeatherForecast WeatherForecast
}

type WeatherForecastCoords struct {
}

type WeatherForecast struct {
}

type Event struct {
	ID             int           `json:"id"`
	StartDate      string        `json:"start_date_local"`
	EndDate        string        `json:"end_date_local"`
	Name           string        `json:"name"`
	Description    string        `json:"description"`
	Type           string        `json:"type"`
	Category       string        `json:"category"`
	Notes          string        `json:"notes"`
	DistanceTarget string        `json:"distance_target"`
	Tags           []string      `json:"tags"`
	Workouts       []interface{} `json:"workouts"`
}

func GetDayInformation(intervalsAPIKey, intervalsUserID, date string) (Day, error) {
	var events []Event

	authToken := base64.StdEncoding.EncodeToString([]byte("API_KEY:" + intervalsAPIKey))

	client := &http.Client{}
	eventsURL := fmt.Sprintf("https://intervals.icu/api/v1/athlete/%s/events?oldest=%s&newest=%s", intervalsUserID, date, date)
	req, err := http.NewRequest("GET", eventsURL, nil)
	if err != nil {
		return Day{}, err
	}
	req.Header.Set("Authorization", "Basic "+authToken)

	resp, err := client.Do(req)
	if err != nil {
		return Day{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return Day{}, fmt.Errorf("error fetching events : %s", string(bodyBytes))
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return Day{}, err
	}

	err = json.Unmarshal(bodyBytes, &events)
	if err != nil {
		return Day{}, err
	}

	wfcoords, err := getWeatherForecastCoords(intervalsUserID, authToken)
	if err != nil {
		return Day{}, err
	}

	wf, err := getWeatherForecast(wfcoords)
	if err != nil {
		return Day{}, err
	}

	return Day{
		Events:          events,
		WeatherForecast: wf,
	}, nil
}

func getWeatherForecastCoords(intervalsUserID, authToken string) (WeatherForecastCoords, error) {
	var wfcoords WeatherForecastCoords

	client := &http.Client{}
	wfCoordsUrl := fmt.Sprintf("https://intervals.icu/api/v1/athlete/%s/weather-forecast", intervalsUserID)

	req, err := http.NewRequest("GET", wfCoordsUrl, nil)
	if err != nil {
		return WeatherForecastCoords{}, err
	}
	req.Header.Set("Authorization", "Basic "+authToken)

	resp, err := client.Do(req)
	if err != nil {
		return WeatherForecastCoords{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return WeatherForecastCoords{}, fmt.Errorf("error fetching weather-forecast : %s", string(bodyBytes))
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return WeatherForecastCoords{}, err
	}

	fmt.Printf("%s\n", resp.Body)
	fmt.Printf("%s\n", bodyBytes)

	err = json.Unmarshal(bodyBytes, &wfcoords)
	if err != nil {
		return WeatherForecastCoords{}, err
	}

	return wfcoords, nil
}

func getWeatherForecast(wfcoords WeatherForecastCoords) (WeatherForecast, error) {
	var wf WeatherForecast

	return wf, nil
}

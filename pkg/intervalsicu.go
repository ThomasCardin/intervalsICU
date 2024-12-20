package pkg

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Day struct {
	Events   []Event
	Forecast Forecast
}

type Forecast struct {
	Forecast []WeatherForecast `json:"forecasts"`
}

type WeatherForecast struct {
	Location string          `json:"location"`
	Daily    []DailyForecast `json:"daily"`
}

type Weather struct {
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

type FeelsLike struct {
	Day   float64 `json:"day"`
	Night float64 `json:"night"`
	Eve   float64 `json:"eve"`
	Morn  float64 `json:"morn"`
	Min   float64 `json:"min"`
	Max   float64 `json:"max"`
}

type DailyForecast struct {
	Date        string    `json:"id"`
	Temperature FeelsLike `json:"feels_like"`
	Weather     []Weather `json:"weather"`
	Snow        float64   `json:"snow"`
	Rain        float64   `json:"rain"`
	Sunrise     string    `json:"sunrise"`
	Sunset      string    `json:"sunset"`
	WindSpeed   float64   `json:"wind_speed"`
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

	forecast, err := getWeatherForecast(intervalsUserID, authToken)
	if err != nil {
		return Day{}, err
	}

	return Day{
		Events:   events,
		Forecast: forecast,
	}, nil
}

func getWeatherForecast(intervalsUserID, authToken string) (Forecast, error) {
	var forecast Forecast

	client := &http.Client{}
	wfCoordsUrl := fmt.Sprintf("https://intervals.icu/api/v1/athlete/%s/weather-forecast", intervalsUserID)

	req, err := http.NewRequest("GET", wfCoordsUrl, nil)
	if err != nil {
		return Forecast{}, err
	}
	req.Header.Set("Authorization", "Basic "+authToken)

	resp, err := client.Do(req)
	if err != nil {
		return Forecast{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return Forecast{}, fmt.Errorf("error fetching weather-forecast : %s", string(bodyBytes))
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return Forecast{}, err
	}

	err = json.Unmarshal(bodyBytes, &forecast)
	if err != nil {
		return Forecast{}, err
	}

	return forecast, nil
}

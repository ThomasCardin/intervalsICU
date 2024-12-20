package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ThomasCardin/intervalsICU/pkg"
	"github.com/ThomasCardin/intervalsICU/templates"
)

const (
	INTERVALS_API_KEY = "INTERVALS_API_KEY"
	INTERVALS_USER_ID = "INTERVALS_USER_ID"

	GOTIFY_URL   = "GOTIFY_URL"
	GOTIFY_TOKEN = "GOTIFY_TOKEN"

	TIME_ZONE = "TIME_ZONE"
)

func main() {
	intervalsAPIKey, found := os.LookupEnv(INTERVALS_API_KEY)
	if !found {
		fmt.Printf("%s not found\n", INTERVALS_API_KEY)
	}

	intervalsUserID, found := os.LookupEnv(INTERVALS_USER_ID)
	if !found {
		fmt.Printf("%s not found\n", INTERVALS_USER_ID)
	}

	gotifyURL, found := os.LookupEnv(GOTIFY_URL)
	if !found {
		fmt.Printf("%s not found\n", GOTIFY_URL)
	}

	gotifyToken, found := os.LookupEnv(GOTIFY_TOKEN)
	if !found {
		fmt.Printf("%s not found\n", GOTIFY_TOKEN)
	}

	timeZone, found := os.LookupEnv(TIME_ZONE)
	if !found {
		fmt.Printf("%s not found\n", TIME_ZONE)
	}

	location, err := time.LoadLocation(timeZone)
	if err != nil {
		fmt.Printf("error loading timezone : %v", err)
	}

	currentTime := time.Now().In(location)
	date := currentTime.Format("2006-01-02")

	events, err := pkg.FetchEventsData(intervalsAPIKey, intervalsUserID, date)
	if err != nil {
		log.Fatalf("error fetching events : %v", err)
	}

	message := templates.FormatEventsMessage(date, events)

	err = pkg.SendMessageToGotify(date, message, gotifyURL, gotifyToken)
	if err != nil {
		log.Fatalf("error sending message to gotify : %v", err)
	}

	fmt.Println("message sent to gotify!")
}

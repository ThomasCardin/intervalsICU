package utils

import (
	"fmt"
	"strings"

	"github.com/ThomasCardin/intervalsICU/pkg"
)

const (
	WORKOUT = "WORKOUT"
	NOTE    = "NOTE"
)

func FormatEventsMessage(date string, day pkg.Day) string {
	var message string

	// Forecast
	found := false
	for _, x := range day.Forecast.Forecast {
		for _, v := range x.Daily {
			if v.Date == date {
				message += fmt.Sprintf("### Weather forecast\n")
				message += fmt.Sprintf("- **Location** : %s\n", x.Location)
				message += fmt.Sprintf("- **Feels like** : Morning (%.2f), Evening (%.2f), Night (%.2f), Day (%.2f)\n", v.Temperature.Morn, v.Temperature.Eve, v.Temperature.Night, v.Temperature.Day)
				message += fmt.Sprintf("- **Description** : %s\n", v.Weather[0].Description)
				message += fmt.Sprintf("- **Wind speed** : %.2f\n", v.WindSpeed)
				message += fmt.Sprintf("- **Rain/Snow** : Rain (%.2f), Snow (%.2f)\n", v.Rain, v.Snow)
				message += fmt.Sprintf("- **Sunrise/Sunset** : Sunrise (%s), Sunset (%s)\n", v.Sunrise, v.Sunset)
			}

			found = true
			break
		}

		if found {
			break
		}
	}

	// Events
	if len(day.Events) > 0 {
		var workouts []string
		var notes []string

		for _, event := range day.Events {
			detail := fmt.Sprintf(
				"- **%s**\n  - Type : %s\n  - Category : %s",
				event.Name,
				event.Type,
				event.Category,
			)

			// Distance is in m
			d := event.Distance / 1000

			detail += fmt.Sprintf("\n  - Description : %s", event.Description)
			detail += fmt.Sprintf("\n  - Distance target : %.2f", d)
			detail += fmt.Sprintf("\n  - Tags : %s", strings.Join(event.Tags, ", "))

			switch event.Category {
			case WORKOUT:
				workouts = append(workouts, detail)
			case NOTE:
				notes = append(notes, detail)
			default:
				workouts = append(workouts, detail)
			}
		}

		if len(workouts) > 0 {
			message += "### Workout :\n\n"
			message += strings.Join(workouts, "\n\n") + "\n\n"
		}

		if len(notes) > 0 {
			message += "### Notes :\n\n"
			message += strings.Join(notes, "\n\n") + "\n\n"
		}
	} else {
		message += "\n_No events for today._\n"
	}

	return message
}

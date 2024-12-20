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

	message += fmt.Sprintf("Weather forecast ")

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
			if event.Description != "" {
				detail += fmt.Sprintf("\n  - Description : %s", event.Description)
			}
			if len(event.Tags) > 0 {
				detail += fmt.Sprintf("\n  - Tags : %s", strings.Join(event.Tags, ", "))
			}

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

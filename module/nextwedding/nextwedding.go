package nextwedding

import (
	"fmt"
	"github.com/gamelost/bot3server/server"
	"time"
)

type NextWeddingService struct{}

func (svc *NextWeddingService) NewService() server.BotHandler {
	return &NextWeddingService{}
}

func (svc *NextWeddingService) Handle(botRequest *server.BotRequest, botResponse *server.BotResponse) {

	botResponse.SetSingleLineResponse(durationToWeddingDate())
}

func durationToWeddingDate() string {

	nowDate := time.Now()
	loc, _ := time.LoadLocation("America/Los_Angeles")
	weddingDate := time.Date(2014, time.September, 6, 16, 0, 0, 0, loc)

	if nowDate.After(weddingDate) {
		return "Ah mon. It be too late, the deed is done!"
	} else {

		duration := weddingDate.Sub(nowDate)
		durationInMinutes := int(duration.Minutes())
		durationDays := durationInMinutes / (24 * 60)
		durationHours := (durationInMinutes - (durationDays * 60 * 24)) / 60
		durationMinutes := durationInMinutes - ((durationDays * 60 * 24) + (durationHours * 60))

		return fmt.Sprintf("Only %d days, %d hours and %d minutes left for natech to come to his senses.  Flee while you can, FLEEEEEE!", durationDays, durationHours, durationMinutes)
	}
}

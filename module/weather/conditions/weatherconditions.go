package conditions

import (
	"errors"
	"fmt"
	"github.com/gamelost/bot3server/module/weather"
	"github.com/gamelost/bot3server/server"
	cache "github.com/pmylund/go-cache"
	"net/url"
	"strings"
	"time"
)

type WeatherConditionsService struct {
	WUAPIKey string
	WUAPIURL string
}

var airportCodeCache *cache.Cache
var weatherConditionsCache *cache.Cache

const stateCityAPICallUrl = "http://api.wunderground.com/api/84a6128036456b39/conditions/q/%s/%s.json"
const cityAPICallUrl = "http://api.wunderground.com/api/84a6128036456b39/conditions/q/%s.json"
const airportAPICallUrl = "http://api.wunderground.com/api/84a6128036456b39/conditions/q/%s.json"

func (svc *WeatherConditionsService) NewService() server.BotHandler {
	return &WeatherConditionsService{}
}

func init() {
	// Create a cache with a default expiration time of 5 minutes, and which
	// purges expired items every 30 seconds
	airportCodeCache = cache.New(0, 0)
	airportCodeCache.Set("LAX", "KCAELSEG4", -1)
	airportCodeCache.Set("SFO", "KCASANFR58", -1)

	// create a cache for stored weather forecasts so we dont waste API calls
	weatherConditionsCache = cache.New(time.Minute*15, time.Minute*5)
}

func (svc *WeatherConditionsService) Handle(botRequest *server.BotRequest, botResponse *server.BotResponse) {

	var err error
	var resp string

	wStr := strings.TrimSpace(botRequest.LineTextWithoutCommand())
	weatherCmd := &weather.WeatherLocation{Location: wStr}

	resp, err = getWeatherConditionsForLocation(weatherCmd.Location)

	if err != nil {
		botResponse.SetSingleLineResponse(err.Error())
	} else {
		botResponse.SetSingleLineResponse(resp)
	}
}

func getWeatherConditionsForLocation(command string) (weatherResponse string, err error) {

	// split if there is a state
	args := strings.SplitN(command, ",", 2)

	// do we have two args?
	if len(args) == 2 {

		state := url.QueryEscape(strings.Replace(strings.TrimSpace(args[1]), " ", "_", -1))
		city := url.QueryEscape(strings.Replace(strings.TrimSpace(args[0]), " ", "_", -1))
		call := fmt.Sprintf(stateCityAPICallUrl, state, city)
		weatherData, callErr := weather.DoWeatherAPICall(call)

		if callErr != nil {
			err = callErr
			return
		}

		if len(weatherData.Response.Results) > 0 {
			err = errors.New("WUG API could not find exact match on this?  Your father was a hamster.")
			return
		}

		weatherResponse = parseWeatherDataIntoResponseString(weatherData)

	} else if len(args) == 1 {

		city := url.QueryEscape(strings.Replace(strings.TrimSpace(args[0]), " ", "_", -1))
		callUrl := fmt.Sprintf(cityAPICallUrl, city)
		weatherData, callErr := weather.DoWeatherAPICall(callUrl)

		if callErr != nil {
			err = callErr
			return
		}

		if len(weatherData.Response.Results) > 0 {
			err = errors.New("Multiple results returned.  Try to specify a state along with your city?")
			return
		}

		weatherResponse = parseWeatherDataIntoResponseString(weatherData)

	} else {
		err = errors.New("Unable to parse location string.  Try it in 'City,State' format, or use an airport code.")
	}

	return
}

func parseWeatherDataIntoResponseString(weatherData *weather.WUAPIResponse) string {

	if len(weatherData.Current_observation.Display_location.Full) == 0 {
		return "Location parameters were ambigious.  Try something more specific.  Like City,[State/Country]."
	}

	return fmt.Sprintf("Weather conditions for %s is %s, temperatures around %s and wind conditions of %s", weatherData.Current_observation.Display_location.Full, weatherData.Current_observation.Weather, weatherData.Current_observation.Temperature_string, weatherData.Current_observation.Wind_string)
}

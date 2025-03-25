package bot_commands

import (
	"alice-go-bot/bot/commands"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/bwmarrin/discordgo"
)

type WeatherResponseCurrentCondition struct {
	ObservationTime string `json:"observation_time"`
	TemperatureC    string `json:"temp_C"`
	TemperatureF    string `json:"temp_F"`
	WeatherCode     string `json:"weatherCode"`
	WeatherIconUrl  []struct {
		Value string `json:"value"`
	} `json:"weatherIconUrl"`
	WeatherDesc []struct {
		Value string `json:"value"`
	} `json:"weatherDesc"`
	WindspeedKmph  string `json:"windspeedKmph"`
	WinddirDegree  string `json:"winddirDegree"`
	Winddir16Point string `json:"winddir16Point"`
	PrecipMM       string `json:"precipMM"`
	Humidity       string `json:"humidity"`
	Visibility     string `json:"visibility"`
	Pressure       string `json:"pressure"`
	Cloudcover     string `json:"cloudcover"`
}

type WeatherResponseNearestArea struct {
	AreaName []struct {
		Value string `json:"value"`
	} `json:"areaName"`
	Country []struct {
		Value string `json:"value"`
	} `json:"country"`
	Region []struct {
		Value string `json:"value"`
	} `json:"region"`
	Latitude   string `json:"latitude"`
	Longitude  string `json:"longitude"`
	Population string `json:"population"`
	WeatherUrl []struct {
		Value string `json:"value"`
	} `json:"weatherUrl"`
}

type WeatherResponseWeather struct {
	Date      string `json:"date"`
	Astronomy []struct {
		Sunrise   string `json:"sunrise"`
		Sunset    string `json:"sunset"`
		Moonphase string `json:"moon_phase"`
		Moonrise  string `json:"moonrise"`
		Moonset   string `json:"moonset"`
	} `json:"astronomy"`
	MaxtempC    string `json:"maxtempC"`
	MintempC    string `json:"mintempC"`
	TotalSnowCm string `json:"totalSnow_cm"`
	SunHour     string `json:"sunHour"`
	UvIndex     string `json:"uvIndex"`
}

type WeatherResponse struct {
	CurrentCondition []WeatherResponseCurrentCondition `json:"current_condition"`
	NearestArea      []WeatherResponseNearestArea      `json:"nearest_area"`
	Weather          []WeatherResponseWeather          `json:"weather"`
}

type WeatherRequest struct {
	Location string
}

func getWeather(request WeatherRequest) (WeatherResponse, error) {
	resp, err := http.Get("https://wttr.in/" + request.Location + "?format=j1")
	if err != nil {
		return WeatherResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return WeatherResponse{}, fmt.Errorf("%v", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return WeatherResponse{}, err
	}

	var weatherResponse WeatherResponse
	err = json.Unmarshal(body, &weatherResponse)
	if err != nil {
		return WeatherResponse{}, fmt.Errorf("`%s`", body)
	}

	return weatherResponse, nil
}

var moon_emojis = map[string]string{
	"New Moon":        "🌑",
	"Waxing Crescent": "🌒",
	"First Quarter":   "🌓",
	"Waxing Gibbous":  "🌔",
	"Full Moon":       "🌕",
	"Waning Gibbous":  "🌖",
	"Last Quarter":    "🌗",
	"Waning Crescent": "🌘",
}

var weather_emojis = map[string]string{
	"Sunny":                                    "☀️",
	"Clear":                                    "☀️",
	"Partly cloudy":                            "⛅",
	"Cloudy":                                   "☁️",
	"Overcast":                                 "🌥️",
	"Mist":                                     "🌫️",
	"Patchy rain possible":                     "🌦️",
	"Patchy snow possible":                     "🌨️",
	"Patchy sleet possible":                    "🌨️",
	"Patchy freezing drizzle possible":         "🌨️",
	"Thundery outbreaks possible":              "🌩️",
	"Blowing snow":                             "🌨️",
	"Blizzard":                                 "🌨️",
	"Fog":                                      "🌫️",
	"Freezing fog":                             "🌫️",
	"Patchy light drizzle":                     "🌧️",
	"Light drizzle":                            "🌧️",
	"Freezing drizzle":                         "🌧️",
	"Heavy freezing drizzle":                   "🌧️",
	"Patchy light rain":                        "🌧️",
	"Light rain":                               "🌧️",
	"Moderate rain at times":                   "🌧️",
	"Moderate rain":                            "🌧️",
	"Heavy rain at times":                      "🌧️",
	"Heavy rain":                               "🌧️",
	"Light freezing rain":                      "🌧️",
	"Moderate or heavy freezing rain":          "🌧️",
	"Light sleet":                              "🌨️",
	"Moderate or heavy sleet":                  "🌨️",
	"Patchy light snow":                        "🌨️",
	"Light snow":                               "🌨️",
	"Patchy moderate snow":                     "🌨️",
	"Moderate snow":                            "🌨️",
	"Patchy heavy snow":                        "🌨️",
	"Heavy snow":                               "🌨️",
	"Ice pellets":                              "🌨️",
	"Light rain shower":                        "🌧️",
	"Moderate or heavy rain shower":            "🌧️",
	"Torrential rain shower":                   "🌧️",
	"Light sleet showers":                      "🌨️",
	"Moderate or heavy sleet showers":          "🌨️",
	"Light snow showers":                       "🌨️",
	"Moderate or heavy snow showers":           "🌨️",
	"Light showers of ice pellets":             "🌨️",
	"Moderate or heavy showers of ice pellets": "🌨️",
	"Patchy light rain with thunder":           "🌩️",
	"Moderate or heavy rain with thunder":      "🌩️",
	"Patchy light snow with thunder":           "🌨️",
	"Moderate or heavy snow with thunder":      "🌨️",
	"Haze":                                     "🌫️",
}

func weatherCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	weather_request := WeatherRequest{
		Location: i.ApplicationCommandData().Options[0].StringValue(),
	}
	weather_response, err := getWeather(weather_request)
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("Failed to fetch weather: %s", err),
			},
		})
		return
	}
	fields := make([]*discordgo.MessageEmbedField, 0)
	for _, currentCondition := range weather_response.CurrentCondition {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   "Temperature",
			Value:  fmt.Sprintf("%s°C", currentCondition.TemperatureC),
			Inline: true,
		})
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   "Weather",
			Value:  fmt.Sprintf("%s %s", weather_emojis[currentCondition.WeatherDesc[0].Value], currentCondition.WeatherDesc[0].Value),
			Inline: true,
		})
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   "Wind",
			Value:  fmt.Sprintf("%s km/h %s", currentCondition.WindspeedKmph, currentCondition.Winddir16Point),
			Inline: true,
		})
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   "Humidity",
			Value:  fmt.Sprintf("%s%%", currentCondition.Humidity),
			Inline: true,
		})
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   "Pressure",
			Value:  fmt.Sprintf("%s hPa", currentCondition.Pressure),
			Inline: true,
		})
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   "Visibility",
			Value:  fmt.Sprintf("%s km", currentCondition.Visibility),
			Inline: true,
		})
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   "Cloud Cover",
			Value:  fmt.Sprintf("%s%%", currentCondition.Cloudcover),
			Inline: true,
		})
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   "Precipitation",
			Value:  fmt.Sprintf("%s mm", currentCondition.PrecipMM),
			Inline: true,
		})
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   "Moon Phase",
			Value:  moon_emojis[weather_response.Weather[0].Astronomy[0].Moonphase],
			Inline: true,
		})
	}
	embed := &discordgo.MessageEmbed{
		Title:       fmt.Sprintf("Weather for %s", weather_response.NearestArea[0].AreaName[0].Value),
		Description: fmt.Sprintf("%s, %s `(Population: %s)`", weather_response.NearestArea[0].AreaName[0].Value, weather_response.NearestArea[0].Country[0].Value, weather_response.NearestArea[0].Population),
		Fields:      fields,
	}
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
		},
	})
}

var WeatherCommand commands.BotCommand = commands.BotCommand{
	Info: discordgo.ApplicationCommand{
		Name:        "weather",
		Description: "Check the weather",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "location",
				Description: "The location to check the weather",
				Required:    true,
			},
		},
	},
	Handler: weatherCommandHandler,
}

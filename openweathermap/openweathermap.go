package openweathermap

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

var (
	ErrCityNotFound = errors.New("city not found")
)

type OpenWeatherMap struct {
	token string
}

func New(token string) *OpenWeatherMap {
	return &OpenWeatherMap{token: token}
}

func (owm *OpenWeatherMap) GetCityCoordinates(city string) (lat float64, lon float64, err error) {
	URL := url.URL{
		Scheme: "https",
		Host:   "api.openweathermap.org",
		Path:   "data/2.5/weather",
	}

	query := URL.Query()
	query.Add("q", city)
	query.Add("appid", owm.token)

	URL.RawQuery = query.Encode()

	resp, err := http.Get(URL.String())
	if err != nil {
		return 0., 0., err
	}
	defer resp.Body.Close()

	buff, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0., 0., errors.New("failed to read response body")
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		if resp.StatusCode == http.StatusNotFound {
			return 0, 0, ErrCityNotFound
		}

		var errorResponse ErrorResponse
		if err := json.Unmarshal(buff, &errorResponse); err != nil {
			return 0., 0., errors.New("unable to unmarshal error response")
		}

		return 0., 0., errors.New(errorResponse.Message)
	}

	type Coordinates struct {
		Lat float64 `json:"lat"`
		Lon float64 `json:"lon"`
	}

	type CurrentWeatherResponse struct {
		Coord Coordinates `json:"coord"`
	}

	var weatherResponse CurrentWeatherResponse
	if err := json.Unmarshal(buff, &weatherResponse); err != nil {
		return 0., 0., errors.New("unable to unmarshal weather response")
	}

	return weatherResponse.Coord.Lat, weatherResponse.Coord.Lon, nil
}

func (owm *OpenWeatherMap) OneCallByCoordinates(lat float64, lon float64) (*OneCallWeatherDataResponse, error) {
	URL := url.URL{
		Scheme: "https",
		Host:   "api.openweathermap.org",
		Path:   "data/2.5/onecall",
	}

	query := URL.Query()
	query.Add("lang", "ru")
	query.Add("appid", owm.token)
	query.Add("lat", fmt.Sprintf("%.1f", lat))
	query.Add("lon", fmt.Sprintf("%.1f", lon))
	query.Add("units", "metric")

	URL.RawQuery = query.Encode()

	resp, err := http.Get(URL.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	buff, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("failed to read response body")
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		var errorResponse ErrorResponse
		if err := json.Unmarshal(buff, &errorResponse); err != nil {
			return nil, errors.New("unable to unmarshal error response")
		}

		return nil, errors.New("error response")
	}

	var weatherResponse OneCallWeatherDataResponse
	if err := json.Unmarshal(buff, &weatherResponse); err != nil {
		return nil, errors.New("unable to unmarshal weather response")
	}

	return &weatherResponse, nil
}

type ErrorResponse struct {
	Code    int    `json:"cod"`
	Message string `json:"message"`
}

type OneCallWeatherDataResponse struct {
	Lat      float64   `json:"lat"`
	Lon      float64   `json:"lon"`
	Timezone string    `json:"timezone"`
	Current  *Current  `json:"current"`
	Hourly   []*Hourly `json:"hourly"`
	Daily    []*Daily  `json:"daily"`
}

type Current struct {
	Dt         int     `json:"dt"`
	Sunrise    int     `json:"sunrise"`
	Sunset     int     `json:"sunset"`
	Temp       float64 `json:"temp"`
	FeelsLike  float64 `json:"feels_like"`
	Pressure   int     `json:"pressure"`
	Humidity   int     `json:"humidity"`
	Uvi        float64 `json:"uvi"`
	Clouds     int     `json:"clouds"`
	Visibility int     `json:"visibility"`
	WindSpeed  float64 `json:"wind_speed"`
	WindDeg    int     `json:"wind_deg"`
	Weather    []*Weather
	Rain       *Rain `json:"rain"`
}

type Hourly struct {
	Dt        int        `json:"dt"`
	Temp      float64    `json:"temp"`
	FeelsLike float64    `json:"feels_like"`
	Pressure  int        `json:"pressure"`
	Humidity  int        `json:"humidity"`
	Clouds    int        `json:"clouds"`
	WindSpeed float64    `json:"wind_speed"`
	WindDeg   int        `json:"wind_deg"`
	Weather   []*Weather `json:"weather"`
	Rain      *Rain      `json:"rain"`
}

type Temperature struct {
	Day     float64 `json:"day"`
	Min     float64 `json:"min"`
	Max     float64 `json:"max"`
	Night   float64 `json:"night"`
	Eve     float64 `json:"eve"`
	Morning float64 `json:"morn"`
}

type FeelsLike struct {
	Day   float64 `json:"day"`
	Night float64 `json:"night"`
	Eve   float64 `json:"eve"`
	Morn  float64 `json:"morn"`
}

type Daily struct {
	Dt          int          `json:"dt"`
	Sunrise     int          `json:"sunrise"`
	Sunset      int          `json:"sunset"`
	Temperature *Temperature `json:"temp"`
	FeelsLike   *FeelsLike   `json:"feels_like"`
	Pressure    int          `json:"pressure"`
	Humidity    int          `json:"humidity"`
	WindSpeed   float64      `json:"wind_speed"`
	WindDeg     int          `json:"wind_deg"`
	Weather     []*Weather   `json:"weather"`
	Clouds      int          `json:"clouds"`
	Rain        float64      `json:"rain"`
	Uvi         float64      `json:"uvi"`
}

type Weather struct {
	ID          int    `json:"id"`
	Main        string `json:"main"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

type Rain struct {
	OneHour float64 `json:"1h"`
}

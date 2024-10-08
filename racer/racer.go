package racer

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"io"
	"math"
	"net/http"
	"sort"
)

type TemperatureData struct {
	Station     string  `json:"station"`
	Temperature float64 `json:"temperature"`
}

type StationAvg struct {
	Station     string  `json:"station"`
	AverageTemp float64 `json:"temperature"`
}

type TemperatureResponse struct {
	RacerID  string       `json:"racerId"`
	Averages []StationAvg `json:"averages"`
}
type TokenRequest struct {
	Token string `json:"token"`
}

type Raced struct {
	ID      string `json:"id"`
	Token   string `json:"token"`
	RacerID string `json:"racerId"`
	Laps    int    `json:"laps,omitempty"`
}

var races = make(map[string]string)

func Hello(c echo.Context) error {
	return c.String(http.StatusOK, "Welcome to the racer API")
}

func Race(c echo.Context) error {
	var tokenReq TokenRequest

	if err := json.NewDecoder(c.Request().Body).Decode(&tokenReq); err != nil {
		fmt.Printf("Error decoding JSON: %v\n", err)
		return c.JSON(http.StatusBadRequest, "failed to decode request body")
	}

	raceID := uuid.New().String()
	racerID := "9e1f0369-3d77-4652-b508-83c4330b2267"

	races[raceID] = tokenReq.Token
	fmt.Printf("Started race with ID: %s, RacerID: %s, Token: %s\n", raceID, racerID, tokenReq.Token)

	return c.JSON(http.StatusCreated, Raced{
		ID:      raceID,
		Token:   tokenReq.Token,
		RacerID: racerID,
	})
}

func RaceLap(c echo.Context) error {
	raceID := c.Param("id")
	if raceID == "" {
		return c.JSON(http.StatusBadRequest, "Race ID is missing")
	}

	if lastToken, ok := races[raceID]; ok {
		bodyBytes, err := io.ReadAll(c.Request().Body)
		if err != nil {
			return c.JSON(http.StatusBadRequest, "failed to read request body")
		}

		newToken := string(bodyBytes)
		races[raceID] = newToken

		return c.JSON(http.StatusOK, Raced{
			RacerID: "9e1f0369-3d77-4652-b508-83c4330b2267",
			Token:   lastToken,
		})
	}

	return c.JSON(http.StatusNotFound, "race not found")
}
func Temperature(c echo.Context) error {
	if c.Request().Header.Get("Content-Encoding") == "gzip" {
		fmt.Println("it has a gzip file")
	} else {
		fmt.Println("it does not contain a gzip file")
	}

	bodyBytes, err := io.ReadAll(c.Request().Body)
	if err != nil {
		fmt.Println("Failed to read request body:", err)
		return c.JSON(http.StatusBadRequest, "Failed to read request body")
	}

	var tempData []TemperatureData
	if err := json.Unmarshal(bodyBytes, &tempData); err != nil {
		fmt.Printf("Failed to decode request body: %v\nBody: %s\n", err, string(bodyBytes))
		return c.JSON(http.StatusBadRequest, "Failed to decode request body")
	}

	stationTemps := make(map[string][]float64)
	for _, data := range tempData {
		stationTemps[data.Station] = append(stationTemps[data.Station], data.Temperature)
	}

	var stationNames []string
	for station := range stationTemps {
		stationNames = append(stationNames, station)
	}
	sort.Strings(stationNames)

	var stationAverages []StationAvg
	for _, station := range stationNames {
		temps := stationTemps[station]
		sum := 0.0
		for _, temp := range temps {
			sum += temp
		}
		avg := sum / float64(len(temps))
		stationAverages = append(stationAverages, StationAvg{
			Station:     station,
			AverageTemp: math.Round(avg*100000) / 100000,
		})
	}

	response := TemperatureResponse{
		RacerID:  "9e1f0369-3d77-4652-b508-83c4330b2267",
		Averages: stationAverages,
	}

	return c.JSON(http.StatusOK, response)
}

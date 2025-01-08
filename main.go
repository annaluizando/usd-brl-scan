package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type USDBRLResponse struct {
	High string `json:"high"`
	Low  string `json:"low"`
}

type Response struct {
	USDBRL USDBRLResponse `json:"USDBRL"`
}

type HistoricalData struct {
	High string `json:"high"`
	Low  string `json:"low"`
}

func main() {
	todaysValue, err := getTodayValueFromAPI()
	if err != nil {
		fmt.Printf(err.Error())
		return
	}

	fmt.Printf("Today's value: $%.3f\n", todaysValue)
	averageSummary("7", todaysValue)
	averageSummary("15", todaysValue)
	averageSummary("30", todaysValue)
}

func averageSummary(numberOfDays string, todaysValue float64) {
	daysData, err := fetchHistoricalData(numberOfDays)
	if err != nil {
		fmt.Printf("Error fetching historical data: %v\n", err)
		return
	}

	daysAverage, err := calculateAverage(daysData)
	if err != nil {
		fmt.Printf("Error calculating thirteen average: %v\n", err)
		return
	}

	Yellow := "\033[33m"
	Green := "\033[32m"
	Red := "\033[31m"
	Reset := "\033[0m"
	Bold := "\033[1m"
	Underline := "\033[4m"

	fmt.Printf(Bold+Underline+"%v-Day"+Reset+" Average: $%.3f\n", numberOfDays, daysAverage)

	if todaysValue > daysAverage {
		fmt.Printf(Green+"Today's value is above the %v-day average.\n"+Reset, numberOfDays)
	} else if todaysValue < daysAverage {
		fmt.Printf(Red+"Today's value is below the %v-day average.\n"+Reset, numberOfDays)
	} else {
		fmt.Printf(Yellow+"Today's value is equal to the %v-day average.\n"+Reset, numberOfDays)
	}
}

func calculateAverage(data []HistoricalData) (float64, error) {
	if len(data) == 0 {
		return 0, fmt.Errorf("no historical data available")
	}

	total := 0.0
	for _, entry := range data {
		high, err := strconv.ParseFloat(entry.High, 64)
		if err != nil {
			return 0, fmt.Errorf("error parsing high value: %v", err)
		}

		low, err := strconv.ParseFloat(entry.Low, 64)
		if err != nil {
			return 0, fmt.Errorf("error parsing low value: %v", err)
		}

		total += (high + low) / 2
	}

	average := total / float64(len(data))
	return average, nil
}

func getTodayValueFromAPI() (float64, error) {
	cotacao_usd_api := "https://economia.awesomeapi.com.br/json/last/USD-BRL"

	res, err := http.Get(cotacao_usd_api)
	if err != nil {
		fmt.Printf(err.Error())
		return 0, err
	}

	defer res.Body.Close()

	var response Response

	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		print(err)
	}

	usdHigh := response.USDBRL.High
	usdLow := response.USDBRL.Low

	convertedHigh, err := strconv.ParseFloat(usdHigh, 64)
	if err != nil {
		panic(err)
	}

	convertedLow, err := strconv.ParseFloat(usdLow, 64)
	if err != nil {
		panic(err)
	}

	valuesAverage := (convertedHigh + convertedLow) / 2

	return valuesAverage, nil
}

func fetchHistoricalData(numberOfDays string) ([]HistoricalData, error) {
	apiURL := fmt.Sprintf("https://economia.awesomeapi.com.br/json/daily/USD-BRL/%s", numberOfDays)

	res, err := http.Get(apiURL)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var response []HistoricalData
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, err
	}

	return response, nil
}

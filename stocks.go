package main

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2/log"
	"github.com/joho/godotenv"
)

func envVariable(key string) string {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return os.Getenv(key)
}

var PoligonPath = envVariable("POLYGON_API_PATH")
var ApiKey = "apiKey=" + envVariable("POLYGON_API_KEY")

var TickerPath = PoligonPath + "v3/reference/tickers"
var DailyValuesPath = PoligonPath + "v1/open-close"

type Stock struct {
	Ticker string `json:"ticker"`
	Name   string `json:"name"`
}

type SearchResult struct {
	Results []Stock `json:"results"`
}

type Values struct {
	Symbol string  `json:"symbol"`
	Open   float64 `json:"open"`
	High   float64 `json:"high"`
	Low    float64 `json:"low"`
	Close  float64 `json:"close"`
}

func Fetch(path string) string {
	resp, err := http.Get(path)
	if err != nil {
		log.Fatal(err)
	}

	if http.StatusOK != resp.StatusCode {
		log.Fatal("Invalid status code: ", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	return string(body)
}

func SearchTicker(ticker string) []Stock {
	body := Fetch(TickerPath + "?" + ApiKey + "&ticker=" + strings.ToUpper(ticker))

	data := SearchResult{}
	json.Unmarshal([]byte(string(body)), &data)

	return data.Results
}

func GetDailyValues(ticker string) Values {
	body := Fetch(DailyValuesPath + "/" + strings.ToUpper(ticker) + "/2023-09-15/?" + ApiKey)

	data := Values{}
	json.Unmarshal([]byte(string(body)), &data)

	return data
}

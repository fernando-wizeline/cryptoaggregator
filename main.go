package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"

	"time"

	"github.com/caarlos0/env/v6"
	"github.com/dghubble/sling"
	"github.com/gin-gonic/gin"

	"ferwizeline.com/cryptoaggregator/types"
)

type environment struct {
	Env                 string        `env:"ENV" envDefault:"production"`
	Port                string        `env:"PORT" envDefault:":8888"`
	HTTPClientTimeout   time.Duration `env:"HTTP_CLIENT_TIMEOUT" envDefault:"15s"`
	DataProviderURL     url.URL       `env:"DATA_PROVIDER_URL" envDefault:"https://stage.bitso.com/api/v3/ticker?book="`
	DataProviderTimeout time.Duration `env:"DATA_PROVIDER_SERVICE_TIMEOUT" envDefault:"30s"`
}

func main() {
	_, cancelFn := context.WithCancel(context.Background())
	defer cancelFn()

	var config environment

	if err := env.Parse(&config); err != nil {
		log.Panic(fmt.Errorf("unable to read config, exiting: %w", err))
	}

	r := gin.Default()
	aggregationsv1 := r.Group("/aggregations/v1")
	aggregationsv1.GET("/aggregate", handleAggregationsGet)

	r.Run(config.Port)
}

func handleAggregationsGet(c *gin.Context) {
	types.JSONResponse(c, func() (any, error) {
		inputLayouts, err := loadInputLayouts()

		if err != nil {
			return nil, err
		}

		return getAggregations(inputLayouts)
	})
}

func loadInputLayouts() (types.InputLayouts, error) {
	b, err := os.ReadFile("fixtures/inputlayout.json")
	if err != nil {
		log.Fatal("Failed to read json from file")
	}

	var il types.InputLayouts
	err = json.NewDecoder(bytes.NewBuffer(b)).Decode(&il)

	if err != nil {
		log.Fatal("unable to parse JSON object")
		return nil, err
	}

	return il, nil

}

func getAggregations(inputLayouts types.InputLayouts) (types.OutputLayouts, error) {

	outputLayouts := types.OutputLayouts{}

	for _, il := range inputLayouts {
		tickerResponse := new(types.TickerResponse)

		_, err := sling.New().Base(fmt.Sprintf("https://stage.bitso.com/api/v3/ticker?book=%s_mxn", il.Component)).Receive(tickerResponse, nil)

		if err != nil {
			return nil, err
		}

		outputLayout := types.Layout{}
		outputLayout.Id = il.Id
		outputLayout.Component = il.Component
		outputLayout.Model.Name = tickerResponse.Payload.Book
		outputLayout.Model.Date = tickerResponse.Payload.CreatedAt
		outputLayout.Model.Price.MXN = tickerResponse.Payload.Last
		outputLayout.Model.Price.USD = tickerResponse.Payload.Last / 21
		outputLayout.Model.TickerSymbol = strings.ToUpper(strings.Split(tickerResponse.Payload.Book, "_")[0])

		outputLayouts = append(outputLayouts, outputLayout)

	}

	return outputLayouts, nil
}

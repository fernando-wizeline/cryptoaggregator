package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"time"

	"github.com/caarlos0/env/v6"
	"github.com/dghubble/sling"
	"github.com/gin-gonic/gin"

	"ferwizeline.com/cryptoaggregator/types"
)

const (
	failureHeader = "x-failure-reason"
)

type environment struct {
	Env                 string        `env:"ENV" envDefault:"production"`
	Port                string        `env:"PORT" envDefault:":8888"`
	HTTPClientTimeout   time.Duration `env:"HTTP_CLIENT_TIMEOUT" envDefault:"15s"`
	DataProviderURL     url.URL       `env:"DATA_PROVIDER_URL" envDefault:"https://stage.bitso.com/api/v3/ticker?book="`
	DataProviderTimeout time.Duration `env:"DATA_PROVIDER_SERVICE_TIMEOUT" envDefault:"30s"`
}

type jsonResponseFunc func() (result any, err error)

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

type RouterError struct {
	StatusCode int                    `json:"code"`
	Message    string                 `json:"message"`
	Details    map[string]interface{} `json:"details,omitempty" swaggerignore:"true"`
	Err        error                  `json:"-"`
}

func (r RouterError) Error() string {
	return r.Message
}

func handleAggregationsGet(c *gin.Context) {
	JSONResponse(c, func() (any, error) {
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

func JSONResponse(c *gin.Context, f jsonResponseFunc) {
	if result, err := f(); err != nil {
		AbortIfErr(c, err)
	} else {
		c.JSON(http.StatusOK, result)
	}
}

func AbortIfErr(c *gin.Context, err error) bool {
	if err != nil {
		_ = c.Error(err) //nolint:errcheck
		var routerErr RouterError
		if errors.As(err, &routerErr) {
			c.Header(failureHeader, routerErr.Message)
			c.AbortWithStatusJSON(routerErr.StatusCode, &routerErr)
		} else {
			c.Header(failureHeader, err.Error())
			routerErr = RouterError{
				StatusCode: http.StatusInternalServerError,
				Message:    err.Error(),
				Err:        err,
			}
			c.AbortWithStatusJSON(http.StatusInternalServerError, &routerErr)
		}
	}

	return err == nil
}

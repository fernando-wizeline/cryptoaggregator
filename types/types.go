package types

import (
	"time"

	"github.com/gin-gonic/gin"
)

type Currency struct {
	Date         time.Time `json:"date"`
	Name         string    `json:"name"`
	TickerSymbol string    `json:"ticker_symbol"`
	Price        struct {
		USD string `json:"usd"`
		MXN string `json:"mxn"`
	} `json:"price"`
}

type Layout struct {
	Id        int      `json:"id"`
	Component string   `json:"component"`
	Model     Currency `json:"model"`
}

type TickerResponse struct {
	Success bool `json:"success"`
	Payload struct {
		High      string    `json:"payload"`
		Last      string    `json:"last"`
		Low       string    `json:"low"`
		Book      string    `json:"book"`
		CreatedAt time.Time `json:"created_at"`
	} `json:"payload"`
}

type AggregatorParams struct {
	Context      *gin.Context
	InputLayouts InputLayouts
}

type FixtureLoaderParams struct {
	Context    *gin.Context
	PathToJSON string
}

type InputLayouts []Layout //A slice of layouts whose Model is empty.

type OutputLayouts []Layout //A slice of layouts whose Model is filled.

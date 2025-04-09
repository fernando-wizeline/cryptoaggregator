package types

import "time"

type Currency struct {
	Date         time.Time `json:"date"`
	Name         string    `json:"name"`
	TickerSymbol string    `json:"ticker_symbol"`
	Price        struct {
		USD int `json:"usd"`
		MXN int `json:"mxn"`
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
		High      int       `json:"payload"`
		Last      int       `json:"last"`
		Low       int       `json:"low"`
		Book      string    `json:"book"`
		CreatedAt time.Time `json:"created_at"`
	} `json:"payload"`
}

type InputLayouts []Layout //A slice of layouts whose Model is empty.

type OutputLayouts []Layout //A slice of layouts whose Model is filled.

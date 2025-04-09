package aggregators

import (
	"fmt"
	"strings"

	"ferwizeline.com/cryptoaggregator/types"
	"github.com/dghubble/sling"
)

type BitsoAggregator struct {
	inputLayouts types.InputLayouts
}

func NewBitsoAggregator(params types.AggregatorParams) *BitsoAggregator {

	return &BitsoAggregator{
		inputLayouts: params.InputLayouts,
	}

}

func (a BitsoAggregator) GetAggregations() (types.OutputLayouts, error) {

	outputLayouts := types.OutputLayouts{}

	for _, il := range a.inputLayouts {
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

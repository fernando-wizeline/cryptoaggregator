package aggregators

import (
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"
	"time"

	"ferwizeline.com/cryptoaggregator/types"
	"github.com/dghubble/sling"
	"github.com/gin-gonic/gin"
)

type BitsoAggregator struct {
	inputLayouts types.InputLayouts
	context      *gin.Context
}

func NewBitsoAggregator(params types.AggregatorParams) *BitsoAggregator {

	return &BitsoAggregator{
		inputLayouts: params.InputLayouts,
		context:      params.Context,
	}

}

func (a BitsoAggregator) GetAggregations() (types.OutputLayouts, error) {
	timer := time.NewTimer(1 * time.Second)
	outputLayouts := types.OutputLayouts{}

	select {
	case <-a.context.Done():
		log.Println("Error when processing request:", a.context.Err())
		return nil, a.context.Err()

	case <-timer.C:
		for _, il := range a.inputLayouts {
			tickerResponse := new(types.TickerResponse)

			rsp, err := sling.New().Base("https://stage.bitso.com/api/v3/ticker?book=btc_mxn").Receive(tickerResponse, nil)
			if err != nil {
				fmt.Println(err)
			}

			body, err := io.ReadAll(rsp.Body)
			if err != nil {
				fmt.Println(err)
			}

			fmt.Println(string(body))
			rsp.Body.Close()

			outputLayout := types.Layout{}
			outputLayout.Id = il.Id
			outputLayout.Component = il.Component
			outputLayout.Model.Name = tickerResponse.Payload.Book
			outputLayout.Model.Date = tickerResponse.Payload.CreatedAt
			outputLayout.Model.Price.MXN = tickerResponse.Payload.Last

			usdPrice, err := strconv.ParseFloat(tickerResponse.Payload.Last, 32)
			if err != nil {
				fmt.Println(err)
			}

			outputLayout.Model.Price.USD = fmt.Sprintf("%f", usdPrice/21)
			outputLayout.Model.TickerSymbol = strings.ToUpper(strings.Split(tickerResponse.Payload.Book, "_")[0])

			outputLayouts = append(outputLayouts, outputLayout)

		}

		return outputLayouts, nil
	}

}

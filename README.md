# cryptoaggregator
Go micro-service that polls cryptocurrency tickers

## Overview
cryptoaggregator provides a Go interface to implement your own aggregator for other crypto ticker services as well as a default aggregator for the Bitso Ticker endpoint.
Similarly, a fixture loader interface is provided to load JSON inputs from any location. The default loader receives a json file located in the fixtures folder.

## How to run
From the root of the project run:
go run main.go
Then, point your browser to http://localhost:8888/aggregations/v1/aggregate and hit enter. After a couple of seconds, the JSON payload with the requestes information will be displayed.
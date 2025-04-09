package api

import "ferwizeline.com/cryptoaggregator/types"

type Aggregator interface {
	GetAggregations() (types.OutputLayouts, error)
}

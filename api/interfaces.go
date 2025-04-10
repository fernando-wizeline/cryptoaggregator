package api

import "ferwizeline.com/cryptoaggregator/types"

type Aggregator interface {
	GetAggregations() (types.OutputLayouts, error)
}

type FixtureLoader interface {
	GetFixture() (types.InputLayouts, error)
}

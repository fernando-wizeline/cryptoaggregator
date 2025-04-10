package loaders

import (
	"bytes"
	"encoding/json"
	"log"
	"os"

	"ferwizeline.com/cryptoaggregator/types"
	"github.com/gin-gonic/gin"
)

type DefaultLoader struct {
	pathToJSON string
	context    *gin.Context
}

func NewDefaultLoader(params types.FixtureLoaderParams) *DefaultLoader {

	return &DefaultLoader{
		pathToJSON: params.PathToJSON,
		context:    params.Context,
	}

}

func (l DefaultLoader) GetFixture() (types.InputLayouts, error) {
	b, err := os.ReadFile(l.pathToJSON)
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

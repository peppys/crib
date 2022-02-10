package estimators

import (
	"fmt"
	"github.com/peppys/crib/internal/zillow"
	"net/http"
)

func DefaultZillowEstimator() func(address string) ([]Estimate, error) {
	return NewZillowEstimator(zillow.NewClient(http.DefaultClient))
}

func NewZillowEstimator(client *zillow.Client) func(address string) ([]Estimate, error) {
	return func(address string) ([]Estimate, error) {
		searchResponse, err := client.SearchProperties(address)
		if err != nil {
			return nil, fmt.Errorf("error looking up address on zillow: %w", err)
		}
		if searchResponse.Results[0].MetaData.Zpid == 0 {
			return nil, fmt.Errorf("could not find property on zillow")
		}
		propertyResponse, err := client.LookupProperty(searchResponse.Results[0].MetaData.Zpid)
		if err != nil {
			return nil, fmt.Errorf("error looking up property on zillow: %w", err)
		}

		return []Estimate{
			{
				Vendor: Zillow,
				Value:  float64(propertyResponse.LookupResults[0].Estimates.Zestimate),
			},
		}, nil
	}
}

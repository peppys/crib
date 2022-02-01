package estimators

import (
	"fmt"
	"github.com/peppys/crib/internal/services/property"
	"github.com/peppys/crib/pkg/zillow"
)

func NewZillowEstimator(client *zillow.Client) func(address string) ([]property.Estimate, error) {
	return func(address string) ([]property.Estimate, error) {
		searchResponse, err := client.SearchProperties(address)
		if err != nil {
			return nil, fmt.Errorf("error looking up address on zillow: %w", err)
		}
		if len(searchResponse.Results) == 0 {
			return nil, fmt.Errorf("could not find property on zillow")
		}
		propertyResponse, err := client.LookupProperty(searchResponse.Results[0].MetaData.Zpid)
		if err != nil {
			return nil, fmt.Errorf("error looking up property on zillow: %w", err)
		}

		return []property.Estimate{
			{
				Vendor: property.Zillow,
				Value:  float64(propertyResponse.LookupResults[0].Estimates.Zestimate),
			},
		}, nil
	}
}

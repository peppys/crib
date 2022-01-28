package estimators

import (
	"fmt"
	"github.com/peppys/crib/internal/services/property"
	"github.com/peppys/crib/pkg/redfin"
	"strings"
)

func NewRedfinEstimator(client *redfin.Client) func(address string) ([]property.Estimate, error) {
	return func(address string) ([]property.Estimate, error) {
		searchResponse, err := client.SearchProperties(address)
		if err != nil {
			return nil, fmt.Errorf("error searching redfin properties: %v", err)
		}
		if !strings.HasPrefix(searchResponse.Payload.ExactMatch.ID, "1_") {
			return nil, fmt.Errorf("todo")
		}
		avmResponse, err := client.GetAutomatedValuationModel(strings.TrimPrefix(searchResponse.Payload.ExactMatch.ID, "1_"))
		if err != nil {
			return nil, fmt.Errorf("error getting redfin avm: %v", err)
		}

		return []property.Estimate{
			{
				Vendor: property.Redfin,
				Value:  avmResponse.Payload.Root.AVMInfo.PredictedValue,
			},
		}, nil
	}
}

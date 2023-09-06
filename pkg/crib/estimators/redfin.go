package estimators

import (
	"fmt"
	"github.com/peppys/crib/internal/redfin"
	"net/http"
	"strings"
)

func DefaultRedfinEstimator() func(address string) ([]Estimate, error) {
	return NewRedfinEstimator(redfin.NewClient(http.DefaultClient))
}

func NewRedfinEstimator(client *redfin.Client) func(address string) ([]Estimate, error) {
	return func(address string) ([]Estimate, error) {
		searchResponse, err := client.SearchProperties(address)
		if err != nil {
			return nil, fmt.Errorf("error searching redfin properties: %v", err)
		}
		if !strings.HasPrefix(searchResponse.Payload.ExactMatch.ID, "1_") {
			return nil, fmt.Errorf("failed to find redfin property ID")
		}
		avmResponse, err := client.GetAutomatedValuationModel(strings.TrimPrefix(searchResponse.Payload.ExactMatch.ID, "1_"))
		if err != nil {
			return nil, fmt.Errorf("error getting redfin avm: %v", err)
		}

		return []Estimate{
			{
				Provider: Redfin,
				Value:    int64(avmResponse.Payload.Root.AVMInfo.PredictedValue * 100),
			},
		}, nil
	}
}

package estimators

import (
	"fmt"
	"github.com/peppys/crib/internal/services/property"
)

func NewBulkEstimator(estimators ...property.Estimator) func(address string) ([]property.Estimate, error) {
	return func(address string) ([]property.Estimate, error) {
		var estimates []property.Estimate

		// TODO - make async
		for _, estimator := range estimators {
			result, err := estimator(address)
			if err != nil {
				return nil, fmt.Errorf("failed running one of the estimators: %w", err)
			}

			estimates = append(estimates, result...)
		}

		return estimates, nil
	}
}

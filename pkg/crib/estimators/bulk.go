package estimators

import (
	"github.com/pkg/errors"
	"sync"
)

func NewBulkEstimator(estimators ...Estimator) func(address string) ([]Estimate, error) {
	return func(address string) ([]Estimate, error) {
		var wg sync.WaitGroup
		var estimatorErrors []error
		var estimates []Estimate

		for _, estimator := range estimators {
			wg.Add(1)

			go func(estimator Estimator) {
				defer wg.Done()
				result, err := estimator(address)
				if err != nil {
					estimatorErrors = append(estimatorErrors, err)
				} else {
					estimates = append(estimates, result...)
				}
			}(estimator)
		}

		wg.Wait()
		if len(estimatorErrors) > 0 {
			err := errors.New("")
			for _, e := range estimatorErrors {
				err = errors.Wrapf(err, "estimator failed: %s", e.Error())
			}

			return nil, err
		}

		return estimates, nil
	}
}

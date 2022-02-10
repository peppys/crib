package crib

import (
	"fmt"
	"github.com/peppys/crib/pkg/crib/estimators"
)

type Crib struct {
	estimate estimators.Estimator
}

type Option func(c *Crib)

func NewCrib(opts ...Option) *Crib {
	c := &Crib{}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

func WithEstimator(e estimators.Estimator) Option {
	return func(c *Crib) {
		c.estimate = e
	}
}

func WithEstimators(e ...estimators.Estimator) Option {
	return func(c *Crib) {
		c.estimate = estimators.NewBulkEstimator(e...)
	}
}

func (c *Crib) Estimate(address string) ([]estimators.Estimate, error) {
	estimate, err := c.estimate(address)
	if err != nil {
		return nil, fmt.Errorf("error while estimating valuation: %w", err)
	}

	return estimate, nil
}

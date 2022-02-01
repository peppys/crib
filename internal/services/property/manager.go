package property

import (
	"fmt"
)

type Manager struct {
	estimate Estimator
}

type Option func(agent *Manager)

type Vendor string

const (
	Zillow Vendor = "ZILLOW"
	Redfin        = "REDFIN"
)

type Estimate struct {
	Vendor Vendor
	Value  float64
}

type Estimator func(string) ([]Estimate, error)

func NewManager(opts ...Option) *Manager {
	m := &Manager{}

	for _, opt := range opts {
		opt(m)
	}

	return m
}

func WithEstimator(e Estimator) Option {
	return func(manager *Manager) {
		manager.estimate = e
	}
}

func (m *Manager) Valuation(address string) ([]Estimate, error) {
	estimate, err := m.estimate(address)
	if err != nil {
		return nil, fmt.Errorf("error while estimating valuation: %w", err)
	}

	return estimate, nil
}

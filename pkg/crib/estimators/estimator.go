package estimators

type Provider string

const (
	Zillow Provider = "zillow"
	Redfin          = "redfin"
)

type Estimate struct {
	Provider Provider `json:"provider"`
	Value    int64    `json:"value"`
}

type Estimator func(string) ([]Estimate, error)

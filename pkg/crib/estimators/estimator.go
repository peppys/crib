package estimators

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

# crib
nice crib

## Installation
Navigate to the [latest release](https://github.com/peppys/crib/releases/latest) and download the binary for your OS.  
Assuming it's mac, here's how to install v2.0.1:
```shell
$ curl -L https://github.com/peppys/crib/releases/download/2.0.1/crib_Darwin_arm64.tar.gz --output crib.tar.gz 

$ tar -xzvf crib.tar.gz
```

## CLI Usage
```shell
$ ./crib value --help
Checks the estimated valuation of your property

Usage:
  cli value [flags]

Flags:
  -a, --address string   Address of your crib
  -f, --format string    Format of output (table/json/csv) (default "table")
  -h, --help             help for value
```

### Machine readable formats
> note: machine readable formats display estimates in `cents` / `pennies`

#### csv
```shell
crib value -a '1443 devlin dr, los angeles, ca' --format csv
 ██████ ██████  ██ ██████  
██      ██   ██ ██ ██   ██ 
██      ██████  ██ ██████  
██      ██   ██ ██ ██   ██ 
 ██████ ██   ██ ██ ██████  

Provider,Estimate                                                                                                                                                                                                                       
redfin,783541324
zillow,756970000
```

#### json
```shell
crib value -a '1443 devlin dr, los angeles, ca' --format json
 ██████ ██████  ██ ██████  
██      ██   ██ ██ ██   ██ 
██      ██████  ██ ██████  
██      ██   ██ ██ ██   ██ 
 ██████ ██   ██ ██ ██████  

[                                                                                                                                                                                                                          
  {
    "provider": "redfin",
    "value": 783541324
  },
  {
    "provider": "zillow",
    "value": 756970000
  }
]
```

json format can be paired with [jq](https://github.com/jqlang/jq) to calculate average value:
```shell
crib value -a '1443 devlin dr, los angeles, ca' --format json | jq 'map(.value) | add/length'
 ██████ ██████  ██ ██████  
██      ██   ██ ██ ██   ██ 
██      ██████  ██ ██████  
██      ██   ██ ██ ██   ██ 
 ██████ ██   ██ ██ ██████  

770255662                         
```

## Library Usage
```go
package main

import (
	"github.com/peppys/crib/pkg/crib"
	"github.com/peppys/crib/pkg/crib/estimators"
	"log"
)

func main() {
	c := crib.New(
		crib.WithEstimators(
			estimators.DefaultZillowEstimator(),
			estimators.DefaultRedfinEstimator(),
		),
	)

	estimates, err := c.Estimate("1443 delvin dr los angeles, ca")
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("found estimates: %+v", estimates)
}
```

## CLI Demo  
![demo](demo.gif)  

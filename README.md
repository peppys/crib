# crib
nice crib

## Install
Navigate to the [latest release](https://github.com/peppys/crib/releases/latest) and download the binary for your OS.  
Assuming it's mac, here's how to install v1.0.1:
```shell
$ curl -L https://github.com/peppys/crib/releases/download/1.0.1/crib_1.0.1_Darwin_arm64.tar.gz --output crib.tar.gz 

$ tar -xzvf crib.tar.gz
```

## Usage
```shell
$ ./crib value --help
Checks the estimated valuation of your property

Usage:
  cli value [flags]

Flags:
  -a, --address string   Address of your crib
  -h, --help             help for value
```
```shell
$ ./crib value -a '123rd Fake St'
 ██████ ██████  ██ ██████
██      ██   ██ ██ ██   ██
██      ██████  ██ ██████
██      ██   ██ ██ ██   ██
 ██████ ██   ██ ██ ██████

Vendor | Estimate
zillow | $748,300.00
redfin | $712,550.92
```

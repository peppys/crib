package main

import "fmt"

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func getVersion() string {
	if version == "dev" {
		return "dev"
	}

	return fmt.Sprintf("%s-%s", version, commit)
}

package main

import (
	"fmt"
)

var (
	// Version is populated at compile time by govvv from ./VERSION
	Version string

	// GitCommit is populated at compile time by govvv.
	GitCommit string

	GitSummary string

	GitBranch string
	BuildDate string

	// GitState is populated at compile time by govvv.
	GitState     string
	BuildTime    string
	BuildMachine string
)

func versionString() string {
	return fmt.Sprintf("%s@%s, build-time:%s, build-machine:%s",
		GitSummary, GitBranch, BuildTime, BuildMachine)
}

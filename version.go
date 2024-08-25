package main

import (
	"fmt"
)

var (
	// Version is populated at compile time by govvv from ./VERSION
	Version string

	CompileInfo string
)

func versionString() string {
	return fmt.Sprintf("compile-info:%s", CompileInfo)
}

package main

import "fmt"

const version = "0.2.0"

var commit = "unknown"

func versionTemplate() string {
	return fmt.Sprintf("version=%s commit=%s\n", version, commit)
}

package build

import (
	"fmt"
)

type BuildInfo struct {
	buildVersion string
	buildDate    string
	buildCommit  string
}

func Create() BuildInfo {
	return BuildInfo{
		buildVersion: "N/A",
		buildDate:    "N/A",
		buildCommit:  "N/A",
	}
}

func (b BuildInfo) Print() {
	fmt.Printf(`
		Build version: %s
		Build date: %s
		Build commit: %s
	`, b.buildVersion, b.buildDate, b.buildCommit)
}

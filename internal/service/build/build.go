package build

import (
	"fmt"
	"strings"
)

type BuildInfo struct {
	buildVersion string
	buildDate    string
	buildCommit  string
}

func (b BuildInfo) Print() {
	template := `
		Build version: <buildVersion>
		Build date: <buildDate>
		Build commit: <buildCommit>
	`

	if b.buildVersion == "" {
		b.buildVersion = "N/A"
	}
	if b.buildDate == "" {
		b.buildDate = "N/A"
	}
	if b.buildCommit == "" {
		b.buildCommit = "N/A"
	}

	template = strings.Replace(template, "<buildVersion>", b.buildVersion, 1)
	template = strings.Replace(template, "<buildDate>", b.buildDate, 1)
	template = strings.Replace(template, "<buildCommit>", b.buildCommit, 1)

	fmt.Println(template)
}

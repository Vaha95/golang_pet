package build

import (
	"fmt"
	"strings"
)

var tmp string = `
	Build version: <buildVersion>
	Build date: <buildDate>
	Build commit: <buildCommit>
`

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
	template := tmp

	template = strings.Replace(template, "<buildVersion>", b.buildVersion, 1)
	template = strings.Replace(template, "<buildDate>", b.buildDate, 1)
	template = strings.Replace(template, "<buildCommit>", b.buildCommit, 1)

	fmt.Println(template)
}

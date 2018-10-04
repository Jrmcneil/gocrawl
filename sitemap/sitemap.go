package sitemap

import (
	"gocrawl/job"
	"strings"
)

type Sitemap struct {
	Result chan string
}

func (sitemap *Sitemap) Build(root job.Job) {
	var sb strings.Builder
	var indentSb strings.Builder
	indentSb.WriteString("")
	buildMap(root, &sb, "", true)
	sitemap.Result <- sb.String()
}

func buildMap(node job.Job, sb *strings.Builder, indent string, last bool) {
	sb.WriteString(indent)
	if last {
		sb.WriteString("\\-")
		indent += "  "
	} else {
		sb.WriteString("|-")
		indent += "| "
	}

	sb.WriteString(node.Address() + "\n")

	<-node.Ready()
    links := node.Links()
    node.ResetLinks()

    length := len(links)
	for i := 0; i < length; i++ {
		buildMap(links[i], sb, indent, i == length - 1)
	}
}

package sitemap

import (
	"gocrawl/job"
	"strings"
)

type Sitemap struct {
	Result chan string
	quit chan bool
}

func (sitemap *Sitemap) Build(root job.Job) {
	var sb strings.Builder
	var indentSb strings.Builder
	indentSb.WriteString("")
	sitemap.buildMap(root, &sb, "", true)
	sitemap.Result <- sb.String()
}

func (sitemap *Sitemap) buildMap(node job.Job, sb *strings.Builder, indent string, last bool) {
		sb.WriteString(indent)
		if last {
			sb.WriteString("\\-")
			indent += "  "
		} else {
			sb.WriteString("|-")
			indent += "| "
		}

		sb.WriteString(node.Address() + "\n")


		<- node.Ready()
		links := node.Links()
		node.ResetLinks()

		length := len(links)
		for i := 0; i < length; i++ {
			sitemap.buildMap(links[i], sb, indent, i == length - 1)
		}
}

func (sitemap *Sitemap) Stop() {
	sitemap.quit <- true
}

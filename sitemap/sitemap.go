package sitemap

import (
    "gocrawl/job"
    "strings"
)

type Sitemap struct {
    Result chan string
}

func (sitemap *Sitemap) Build(root job.Job) {
    go func() {
        var sb strings.Builder
        var indentSb strings.Builder
        indentSb.WriteString("")
        buildMap(root, &sb, &indentSb, true)
        sitemap.Result <- sb.String()
    } ()
}

func buildMap(node job.Job, sb *strings.Builder, indentSb *strings.Builder, last bool) {
    sb.WriteString(indentSb.String())

    if last {
        sb.WriteString("\\-")
        indentSb.WriteString("  ")
    } else {
        sb.WriteString("|-")
        indentSb.WriteString("| ")
    }
    sb.WriteString(node.Address() + "\n")

    <- node.Ready()

    for i := 0; i < len(node.Links()); i++ {
        buildMap(node.Links()[i], sb, indentSb, i == len(node.Links()) - 1)
    }
}
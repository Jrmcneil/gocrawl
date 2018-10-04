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
        buildMap(root, &sb, "", true)
        sitemap.Result <- sb.String()
    } ()
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

   <- node.Ready()

   for i := 0; i < len(node.Links()); i++ {
       buildMap(node.Links()[i], sb, indent, i == len(node.Links()) - 1)
   }
}

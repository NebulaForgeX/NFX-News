package htmlreport

import (
	"html"
	"strconv"
	"strings"
	"time"
)

type Item struct {
	Title      string
	URL        string
	SourceID   string
	Group      string
	IsNew      bool
}

func Render(title, mode string, generatedAt time.Time, items []Item) string {
	var b strings.Builder
	b.WriteString("<!DOCTYPE html><html lang=\"zh\"><head><meta charset=\"utf-8\">")
	b.WriteString("<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\">")
	b.WriteString("<title>")
	b.WriteString(html.EscapeString(title))
	b.WriteString("</title><style>")
	b.WriteString("body{font-family:ui-sans-serif,system-ui,sans-serif;margin:24px;background:#0b0d10;color:#e8eaed}")
	b.WriteString("h1{font-size:1.4rem} .meta{opacity:.7;margin-bottom:1.5rem}")
	b.WriteString("a{color:#8ab4f8;text-decoration:none} li{margin:.4rem 0}")
	b.WriteString(".new{color:#fbbc04} .group{opacity:.6;font-size:.85rem}")
	b.WriteString("</style></head><body>")
	b.WriteString("<h1>")
	b.WriteString(html.EscapeString(title))
	b.WriteString("</h1><div class=\"meta\">mode=")
	b.WriteString(html.EscapeString(mode))
	b.WriteString(" · ")
	b.WriteString(html.EscapeString(generatedAt.Format(time.RFC3339)))
	b.WriteString(" · ")
	b.WriteString(strconv.Itoa(len(items)))
	b.WriteString(" items</div><ol>")
	for _, it := range items {
		b.WriteString("<li>")
		if it.IsNew {
			b.WriteString("<span class=\"new\">NEW </span>")
		}
		b.WriteString("<a href=\"")
		b.WriteString(html.EscapeString(it.URL))
		b.WriteString("\">")
		b.WriteString(html.EscapeString(it.Title))
		b.WriteString("</a> <span class=\"group\">")
		b.WriteString(html.EscapeString(it.Group))
		if it.SourceID != "" {
			b.WriteString(" · ")
			b.WriteString(html.EscapeString(it.SourceID))
		}
		b.WriteString("</span></li>")
	}
	b.WriteString("</ol></body></html>")
	return b.String()
}

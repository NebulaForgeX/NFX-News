package channels

import (
	"encoding/json"
	"fmt"
	"strings"

	"nfxnews/events"
)

type Item struct {
	Title    string
	URL      string
	SourceID string
	Group    string
	IsNew    bool
}

func ParseItems(payloadJSON string) []Item {
	if payloadJSON == "" {
		return nil
	}
	var raw map[string]any
	if err := json.Unmarshal([]byte(payloadJSON), &raw); err != nil {
		return nil
	}
	list, _ := raw["items"].([]any)
	out := make([]Item, 0, len(list))
	for _, it := range list {
		row, ok := it.(map[string]any)
		if !ok {
			continue
		}
		item := Item{}
		item.Title, _ = row["title"].(string)
		item.URL, _ = row["url"].(string)
		item.SourceID, _ = row["source_id"].(string)
		item.Group, _ = row["group"].(string)
		item.IsNew, _ = row["is_new"].(bool)
		if item.Title == "" {
			continue
		}
		out = append(out, item)
	}
	return out
}

func ReportBody(ev events.ReportGeneratedEvent) string {
	items := ParseItems(ev.Payload)
	var b strings.Builder
	title := ev.Title
	if title == "" {
		title = "News report"
	}
	fmt.Fprintf(&b, "📊 %s\n", title)
	fmt.Fprintf(&b, "mode=%s · items=%d\n\n", ev.Mode, ev.ItemCount)
	if len(items) == 0 {
		b.WriteString("(no matched titles)\n")
		return b.String()
	}
	currentGroup := ""
	for _, it := range items {
		if it.Group != "" && it.Group != currentGroup {
			currentGroup = it.Group
			fmt.Fprintf(&b, "【%s】\n", currentGroup)
		}
		marker := "•"
		if it.IsNew {
			marker = "🆕"
		}
		fmt.Fprintf(&b, "%s %s\n", marker, it.Title)
		if it.URL != "" {
			fmt.Fprintf(&b, "  %s\n", it.URL)
		}
		if it.SourceID != "" {
			fmt.Fprintf(&b, "  %s\n", it.SourceID)
		}
		b.WriteByte('\n')
	}
	return b.String()
}

func SplitBatches(text string, maxBytes int) []string {
	if maxBytes <= 0 {
		maxBytes = 4000
	}
	if len([]byte(text)) <= maxBytes {
		return []string{text}
	}
	lines := strings.Split(text, "\n")
	var batches []string
	var cur strings.Builder
	for _, line := range lines {
		candidate := line + "\n"
		if cur.Len()+len(candidate) > maxBytes && cur.Len() > 0 {
			batches = append(batches, strings.TrimRight(cur.String(), "\n"))
			cur.Reset()
		}
		if len(candidate) > maxBytes {
			batches = append(batches, line)
			continue
		}
		cur.WriteString(candidate)
	}
	if cur.Len() > 0 {
		batches = append(batches, strings.TrimRight(cur.String(), "\n"))
	}
	if len(batches) == 0 {
		return []string{text}
	}
	return batches
}

func MaxBytes(kind string, cfg map[string]any) int {
	if n := cfgInt(cfg, "batch_size"); n > 0 {
		return n
	}
	switch kind {
	case "dingtalk":
		return 20000
	case "feishu":
		return 29000
	case "bark":
		return 3600
	case "ntfy":
		return 3800
	case "slack":
		return 4000
	default:
		return 4000
	}
}

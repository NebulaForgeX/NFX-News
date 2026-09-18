package frequency

import (
	"os"
	"strconv"
	"strings"
	"unicode"
)

type Group struct {
	Required []string
	Normal   []string
	GroupKey string
	MaxCount int
	Name     string
}

type Dictionary struct {
	Groups  []Group
	Exclude []string
}

func LoadFile(path string) (Dictionary, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Dictionary{}, err
	}
	return Parse(string(data)), nil
}

func Parse(content string) Dictionary {
	blocks := strings.Split(content, "\n\n")
	var dict Dictionary
	for _, block := range blocks {
		lines := strings.Split(block, "\n")
		var required, normal []string
		maxCount := 0
		name := ""
		for _, raw := range lines {
			word := strings.TrimSpace(raw)
			if word == "" {
				continue
			}
			if strings.HasPrefix(word, "#") {
				if name == "" {
					name = strings.TrimSpace(strings.TrimPrefix(word, "#"))
					name = strings.Trim(name, "= ")
				}
				continue
			}
			switch {
			case strings.HasPrefix(word, "@"):
				if n, err := strconv.Atoi(word[1:]); err == nil && n > 0 {
					maxCount = n
				}
			case strings.HasPrefix(word, "!"):
				dict.Exclude = append(dict.Exclude, strings.TrimPrefix(word, "!"))
			case strings.HasPrefix(word, "+"):
				required = append(required, strings.TrimPrefix(word, "+"))
			default:
				normal = append(normal, word)
			}
		}
		if len(required) == 0 && len(normal) == 0 {
			continue
		}
		key := strings.Join(normal, " ")
		if key == "" {
			key = strings.Join(required, " ")
		}
		dict.Groups = append(dict.Groups, Group{Required: required, Normal: normal, GroupKey: key, MaxCount: maxCount, Name: name})
	}
	return dict
}

func MatchTitle(title string, dict Dictionary) (group Group, ok bool) {
	t := fold(title)
	for _, ex := range dict.Exclude {
		if strings.Contains(t, fold(ex)) {
			return Group{}, false
		}
	}
	if len(dict.Groups) == 0 {
		return Group{GroupKey: "all"}, true
	}
	for _, g := range dict.Groups {
		if matchGroup(t, g) {
			return g, true
		}
	}
	return Group{}, false
}

func matchGroup(title string, g Group) bool {
	for _, w := range g.Required {
		if !strings.Contains(title, fold(w)) {
			return false
		}
	}
	if len(g.Normal) == 0 {
		return len(g.Required) > 0
	}
	for _, w := range g.Normal {
		if strings.Contains(title, fold(w)) {
			return true
		}
	}
	return len(g.Required) > 0
}

func fold(s string) string {
	return strings.Map(func(r rune) rune {
		return unicode.ToLower(r)
	}, s)
}

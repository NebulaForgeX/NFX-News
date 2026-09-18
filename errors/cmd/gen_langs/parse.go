package main

import (
	"os"
	"path/filepath"
	"regexp"
)

var (
	blockCommentRe = regexp.MustCompile(`/\*\*?([\s\S]*?)\*/`)
	codeLineRe     = regexp.MustCompile(`(?m)^!\s*([A-Z0-9_]+)\s*$`)
	langLineRe     = regexp.MustCompile(`(?m)^\*\s*([a-z]{2,})\s*<\s*([^>]*)\s*>?\s*$`)
	prodLangLineRe = regexp.MustCompile(`(?m)^\*p\*\s*([a-z]{2,})\s*<\s*([^>]*)\s*>?\s*$`)
)

// Collect walks srcPath and parses all .go files, returning lang -> code -> text.
// When prod is true, *p*lang<...> lines override the default *lang<...> text for that locale.
func Collect(srcPath string, prod bool) (map[string]map[string]string, error) {
	byLang := make(map[string]map[string]string)
	err := filepath.Walk(srcPath, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || filepath.Ext(path) != ".go" {
			return err
		}
		collectFromFile(path, prod, byLang)
		return nil
	})
	return byLang, err
}

func collectFromFile(path string, prod bool, byLang map[string]map[string]string) {
	data, err := os.ReadFile(path)
	if err != nil {
		return
	}
	blocks := blockCommentRe.FindAllStringSubmatch(string(data), -1)
	for _, m := range blocks {
		if len(m) < 2 {
			continue
		}
		block := m[1]
		var currentCode string
		devByLang := make(map[string]string)
		prodByLang := make(map[string]string)
		flush := func() {
			if currentCode == "" {
				return
			}
			mergeCodeTranslations(currentCode, devByLang, prodByLang, prod, byLang)
			currentCode = ""
			devByLang = make(map[string]string)
			prodByLang = make(map[string]string)
		}
		for _, line := range splitLines(block) {
			if codeLineRe.MatchString(line) {
				flush()
				sm := codeLineRe.FindStringSubmatch(line)
				if len(sm) >= 2 {
					currentCode = sm[1]
				}
				continue
			}
			if currentCode == "" {
				continue
			}
			if prodLangLineRe.MatchString(line) {
				sm := prodLangLineRe.FindStringSubmatch(line)
				if len(sm) >= 3 {
					prodByLang[sm[1]] = sm[2]
				}
				continue
			}
			if langLineRe.MatchString(line) {
				sm := langLineRe.FindStringSubmatch(line)
				if len(sm) >= 3 {
					devByLang[sm[1]] = sm[2]
				}
			}
		}
		flush()
	}
}

func mergeCodeTranslations(
	code string,
	devByLang map[string]string,
	prodByLang map[string]string,
	prod bool,
	byLang map[string]map[string]string,
) {
	langs := make(map[string]struct{}, len(devByLang)+len(prodByLang))
	for lang := range devByLang {
		langs[lang] = struct{}{}
	}
	for lang := range prodByLang {
		langs[lang] = struct{}{}
	}

	for lang := range langs {
		text := devByLang[lang]
		if prod {
			if override, ok := prodByLang[lang]; ok && override != "" {
				text = override
			}
		}
		if text == "" {
			continue
		}
		if byLang[lang] == nil {
			byLang[lang] = make(map[string]string)
		}
		byLang[lang][code] = text
	}
}

func splitLines(s string) []string {
	var out []string
	start := 0
	for i := 0; i <= len(s); i++ {
		if i == len(s) || s[i] == '\n' {
			out = append(out, trimSpace(s[start:i]))
			start = i + 1
		}
	}
	return out
}

func trimSpace(s string) string {
	i := 0
	for i < len(s) && (s[i] == ' ' || s[i] == '\t' || s[i] == '\r') {
		i++
	}
	j := len(s)
	for j > i && (s[j-1] == ' ' || s[j-1] == '\t' || s[j-1] == '\r') {
		j--
	}
	return s[i:j]
}

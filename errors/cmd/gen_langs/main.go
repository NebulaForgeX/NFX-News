// Run from repo root:
//
//	go run ./errors/cmd/gen_langs/ ./errors          # dev (default)
//	go run ./errors/cmd/gen_langs/ ./errors --prod   # production-safe messages (*p* overrides)
package main

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	srcDir   = "src"
	langsDir = "langs"
)

func main() {
	useProd, baseDir := parseArgs(os.Args[1:])

	srcPath := filepath.Join(baseDir, srcDir)
	outPath := filepath.Join(baseDir, langsDir)

	if err := os.RemoveAll(outPath); err != nil {
		fmt.Fprintf(os.Stderr, "clean %s: %v\n", outPath, err)
		os.Exit(1)
	}
	if err := os.MkdirAll(outPath, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "mkdir %s: %v\n", outPath, err)
		os.Exit(1)
	}

	byLang, err := Collect(srcPath, useProd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "walk: %v\n", err)
		os.Exit(1)
	}

	if err := WriteLangFiles(outPath, byLang); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	mode := "dev"
	if useProd {
		mode = "prod"
	}
	fmt.Printf("generated langs (%s mode)\n", mode)
}

// parseArgs accepts optional --prod / --dev flags in any order and an optional base directory (defaults to errors/).
func parseArgs(args []string) (prod bool, baseDir string) {
	baseDir = "."
	prodSet := false
	for _, arg := range args {
		switch arg {
		case "--prod":
			prod = true
			prodSet = true
		case "--dev":
			prod = false
			prodSet = true
		default:
			baseDir = arg
		}
	}
	if !prodSet {
		prod = false
	}
	if _, err := os.Stat(filepath.Join(baseDir, srcDir)); os.IsNotExist(err) {
		if baseDir == "." {
			baseDir = "errors"
		}
	}
	return prod, baseDir
}

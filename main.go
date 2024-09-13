package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Find files with extensions `.tsx` and `.js`
func findFiles(root string) ([]string, error) {
	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && (filepath.Ext(path) == ".tsx" || filepath.Ext(path) == ".js") {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

// Parse file and extract exported React components using a regex (simplified)
func extractExports(filePath string) ([]string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var exports []string
	// Simplified regex to find exported components
	// re := regexp.MustCompile(`export\s+(default|const|class|function)\s+([A-Z].+)`)
	re := regexp.MustCompile(`export\s+default\s+([A-Z]\w+|(function|const)\s+([A-Z]+\w+))`)
	matches := re.FindAllStringSubmatch(string(content), -1)

	for _, match := range matches {
		if len(match) > 1 {
			len := 0
			for _, sub := range match {
				if sub != "" {
					len++
				}
			}
			exports = append(exports, match[len-1])
		}
	}

	return exports, nil
}

func main() {
	rootDir := os.Args[1]

	files, err := findFiles(rootDir)
	if err != nil {
		fmt.Println("Error finding files:", err)
		return
	}

	for _, file := range files {
		components, err := extractExports(file)
		if err != nil {
			fmt.Println("Error parsing file:", file, err)
			continue
		}
		if len(components) == 0 {
			continue
		}
		fmt.Printf(
			"%s: %v\n",
			strings.ReplaceAll(file, "/Users/souria-sakyhocquenghem/Documents/qeeps/front-end-app/", ""),
			components,
		)
	}
}

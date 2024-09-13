package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
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
	re := regexp.MustCompile(`export\s+(?:default\s+)?([A-Z]\w+|(function|const)\s+([A-Z]+\w+))`)
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

func mapExportedComponents(files []string) map[string][]string {
	m := make(map[string][]string, len(files))
	for _, file := range files {
		components, err := extractExports(file)
		if err != nil {
			fmt.Println("Error parsing file:", file, err)
			continue
		}
		if len(components) == 0 {
			continue
		}
		m[file] = components
	}
	return m
}

func extractImports(filePath string) ([]string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var imports []string
	// Simplified regex to find imported components
	re := regexp.MustCompile(`import\s+(?:\{\s+)?([A-Z]+\w+|[A-Z]+\w+)(?:\s+\})?\s+from\s+\"(\@/|\.+)[a-zA-Z/_]+\"`)
	matches := re.FindAllStringSubmatch(string(content), -1)

	for _, match := range matches {
		if len(match) > 1 {
			imports = append(imports, match[1])
		}
	}

	return imports, nil
}

func mapImportedComponents(files []string) map[string][]string {
	m := make(map[string][]string, len(files))
	for _, file := range files {
		imports, err := extractImports(file)
		if err != nil {
			fmt.Println("Error parsing file:", file, err)
			continue
		}
		if len(imports) != 0 {
			m[file] = imports
		}
	}
	return m
}

func main() {
	if len(os.Args) != 3 || (os.Args[1] != "import" && os.Args[1] != "export") {
		log.Fatal("usage: <import|export> <path>\n")
	}
	rootDir := os.Args[2]

	files, err := findFiles(rootDir)
	if err != nil {
		fmt.Println("Error finding files:", err)
		return
	}
	var m map[string][]string
	if os.Args[1] == "import" {
		m = mapImportedComponents(files)
	} else if os.Args[1] == "export" {
		m = mapExportedComponents(files)
	}
	for file, components := range m {
		if len(components) > 0 {
			fmt.Printf(
				"%s: %v\n",
				file,
				components,
			)
		}
	}
}

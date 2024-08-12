package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// createDirStructure creates directories and files from the input structure
func createDirStructure(baseDir, structure string, dryRun bool) error {
	scanner := bufio.NewScanner(strings.NewReader(structure))
	//currentDir := baseDir
	fileStack := []string{}
	previousContext := 0
	previousDirContinue := false
	for scanner.Scan() {
		line := scanner.Text()
		item := strings.TrimSpace(line)
		localPrevDirCon := false
		if strings.Contains(line, "├──") || strings.Contains(line, "└──") {
			currentContext := strings.Count(line, "│")
			if strings.Contains(line, "├──") {
				_, item, _ = strings.Cut(line, "├──")
				localPrevDirCon = true
				fileStack = fileStack[:currentContext+1]
			}
			if strings.Contains(line, "└──") {
				_, item, _ = strings.Cut(line, "└──")
				if currentContext < previousContext || (currentContext == previousContext && previousDirContinue) {
					fileStack = fileStack[:currentContext+1]
				}
			}
			previousDirContinue = localPrevDirCon
			item = strings.TrimSpace(item)
			fileStack = append(fileStack, item)
			previousContext = currentContext
			path := strings.Join(fileStack, "")
			path = filepath.Join(baseDir, path)
			// Directory or file
			if strings.HasSuffix(line, "/") {
				// It's a directory
				fmt.Println(path)
				if !dryRun {
					if err := os.MkdirAll(path, 0755); err != nil {
						return err
					}
				}
			} else {
				// It's a file
				if !dryRun {
					dir := filepath.Dir(path)
					if err := os.MkdirAll(dir, 0755); err != nil {
						return err
					}
					file, err := os.Create(path)
					if err != nil {
						return err
					}
					if err = file.Close(); err != nil {
						return err
					}
				}
				fmt.Println(path)
			}
		} else {
			fileStack = append(fileStack, item)
		}
	}

	return nil
}

func main() {
	// Define command-line flags
	structureFile := flag.String("file", "", "Path to the file containing the directory structure")
	targetDir := flag.String("target", ".", "Target directory where the structure will be created")
	dryRun := flag.Bool("dry-run", true, "By default do a dry run rather than creating directory structure")
	flag.Parse()

	// Check if the structure file path is provided
	if *structureFile == "" {
		fmt.Println("Error: Please provide a path to the structure file using the -file flag.")
		os.Exit(1)
	}

	// Read the structure from the file
	content, err := os.ReadFile(*structureFile)
	if err != nil {
		fmt.Println("Error reading the structure file:", err)
		os.Exit(1)
	}

	// Create the directory structure
	err = createDirStructure(*targetDir, string(content), *dryRun)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Directory structure created successfully.")
	}
}

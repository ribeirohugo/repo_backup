package backup

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

func LoadFromFile() ([]string, error) {
	var repos []string

	// Define flag
	filePath := flag.String("f", "", "Path to file containing repository URLs (one per line)")
	flag.Parse()

	if *filePath == "" {
		fmt.Println("Usage: go run main.go -f repos.txt")
		os.Exit(1)
	}

	// Open the file
	file, err := os.Open(*filePath)
	if err != nil {
		return repos, fmt.Errorf("Error opening file: %v\n", err)
	}
	defer file.Close()

	// Read each line and clone
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		repo := strings.TrimSpace(scanner.Text())
		if repo == "" || strings.HasPrefix(repo, "#") {
			continue // skip empty lines and comments
		}

		repos = append(repos, repo)
	}

	return repos, nil
}

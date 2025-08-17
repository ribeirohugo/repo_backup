package backup

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

func Load() ([]string, error) {
	// Define flag
	filePath := flag.String("f", "", "Path to file containing repository URLs (one per line)")
	flag.Parse()

	if *filePath != "" {
		return LoadFromFile(*filePath)
	}

	return LoadFromArgs()
}

func LoadFromFile(filePath string) ([]string, error) {
	var repos []string

	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return repos, fmt.Errorf("error opening file: %v", err)
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

func LoadFromArgs() ([]string, error) {
	var repos []string

	for _, repoURL := range os.Args[1:] {
		if repoURL != "" {
			repos = append(repos, repoURL)
		}
	}

	return repos, nil
}

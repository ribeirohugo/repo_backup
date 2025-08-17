package main

import (
	"fmt"
	"os"

	"github.com/ribeirohugo/repo_backup/pkg/backup"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <repo1> <repo2> ...")
		return
	}

	var clonedDirs []string
	for _, repoURL := range os.Args[1:] {
		dir, err := backup.Clone(repoURL)
		if err != nil {
			fmt.Printf("Error cloning %s: %v\n", repoURL, err)
			continue
		}
		clonedDirs = append(clonedDirs, dir)
	}

	if len(clonedDirs) == 0 {
		fmt.Println("No repositories cloned successfully. Exiting.")
		return
	}

	zipName, err := backup.Zip(clonedDirs)
	if err != nil {
		fmt.Printf("Error creating zip: %v\n", err)
		return
	}
	fmt.Printf("Repositories zipped into %s\n", zipName)

	if err := backup.Remove(clonedDirs); err != nil {
		fmt.Printf("Error removing repos: %v\n", err)
		return
	}
	fmt.Println("Cleanup complete.")
}

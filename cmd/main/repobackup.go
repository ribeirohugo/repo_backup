package main

import (
	"fmt"
	"os"

	"github.com/ribeirohugo/repo_backup/pkg/backup"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go [-f file.txt] <repo1> <repo2> ...")
		return
	}

	repos, err := backup.Load()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if len(repos) == 0 {
		fmt.Println("No repositories cloned successfully. Exiting.")
		return
	}

	var cloneDirs []string
	for _, repoURL := range repos {
		dir, err := backup.Clone(repoURL)
		if err != nil {
			fmt.Printf("Error cloning %s: %v\n", repoURL, err)
			continue
		}
		cloneDirs = append(cloneDirs, dir)
	}

	zipName, err := backup.Zip(cloneDirs)
	if err != nil {
		fmt.Printf("Error creating zip: %v\n", err)
		return
	}
	fmt.Printf("Repositories zipped into %s\n", zipName)

	if err := backup.Remove(cloneDirs); err != nil {
		fmt.Printf("Error removing repos: %v\n", err)
		return
	}
	fmt.Println("Cleanup complete.")
}

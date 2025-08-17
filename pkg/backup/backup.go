package backup

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// Clone clones a git repository into the current directory
func Clone(repoURL string) (string, error) {
	parts := strings.Split(repoURL, "/")
	name := strings.TrimSuffix(parts[len(parts)-1], ".git")

	fmt.Printf("Cloning %s...\n", repoURL)
	cmd := exec.Command("git", "clone", repoURL)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return name, nil
}

// Zip compresses given repo directories into a zip file with current date
func Zip(dirs []string) (string, error) {
	date := time.Now().Format("2006-01-02")
	zipName := fmt.Sprintf("%s.zip", date)

	fmt.Printf("Creating archive %s...\n", zipName)
	zipFile, err := os.Create(zipName)
	if err != nil {
		return "", err
	}
	defer zipFile.Close()

	archive := zip.NewWriter(zipFile)
	defer archive.Close()

	for _, dir := range dirs {
		err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}

			relPath, err := filepath.Rel(filepath.Dir(dir), path)
			if err != nil {
				return err
			}

			f, err := archive.Create(relPath)
			if err != nil {
				return err
			}

			fsFile, err := os.Open(path)
			if err != nil {
				return err
			}
			defer fsFile.Close()

			_, err = io.Copy(f, fsFile)
			return err
		})
		if err != nil {
			return "", err
		}
	}
	return zipName, nil
}

// Remove deletes the cloned repo directories
func Remove(dirs []string) error {
	for _, dir := range dirs {
		fmt.Printf("Removing %s...\n", dir)
		if err := os.RemoveAll(dir); err != nil {
			return err
		}
	}
	return nil
}

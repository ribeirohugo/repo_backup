package backup

import (
	"archive/zip"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --- Mock setup ---

var fakeExecCommand = func(command string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestHelperProcess", "--", command}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
	return cmd
}

var fakeExecCommandError = func(command string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestHelperProcessError", "--", command}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS_ERROR=1"}
	return cmd
}

// --- Tests ---

func TestClone(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		old := ExecCommand
		ExecCommand = fakeExecCommand
		defer func() { ExecCommand = old }()

		repo := "https://github.com/example/repo.git"
		name, err := Clone(repo)

		assert.NoError(t, err)
		assert.Equal(t, "repo", name)
	})

	t.Run("failure", func(t *testing.T) {
		old := ExecCommand
		ExecCommand = fakeExecCommandError
		defer func() { ExecCommand = old }()

		repo := "https://github.com/example/repo.git"
		name, err := Clone(repo)

		assert.Error(t, err)
		assert.Empty(t, name)
	})
}

func TestZip(t *testing.T) {
	t.Run("creates zip with files", func(t *testing.T) {
		tmp := t.TempDir()

		// Create a dummy repo dir with a file
		repoDir := filepath.Join(tmp, "repo")
		assert.NoError(t, os.Mkdir(repoDir, 0o755))

		filePath := filepath.Join(repoDir, "file.txt")
		assert.NoError(t, os.WriteFile(filePath, []byte("hello"), 0o644))

		zipName, err := Zip([]string{repoDir})
		assert.NoError(t, err)
		defer func(path string) {
			err := os.RemoveAll(path)
			require.NoError(t, err)
		}(repoDir)

		// Ensure zip file exists
		_, err = os.Stat(zipName)
		require.NoError(t, err)
		defer func() {
			err = os.RemoveAll(zipName)
			require.NoError(t, err)
		}()

		// Verify zip content
		r, err := zip.OpenReader(zipName)
		assert.NoError(t, err)
		defer func(r *zip.ReadCloser) {
			err := r.Close()
			require.NoError(t, err)
		}(r)

		found := false
		for _, f := range r.File {
			if strings.HasSuffix(f.Name, "file.txt") {
				rc, _ := f.Open()
				content, _ := io.ReadAll(rc)
				_ = rc.Close()
				assert.Equal(t, "hello", string(content))
				found = true
			}
		}
		assert.True(t, found, "expected file.txt inside zip")
	})

	t.Run("returns error if dir not exists", func(t *testing.T) {
		_, err := Zip([]string{"does-not-exist"})
		assert.Error(t, err)
	})
}

func TestRemove(t *testing.T) {
	t.Run("removes dirs", func(t *testing.T) {
		tmp := t.TempDir()
		repoDir := filepath.Join(tmp, "repo")
		assert.NoError(t, os.Mkdir(repoDir, 0o755))

		// Ensure dir exists before remove
		_, err := os.Stat(repoDir)
		assert.NoError(t, err)

		// Remove it
		err = Remove([]string{repoDir})
		assert.NoError(t, err)

		_, err = os.Stat(repoDir)
		assert.Error(t, err, "expected dir to be removed")
	})

	t.Run("ignores already removed dirs", func(t *testing.T) {
		err := Remove([]string{"not-exist"})
		assert.NoError(t, err)
	})
}

// --- Fake helper processes for mocks ---

func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	os.Exit(0) // simulate git success
}

func TestHelperProcessError(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS_ERROR") != "1" {
		return
	}
	os.Exit(1) // simulate git failure
}

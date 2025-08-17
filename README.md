# Repository Backup

The *Repository Backup* allows to automatize the backup process of hosted repositories,
by cloning it and compressing into a file.

## 1. How to run
`repobackup` requires a list of repositories that should be inserted or through the arguments, or through a file.

Using command arguments:

```
repobackup <repo1> <repo2>
```

You can run using a file, that holds a list of repositories separated line break.

```
repobackup -f repos.txt
```

### 1.1. On Windows

Example of compilation for Windows x64 b

`$env:GOOS = "windows"; $env:GOARCH = "amd64"; go build -o repobackup.exe cmd/main/repobackup.go`

Example running 

`.\repobackup.exe git@github.com:ribeirohugo/repo_backup.git`

### 1.1. Compiling on Windows

Example of compilation for Windows x64:

`$env:GOOS = "windows"; $env:GOARCH = "amd64"; go build -o repobackup.exe cmd/main/repobackup.go`

Example of compilation for Windows x64 file loader:

`$env:GOOS = "windows"; $env:GOARCH = "amd64"; go build -o repobackup.exe cmd/file/repobackup.go`

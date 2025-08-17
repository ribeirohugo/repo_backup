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

Example of running `repobackup` with file flag:

`.\repobackup -f repos.txt`

Example of compilation for Windows x64:

`$env:GOOS = "windows"; $env:GOARCH = "amd64"; go build -o repobackup.exe cmd/main/repobackup.go`

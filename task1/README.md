# GitHub Repository Info

CLI utility for retrieving information about GitHub repositories. Displays the name, description, star count, fork count, creation date, and repository URL.

## System Installation
Make sure **`$GOPATH/bin`** (usually `~/go/bin`) is added to your `PATH`.

```bash
git clone https://github.com/IliaSotnikov2005/golang-course.git
cd golang-course/task1
```

```bash
go install ./cmd/repoViewer
```

## Build
```bash
git clone https://github.com/IliaSotnikov2005/golang-course.git
cd golang-course/task1
```

```bash
go build ./cmd/repoViewer
```

## Examples
```bash
# Using repository URL
repoViewer https://github.com/golang/go

# Using URL with .git
repoViewer https://github.com/golang/go.git

# Using owner/repo format
repoViewer golang/go

# Using the -repo flag
repoViewer -repo golang/go
```
## Running URL shortener locally

### Prerequisites:
1. Golang
2. `GOPATH` is set
3. Install [dep](https://github.com/golang/dep/)

### Build/Run the URL shortener:
1. Clone the repository.
2. Run `dep ensure` for getting dependencies.
2. Run `go build main.go` which will generate an eexecutable file of the format depending on the host OS.
3. Run the generated executable file:
    1. For windows, run `./main.exe`
    2. For macOS, linux, run `./main`


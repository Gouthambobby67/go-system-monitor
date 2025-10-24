# Build instructions

This project is a Go terminal application. The preferred way to build it locally is using the included build script.

Requirements
- Go 1.24.x (the project declares `go 1.24.2` in `go.mod`)

Build
```bash
# from project root
chmod +x build.sh   # only needed once
./build.sh
```

This runs `go mod tidy` and `go build -o sysmon main.go`. The produced binary is `./sysmon`.

Run
```bash
# show version (non-interactive)
./sysmon -version

# run the interactive TUI
./sysmon
```

Notes
- The TUI uses alternate screen buffer and will capture the terminal. Press `q` or Ctrl+C to exit.
- If you need to cross-compile, set `GOOS` and `GOARCH` environment variables before `go build`.

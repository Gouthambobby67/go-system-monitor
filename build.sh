#!/bin/bash

# Build script for Go System Monitor

echo "Building Go System Monitor..."

# Ensure modules are up to date
go mod tidy

# Build for the current platform
go build -o sysmon main.go

# Make executable
chmod +x sysmon

echo "Build complete! The executable is named 'sysmon'"
echo "Run ./sysmon to start the system monitor"

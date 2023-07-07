#!/bin/bash

echo "Welcome to SSH Pod Log Downloader installation"

# Detect the operating system
case "$(uname -s)" in
    Linux*)
        os="linux"
        ;;
    Darwin*)
        os="macos"
        ;;
    CYGWIN*|MINGW32*|MSYS*|MINGW*)
        os="windows"
        ;;
    *)
        echo "Unsupported operating system"
        exit 1
        ;;
esac

# Build and run the executable based on the operating system
case "$os" in
    linux)
        go build -o cmd/sshpodlog-linux
        ;;
    macos)
        go build -o sshpodlog-macos
        ;;
    windows)
        go build -o sshpodlog-windows.exe
        ;;
esac
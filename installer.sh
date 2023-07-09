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
        executable="sshpodlog-linux"
        go build -o "$executable" cmd/sshpodlog/*
        ;;
    macos)
        executable="sshpodlog-macos"
        go build -o "$executable" cmd/sshpodlog/*
        ;;
    windows)
        executable="sshpodlog-windows.exe"
        go build -o "$executable" cmd/sshpodlog/*
        ;;
esac

echo "Compiled code created"
echo "USAGE: ./$executable -server <ip address> -port(default 22) <port number> -username <username> -key(optional) <file path to the private key>"
echo "optional arguments are port, cluster, key"

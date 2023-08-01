#!/bin/bash

echo "Welcome to SSH Pod Log Downloader installation"

# Prompt the user for the operating system choice
echo "Please choose the target operating system:"
echo "1) Linux"
echo "2) macOS"
echo "3) Windows"
read -p "Enter the number corresponding to your choice: " os_choice

# Define the operating system based on the user's choice
case "$os_choice" in
    1)
        os="linux"
        ;;
    2)
        os="macos"
        ;;
    3)
        os="windows"
        ;;
    *)
        echo "Invalid choice. Supported choices are 1, 2, or 3."
        exit 1
        ;;
esac

# Build and run the executable based on the operating system
case "$os" in
    linux)
        executable="sshpodlog-linux"
        GOOS=linux GOARCH=amd64 go build -o "$executable" cmd/sshpodlog/*
        ;;
    macos)
        executable="sshpodlog-macos"
        GOOS=darwin GOARCH=amd64 go build -o "$executable" cmd/sshpodlog/*
        ;;
    windows)
        executable="sshpodlog-windows.exe"
        GOOS=windows GOARCH=amd64 go build -o "$executable" cmd/sshpodlog/*
        ;;
esac

echo "Compiled code created"
echo "USAGE: ./$executable -server <ip address> -port(default 22) <port number> -username <username> -cluster <cluster-name> -key(optional) <file path to the private key>"
echo "Optional arguments are port, cluster, key"


## SSH Pod Log Downloader

This is a command-line tool written in Go that connects to an SSH server, interacts with Kubernetes using `kubectl`, and downloads logs from a specified pod.

## Installation

1. Clone the repository:

   ```shell
   git clone https://github.com/ipsegs/sshpodlog.git
   ```

2. Navigate to the project directory:
   ```shell
   cd sshpodlog
   ```

3. Run the installation script:
   ```shell
   ./installer.sh
   ```

   The installation script will automatically detect the target operating system for building the executable (Linux, macOS, or Windows). After selecting the operating system, the script will build the executable accordingly.

4. Once the installation is complete, you will find the built executable file in the project directory.

## Another Installation Approach
 Download the Binary from the releases Assets and use right away. 

## Manual Installation Approach
 Confirm the operating system and architecture(amd or arm) you are compiling for.

1. Compiling for Linux system
   ```shell
   GOOS=linux GOARCH=amd64 go build -o sshpodlog-linux cmd/sshpodlog/*
   ```
2. Compiling for MacOS system
   ```shell
   GOOS=darwin GOARCH=amd64 go build -o sshpodlog-macos cmd/sshpodlog/*
   ```
3. Compiling for Windows system.
   ```shell
   GOOS=windows GOARCH=amd64 go build -o sshpodlog-windows.exe cmd/sshpodlog/*
   ```

## Usage

To use the SSH Pod Log Downloader, follow these steps:

1. Open a terminal or command prompt.

2. Navigate to the directory where the SSH Pod Log Downloader executable is located.

3. Run the executable with the desired command-line arguments. For example:
   ```shell
   sshpodlog-binary -server 192.168.1.100 -port 22 -username myuser -cluster production -key ~/.ssh/id_rsa
   ```

   Replace the arguments with the appropriate values for your SSH server configuration.

## Contributing

Contributions are welcome! If you find any issues or have suggestions for improvement, please open an issue or submit a pull request.

## License

This project is licensed under the [MIT License](LICENSE).
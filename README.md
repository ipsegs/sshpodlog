
# SSH Pod Log Downloader

This is a command-line tool written in Go that connects to an SSH server, interacts with Kubernetes using `kubectl`, and downloads logs from a specified pod.

## Binary Download to start using right away
Download the Binary from the Releases Assets to start using right away.

## Installation For Each Operating System using Bash script
**Note:** Go must be installed on the system.

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

   The installation script will ask for the target operating system for building the executable (Linux, macOS, or Windows). After selecting the operating system, the script will build the executable accordingly.

4. Once the installation is complete, you will find the built executable file in the project directory.

## Manual Installation Approach
 Confirm the operating system and architecture(amd or arm) you are compiling for and go has to be installed

1. Compiling for Linux system
   ```shell
   GOOS=linux GOARCH=amd64 go build ./...
   ```
2. Compiling for MacOS system
   ```shell
   GOOS=darwin GOARCH=amd64 go build ./...
   ```
3. Compiling for Windows system.
   ```shell
   GOOS=windows GOARCH=amd64 go build ./...
   ```

## Usage

To use the SSH Pod Log Downloader, follow these steps:

1. Open a terminal or command prompt.

2. Navigate to the directory where the SSH Pod Log Downloader executable is located.

3. Run the executable with the desired command-line arguments. For example:
   ```shell
   sshpodlog-windows.exe --server 192.168.1.100 --port 22 --username myuser --cluster production --key ~/.ssh/id_rsa
   ```

   Replace the arguments with the appropriate values for your SSH server configuration.

4. Alternatively, you can use a config file to load configuration settings. Provide the `-f` or `--from-file` flag followed by the path to the config file, the config file extension can be (yaml,json,toml e.t.c.). For example:
   ```shell
   sshpodlog --from-file config.yaml
   ```

   The following flags can be used:

   - `--server` (`-s`): The SSH server address.
   - `--username` (`-u`): The SSH username.
   - `--cluster` (`-c`): The Kubernetes context switch.
   - `--port` (`-p`): The SSH port.
   - `--key` (`-k`): The path to the SSH private key file.
   - `--from-file` (`-f`): The configuration file

## Contributing

Contributions are welcome! If you find any issues or have suggestions for improvement, please open an issue or submit a pull request.

## License

This project is licensed under the [MIT License](LICENSE).
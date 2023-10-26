# SSH Pod Log Downloader

This Go command-line tool, powered by Cobra and Viper, has been enhanced to provide new features and a more streamlined user experience. The codebase has been refactored to make it easier to configure and use. Now, all flags can be set either from a configuration file or directly on the command line. Additionally, it introduces different subcommands for various functionalities.

## Installation

### Binary Download

You can download the binary executable directly from the Releases Assets to start using it immediately.

### Manual Installation

If you have Go installed, you can also install SSH Pod Log Downloader using the following command:

```shell
go install github.com/ipsegs/sshpodlog@latest
```

This command will fetch and install the latest version of SSH Pod Log Downloader directly into your Go binary directory, making it accessible from your command line.

### Building from Source

If you prefer building from source, follow these steps:

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

   The script will prompt you to select the target operating system (Linux, macOS, or Windows) for building the executable. After making your selection, the script will build the executable accordingly. Once the installation is complete, you will find the built executable file in the project directory.

4. If you prefer to manually build the executable, ensure you have Go installed and follow these steps based on your target operating system:

   - Compiling for Linux:

     ```shell
     GOOS=linux GOARCH=amd64 go build ./...
     ```

   - Compiling for macOS:

     ```shell
     GOOS=darwin GOARCH=amd64 go build ./...
     ```

   - Compiling for Windows:

     ```shell
     GOOS=windows GOARCH=amd64 go build ./...
     ```

## Usage

The SSH Pod Log Downloader now offers three distinct subcommands, each serving a specific purpose:

1. **File Subcommand**: This subcommand helps you print the logs into a file and send the logs from the server to a local client server. Example usage:

   ```shell
   sshpodlog file -f config-file.yaml
   ```

2. **Terminal Subcommand**: Use this subcommand to print the logs directly to the terminal. Example usage:

   ```shell
   sshpodlog terminal -f config-file.yaml
   ```

3. **Filter Subcommand**: The filter subcommand allows you to filter logs based on a string and print them to the terminal. Example usage:

   ```shell
   sshpodlog filter -r ERROR -f config-file.yaml
   ```
4. **Tail Subcommand**: The tail subcommand allows you to watch live logs as they come in to the terminal. 
Example usage:

```shell
sshpodlog tail -f config-file.yaml
```


   Note that the configuration file can be in YAML, JSON, TOML, or other supported formats.

Here's how the configuration file (e.g., `config-file.yaml`) looks in YAML format:

```yaml
# SSH Pod Log Downloader Configuration
server: 192.168.1.100
port: 22
username: myuser
cluster: production
key: ~/.ssh/id_rsa
```

The following flags can be used:

- `--server` (`-s`): The SSH server address.
- `--username` (`-u`): The SSH username.
- `--cluster` (`-c`): The Kubernetes context switch.
- `--port` (`-p`): The SSH port.
- `--key` (`-k`): The path to the SSH private key file.
- `--from-file` (`-f`): The path to the configuration file.

## Contributing

Contributions are welcome! If you encounter any issues or have suggestions for improvement, please open an issue or submit a pull request.

## License

This project is licensed under the [MIT License](LICENSE).
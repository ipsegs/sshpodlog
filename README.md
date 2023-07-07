```markdown
# SSH Pod Log Downloader

This is a command-line tool written in Go that connects to an SSH server, interacts with Kubernetes using `kubectl`, and downloads logs from a specified pod.

## Installation

1. Clone the repository:
   ```shell
   git clone https://github.com/yourusername/ssh-pod-log-downloader.git
   ```

2. Navigate to the project directory:
   ```shell
   cd ssh-pod-log-downloader
   ```

3. Run the installation script:
   ```shell
   ./install.sh
   ```

   The installation script will prompt you to choose the target operating system for building the executable (Linux, macOS, or Windows). After selecting the operating system, the script will build the executable accordingly.

4. Once the installation is complete, you will find the built executable file in the project directory.

## Usage

To use the SSH Pod Log Downloader, follow these steps:

1. Open a terminal or command prompt.

2. Navigate to the directory where the SSH Pod Log Downloader executable is located.

3. Run the executable with the desired command-line arguments. For example:
   ```shell
   ./ssh-pod-log-downloader -server 192.168.1.100 -port 22 -username myuser -cluster production -key ~/.ssh/id_rsa
   ```

   Replace the arguments with the appropriate values for your SSH server configuration.

## Contributing

Contributions are welcome! If you find any issues or have suggestions for improvement, please open an issue or submit a pull request.

## License

This project is licensed under the [MIT License](LICENSE).
```

Feel free to customize the README.md file further if needed, adding any additional information or sections that you think are relevant.
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/pkg/sftp"
	"github.com/schollz/progressbar/v3"
	"golang.org/x/crypto/ssh"
	"golang.org/x/term"
)

func main() {
	server := flag.String("server", "", "<usage> -server <ip address or hostname>")
	port := flag.Int("port", 22, "<usage -port <port number>")
	username := flag.String("username", "", "<usege -username <username>")
	kubectlClusterSwitch := flag.String("cluster", "default", "usage -context <cluster>")
	flag.Parse()

	if *server == "" {
		log.Fatal("Usage: -server <ip address>")
		return
	}

	if *username == "" {
		log.Fatal("Usage: -username <username>")
		return
	}

	fmt.Print("Enter Password: ")
	password, _ := readPassword()
	if password == nil {
		log.Fatalf("Please enter a password\n")
		return
	}

	config := &ssh.ClientConfig{
		User: *username,
		Auth: []ssh.AuthMethod{
			ssh.Password(string(password)),
		},
		Timeout:         5 * time.Second,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	conn, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", *server, *port), config)
	if err != nil {
		log.Fatalf("Error: Cannot connect to the server %v", err)
		return
	}
	defer conn.Close()

	session, err := conn.NewSession()
	if err != nil {
		log.Fatalf("Error: SSH session failed %v", err)
		return
	}
	defer session.Close()

	fmt.Print("Namespace: ")
	namespace, err := readInput()
	if err != nil {
		log.Fatalf("Namespace does not exist: %v", err)
		return
	}

	contextSwitch := fmt.Sprintf("kubectl config use-context %s", *kubectlClusterSwitch)
	session.Output(contextSwitch)
	fmt.Printf("Context has been switched to %s\n", *kubectlClusterSwitch)

	session, err = conn.NewSession()
	if err != nil {
		log.Fatalf("Error: SSH session failed %v", err)
		return
	}
	defer session.Close()

	listPods := fmt.Sprintf("kubectl get po -n %s -o jsonpath='{.items[*].metadata.name}'", namespace)
	podList, err := session.Output(listPods)
	if err != nil {
		log.Fatalf("Command failed to run in SSH Session: %v", err)
	}
	fmt.Println(string(podList))

	fmt.Print("Enter pod name: ")
	podName, err := readInput()
	if err != nil {
		log.Fatalf("Pod does not exist: %v", err)
	}
	session.Close()

	newSession, err := conn.NewSession()
	if err != nil {
		log.Fatalf("Unable to create second SSH connection: %v", err)
	}
	defer newSession.Close()
	logFileName := podName + ".txt"
	getPodLogs := fmt.Sprintf("kubectl logs %s -n %s > %s", podName, namespace, logFileName)
	_, err = newSession.CombinedOutput(getPodLogs)
	if err != nil {
		log.Fatal("Failed to run second command in second SSH connection", err)
	}

	sftpClient, err := sftp.NewClient(conn)
	if err != nil {
		log.Println("Failed to create SFTP client:", err)
		return
	}
	defer sftpClient.Close()

	var homeDir string
	if runtime.GOOS == "windows" {
		homeDrive := os.Getenv("HOMEDRIVE")
		homePath := os.Getenv("HOMEPATH")
		homeDir = filepath.Join(homeDrive, homePath)
	} else {
		homeDir, err = os.UserHomeDir()
		if err != nil {
			log.Fatalf("Failed to get user's home directory: %v", err)
		}
	}

	var localFilePath string
	if runtime.GOOS == "windows" {
		downloadFolder := filepath.Join(homeDir, "Downloads")
		localFilePath = filepath.Join(downloadFolder, logFileName)
	} else {
		localFilePath = filepath.Join(homeDir, logFileName)
	}
	fmt.Println("Location of file on local directory", localFilePath)

	remoteFile, err := sftpClient.Open(logFileName)
	if err != nil {
		log.Println("Failed to open remote file:", err)
		return
	}
	defer remoteFile.Close()

	localFile, err := os.Create(localFilePath)
	if err != nil {
		log.Println("Failed to create the local file:", err)
		return
	}
	defer localFile.Close()

	// Get the file size
	remoteFileInfo, _ := remoteFile.Stat()
	fileSize := remoteFileInfo.Size()

	// Create a progress bar
	bar := progressbar.DefaultBytes(fileSize, "copying")

	_, err = io.Copy(localFile, remoteFile)
	if err != nil {
		log.Println("Error copying file:", err)
		return
	}

	bar.Finish()
	filesizeToKb := float64(fileSize/1024)
	fmt.Printf("Copied %.2f kilobytes content.\n", filesizeToKb)
}

func readPassword() ([]byte, error) {
	password, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		log.Fatal(err)
	}
	return password, err
}

func readInput() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	input = strings.TrimSpace(input)
	return input, err
}

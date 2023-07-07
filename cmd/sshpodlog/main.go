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
	kubectlClusterSwitch := flag.String("cluster", "default", "usage -cluster <cluster>")
	privateKey := flag.String("key", "", "usage -key <path to the private key file>")
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
		log.Fatalf("Please enter a password")
		return
	}

	config := &ssh.ClientConfig{
		User: *username,
		Auth: []ssh.AuthMethod{
			ssh.Password(string(password)),
		},
		//Auth: []ssh.AuthMethod {},
		Timeout:         5 * time.Second,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	// var password byte
	// fmt.Print("Enter Password: ")
	// if password != 0 {
	// 	//fmt.Print("Enter Password: ")
	// 	password, _ := readPassword()
	// 	if password == nil {
	// 	log.Fatalf("Please enter a password")
	// 	return
	// }
	// 	config.Auth = append(config.Auth, ssh.Password(string(password)))
	// }

	if *privateKey != "" {
		file, err := os.Open(*privateKey)
		if err != nil {
			log.Printf("Unable to open file path: %v", err)
			return
		}
		defer file.Close()

		privateKeyBytes, err := io.ReadAll(file)
		if err != nil {
			log.Printf("unable to read file: %v", err)
			return
		}
		key, err := ssh.ParsePrivateKey(privateKeyBytes)
		if err != nil {
			log.Printf("Failed to parse private key: %v", err)
			return
		}
		config.Auth = append(config.Auth, ssh.PublicKeys(key))
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

	//	fmt.Print("Namespace: ")
	//	namespace, err := readInput()
	//	if err != nil {
	//		log.Fatalf("Namespace does not exist: %v", err)
	//		return
	//	}

	fmt.Println()

	contextSwitch := fmt.Sprintf("kubectl config use-context %s\n", *kubectlClusterSwitch)
	session.Output(contextSwitch)
	fmt.Printf("in %s cluster\n", *kubectlClusterSwitch)

	session, err = conn.NewSession()
	if err != nil {
		log.Fatalf("Error: SSH session failed %v", err)
		return
	}
	defer session.Close()

	// listPods := fmt.Sprintf("kubectl get po -n %s -o jsonpath='{.items[*].metadata.name}'", namespace)
	// podList, err := session.Output(listPods)
	// if err != nil {
	// 	log.Fatalf("Command failed to run in SSH Session: %v", err)
	// }
	// fmt.Println(string(podList))
	var namespace string
	for {
		fmt.Print("Enter the namespace: ")
		namespace, err = readInput()
		if err != nil {
			log.Printf("The namespace does not exist: %v\n", err)
			return
		}

		session, err = conn.NewSession()
		if err != nil {
			log.Printf("Failed to start the session connection: %v\n", err)
			return
		}

		defer session.Close()

		checkNamespace := fmt.Sprintln("kubectl get namespace", namespace)
		_, err = session.CombinedOutput(checkNamespace)
		if err == nil {
			break
		}

		log.Println("Error: Namespace does not exist")

		session, err = conn.NewSession()
		if err != nil {
			log.Printf("Unable to start another session connection: %v\n", err)
			return
		}

		defer session.Close()

		fmt.Println("Available namespaces:")
		namespaceList, _ := session.CombinedOutput(fmt.Sprintln("kubectl get ns -o jsonpath='{.items[*].metadata.name}'"))
		fmt.Println(string(namespaceList))
	}

	session, err = conn.NewSession()
	if err != nil {
		log.Fatalf("Unable to start the session connection: %v", err)
		return
	}

	defer session.Close()

	listPods := fmt.Sprintf("kubectl get po -n %s -o jsonpath='{.items[*].metadata.name}'", namespace)
	pods, err := session.Output(listPods)
	if err != nil {
		log.Fatalf("Unable to list pods: %v", err)
		return
	}
	if len(pods) == 0 {
		log.Println("There are no pods in this namespace")
		return
	}
	fmt.Println(string(pods))

	fmt.Print("Enter pod name: ")
	podName, err := readInput()
	if err != nil {
		log.Fatalf("Pod does not exist: %v", err)
	}
	session.Close()

	session, err = conn.NewSession()
	if err != nil {
		log.Fatalf("Unable to create second SSH connection: %v", err)
	}
	defer session.Close()

	logFileName := podName + ".txt"
	//logFilePath := "logs/" + logFileName
	getPodLogs := fmt.Sprintf("kubectl logs %s -n %s > %s", podName, namespace, logFileName)
	_, err = session.CombinedOutput(getPodLogs)
	if err != nil {
		log.Fatalf("Failed to run second command in second SSH connection: %v", err)
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
			log.Printf("Failed to get user's home directory: %v", err)
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

	remoteFileInfo, _ := remoteFile.Stat()
	fileSize := remoteFileInfo.Size()

	bar := progressbar.DefaultBytes(fileSize, "copying to local")

	_, err = io.Copy(localFile, remoteFile)
	if err != nil {
		log.Println("Error copying file:", err)
		return
	}

	bar.Finish()
	filesizeToKb := float64(fileSize / 1024)
	fmt.Printf("Copied %.2f kilobytes content.\n", filesizeToKb)

	session, err = conn.NewSession()
	if err != nil {
		log.Printf("Error: SSH connection cannot be established: %v\n", err)
		return
	}
	defer session.Close()

	rmLogFile := fmt.Sprintf("rm %s", logFileName)
	_, err = session.CombinedOutput(rmLogFile)
	if err != nil {
		log.Printf("Error: command can't be ran: %v", err)
	}
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

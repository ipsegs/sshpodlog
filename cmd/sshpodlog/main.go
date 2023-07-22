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

type Config struct{
	Server string
	Port int
	Username string
	KctlCtxSwitch string
	PrivateKey string
}

type Application struct{
	InfoLog *log.Logger
	ErrorLog *log.Logger
	Config Config
}

func main() {

	var cfg Config
	//command line arguments
	flag.StringVar(&cfg.Server, "server", "", "usage -server <ip address or hostname>")
	flag.IntVar(&cfg.Port, "port", 22, "usage -port <port number>")
	flag.StringVar(&cfg.Username,"username", "", "usege -username <username>")
	flag.StringVar(&cfg.KctlCtxSwitch, "cluster", "default", "usage -cluster <cluster>")
	flag.StringVar(&cfg.PrivateKey, "key", "", "usage -key <path to the private key file>")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	app := &Application{
		InfoLog: infoLog,
		ErrorLog: errorLog,
		Config: cfg,
	}

	//input validation
	if cfg.Server == "" {
		errorLog.Println("Usage: -server <ip address>")
		return
	}

	if cfg.Username == "" {
		errorLog.Println("Usage: -username <username>")
		return
	}

	fmt.Print("Enter Password: ")
	password, _ := app.readPassword()
	if password == nil {
		errorLog.Println("Please enter a password")
		return
	}

	config := &ssh.ClientConfig{
		User: cfg.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(string(password)),
		},
		//Auth: []ssh.AuthMethod {},
		Timeout:         5 * time.Second,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// Load private key if provided
	if cfg.PrivateKey != "" {
		file, err := os.Open(cfg.PrivateKey)
		if err != nil {
			errorLog.Printf("Unable to open file path: %v", err)
			return
		}
		defer file.Close()

		privateKeyBytes, err := io.ReadAll(file)
		if err != nil {
			errorLog.Printf("unable to read file: %v", err)
			return
		}
		key, err := ssh.ParsePrivateKey(privateKeyBytes)
		if err != nil {
			errorLog.Printf("Failed to parse private key: %v", err)
			return
		}
		config.Auth = append(config.Auth, ssh.PublicKeys(key))
	}

	//SSH Server connection
	//conn, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", cfg.Server, cfg.Port), config)
	conn, err := ssh.Dial("tcp", app.fmtSprint(), config)
	if err != nil {
		errorLog.Fatalf("Error: Cannot connect to the server %v", err)
		return
	}
	defer conn.Close()

	session, err := conn.NewSession()
	if err != nil {
		log.Fatalf("Error: SSH session failed %v", err)
		return
	}
	defer session.Close()

	fmt.Println()

	//switch kubernetes context, Default is for the default namespace
	contextSwitch := fmt.Sprintf("kubectl config use-context %s\n", cfg.KctlCtxSwitch)
	session.Output(contextSwitch)
	fmt.Printf("in %s cluster\n", cfg.KctlCtxSwitch)

	session, err = conn.NewSession()
	if err != nil {
		errorLog.Fatalf("Error: SSH session failed %v", err)
		return
	}
	defer session.Close()

	//Get kubernetes namespace from user
	var namespace string
	for {
		session, err = conn.NewSession()
		if err != nil {
			errorLog.Printf("Unable to start another session connection: %v\n", err)
			return
		}

		defer session.Close()

		fmt.Println("Available namespaces:")
		namespaceList, _ := session.CombinedOutput(fmt.Sprintln("kubectl get ns -o jsonpath='{.items[*].metadata.name}'"))
		fmt.Println(string(namespaceList))
		
		fmt.Print("Enter the namespace: ")
		namespace, err = app.readInput()
		if err != nil {
			errorLog.Printf("The namespace does not exist: %v\n", err)
			return
		}

		session, err = conn.NewSession()
		if err != nil {
			errorLog.Printf("Failed to start the session connection: %v\n", err)
			return
		}

		defer session.Close()

		checkNamespace := fmt.Sprintln("kubectl get namespace", namespace)
		_, err = session.CombinedOutput(checkNamespace)
		if err == nil {
			break
		}

		errorLog.Println("Error: Namespace does not exist")
	}

	session, err = conn.NewSession()
	if err != nil {
		errorLog.Fatalf("Unable to start the session connection: %v", err)
		return
	}

	defer session.Close()

	//List kubernetes pod within the namespace
	listPods := fmt.Sprintf("kubectl get po -n %s -o jsonpath='{.items[*].metadata.name}'", namespace)
	pods, err := session.Output(listPods)
	if err != nil {
		errorLog.Fatalf("Unable to list pods: %v", err)
		return
	}
	if len(pods) == 0 {
		errorLog.Println("There are no pods in this namespace")
		return
	}
	fmt.Println(string(pods))

	//Enter pod name from the list provided above
	fmt.Print("Enter pod name: ")
	podName, err := app.readInput()
	if err != nil {
		errorLog.Fatalf("Pod does not exist: %v", err)
	}
	session.Close()

	session, err = conn.NewSession()
	if err != nil {
		errorLog.Fatalf("Unable to create second SSH connection: %v", err)
	}
	defer session.Close()

	// create file name from the pod, using .txt, it can be .log
	logFileName := podName + ".txt"
	getPodLogs := fmt.Sprintf("kubectl logs %s -n %s > %s", podName, namespace, logFileName)
	_, err = session.CombinedOutput(getPodLogs)
	if err != nil {
		errorLog.Fatalf("Failed to run second command in second SSH connection: %v", err)
	}

	//create sftp connection
	sftpClient, err := sftp.NewClient(conn)
	if err != nil {
		errorLog.Println("Failed to create SFTP client:", err)
		return
	}
	defer sftpClient.Close()

	//get home directory of the local server
	var homeDir string
	if runtime.GOOS == "windows" {
		homeDrive := os.Getenv("HOMEDRIVE")
		homePath := os.Getenv("HOMEPATH")
		homeDir = filepath.Join(homeDrive, homePath)

	} else {
		homeDir, err = os.UserHomeDir()
		if err != nil {
			errorLog.Printf("Failed to get user's home directory: %v", err)
		}
	}

	//get file path to save the file into the local machin
	var localFilePath string
	if runtime.GOOS == "windows" {
		downloadFolder := filepath.Join(homeDir, "Downloads")
		localFilePath = filepath.Join(downloadFolder, logFileName)
	} else {
		localFilePath = filepath.Join(homeDir, logFileName)
	}
	infoLog.Println("Location of file on local directory", localFilePath)

	// Open remote file
	remoteFile, err := sftpClient.Open(logFileName)
	if err != nil {
		log.Println("Failed to open remote file:", err)
		return
	}
	defer remoteFile.Close()

	//create the file name in the local machine
	localFile, err := os.Create(localFilePath)
	if err != nil {
		errorLog.Println("Failed to create the local file:", err)
		return
	}
	defer localFile.Close()

	//file size of the remote file
	remoteFileInfo, _ := remoteFile.Stat()
	fileSize := remoteFileInfo.Size()

	bar := progressbar.DefaultBytes(fileSize, "copying to local")

	//copy the file from remote to local
	_, err = io.Copy(localFile, remoteFile)
	if err != nil {
		errorLog.Println("Error copying file:", err)
		return
	}

	bar.Finish()
	filesizeToKb := fileSize / 1024
	fmt.Printf("Copied %d kilobytes content.\n", filesizeToKb)

	session, err = conn.NewSession()
	if err != nil {
		errorLog.Printf("Error: SSH connection cannot be established: %v\n", err)
		return
	}
	defer session.Close()

	//remove log file from the remote to reduce excesses
	rmLogFile := fmt.Sprintf("rm %s", logFileName)
	_, err = session.CombinedOutput(rmLogFile)
	if err != nil {
		errorLog.Printf("Error: command can't be ran: %v", err)
	}
}

// function to input password without showing it on the terminal
func (app *Application) readPassword() ([]byte, error) {
	password, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return password, err
}

// input value but remove spaces and any unnecessary input that can be present.
func (app *Application) readInput() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		log.Println(err)
		return "", err
	}
	input = strings.TrimSpace(input)
	return input, err
}

func (app *Application) fmtSprint() string {
	return fmt.Sprintf("%s:%d", app.Config.Server, app.Config.Port)
}
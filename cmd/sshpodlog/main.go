package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"golang.org/x/crypto/ssh"
)

type Config struct {
	Server        string
	Port          int
	Username      string
	KctlCtxSwitch string
	PrivateKey    string
}

type Application struct {
	InfoLog   *log.Logger
	ErrorLog  *log.Logger
	Config    Config
	Namespace string
}

func main() {

	var cfg Config
	//command line arguments
	flag.StringVar(&cfg.Server, "server", "", "usage -server <ip address or hostname>")
	flag.IntVar(&cfg.Port, "port", 22, "usage -port <port number>")
	flag.StringVar(&cfg.Username, "username", "", "usege -username <username>")
	flag.StringVar(&cfg.KctlCtxSwitch, "cluster", "default", "usage -cluster <cluster>")
	flag.StringVar(&cfg.PrivateKey, "key", "", "usage -key <path to the private key file>")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	app := &Application{
		InfoLog:  infoLog,
		ErrorLog: errorLog,
		Config:   cfg,
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

	//Password Auth
	fmt.Print("Enter Password: ")
	password, _ := app.readPassword() 
	if password == nil {
		errorLog.Println("Please enter a password")
		return
	}

	//configure ssh client information
	sshConfig := &ssh.ClientConfig{
		User: cfg.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(string(password)),
		},
		Timeout:         5 * time.Second,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// Load private key Auth if provided
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
		sshConfig.Auth = append(sshConfig.Auth, ssh.PublicKeys(key))
	}

	//create SSH connection.
	conn, err := ssh.Dial("tcp", app.fmtSprint(), sshConfig)
	if err != nil {
		errorLog.Fatalf("Error: Cannot connect to the server %v", err)
		return
	}
	defer conn.Close()

	// This is to switch kubernetes context if the jumper is connected to multi-clusters
	if err = app.switchContext(conn); err != nil {
		errorLog.Fatalf("Unable to switch Context, %v", err)
	}

	//Get namespace
	namespace, err := app.getNamespace(conn)
	if err != nil {
		app.ErrorLog.Printf("Unable to get namespace: %v\n", err)
		return
	}

	//get pod name
	logFileName, err := app.podInfo(conn, namespace)
	if err != nil {
		errorLog.Fatalf("Unable to create get Pod: %v \n", err)
	}

	// Copy file from remote Server to Local
	if err = app.sftpClientCopy(conn, logFileName); err != nil {
		errorLog.Fatalf("Unable to copy file: %v", err)
		return
	}

	// Remove the file from the remote server
	if err = app.rmFile(conn, logFileName); err != nil {
		errorLog.Fatalf("Unable to remove file: %v", err)
		return
	}

}

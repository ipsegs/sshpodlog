package pkg

import (
	//"flag"

	"log"
	"os"

	data "github.com/ipsegs/sshpodlog/internal"
	"golang.org/x/crypto/ssh"
)

// Struct used to create dependency injection
type Application struct {
	Cfg data.Config
	App data.Application
}


func Sshpodlog(server, username, kctlCtxSwitch, privateKey string, port int) {

	cfgCheck := &data.Config{
		Server:        server,
		Port:          port,
		Username:      username,
		KctlCtxSwitch: kctlCtxSwitch,
		PrivateKey:    privateKey,
	}

	app := &Application{
		Cfg: *cfgCheck,
	}

	app.App.InfoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.App.ErrorLog = log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	//input validation
	if server == "" {
		app.App.ErrorLog.Println("Usage: -server <ip address>")
		return
	}

	if username == "" {
		app.App.ErrorLog.Println("Usage: -username <username>")
		return
	}

	//create SSH Configuration
	sshConfig, err := app.sshConfigInfo()
	if err != nil {
		app.App.ErrorLog.Fatalf("Cannot connect to the server: %v", err)
		return
	}

	//create SSH connection.
	conn, err := ssh.Dial("tcp", app.fmtSprint(), sshConfig)
	if err != nil {
		app.App.ErrorLog.Fatalf("Cannot connect to the server: %v", err)
		return
	}
	defer conn.Close()

	// This is to switch kubernetes context if the jumper is connected to multi-clusters
	if err = app.switchContext(conn); err != nil {
		app.App.ErrorLog.Fatalf("Unable to switch Kubernetes Context, %v", err)
	}

	//Get namespace
	namespace, err := app.getNamespace(conn)
	if err != nil {
		app.App.ErrorLog.Printf("Unable to get namespace: %v\n", err)
		return
	}

	//get pod name
	err = app.listPodsinNamespace(conn, namespace)
	if err != nil {
		app.App.ErrorLog.Fatalf("Unable to get pod logs: %v \n", err)
	}

	//get filename from pods
	podName, err := app.getpodName(conn, namespace)
	if err != nil {
		app.App.ErrorLog.Fatalf("Unable to get pod name: %v \n", err)
	}

	logFileName, err := app.getlogFileNameFromPodName(conn, namespace, podName)
	if err != nil {
		app.App.ErrorLog.Fatalf("Unable to get log file name from pods: %v \n", err)
	}

	// Copy file from remote Server to Local
	if err = app.sftpClientCopy(conn, logFileName); err != nil {
		app.App.ErrorLog.Fatalf("Unable to copy file: %v", err)
		return
	}

	// Remove the file from the remote server
	if err = app.rmFile(conn, logFileName); err != nil {
		app.App.ErrorLog.Fatalf("Unable to remove file: %v", err)
		return
	}

}

package pkg

import (
	"log"
	"os"

	data "github.com/ipsegs/sshpodlog/internal"
	"golang.org/x/crypto/ssh"
)

// Struct used to create dependency injection
type Application struct {
	Cfg           data.ClientConfig
	App           data.Application
	KctlCtxSwitch string
}

func Sshpodlog(server, username, kctlCtxSwitch, privateKey string, port int) *ssh.Client {

	cfgCheck := &data.ClientConfig{
		Server:        server,
		Port:          port,
		Username:      username,
		PrivateKey:    privateKey,
		KctlCtxSwitch: kctlCtxSwitch,
	}

	app := &Application{
		Cfg: *cfgCheck,
	}

	app.App.InfoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.App.ErrorLog = log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	//input validation
	if server == "" {
		app.App.ErrorLog.Println("Usage: -server <ip address>")
		return nil
	}

	if username == "" {
		app.App.ErrorLog.Println("Usage: -username <username>")
		return nil
	}

	//create SSH Configuration and SSH connection
	conn, err := app.SshConnectConfig()
	if err != nil {
		app.App.ErrorLog.Fatalf("Cannot connect to the server: %v", err)
		return nil
	}

	// This is to switch kubernetes context if the jumper is connected to multi-clusters
	if err = app.SwitchContext(conn); err != nil {
		app.App.ErrorLog.Fatalf("Unable to switch Kubernetes Context, %v", err)
	}

	return conn

}

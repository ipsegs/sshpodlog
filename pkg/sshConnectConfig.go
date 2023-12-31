package pkg

import (
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"golang.org/x/crypto/ssh"
)

func (app *Application) SshConnectConfig() (*ssh.Client, error) {
	//Password Auth
	fmt.Print("Enter Password: ")
	password, _ := app.readPassword()
	if password == nil {
		app.App.ErrorLog.Println("Please enter a password")
		return nil, errors.New("please input password")
	}

	//configure ssh client information
	sshConfig := &ssh.ClientConfig{
		User: app.Cfg.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(string(password)),
		},
		Timeout:         5 * time.Second,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// Load private key Auth if provided
	if app.Cfg.PrivateKey != "" {
		file, err := os.Open(app.Cfg.PrivateKey)
		if err != nil {
			app.App.ErrorLog.Printf("Unable to open file path: %v", err)
			return nil, err
		}
		defer file.Close()

		privateKeyBytes, err := io.ReadAll(file)
		if err != nil {
			app.App.ErrorLog.Printf("unable to read file: %v", err)
			return nil, err
		}
		key, err := ssh.ParsePrivateKey(privateKeyBytes)
		if err != nil {
			app.App.ErrorLog.Printf("Failed to parse private key: %v", err)
			return nil, err
		}
		sshConfig.Auth = append(sshConfig.Auth, ssh.PublicKeys(key))
	}
	//fmt.Println("testing")
	//create SSH connection.
	conn, err := ssh.Dial("tcp", app.IpPort(), sshConfig)
	if err != nil {
		app.App.ErrorLog.Printf("Cannot connect to the server: %v", err)
		return nil, err
	}

	return conn, nil
}

package main

import (
	"fmt"
	"io"
	"os"
	"time"

	"golang.org/x/crypto/ssh"
)

func (app *Application) sshConfigInfo() (*ssh.ClientConfig, error) {
	//Password Auth
	fmt.Print("Enter Password: ")
	password, _ := app.readPassword()
	if password == nil {
		app.ErrorLog.Println("Please enter a password")
		return nil, nil
	}

	//configure ssh client information
	sshConfig := &ssh.ClientConfig{
		User: app.Config.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(string(password)),
		},
		Timeout:         5 * time.Second,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// Load private key Auth if provided
	if app.Config.PrivateKey != "" {
		file, err := os.Open(app.Config.PrivateKey)
		if err != nil {
			app.ErrorLog.Printf("Unable to open file path: %v", err)
			return nil, err
		}
		defer file.Close()

		privateKeyBytes, err := io.ReadAll(file)
		if err != nil {
			app.ErrorLog.Printf("unable to read file: %v", err)
			return nil, err
		}
		key, err := ssh.ParsePrivateKey(privateKeyBytes)
		if err != nil {
			app.ErrorLog.Printf("Failed to parse private key: %v", err)
			return nil, err
		}
		sshConfig.Auth = append(sshConfig.Auth, ssh.PublicKeys(key))
	}

	return sshConfig, nil
}

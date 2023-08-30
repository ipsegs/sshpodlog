package pkg

import (
	"fmt"

	"golang.org/x/crypto/ssh"
)

func (app *Application) getpodName(conn *ssh.Client, namespace string) (string, error) {
	//Enter pod name from the list provided
	fmt.Print("Enter pod name: ")
	podName, _ := app.readInput()

	session, err := conn.NewSession()
	if err != nil {
		app.App.ErrorLog.Printf("Unable to start session: %v\n", err)
	}
	defer session.Close()

	return podName, err
}
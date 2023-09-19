package pkg

import (
	"fmt"

	"golang.org/x/crypto/ssh"
)

func (app *Application) GetpodName(conn *ssh.Client, namespace string) (string, error) {
		//Enter pod name from the list provided
		fmt.Print("Enter pod name: ")
		podName, err := app.readInput()
		if err != nil {
			app.App.ErrorLog.Println("unable to read pod name input")
		}
		
		session, err := conn.NewSession()
		if err != nil {
			app.App.ErrorLog.Printf("Unable to start session: %v\n", err)
		}
		defer session.Close()
	

	return podName, err
}

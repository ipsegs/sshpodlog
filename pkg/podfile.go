package pkg

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/ssh"
)

func (app *Application) podFile(conn *ssh.Client, namespace string) (string, error) {
	//Enter pod name from the list provided above
	fmt.Print("Enter pod name: ")
	podName, _ := app.readInput()

	session, err := conn.NewSession()
	if err != nil {
		app.App.ErrorLog.Printf("Unable to start session: %v\n", err)
	}
	defer session.Close()

	// create file name from the pod name with a .log extension
	logFileName := podName + ".log"
	getPodLogs := fmt.Sprintf("kubectl logs %s -n %s > %s", podName, namespace, logFileName)
	_, err = session.CombinedOutput(getPodLogs)
	if err != nil {
		app.App.ErrorLog.Printf("Failed to get pod logs")
		return "", errors.New("failed to get pod logs")
	}
	return logFileName, err
}

package pkg

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/ssh"
)

func (app *Application) GetlogFileNameFromPodName(conn *ssh.Client, namespace string) (string, error) {
	podName, err := app.GetpodName(conn, namespace)
	if err != nil {
		app.App.ErrorLog.Printf("Unable to get pod name: %v\n", err)
	}
	session, err := conn.NewSession()
	if err != nil {
		fmt.Printf("Unable to start session: %v\n", err)
	}
	defer session.Close()
	logFileName := podName + ".log"
	getPodLogs := fmt.Sprintf("kubectl logs %s -n %s > %s", podName, namespace, logFileName)
	_, err = session.CombinedOutput(getPodLogs)
	if err != nil {
		app.App.ErrorLog.Printf("Failed to get pod logs")
		return "", errors.New("failed to get pod logs")
	}

	return logFileName, err
}

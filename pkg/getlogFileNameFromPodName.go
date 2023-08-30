package pkg

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/ssh"
)

func (app *Application) getlogFileNameFromPodName(conn *ssh.Client, namespace string, podName string) (string, error) {
	session, err := conn.NewSession()
	if err != nil {
		app.App.ErrorLog.Printf("Unable to start session: %v\n", err)
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

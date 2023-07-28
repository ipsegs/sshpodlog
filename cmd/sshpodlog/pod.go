package main

import (
	"fmt"

	"golang.org/x/crypto/ssh"
)

func (app *Application) podInfo(conn *ssh.Client, namespace string) (string, error) {
	

	session, err := conn.NewSession()
	if err != nil {
		app.ErrorLog.Printf("Unable to start the session connection: %v\n", err)
		return "", nil
	}
	defer session.Close()

	//List kubernetes pod within the namespace
	listPods := fmt.Sprintf("kubectl get po -n %s -o jsonpath='{.items[*].metadata.name}'", namespace)
	pods, err := session.Output(listPods)
	if err != nil {
		app.ErrorLog.Printf("Unable to list pods: %v\n", err)
		return "", nil
	}
	if len(pods) == 0 {
		app.ErrorLog.Println("There are no pods in this namespace")
		return "", nil
	}
	fmt.Println(string(pods))

	//Enter pod name from the list provided above
	fmt.Print("Enter pod name: ")
	podName, _ := app.readInput()

	session, err = conn.NewSession()
	if err != nil {
		app.ErrorLog.Printf("Unable to start the session connection: %v\n", err)
	}
	defer session.Close()

	// create file name from the pod, using .txt, it can be .log
	logFileName := podName + ".txt"
	getPodLogs := fmt.Sprintf("kubectl logs %s -n %s > %s", podName, namespace, logFileName)
	_, err = session.CombinedOutput(getPodLogs)
	if err != nil {
		app.ErrorLog.Printf("Failed to run second command in second SSH connection: %v\n", err)
		return "", err
	}

	return logFileName, err

}

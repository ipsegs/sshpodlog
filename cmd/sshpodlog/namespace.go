package main

import (
	"fmt"

	"golang.org/x/crypto/ssh"
)

func (app *Application) getNamespace(session *ssh.Session, conn *ssh.Client) (string, error) {
	var namespace string
	var err error

	fmt.Println("Available namespaces:")
	namespaceList, _ := session.CombinedOutput(fmt.Sprintln("kubectl get ns -o jsonpath='{.items[*].metadata.name}'"))
	fmt.Println(string(namespaceList))
	for {
		//Get Namespace from User
		fmt.Print("Enter the namespace: ")
		namespace, err = app.readInput()
		if err != nil {
			app.ErrorLog.Printf("The namespace does not exist: %v\n", err)
			return "", err
		}

		session, err = conn.NewSession()
		if err != nil {
			app.ErrorLog.Printf("Error: SSH connection cannot be established: %v\n", err)
			return "", err
		}
		defer session.Close()

		checkNamespace := fmt.Sprintln("kubectl get namespace", namespace)
		_, err = session.CombinedOutput(checkNamespace)
		if err == nil {
			break
		}
	}
	return namespace, err
}

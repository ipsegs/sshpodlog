package pkg

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/ssh"
)

func (app *Application) getNamespace(conn *ssh.Client) (string, error) {
	var namespace string
	var err error

	session, err := conn.NewSession()
	if err != nil {
		app.App.ErrorLog.Fatalf("Error: SSH session failed %v", err)
		return "", err
	}
	defer session.Close()

	fmt.Println("Available namespaces:")
	namespaceList, _ := session.CombinedOutput(fmt.Sprintln("kubectl get ns -o jsonpath='{.items[*].metadata.name}'"))
	fmt.Println(string(namespaceList))
	for {
		//Get Namespace from User
		fmt.Print("Enter the namespace: ")
		namespace, err = app.readInput()
		if namespace == "" {
			app.App.ErrorLog.Printf("No namespace provided\n")
			continue
			//return "", errors.New("please input namespace")
		}
		if err != nil {
			app.App.ErrorLog.Printf("The namespace does not exist: %v\n", err)
			return "", errors.New("namespace does not exist")
		}

		session, err = conn.NewSession()
		if err != nil {
			app.App.ErrorLog.Printf("Error: SSH connection cannot be established: %v\n", err)
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

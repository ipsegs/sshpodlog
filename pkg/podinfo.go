package pkg

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/ssh"
)

func (app *Application) podInfo(conn *ssh.Client, namespace string) (error) {

	session, err := conn.NewSession()
	if err != nil {
		app.App.ErrorLog.Printf("Unable to start the session connection: %v\n", err)
		return  nil
	}
	defer session.Close()

	//List kubernetes pod within the namespace
	listPods := fmt.Sprintf("kubectl get po -n %s -o jsonpath='{.items[*].metadata.name}'", namespace)
	pods, err := session.Output(listPods)
	if err != nil {
		app.App.ErrorLog.Printf("Unable to list pods: %v\n", err)
		return nil
	}
	if len(pods) == 0 {
		app.App.ErrorLog.Println("There are no pods in the namespace:", namespace)
		return errors.New("pods unavailable in namespace")
	}
	fmt.Println(string(pods))
	 return err
}


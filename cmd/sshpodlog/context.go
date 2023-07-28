package main

import (
	"fmt"

	"golang.org/x/crypto/ssh"
)

func (app *Application) switchContext(conn *ssh.Client) error {
	session, err := conn.NewSession()
	if err != nil {
		app.ErrorLog.Printf("Error: SSH session failed %v\n", err)
		return err
	}
	defer session.Close()
	
	fmt.Println()
	//switch kubernetes context, if no namespace argument is given, it uses the current context
	contextSwitch := fmt.Sprintf("kubectl config use-context %s\n", app.Config.KctlCtxSwitch)
	session.Output(contextSwitch)
	fmt.Printf("in %s cluster\n", app.Config.KctlCtxSwitch)
	return err
}

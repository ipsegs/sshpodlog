package pkg

import (
	"errors"
	"fmt"
	"golang.org/x/crypto/ssh"
)

func (app *Application) switchContext(conn *ssh.Client) error {
	session, err := conn.NewSession()
	if err != nil {
		app.App.ErrorLog.Printf("Error: SSH session failed: %v\n", err)
		return errors.New("SSH Session failed")
	}
	defer session.Close()

	fmt.Println()
	//switch kubernetes context, if no namespace argument is given, it uses the current context
	contextSwitch := fmt.Sprintf("kubectl config use-context %s\n", app.Cfg.KctlCtxSwitch)
	session.Output(contextSwitch)
	fmt.Printf("in %s cluster\n", app.Cfg.KctlCtxSwitch)
	return err
}

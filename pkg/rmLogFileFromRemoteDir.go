package pkg

import (
	"fmt"

	"golang.org/x/crypto/ssh"
)

func (app *Application) RmFile(conn *ssh.Client, logFileName string) error {

	session, err := conn.NewSession()
	if err != nil {
		app.App.ErrorLog.Printf("Error: SSH connection cannot be established: %v \n", err)
	}
	defer session.Close()

	//remove log file from the remote after copy
	rmLogFile := fmt.Sprintf("rm %s", logFileName)
	_, err = session.CombinedOutput(rmLogFile)
	if err != nil {
		app.App.ErrorLog.Printf("Error: File cannot be removed: %v", err)
		return err
	}
	return err
}

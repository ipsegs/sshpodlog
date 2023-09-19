package terminal

import (
	"fmt"

	"github.com/ipsegs/sshpodlog/pkg"
	"golang.org/x/crypto/ssh"
)



func ShowLogsInTerminal(conn *ssh.Client) error {
	inst := &pkg.Application{}

	namespace, err := inst.GetNamespace(conn)
	if err != nil {
		inst.App.ErrorLog.Printf("Unable to get correct Namespace: %v \n", err)
		return err
	}

	// Get pod list in the specific namespace
	err = inst.ListPodsinNamespace(conn, namespace)
	if err != nil {
		inst.App.ErrorLog.Printf("Unable to get pod logs: %v \n", err)
		return err
	}
	podName, err := inst.GetpodName(conn, namespace)
	if err != nil {
		inst.App.ErrorLog.Printf("Unable to get pod name: %v \n", err)
		return err
	}

	session, err := conn.NewSession()
	if err != nil {
		inst.App.ErrorLog.Printf("Error: SSH connection cannot be established: %v \n", err)
		return err
	}
	defer session.Close()

	getPodLogs := fmt.Sprintf("kubectl logs %s -n %s", podName, namespace)
	printPodLogs, err := session.CombinedOutput(getPodLogs)
	if err != nil {
		inst.App.ErrorLog.Printf("Incorrect pod input: %v \n", err )
		return err
	}
	fmt.Println(string(printPodLogs))
	return err
}

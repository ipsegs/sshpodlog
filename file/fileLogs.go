package file

import (
	"github.com/ipsegs/sshpodlog/pkg"
	"golang.org/x/crypto/ssh"
)

func GetLogsInFile(conn *ssh.Client) error {
	inst := &pkg.Application{}

	namespace, err := inst.GetNamespace(conn)
	if err != nil {
		inst.App.ErrorLog.Printf("Unable to get Namespace: %v \n", err)
		return err
	}

	// Get pod list in the specific namespace
	err = inst.ListPodsinNamespace(conn, namespace)
	if err != nil {
		inst.App.ErrorLog.Printf("Unable to get pod logs: %v \n", err)
		return err
	}

	session, err := conn.NewSession()
	if err != nil {
		inst.App.ErrorLog.Printf("Error: SSH connection cannot be established: %v \n", err)
		return err
	}
	defer session.Close()

	// Get log file name using the pod name
	logFileName, err := inst.GetlogFileNameFromPodName(conn, namespace)
	if err != nil {
		inst.App.ErrorLog.Printf("Unable to get log file name from pods: %v \n", err)
		return err
	}

	if err = inst.CopyFileFromRemoteToLocal(conn, logFileName); err != nil {
		inst.App.ErrorLog.Printf("Unable to copy file from remote to Locale: %v", err)
		return err
	}

	// Remove the file from the remote server
	if err = inst.RmFile(conn, logFileName); err != nil {
		inst.App.ErrorLog.Printf("Unable to remove file: %v", err)
		return err
	}
	return err
}

package tail

import (
	"bufio"
	"fmt"
	"sync"

	"github.com/ipsegs/sshpodlog/pkg"
	"golang.org/x/crypto/ssh"
)

func TailLogsInTerminal(conn *ssh.Client) error {
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

	getPodLogs := fmt.Sprintf("kubectl logs -f %s -n %s", podName, namespace)

	// Set up standard output, and error streams
	stdout, err := session.StdoutPipe()
	if err != nil {
		inst.App.ErrorLog.Printf("Error: Unable to setup stdout: %v \n", err)
		return err
	}

	stderr, err := session.StderrPipe()
	if err != nil {
		inst.App.ErrorLog.Printf("Error: Unable to setup stderr: %v \n", err)
		return err
	}

	if err := session.Start(getPodLogs); err != nil {
		inst.App.ErrorLog.Printf("Error: Unable to start command: %v \n", err)
		return err
	}

	var wg sync.WaitGroup
	wg.Add(1)
	// Read and process output from standard output and error streams
	go func() {
		defer wg.Done()

		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			line := scanner.Text()
			fmt.Println(line) // Process each line of the command output here
		}

		scannerErr := bufio.NewScanner(stderr)
		for scannerErr.Scan() {
			line := scannerErr.Text()
			fmt.Println(line) // Process each line of the command error output here
		}
	}()

	// sigCh := make(chan os.Signal, 1)
	// signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	// <-sigCh

	wg.Wait()
	return nil
}
package terminal

import (
	"bufio"
	"bytes"
	"fmt"
	"sync"

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

	var wg sync.WaitGroup
	wg.Add(1)

	getPodLogs := fmt.Sprintf("kubectl logs %s -n %s", podName, namespace)

	printOutput := func(output []byte) {
		scanner := bufio.NewScanner(bytes.NewReader(output))
		for scanner.Scan() {
			line := scanner.Text()
			fmt.Println(line) // Process each line of the command output here
		}
		if err := scanner.Err(); err != nil {
			inst.App.ErrorLog.Printf("Error reading from command output: %v\n", err)
		}
	}

	go func() {
		defer wg.Done()

		printPodLogs, err := session.CombinedOutput(getPodLogs)
		if err != nil {
			inst.App.ErrorLog.Printf("Error running command: %v\n", err)
			return
		}

		printOutput(printPodLogs)
	}()

	wg.Wait()
	return nil

	// var wg sync.WaitGroup
	// wg.Add(1)
	// go func() {
	// defer wg.Done()
	// getPodLogs := fmt.Sprintf("kubectl logs %s -n %s", podName, namespace)
	// printPodLogs, err := session.CombinedOutput(getPodLogs)
	// if err != nil {
	// 	inst.App.ErrorLog.Printf("Incorrect pod input: %v \n", err )
	// 	return
	// }

	// fmt.Println(string(printPodLogs))
	// }()
	// wg.Wait()
	// return err
}

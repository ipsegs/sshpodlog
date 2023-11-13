package tail

import (
	"fmt"
	"io"
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
	const bufferSize = 1024
	processStream := func(reader io.Reader, wg *sync.WaitGroup) {
		defer wg.Done()

		buffer := make([]byte, bufferSize)

		for {
			bytesRead, err := reader.Read(buffer)
			if err != nil {
				if err == io.EOF {
					return
				} else {
					fmt.Printf("Error reading stream: %v\n", err)
				}
				return
			}
			fmt.Print(string(buffer[:bytesRead]))
		}
	}
	wg := &sync.WaitGroup{}
	wg.Add(2)
	//wg.Add(1)
	go processStream(stdout, wg)
	go processStream(stderr, wg)

	// Wait for the command to finish
	err = session.Wait()
	if err != nil {

		return err
	}
	// var wg sync.WaitGroup
	// wg.Add(1)
	// // Read and process output from standard output and error streams
	// go func() {
	// 	defer wg.Done()

	// 	scanner := bufio.NewScanner(stdout)
	// 	for scanner.Scan() {
	// 		line := scanner.Text()
	// 		fmt.Println(line) 
	// 	}

	// 	scannerErr := bufio.NewScanner(stderr)
	// 	for scannerErr.Scan() {
	// 		line := scannerErr.Text()
	// 		fmt.Println(line) 
	// 	}
	// }()

	wg.Wait()
	return nil
}

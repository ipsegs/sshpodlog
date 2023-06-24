package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"syscall"
	"time"

	"golang.org/x/crypto/ssh"
	"golang.org/x/term"
)

func main() {

	server := flag.String("server", "", "<usage> -server <ip address or hostname>")
	port := flag.Int("port", 22, "<usage -port <port number>")
	username := flag.String("username", "", "<usege -username <username>")
	flag.Parse()

	if *server == "" {
		log.Fatal("Usage: -server <ip address>  ")
		return
	}

	if *username == "" {
		log.Fatal("Usage: -username <username>")
		return
	}

	fmt.Print("Enter Password: ")
	password, _ := readPassword()
	if password == nil {
		log.Fatalf("Please Enter a password\n")
		return
	}

	config := &ssh.ClientConfig{
		User: *username,
		Auth: []ssh.AuthMethod{
			ssh.Password(string(password)),
		},
		Timeout:         5 * time.Second,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	conn, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", *server, *port), config)
	if err != nil {
		log.Fatalf("Error: Cannot connect to the server %v", err)
		return
	}
	defer conn.Close()

	session, err := conn.NewSession()
	if err != nil {
		log.Fatalf("Error: SSH session failed %v", err)
	}
	defer session.Close()

	fmt.Print("Namespace: ")
	namespace, err := readInput()
	if err != nil {
		log.Fatalf("namespace does not exist: %v", err)
		return
	}

	listPods := fmt.Sprintf("kubectl get po -n %s -o jsonpath='{.items[*].metadata.name}'", namespace)
	podList, err := session.Output(listPods)
	if err != nil {
		log.Fatalf("command failed to run in SSH Session: %v", err)
	}
	fmt.Println(string(podList))

	fmt.Print("Enter pod name: ")
	podName, err := readInput()
	if err != nil {
		log.Fatalf("pod does not exist: %v", err)
	}
	session.Close()

	newSession, err := conn.NewSession()
	if err != nil {
		log.Fatalf("Unable to create second ssh connection: %v", err)
	}
	defer newSession.Close()
	getPodLogs := fmt.Sprintf("kubectl logs %s -n %s", podName, namespace)
	// podLog, err := newSession.CombinedOutput(getPodLogs)
	// if err != nil {
	// 	log.Fatal("Failed to run second command in second ssh connection", err)
	// }
	stdout, err := newSession.StdoutPipe()
	if err != nil {
		log.Fatalf("Failed to create stdout pipe: %v", err)
	}

	err = newSession.Start(getPodLogs)
	if err != nil {
		log.Fatalf("Failed to start command execution: %v", err)
	}

	podLog, err := ioutil.ReadAll(stdout)
	if err != nil {
		log.Fatalf("Failed to read command output: %v", err)
	}

	err = newSession.Wait()
	if err != nil {
		log.Fatalf("Command execution failed: %v", err)
	}

	fmt.Println(string(podLog))
}

func readPassword() ([]byte, error) {
	password, _ := term.ReadPassword(int(syscall.Stdin))
	return password, nil
}

func readInput() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	input = strings.TrimSpace(input)
	return input, err
}

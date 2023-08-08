package pkg

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	//"syscall"

	"golang.org/x/term"
)

// function to input password without showing it on the terminal
func (app *Application) readPassword() ([]byte, error) {
	// password, err := term.ReadPassword(int(syscall.Stdin))
	password, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	
	return password, err
}

// input value but remove spaces and any unnecessary input that can be present.
func (app *Application) readInput() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		log.Println(err)
		return "", err
	}
	input = strings.TrimSpace(input)
	return input, err
}

func (app *Application) fmtSprint() string {
	return fmt.Sprintf("%s:%d", app.Cfg.Server, app.Cfg.Port)
}

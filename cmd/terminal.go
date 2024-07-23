/*
Copyright © 2023 NAME HERE Adebayosegun174@gmail.com
*/
package cmd

import (
	"log"
	"strconv"

	"github.com/ipsegs/sshpodlog/pkg"
	"github.com/ipsegs/sshpodlog/terminal"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// terminalCmd represents the terminal command
var terminalCmd = &cobra.Command{
	Use:   "terminal",
	Short: "Outputs logs into terminal",
	Run:   terminalFunction,
}

func init() {
	rootCmd.AddCommand(terminalCmd)
}

func terminalFunction(cmd *cobra.Command, args []string) {
	if cmd.Flag("from-file").Changed {
		flags.server = viper.GetString("server")
		flags.username = viper.GetString("username")
		flags.kctlCtxSwitch = viper.GetString("cluster")
		flags.privateKey = viper.GetString("key")
		flags.port = viper.GetInt("port")
	} else {
		// Use the values from command-line flags if provided
		if ServerFlag := cmd.Flag("server"); ServerFlag != nil && ServerFlag.Changed {
			flags.server = ServerFlag.Value.String()
		}
		if UsernameFlag := cmd.Flag("username"); UsernameFlag != nil && UsernameFlag.Changed {
			flags.username = UsernameFlag.Value.String()
		}
		if KctlCtxSwitchFlag := cmd.Flag("cluster"); KctlCtxSwitchFlag != nil && KctlCtxSwitchFlag.Changed {
			flags.kctlCtxSwitch = KctlCtxSwitchFlag.Value.String()
		}
		if PrivateKeyFlag := cmd.Flag("key"); PrivateKeyFlag != nil && PrivateKeyFlag.Changed {
			flags.privateKey = PrivateKeyFlag.Value.String()
		}
		if PortFlag := cmd.Flag("port"); PortFlag != nil && PortFlag.Changed {
			port, err := strconv.Atoi(PortFlag.Value.String())
			if err != nil {
				log.Printf("Error parsing port: %v", err)
			}
			flags.port = port
		}
	}

	//Check if "cluster flag" is empty, and if so, set it to "current"
	if flags.kctlCtxSwitch == "" {
		flags.kctlCtxSwitch = "current"
	}

	conn, err := pkg.Sshpodlog(flags.server, flags.username, flags.kctlCtxSwitch, flags.privateKey, flags.port)
	if err != nil {
		return
	}
	terminal.ShowLogsInTerminal(conn)
}

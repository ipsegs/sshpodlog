/*
Copyright © 2023 NAME HERE Adebayosegun174@gmail.com
*/
package cmd

import (
	"log"
	"strconv"

	"github.com/ipsegs/sshpodlog/file"
	"github.com/ipsegs/sshpodlog/pkg"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// fileCmd represents the file command
var fileCmd = &cobra.Command{
	Use:   "file",
	Short: "A subcommand that outputs kubernetes pod logs into a file in the client's computer",
	Run:   fileFunction,
}

func init() {
	rootCmd.AddCommand(fileCmd)
}

func fileFunction(cmd *cobra.Command, args []string) {
	if cmd.Flag("from-file").Changed {
		flags.Server = viper.GetString("server")
		flags.Username = viper.GetString("username")
		flags.KctlCtxSwitch = viper.GetString("cluster")
		flags.PrivateKey = viper.GetString("key")
		flags.Port = viper.GetInt("port")
	} else {
		// Use the values from command-line flags if provided
		if ServerFlag := cmd.Flag("server"); ServerFlag != nil && ServerFlag.Changed {
			flags.Server = ServerFlag.Value.String()
		}
		if UsernameFlag := cmd.Flag("username"); UsernameFlag != nil && UsernameFlag.Changed {
			flags.Username = UsernameFlag.Value.String()
		}
		if KctlCtxSwitchFlag := cmd.Flag("cluster"); KctlCtxSwitchFlag != nil && KctlCtxSwitchFlag.Changed {
			flags.KctlCtxSwitch = KctlCtxSwitchFlag.Value.String()
		}
		if PrivateKeyFlag := cmd.Flag("key"); PrivateKeyFlag != nil && PrivateKeyFlag.Changed {
			flags.PrivateKey = PrivateKeyFlag.Value.String()
		}
		if PortFlag := cmd.Flag("port"); PortFlag != nil && PortFlag.Changed {
			port, err := strconv.Atoi(PortFlag.Value.String())
			if err != nil {
				log.Printf("Error parsing port: %v", err)
			}
			flags.Port = port
		}
	}

	//Check if "cluster flag" is empty, and if so, set it to "current"
	if flags.KctlCtxSwitch == "" {
		flags.KctlCtxSwitch = "current"
	}
	conn := pkg.Sshpodlog(flags.Server, flags.Username, flags.KctlCtxSwitch, flags.PrivateKey, flags.Port)
	file.GetLogsInFile(conn)
}

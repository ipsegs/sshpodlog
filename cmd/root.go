/*
Copyright © 2023 Adebayosegun174@gmail.com
*/
package cmd

import (
	"log"
	"strconv"

	"github.com/ipsegs/sshpodlog/pkg"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var flags struct {
	server        string
	username      string
	kctlCtxSwitch string
	privateKey    string
	fromFile      string
	port          int
}

var rootCmd = &cobra.Command{
	Use:   "sshpodlog",
	Short: "A tool to manage Kubernetes pod logs",
	Long:  `A CLI tool designed for efficient management and retrieval of logs from Kubernetes pods.`,
	Run:   defaultFunction,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)

	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Global flags
	rootCmd.PersistentFlags().StringVarP(&flags.server, "server", "s", "", "SSH server address")
	rootCmd.PersistentFlags().StringVarP(&flags.username, "username", "u", "", "SSH username")
	rootCmd.PersistentFlags().StringVarP(&flags.kctlCtxSwitch, "cluster", "c", "current", "kubectl context switch")
	rootCmd.PersistentFlags().StringVarP(&flags.privateKey, "key", "k", "", "SSH private key path")
	rootCmd.PersistentFlags().IntVarP(&flags.port, "port", "p", 22, "SSH port")
	rootCmd.PersistentFlags().StringVarP(&flags.fromFile, "from-file", "f", "default", "Configuration properties file")
}

func initConfig() {
	if flags.fromFile != "" {
		viper.SetConfigFile(flags.fromFile)
		viper.ReadInConfig()
	}
}

func defaultFunction(cmd *cobra.Command, args []string) {
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

	// Check if "cluster flag" is empty, and if so, set it to "current"
	if flags.kctlCtxSwitch == "" {
		flags.kctlCtxSwitch = "current"
	}

	pkg.Sshpodlog(flags.server, flags.username, flags.kctlCtxSwitch, flags.privateKey, flags.port)
}

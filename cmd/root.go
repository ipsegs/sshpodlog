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
	Server        string
	Username      string
	KctlCtxSwitch string
	PrivateKey    string
	Port          int
	FromFile      string
}

var rootCmd = &cobra.Command{
	Use:   "sshpodlog",
	Short: "A tool to manage Kubernetes pod logs",
	Long:  `A CLI tool designed for efficient management and retrieval of logs from Kubernetes pods. ...`,
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
	rootCmd.PersistentFlags().StringVarP(&flags.Server, "server", "s", "", "SSH server address")
	rootCmd.PersistentFlags().StringVarP(&flags.Username, "username", "u", "", "SSH username")
	rootCmd.PersistentFlags().StringVarP(&flags.KctlCtxSwitch, "cluster", "c", "current", "kubectl context switch")
	rootCmd.PersistentFlags().StringVarP(&flags.PrivateKey, "key", "k", "", "SSH private key path")
	rootCmd.PersistentFlags().IntVarP(&flags.Port, "port", "p", 22, "SSH port")
	rootCmd.PersistentFlags().StringVarP(&flags.FromFile, "from-file", "f", "default", "Configuration properties file")
}

func initConfig() {
	if flags.FromFile != "" {
		viper.SetConfigFile(flags.FromFile)
		viper.ReadInConfig()
	}
}

func defaultFunction(cmd *cobra.Command, args []string) {
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

	pkg.Sshpodlog(flags.Server, flags.Username, flags.KctlCtxSwitch, flags.PrivateKey, flags.Port)
}

/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/ipsegs/sshpodlog/pkg"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "sshpodlog",
	Short: "A tool to manage Kubernetes pod logs",
	Long:  `A CLI tool designed for efficient management and retrieval of logs from Kubernetes pods. This tool operates by utilizing SSH connections to establish access through a dedicated jump server, also known as a jump host or Bastion host. This jump server acts as an intermediary node to ensure secure access to the target environment. The jump server is equipped with a preinstalled instance of kubectl, the Kubernetes command-line tool. This strategic setup empowers the CLI to simplify the process of fetching logs from Kubernetes clusters Through the utilization of SSH connections and the jump server, the utility offers a secure gateway into the Kubernetes environment. By leveraging the integrated "kubectl", users can seamlessly extract logs from various pods within the cluster. This approach streamlines the process of log retrieval, ultimately enhancing the overall management and troubleshooting capabilities within the Kubernetes ecosystem.`,
	Run:   defaultFunction,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var (
	Server        string
	Username      string
	KctlCtxSwitch string
	PrivateKey    string
	Port          int
)

func init() {

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVarP(&Server, "server", "s", "", "usage")
	rootCmd.PersistentFlags().StringVarP(&Username, "username", "u", "", "usage")
	rootCmd.PersistentFlags().StringVarP(&KctlCtxSwitch, "cluster", "c", "current", "usage")
	rootCmd.PersistentFlags().StringVarP(&PrivateKey, "key", "k", "", "usage")
	rootCmd.PersistentFlags().IntVarP(&Port, "port", "p", 22, "usage")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

//The rootCMD that will be run which is the default function being run if no subcommand is added.
func defaultFunction(cmd *cobra.Command, args []string) {
	pkg.Sshpodlog(Server, Username, KctlCtxSwitch, PrivateKey, Port)

}

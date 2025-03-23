package main

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/yuansheng0111/MultiSSH/internal/ssh"
)

var (
	address  []string
	username []string
	password []string
	keyPath  []string
	command  []string
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "ssh-tool",
		Short: "A tool for SSH operations",
		Long:  `An SSH utility tool that allows connecting to multiple hosts and running commands.`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if len(address) == 0 {
				log.Fatal("Error: at least one host address is required")
			}
		},
	}

	rootCmd.Flags().StringArrayVarP(&address, "address", "a", []string{}, "SSH host addresses")
	rootCmd.Flags().StringArrayVarP(&username, "user", "u", []string{}, "SSH usernames")
	rootCmd.Flags().StringArrayVarP(&password, "password", "p", []string{}, "SSH passwords")
	rootCmd.Flags().StringArrayVarP(&keyPath, "key", "k", []string{}, "SSH private key paths")
	rootCmd.Flags().StringArrayVarP(&command, "cmd", "c", []string{}, "SSH commands to execute")

	rootCmd.SilenceUsage = true
	rootCmd.SilenceErrors = true
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}

	if len(address) == 0 || len(username) == 0 || len(command) == 0 {
		log.Fatal("Please provide host address, username, and command")
	} else if len(address) != len(username) || len(username) != len(command) {
		log.Fatal("Please provide the same number of host addresses, usernames, and commands")
	}

	if len(password) == 0 && len(keyPath) == 0 {
		log.Fatal("Please provide either password or key path")
	} else if len(password) != 0 && len(keyPath) != 0 {
		log.Fatal("Please provide either password or key path, but not both")
	}

	hosts := []ssh.HostInfo{}
	for id := range address {
		host := ssh.HostInfo{
			Address: address[id] + ":22",
			User:    username[id],
			Command: command[id],
		}

		if len(password) > 0 {
			host.Password = password[id]
		} else {
			host.KeyPath = keyPath[id]
		}

		hosts = append(hosts, host)
	}

	ssh.ExecuteCommandOnHosts(hosts)
}

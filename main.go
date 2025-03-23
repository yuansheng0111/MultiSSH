package main

import (
	"github.com/spf13/cobra"
	"github.com/yuansheng0111/MultiSSH/cmd"
	"github.com/yuansheng0111/MultiSSH/internal/ssh"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "ssh-tool",
		Short: "A tool for SSH operations",
		Long:  `An SSH utility tool that allows connecting to multiple hosts and running commands.`,
	}

	config := cmd.NewConfig()
	config.ParseFlags(rootCmd)
	config.CheckFlags()

	hosts := []ssh.HostInfo{}
	for id := range config.Address {
		host := ssh.HostInfo{
			Address: config.Address[id] + ":22",
			User:    config.Username[id],
			Command: config.Command[id],
		}

		if len(config.Password) > 0 {
			host.Password = config.Password[id]
		} else {
			host.KeyPath = config.KeyPath[id]
		}

		hosts = append(hosts, host)
	}

	ssh.ExecuteCommandOnHosts(hosts)
}

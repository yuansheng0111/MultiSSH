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

	hosts, _ := ssh.BuildHosts(config)
	ssh.ExecuteCommandOnHosts(hosts)
}

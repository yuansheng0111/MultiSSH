package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

type Config struct {
	Address  []string
	Username []string
	Password []string
	KeyPath  []string
	Command  []string
}

func NewConfig() *Config {
	return &Config{}
}

func (c *Config) ParseFlags(rootCmd *cobra.Command) {
	rootCmd.Flags().StringArrayVarP(&c.Address, "address", "a", []string{}, "SSH host addresses")
	rootCmd.Flags().StringArrayVarP(&c.Username, "user", "u", []string{}, "SSH usernames")
	rootCmd.Flags().StringArrayVarP(&c.Password, "password", "p", []string{}, "SSH passwords")
	rootCmd.Flags().StringArrayVarP(&c.KeyPath, "key", "k", []string{}, "SSH private key paths")
	rootCmd.Flags().StringArrayVarP(&c.Command, "cmd", "c", []string{}, "SSH commands to execute")

	rootCmd.SilenceUsage = true
	rootCmd.SilenceErrors = true
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func (c *Config) CheckFlags() {
	if len(c.Address) == 0 || len(c.Address) == 0 || len(c.Command) == 0 {
		log.Fatal("Please provide host address, username, and command")
	} else if len(c.Address) != len(c.Username) || len(c.Username) != len(c.Command) {
		log.Fatal("Please provide the same number of host addresses, usernames, and commands")
	}

	if len(c.Password) == 0 && len(c.KeyPath) == 0 {
		log.Fatal("Please provide either password or key path")
	} else if len(c.Password) != 0 && len(c.KeyPath) != 0 {
		log.Fatal("Please provide either password or key path, but not both")
	}
}

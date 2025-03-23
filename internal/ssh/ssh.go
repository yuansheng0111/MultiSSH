package ssh

import (
	"fmt"
	"os"
	"sync"

	"golang.org/x/crypto/ssh"
)

// HostInfo stores connection details for each host
type HostInfo struct {
	Address  string
	User     string
	Password string
	KeyPath  string // Path to private key
	Command  string
}

// ExecuteCommandOnHosts runs a command on multiple hosts concurrently
func ExecuteCommandOnHosts(hosts []HostInfo) map[string]string {
	results := make(map[string]string) // Store output
	var wg sync.WaitGroup
	var mu sync.Mutex // Prevent race conditions

	for _, host := range hosts {
		wg.Add(1)
		go func(h HostInfo) {
			defer wg.Done()
			output, err := runSSHCommand(h)
			if err != nil {
				output = fmt.Sprintf("Error: %v", err)
			}
			mu.Lock()
			results[h.Address] = output
			mu.Unlock()
		}(host)
	}

	wg.Wait() // Wait for all goroutines to finish
	return results
}

// runSSHCommand handles SSH connection and execution
func runSSHCommand(host HostInfo) (string, error) {
	// Setup SSH config
	var authMethods []ssh.AuthMethod
	if host.KeyPath != "" {
		privateKeyBytes, err := os.ReadFile(host.KeyPath)
		if err != nil {
			return "", fmt.Errorf("failed to read SSH key: %w", err)
		} else {
			signer, err := ssh.ParsePrivateKey(privateKeyBytes)
			if err != nil {
				return "", fmt.Errorf("failed to parse SSH key: %w", err)
			}
			authMethods = []ssh.AuthMethod{ssh.PublicKeys(signer)}
		}
	} else {
		authMethods = []ssh.AuthMethod{ssh.Password(host.Password)}
	}

	config := &ssh.ClientConfig{
		User:            host.User,
		Auth:            authMethods,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	fmt.Println(host.Command)
	// Connect to SSH
	client, err := ssh.Dial("tcp", host.Address, config)
	if err != nil {
		fmt.Printf("failed to connect to %s: %w\n", host.Address, err)
		return "", fmt.Errorf("failed to connect to %s: %w", host.Address, err)
	} else {
		fmt.Printf("Connected to %s\n", host.Address)
	}
	defer client.Close()

	// Create session
	session, err := client.NewSession()
	if err != nil {
		return "", fmt.Errorf("failed to create session: %w", err)
	}
	defer session.Close()

	// Run command
	output, err := session.CombinedOutput(host.Command)
	if err != nil {
		return "", fmt.Errorf("command execution failed: %w", err)
	}

	fmt.Println(string(output))
	return string(output), nil
}

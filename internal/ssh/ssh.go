package ssh

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"

	"gopkg.in/yaml.v3"

	"github.com/yuansheng0111/MultiSSH/cmd"
	"golang.org/x/crypto/ssh"
)

// HostInfo stores connection details for each host
type HostInfo struct {
	FileName         string
	UploadFilePath   string
	DownloadFilePath string
	Address          string
	User             string
	Password         string
	KeyPath          string // Path to private key
	Command          string
}

func BuildHostsFromConfigFile(configFile string) ([]HostInfo, error) {
	var config map[string]interface{}
	var hosts []HostInfo

	if strings.HasSuffix(configFile, ".json") {
		jsonData, err := os.ReadFile(configFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
		err = json.Unmarshal(jsonData, &config)
		if err != nil {
			return nil, fmt.Errorf("failed to parse config file: %w", err)
		}
	} else if strings.HasSuffix(configFile, ".yaml") {
		yamlData, err := os.ReadFile(configFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}

		err = yaml.Unmarshal(yamlData, &config)
		if err != nil {
			return nil, fmt.Errorf("failed to parse config file: %w", err)
		}
	} else {
		return nil, fmt.Errorf("invalid config file format: %s", configFile)
	}

	hostsConfig := config["hosts"].([]interface{})
	for _, host := range hostsConfig {
		hostInfo := host.(map[string]interface{})
		host := HostInfo{}
		host.Address = hostInfo["address"].(string) + ":22"
		host.User = hostInfo["user"].(string)
		if hostInfo["password"] != nil {
			host.Password = hostInfo["password"].(string)
		}
		if hostInfo["key"] != nil {
			host.KeyPath = hostInfo["key"].(string)
		}
		if hostInfo["uploadfilepath"] != nil {
			host.UploadFilePath = hostInfo["uploadfilepath"].(string)
		}
		if hostInfo["downloadfilepath"] != nil {
			host.DownloadFilePath = hostInfo["downloadfilepath"].(string)
		}
		if hostInfo["command"] != nil {
			host.Command = hostInfo["command"].(string)
		}
		if hostInfo["filename"] != nil {
			host.FileName = hostInfo["filename"].(string)
		}
		hosts = append(hosts, host)
	}

	return hosts, nil
}

func BuildHosts(config *cmd.Config) ([]HostInfo, error) {
	if config.FilePath != "" {
		return BuildHostsFromConfigFile(config.FilePath)
	}

	hosts := []HostInfo{}
	for id := range config.Address {
		host := HostInfo{
			Address: config.Address[id] + ":22",
			User:    config.Username[id],
		}

		if len(config.Password) > 0 {
			host.Password = config.Password[id]
		} else {
			host.KeyPath = config.KeyPath[id]
		}

		if config.UploadFilePath != "" {
			host.UploadFilePath = config.UploadFilePath
			host.FileName = config.UploadFilePath
		}
		if config.DownloadFilePath != "" {
			host.DownloadFilePath = config.DownloadFilePath
			host.FileName = config.DownloadFilePath
		}
		if config.Command != nil && len(config.Command) > 0 {
			host.Command = config.Command[id]
		}
		if config.FileName != "" {
			host.FileName = config.FileName
		}

		hosts = append(hosts, host)
	}
	return hosts, nil
}

// ExecuteCommandOnHosts runs a command on multiple hosts concurrently
func ExecuteCommandOnHosts(hosts []HostInfo) map[string]string {
	results := make(map[string]string) // Store output
	var wg sync.WaitGroup
	var mu sync.Mutex // Prevent race conditions

	for _, host := range hosts {
		wg.Add(1)
		go func(host HostInfo) {
			defer wg.Done()
			output, err := runSSHCommand(host)
			if err != nil {
				output = fmt.Sprintf("Error: %v", err)
			}
			mu.Lock()
			results[host.Address] = output
			mu.Unlock()
		}(host)
	}

	wg.Wait() // Wait for all goroutines to finish
	return results
}

// runSSHCommand handles SSH connection and execution
func runSSHCommand(host HostInfo) (string, error) {
	client, err := NewSSHClient(host)
	if err != nil {
		return "", fmt.Errorf("failed to create SSH client: %w", err)
	}
	defer client.Close()

	if host.UploadFilePath != "" {
		err := UploadFile(client, host.UploadFilePath, host.FileName)
		if err != nil {
			return "", fmt.Errorf("failed to upload file: %w", err)
		}
		return "", nil
	}

	if host.DownloadFilePath != "" {
		err := DownloadFile(client, host.DownloadFilePath, host.FileName)
		if err != nil {
			return "", fmt.Errorf("failed to download file: %w", err)
		}
		return "", nil
	}

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

	return string(output), nil
}

func NewSSHClient(host HostInfo) (*ssh.Client, error) {
	var authMethods []ssh.AuthMethod

	if host.KeyPath != "" {
		privateKeyBytes, err := os.ReadFile(host.KeyPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read SSH key: %w", err)
		}
		signer, err := ssh.ParsePrivateKey(privateKeyBytes)
		if err != nil {
			return nil, fmt.Errorf("failed to parse SSH key: %w", err)
		}
		authMethods = []ssh.AuthMethod{ssh.PublicKeys(signer)}
	} else {
		authMethods = []ssh.AuthMethod{ssh.Password(host.Password)}
	}

	config := &ssh.ClientConfig{
		User:            host.User,
		Auth:            authMethods,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", host.Address, config)
	if err != nil {
		return nil, fmt.Errorf("failed to dial SSH: %w", err)
	} else {
		fmt.Printf("Connected to %s\n", host.Address)
	}
	return client, nil
}

# MultiSSH - Execute Commands on Multiple Remote Hosts

## Overview
MultiSSH is a powerful and lightweight tool that enables users to run commands, transfer files, and clone directories on multiple remote hosts simultaneously over SSH. It provides both CLI and GUI interfaces, making it suitable for automation, DevOps, and system administration tasks.

## Features Done
- **Run Commands on Multiple Hosts**: Execute the same command on multiple servers concurrently.
- **Parallel Execution**: Uses Goroutines to handle multiple SSH connections efficiently.
- **Secure Authentication**: Supports both password and private key authentication.
- **Configuration Management**: Supports JSON and YAML configuration.
- **File Transfer**: Upload single file through command line arguments.

## Features To Be Done
- **File Transfer**: Upload or download files and directories through command line arguments and config file.
- **Container Deployment**: Easily deployable on Docker and Kubernetes.

## Usage
Recommend to use
```shell
go run main.go -f <config_file>
```
- Set up either absolute private key path or password for each host in config file.

Alternatively, use
```shell
go run main.go -a <address> -u <username> -p <password> -k <key_path> -c <command>
```
- For multiple hosts, specify each address with its own flag: `--address <address1> --address <address2> ...`.
- Use only one authentication method for all hosts: either password or private key.

## Roadmap
- [ ] GUI interface
- [ ] Docker/Kubernetes integration
- [ ] Logging & monitoring with Prometheus/Grafana
- [ ] Windows support

## License
MIT License

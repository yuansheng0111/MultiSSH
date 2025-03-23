# MultiSSH - Execute Commands on Multiple Remote Hosts

## Overview
MultiSSH is a powerful and lightweight tool that enables users to run commands, transfer files, and clone directories on multiple remote hosts simultaneously over SSH. It provides both CLI and GUI interfaces, making it suitable for automation, DevOps, and system administration tasks.

## Features Done
- **Run Commands on Multiple Hosts**: Execute the same command on multiple servers concurrently.
- **Parallel Execution**: Uses Goroutines to handle multiple SSH connections efficiently.
- **Secure Authentication**: Supports both password and private key authentication.

## Features To Be Done
- **File Transfer**: Upload or download files across multiple remote hosts.
- **Configuration Management**: Supports JSON/YAML-based configuration.
- **CLI and GUI Support**: Offers a command-line tool and an optional web-based UI.
- **Container Deployment**: Easily deployable on Docker and Kubernetes.

## Usage
```shell
go run main.go -a <address> -u <username> -p <password> -k <key_path> -c <command>
```
- For multiple hosts, specify each address with its own flag: `--address <address1> --address <address2> ...`.
- Use only one authentication method for all hosts: either password or private key.

## Roadmap
- [ ] Web-based GUI using React/Wails
- [ ] Docker/Kubernetes integration
- [ ] Logging & monitoring with Prometheus/Grafana
- [ ] Windows support

## License
MIT License

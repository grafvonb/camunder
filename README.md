# Camunder
Camunder is a CLI tool to interact with Camunda 8.

> **Disclaimer:** This project is currently in an **early development stage**.  
> Features, APIs, and behavior may change without notice. Use at your own risk.

## Overview

**Camunder** is a command-line interface (CLI) tool designed to interact with [Camunda 8](https://camunda.com/platform/), a popular workflow and decision automation platform.
It provides a convenient way to manage and monitor Camunda 8 resources directly from the terminal.

## Features implemented

- get cluster topology: `$ camunder get cluster-topology`
- get process definitions: `$ camunder get process-definition`
- 
## Features planned

- run process instances (bulk): `$ camunder run process-instance --bpmn-process-id <bpmnProcessId> --variables <key1=value1,key2=value2,...> --count <number>`
- get process instances: `$ camunder get process-instance`
- cancel process instances (bulk): `$ camunder cancel process-instance --ids <id1,id2,...>`
- delete process instances (bulk): `$ camunder delete process-instance --ids <id1,id2,...>`
- ...and more to come!

## Configuration

Camunder can be configured using command-line flags, environment variables or a configuration file. 
Default configuration file locations are application path or user config directory:
- `$HOME/.config/camunder/config.yaml` (Linux)
- `$HOME/Library/Application Support/camunder/config.yaml` (macOS)
- `%APPDATA%\camunder\config.yaml` (Windows)

## Usage 
```
Camunder is a CLI tool to interact with Camunda 8.

Usage:
  camunder [flags]
  camunder [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  get         List resources of a defined type e.g. cluster-topology, process-definition, process-instance etc.
  help        Help about any command

Flags:
      --camunda8-base-url string   Camunda 8 API base URL
      --camunda8-token string      Camunda 8 API bearer token
      --config string              path to config file
  -h, --help                       help for camunder
      --operate-base-url string    Operate API base URL
      --operate-token string       Operate API bearer token
      --tasklist-base-url string   Tasklist API base URL
      --tasklist-token string      Tasklist API bearer token
      --timeout duration           HTTP timeout (e.g. 10s, 1m)

Use "camunder [command] --help" for more information about a command.
```

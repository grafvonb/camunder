# Camunder
Camunder is a CLI tool to interact with Camunda 8.

> **Disclaimer:** This project is currently in an **early development stage**.  
> Features, APIs, and behavior may change without notice. Use at your own risk.

## Overview & Motivation

**Camunder** is a command-line interface (CLI) tool designed to interact with [Camunda 8](https://camunda.com/platform/),
a popular workflow and decision automation platform.
It provides a convenient way to manage and monitor Camunda 8 resources directly from the terminal.

Standard Camunda 8 tools like [Operate](https://docs.camunda.io/docs/components/operate/overview/) and [Tasklist](https://docs.camunda.io/docs/components/tasklist/overview/) are web-based applications. 
While they offer comprehensive features, there are scenarios where a CLI tool can be more efficient, 
especially for automation, scripting, and quick interactions.

Camunder aims to fill this gap by providing alongside simple commands like `get`, `cancel` or `delete` 
also bulk operations to manage multiple resources at once and composed commands to chain multiple operations together 
like cancelling prior to deleting process instances.

## Features implemented

- get cluster topology: `$ camunder get cluster-topology`
- get process definitions: `$ camunder get process-definition`
- delete single process instance with cancel option: `$ camunder delete process-instance --key <key>`

## Features planned

- run process instances (bulk): `$ camunder run process-instance --bpmn-process-id <bpmnProcessId> --variables <key1=value1,key2=value2,...> --count <number>`
- get process instances: `$ camunder get process-instance`
- cancel process instance (bulk): `$ camunder cancel process-instance --bulk <key1,key2,...>`
- delete process instance (bulk): `$ camunder delete process-instance --bulk <key1,key2,...>`
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
  delete      Delete a resource of a given type by its key. Supported resource types are: process-instance (pi)
  get         List resources of a defined type. Supported resource types are: cluster-topology (ct), process-definition (pd)
  help        Help about any command
  version     Print version information

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

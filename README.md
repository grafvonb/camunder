# Camunder
Camunder is a CLI tool to interact with Camunda 8.

> **Disclaimer:** This project is currently in an **early development stage**.  
> Features, APIs, and behavior may change without notice. Use at your own risk.

## Overview

**Camunder** is a command-line interface (CLI) tool designed to interact with [Camunda 8](https://camunda.com/platform/), a popular workflow and decision automation platform.
It provides a convenient way to manage and monitor Camunda 8 resources directly from the terminal.

## Features implemented

- get cluster topology: `camunder get cluster-topology`

## Features planned

- get process instances: `camunder get process-instances`
- get process definitions: `camunder get process-definitions`
- cancel process instances (bulk): `camunder cancel process-instances --ids <id1,id2,...>`
- delete process instances (bulk): `camunder delete process-instances --ids <id1,id2,...>`
- ...and more to come!

## Usage 
```
Camunder is a CLI tool to interact with Camunda 8.

Usage:
  camunder [flags]
  camunder [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  get         List resources of a defined type e.g. cluster-topology, process-instances etc.
  help        Help about any command

Flags:
      --base-url string    API base URL
      --config string      Path to config file
  -h, --help               help for camunder
      --timeout duration   HTTP timeout (e.g. 10s, 1m)
      --token string       API bearer token

Use "camunder [command] --help" for more information about a command.
```

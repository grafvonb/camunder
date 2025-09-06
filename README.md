<p align="center">
<img src="./docs/logo/camunder-logo-orange-background_170x170.png" alt="camunder logo" style="border-radius: 12px;" />
</p>

# Camunder

Camunder is a CLI tool to interact with [Camunda 8](https://camunda.com/platform/).

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

- get process instances for a process definition id: `$ camunder get process-instance --bpmn-process-id <id>`
- get process definitions: `$ camunder get process-definition`
- get cluster topology: `$ camunder get cluster-topology`
- delete single process instance with cancel option: `$ camunder delete process-instance --key <key>`

## Features planned

- run process instances (bulk)
- cancel process instance (bulk)
- delete process instance (bulk)
- ...and more to come!

## Configuration

Camunder can be configured using command-line flags, environment variables or a configuration file. 
Default configuration file locations are application path or user config directory:
- `$HOME/.config/camunder/config.yaml` (Linux)
- `$HOME/Library/Application Support/camunder/config.yaml` (macOS)
- `%APPDATA%\camunder\config.yaml` (Windows)

## Usage 
### Help Output
```
Usage:
  camunder [flags]
  camunder [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  delete      Delete a resource of a given type by its key. Supported resource types are: process-instance (pi)
  get         List resources of a defined type. Supported resource types are: cluster-topology (ct), process-definition (pd), process-instance (pi)
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
### Camunder in Action
In the following example Camunder is used to list process definitions, get process instances for a specific process definition id,
and tries to delete a process instance by its key. As the process instance is still active, the deletion fails, so the process instance is cancelled first and then deleted successfully.
```bash
# get process definitions
$ camunder get pd
{
  "items": [
    {
      "bpmnProcessId": "DoSomethingUserTaskProcess2ID",
      "key": 2251799813685251,
      "name": "Do Something User Task Process",
      "tenantId": "\u003cdefault\u003e",
      "version": 1,
      "versionTag": "v1.0.1"
    },
    {
      "bpmnProcessId": "DoSomethingUserTaskProcessID",
      "key": 2251799813685252,
      "name": "Do Something User Task Process",
      "tenantId": "\u003cdefault\u003e",
      "version": 1,
      "versionTag": "v1.0.1"
    }
  ],
  "sortValues": [
    "2251799813685252"
  ],
  "total": 2
}
# get process instances for a specific process definition id
$ camunder get pi --bpmn-process-id=DoSomethingUserTaskProcess2ID
{
  "items": [
    {
      "bpmnProcessId": "DoSomethingUserTaskProcess2ID",
      "incident": false,
      "key": 2251799813685391,
      "processDefinitionKey": 2251799813685251,
      "processVersion": 1,
      "processVersionTag": "v1.0.1",
      "startDate": "2025-09-06T06:48:08.625+0000",
      "state": "ACTIVE",
      "tenantId": "\u003cdefault\u003e"
    },
    {
      "bpmnProcessId": "DoSomethingUserTaskProcess2ID",
      "incident": false,
      "key": 2251799813685415,
      "processDefinitionKey": 2251799813685251,
      "processVersion": 1,
      "processVersionTag": "v1.0.1",
      "startDate": "2025-09-06T06:52:13.069+0000",
      "state": "ACTIVE",
      "tenantId": "\u003cdefault\u003e"
    }
  ],
  "sortValues": [
    2251799813685415
  ],
  "total": 2
}
# try to delete a process instance by its key
$ camunder delete pi --key 2251799813685391
Trying to delete process instance with key 2251799813685391...
Error deleting process instance with key 2251799813685391: unexpected status 400: {"status":400,"message":"Process instances needs to be in one of the states [COMPLETED, CANCELED]","instance":"d24bb589-6c62-4985-a5cb-0712d1e31152","type":"Invalid request"}
# as the process instance is still active, try to cancel it first and then delete it
$ camunder delete pi --key 2251799813685391 -c
Process instance with key 2251799813685391 not in state COMPLETED or CANCELED, cancelling it first...
Trying to cancel process instance with key 2251799813685391...
Process instance with key 2251799813685391 was successfully cancelled
Waiting for process instance with key 2251799813685391 to be cancelled by workflow engine...
{
  "deleted": 1,
  "message": "Process instance and dependant data deleted for key '2251799813685391'"
}
```
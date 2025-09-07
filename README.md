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
- cancel single process instance: `$ camunder cancel process-instance --key <key>`

## Features planned

- run process instances (bulk)
- cancel process instance (bulk)
- delete process instance (bulk)
- ...and more to come!

## Supported Camunda 8 APIs

- 8.7.x

## Configuration

Camunder uses [Viper](https://github.com/spf13/viper) under the hood.\
Configuration values can come from:

-   **Flags** (`--auth-client-id=...`)\
-   **Environment variables** (`CAMUNDER_AUTH_CLIENT_ID=...`)\
-   **Config file** (YAML)\
-   **Defaults** (hardcoded fallbacks)

### Precedence

When multiple sources define the same setting, the **highest-priority value wins**:

| Priority    | Source             | Example                          |
|-------------|--------------------|----------------------------------|
| 1 (highest) | Command-line flags | `--auth-client-id=cli-id`        |
| 2           | Environment vars   | `CAMUNDER_AUTH_CLIENT_ID=env-id` |
| 3           | Config file (YAML) | `auth.client_id: file-id`        |
| 4 (lowest)  | Defaults           | `http.timeout: "30s"` (built-in) |

### Default configuration file locations

When searching for a config file, Camunder checks these paths in order and uses the first one it finds:

| Priority | Location                                        | Notes                                                                                                                              |
|----------|-------------------------------------------------|------------------------------------------------------------------------------------------------------------------------------------|
| 1        | `./config.yaml`                                 | Current working directory                                                                                                          |
| 2        | `$XDG_CONFIG_HOME/camunder/config.yaml`         | Skipped if `$XDG_CONFIG_HOME` is not set                                                                                           |
| 3        | `$HOME/.config/camunder/config.yaml`            | XDG default on Linux/macOS                                                                                                         |
| 4        | `$HOME/.camunder/config.yaml`                   | Legacy fallback                                                                                                                    |
| 5        | `%AppData%\camunder\config.yaml` (Windows only) | `%AppData%` usually expands to `C:\Users\<User>\AppData\Roaming`<br>Example: `C:\Users\Alice\AppData\Roaming\camunder\config.yaml` |

### File format

Config files must be **YAML**. Example:

``` yaml
auth:
  # OAuth token endpoint
  token_url: "http://localhost:18080/auth/realms/camunda-platform/protocol/openid-connect"

  # Client credentials (use env vars if possible)
  client_id: "camunder"
  client_secret: "*******"

  # Scopes as key:value pairs (names -> scope strings)
  # Do not define if not in use or empty
  scopes:
    camunda8_api: "profile"
    operate_api: "profile"
    tasklist_api: "profile"

http:
  # Go duration string (e.g., 10s, 1m, 2m30s)
  timeout: "30s"

apis:
  # Base URLs for your endpoints
  camunda8_api:
    base_url: "http://localhost:8080/v2"
  operate_api:
    base_url: "http://localhost:8081/v1"
  tasklist_api:
    base_url: "http://localhost:8082/v1"
```

### Environment variables

Each config key can also be set via environment variable.\
The prefix is `CAMUNDER_`, and nested keys are joined with `_`. For
example:

-   `CAMUNDER_AUTH_CLIENT_ID`
-   `CAMUNDER_AUTH_CLIENT_SECRET`
-   `CAMUNDER_HTTP_TIMEOUT`

### Security note

Sensitive fields such as `auth.client_secret` are **always masked** when
the configuration is printed (e.g. with `--show-config`) or logged.\
The raw values are still loaded and used internally, but they will never
appear in output.

### Example: Show effective configuration

You can inspect the effective configuration (after merging defaults,
config file, env vars, and flags) with:

```bash
camunder --show-config
```

Example output:

```json
{
  "Config": "",
  "Auth": {
    "TokenURL": "http://localhost:18080/auth/realms/camunda-platform/protocol/openid-connect",
    "ClientID": "******",
    "ClientSecret": "******",
    "Scopes": {
      "camunda8_api": "profile",
      "operate_api": "profile",
      "tasklist_api": "profile"
    }
  },
  "APIs": {
    "Camunda8": {
      "Key": "camunda8_api",
      "BaseURL": "http://localhost:8086/v2"
    },
    "Operate": {
      "Key": "operate_api",
      "BaseURL": "http://localhost:8081/v1"
    },
    "Tasklist": {
      "Key": "tasklist_api",
      "BaseURL": "http://localhost:8082/v1"
    }
  },
  "HTTP": {
    "Timeout": "30s"
  }
}
```
## Usage 
### Help Output
```
Usage:
  camunder [flags]
  camunder [command]

Available Commands:
  cancel      Cancel a resource of a given type by its key. Supported resource types are: process-instance (pi)
  completion  Generate the autocompletion script for the specified shell
  delete      Delete a resource of a given type by its key. Supported resource types are: process-instance (pi)
  get         List resources of a defined type. Supported resource types are: cluster-topology (ct), process-definition (pd), process-instance (pi)
  help        Help about any command
  version     Print version information

Flags:
      --auth-client-id string        auth client ID
      --auth-client-secret string    auth client secret
      --auth-scopes stringToString   auth scopes as key=value (repeatable or comma-separated) (default [])
      --auth-token-url string        auth token URL
      --camunda8-base-url string     Camunda8 API base URL
      --config string                path to config file
  -h, --help                         help for camunder
      --http-timeout string          HTTP timeout (Go duration, e.g. 30s)
      --operate-base-url string      Operate API base URL
  -q, --quiet                        suppress output, use exit code only
      --show-config                  print effective config (secrets redacted)
      --tasklist-base-url string     Tasklist API base URL

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

  "message": "Process instance and dependant data deleted for key '2251799813685391'"
}
```
Copyright © 2025 Adam Bogdan Boczek | [boczek.info](https://boczek.info)
<p align="center">
<img src="./docs/logo/camunder-logo-orange-background_170x170.png" alt="camunder logo" style="border-radius: 12px;" />
</p>

# Camunder – a CLI for Camunda 8

**Camunder** is a CLI (command-line interface) for Camunda 8 that gives developers and operators faster, scriptable management of Camunda resources.
It complements Camunda's [Operate](https://docs.camunda.io/docs/components/operate/overview/) and [Tasklist](https://docs.camunda.io/docs/components/tasklist/overview/) by enabling automation, bulk operations, and integration into existing workflows and pipelines.

While Operate and Tasklist cover most use cases via web interfaces, a CLI can be more efficient for automation, scripting, and quick operational tasks.  

**Camunder** fills this gap with commands such as `get`, `cancel`, and `delete`, 
as well as specialized use cases like *[deleting active process instances by canceling it first](#deleting-an-active-process-instance-by-cancelling-it-first)* 
or *[finding process instances with orphan parent process instances](#finding-process-instances-with-orphan-parent-process-instances)*, 
which simplify recurring administration and maintenance of Camunda 8 process instances. 

See [Camunder in Action](#camunder-in-action) for more examples.
 
## Highlights

Camunder simplifies various tasks related to Camunda 8, including these special use cases:
- deleting active process instances by cancelling them first,
- finding process instances with orphan parent process instances,
- recursive search process instances with parent-child relationships:
  - list all child process instances of a given process instance,
  - list path from a given process instance to its root ancestor (top-level parent),
  - list the entire family (parent, grandparent, ...) of a given process instance (traverse up and down the tree),
- bulk cancelling or deleting process instances (not implemented yet),
- and more to come...

## Supported Camunda 8 APIs

- 8.7.x

## Configuration

Camunder uses [Viper](https://github.com/spf13/viper) under the hood.
Configuration values can come from:

-   **Flags** (`--auth-client-id=...`)
-   **Environment variables** (`CAMUNDER_AUTH_CLIENT_ID=...`)
-   **Config file** (YAML)
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
app:
  backoff:
    strategy: exponential
    initial_delay: 500ms
    max_delay: 8s
    max_retries: 0
    multiplier: 2.0
    timeout: 2m

auth:
  # OAuth token endpoint
  token_url: "http://localhost:18080/auth/realms/camunda-platform/protocol/openid-connect"

  # Client credentials (use env vars if possible)
  client_id: "camunder"
  client_secret: ""

  # Scopes as key:value pairs (names -> scope strings)
  # Do not define if not in use or empty
  scopes:
    camunda_api: "profile"
    operate_api: "profile"
    tasklist_api: "profile"

http:
  # Go duration string (e.g., 10s, 1m, 2m30s)
  timeout: "30s"

apis:
  # Base URLs for your endpoints
  camunda_api:
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
$ camunder --show-config
config loaded: /Users/adam.boczek/.camunder/config.yaml
{
  "Config": "",
  "App": {
    "Tenant": "",
    "Backoff": {
      "Strategy": "exponential",
      "InitialDelay": 500000000,
      "MaxDelay": 8000000000,
      "MaxRetries": 0,
      "Multiplier": 2,
      "Timeout": 120000000000
    }
  },
  "Auth": {
    "TokenURL": "http://localhost:18080/auth/realms/camunda-platform/protocol/openid-connect",
    "ClientID": "******",
    "ClientSecret": "******",
    "Scopes": {
      "camunda_api": "profile",
      "operate_api": "profile",
      "tasklist_api": "profile"
    }
  },
  "APIs": {
    "Camunda": {
      "Key": "camunda_api",
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

```bash
$ camunder --help
Camunder is a CLI tool to interact with Camunda 8.

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
  walk        Traverse (walk) the parent/child graph process instances.

Flags:
      --auth-client-id string        auth client ID
      --auth-client-secret string    auth client secret
      --auth-scopes stringToString   auth scopes as key=value (repeatable or comma-separated) (default [])
      --auth-token-url string        auth token URL
      --camunda-base-url string     Camunda8 API base URL
      --config string                path to config file
  -h, --help                         help for camunder
      --http-timeout string          HTTP timeout (Go duration, e.g. 30s)
      --operate-base-url string      Operate API base URL
      --quiet                        suppress output, use exit code only
      --show-config                  print effective config (secrets redacted)
      --tasklist-base-url string     Tasklist API base URL
      --tenant string                default tenant ID

Use "camunder [command] --help" for more information about a command.
```
### Camunder in Action

#### Deleting an active process instance by cancelling it first
```bash
$ camunder get pi --bpmn-process-id=C87SimpleUserTask_Process --one-line
found: 12
2251799813685511 dev01 C87SimpleUserTask_Process v1/v1.0.0 ACTIVE s:2025-09-09T12:14:30.380+0000 i:false
2251799813685518 dev01 C87SimpleUserTask_Process v1/v1.0.0 ACTIVE s:2025-09-09T12:14:36.618+0000 i:false
2251799813685525 dev01 C87SimpleUserTask_Process v1/v1.0.0 ACTIVE s:2025-09-09T12:14:40.675+0000 i:false
2251799813685532 dev01 C87SimpleUserTask_Process v1/v1.0.0 ACTIVE s:2025-09-09T12:14:44.338+0000 i:false
2251799813685541 dev01 C87SimpleUserTask_Process v2 ACTIVE s:2025-09-09T16:11:52.976+0000 i:false
2251799813685548 dev01 C87SimpleUserTask_Process v2 ACTIVE s:2025-09-09T16:13:30.653+0000 i:false
2251799813685556 dev01 C87SimpleUserTask_Process v2 ACTIVE s:2025-09-09T16:13:53.060+0000 i:false
2251799813685571 dev01 C87SimpleUserTask_Process v2 CANCELED s:2025-09-09T19:10:28.190+0000 e:2025-09-10T20:36:44.990+0000 p:2251799813685566 i:false
2251799813685582 dev01 C87SimpleUserTask_Process v2 CANCELED s:2025-09-09T19:10:33.364+0000 e:2025-09-09T20:27:55.530+0000 p:2251799813685577 i:false
2251799813685595 dev01 C87SimpleUserTask_Process v2 ACTIVE s:2025-09-09T22:01:27.621+0000 p:2251799813685590 i:false
2251799813685606 dev01 C87SimpleUserTask_Process v2 ACTIVE s:2025-09-09T22:01:33.533+0000 p:2251799813685601 i:false
2251799813685614 dev01 C87SimpleUserTask_Process v3 ACTIVE s:2025-09-10T10:03:12.700+0000 i:false
$ camunder get pi --bpmn-process-id=C87SimpleUserTask_Process --one-line --state=active
filter: state=active
found: 10
2251799813685511 dev01 C87SimpleUserTask_Process v1/v1.0.0 ACTIVE s:2025-09-09T12:14:30.380+0000 i:false
2251799813685518 dev01 C87SimpleUserTask_Process v1/v1.0.0 ACTIVE s:2025-09-09T12:14:36.618+0000 i:false
2251799813685525 dev01 C87SimpleUserTask_Process v1/v1.0.0 ACTIVE s:2025-09-09T12:14:40.675+0000 i:false
2251799813685532 dev01 C87SimpleUserTask_Process v1/v1.0.0 ACTIVE s:2025-09-09T12:14:44.338+0000 i:false
2251799813685541 dev01 C87SimpleUserTask_Process v2 ACTIVE s:2025-09-09T16:11:52.976+0000 i:false
2251799813685548 dev01 C87SimpleUserTask_Process v2 ACTIVE s:2025-09-09T16:13:30.653+0000 i:false
2251799813685556 dev01 C87SimpleUserTask_Process v2 ACTIVE s:2025-09-09T16:13:53.060+0000 i:false
2251799813685595 dev01 C87SimpleUserTask_Process v2 ACTIVE s:2025-09-09T22:01:27.621+0000 p:2251799813685590 i:false
2251799813685606 dev01 C87SimpleUserTask_Process v2 ACTIVE s:2025-09-09T22:01:33.533+0000 p:2251799813685601 i:false
2251799813685614 dev01 C87SimpleUserTask_Process v3 ACTIVE s:2025-09-10T10:03:12.700+0000 i:false
$ camunder delete pi --key 2251799813685511
trying to delete process instance with key 2251799813685511...
Error deleting process instance with key 2251799813685511: unexpected status 400: {"status":400,"message":"Process instances needs to be in one of the states [COMPLETED, CANCELED]","instance":"dae2c2ce-58dd-4396-a948-4d57463168ed","type":"Invalid request"}
$ camunder delete pi --key 2251799813685511 --cancel
trying to delete process instance with key 2251799813685511...
process instance with key 2251799813685511 not in state COMPLETED or CANCELED, cancelling it first...
trying to cancel process instance with key 2251799813685511...
process instance with key 2251799813685511 was successfully cancelled
waiting for process instance with key 2251799813685511 to be cancelled by workflow engine...
process instance "2251799813685511" currently in state "ACTIVE"; waiting...
process instance "2251799813685511" currently in state "ACTIVE"; waiting...
process instance "2251799813685511" currently in state "ACTIVE"; waiting...
process instance "2251799813685511" currently in state "ACTIVE"; waiting...
process instance "2251799813685511" currently in state "ACTIVE"; waiting...
process instance "2251799813685511" reached desired state "CANCELED"
process instance with key 2251799813685511 was successfully deleted
{
  "deleted": 1,
  "message": "Process instance and dependant data deleted for key '2251799813685511'"
}
```
#### Finding process instances with orphan parent process instances
```bash
$ camunder get pi --bpmn-process-id=C87SimpleUserTask_Process --one-line
found: 12
2251799813685511 dev01 C87SimpleUserTask_Process v1/v1.0.0 ACTIVE s:2025-09-09T12:14:30.380+0000 i:false
2251799813685518 dev01 C87SimpleUserTask_Process v1/v1.0.0 ACTIVE s:2025-09-09T12:14:36.618+0000 i:false
2251799813685525 dev01 C87SimpleUserTask_Process v1/v1.0.0 ACTIVE s:2025-09-09T12:14:40.675+0000 i:false
2251799813685532 dev01 C87SimpleUserTask_Process v1/v1.0.0 ACTIVE s:2025-09-09T12:14:44.338+0000 i:false
2251799813685541 dev01 C87SimpleUserTask_Process v2 ACTIVE s:2025-09-09T16:11:52.976+0000 i:false
2251799813685548 dev01 C87SimpleUserTask_Process v2 ACTIVE s:2025-09-09T16:13:30.653+0000 i:false
2251799813685556 dev01 C87SimpleUserTask_Process v2 ACTIVE s:2025-09-09T16:13:53.060+0000 i:false
2251799813685571 dev01 C87SimpleUserTask_Process v2 CANCELED s:2025-09-09T19:10:28.190+0000 e:2025-09-10T20:36:44.990+0000 p:2251799813685566 i:false
2251799813685582 dev01 C87SimpleUserTask_Process v2 CANCELED s:2025-09-09T19:10:33.364+0000 e:2025-09-09T20:27:55.530+0000 p:2251799813685577 i:false
2251799813685595 dev01 C87SimpleUserTask_Process v2 ACTIVE s:2025-09-09T22:01:27.621+0000 p:2251799813685590 i:false
2251799813685606 dev01 C87SimpleUserTask_Process v2 ACTIVE s:2025-09-09T22:01:33.533+0000 p:2251799813685601 i:false
2251799813685614 dev01 C87SimpleUserTask_Process v3 ACTIVE s:2025-09-10T10:03:12.700+0000 i:false
$ camunder get pi --bpmn-process-id=C87SimpleUserTask_Process --one-line --children-only
filter: children-only=true
found: 4
2251799813685571 dev01 C87SimpleUserTask_Process v2 CANCELED s:2025-09-09T19:10:28.190+0000 e:2025-09-10T20:36:44.990+0000 p:2251799813685566 i:false
2251799813685582 dev01 C87SimpleUserTask_Process v2 CANCELED s:2025-09-09T19:10:33.364+0000 e:2025-09-09T20:27:55.530+0000 p:2251799813685577 i:false
2251799813685595 dev01 C87SimpleUserTask_Process v2 ACTIVE s:2025-09-09T22:01:27.621+0000 p:2251799813685590 i:false
2251799813685606 dev01 C87SimpleUserTask_Process v2 ACTIVE s:2025-09-09T22:01:33.533+0000 p:2251799813685601 i:false
$ camunder get pi --bpmn-process-id=C87SimpleUserTask_Process --one-line --orphan-parents-only
filter: orphan-parents-only=true
found: 2
2251799813685571 dev01 C87SimpleUserTask_Process v2 CANCELED s:2025-09-09T19:10:28.190+0000 e:2025-09-10T20:36:44.990+0000 p:2251799813685566 i:false
2251799813685582 dev01 C87SimpleUserTask_Process v2 CANCELED s:2025-09-09T19:10:33.364+0000 e:2025-09-09T20:27:55.530+0000 p:2251799813685577 i:false
```
#### Listing process instances for a specific process definition (model) and its first version
```bash
$ camunder get pi --bpmn-process-id=C87SimpleUserTask_Process --one-line
found: 11
2251799813685518 dev01 C87SimpleUserTask_Process v1/v1.0.0 ACTIVE s:2025-09-09T12:14:36.618+0000 i:false
2251799813685525 dev01 C87SimpleUserTask_Process v1/v1.0.0 ACTIVE s:2025-09-09T12:14:40.675+0000 i:false
2251799813685532 dev01 C87SimpleUserTask_Process v1/v1.0.0 ACTIVE s:2025-09-09T12:14:44.338+0000 i:false
2251799813685541 dev01 C87SimpleUserTask_Process v2 ACTIVE s:2025-09-09T16:11:52.976+0000 i:false
2251799813685548 dev01 C87SimpleUserTask_Process v2 ACTIVE s:2025-09-09T16:13:30.653+0000 i:false
2251799813685556 dev01 C87SimpleUserTask_Process v2 ACTIVE s:2025-09-09T16:13:53.060+0000 i:false
2251799813685571 dev01 C87SimpleUserTask_Process v2 CANCELED s:2025-09-09T19:10:28.190+0000 e:2025-09-10T20:36:44.990+0000 p:2251799813685566 i:false
2251799813685582 dev01 C87SimpleUserTask_Process v2 CANCELED s:2025-09-09T19:10:33.364+0000 e:2025-09-09T20:27:55.530+0000 p:2251799813685577 i:false
2251799813685595 dev01 C87SimpleUserTask_Process v2 ACTIVE s:2025-09-09T22:01:27.621+0000 p:2251799813685590 i:false
2251799813685606 dev01 C87SimpleUserTask_Process v2 ACTIVE s:2025-09-09T22:01:33.533+0000 p:2251799813685601 i:false
2251799813685614 dev01 C87SimpleUserTask_Process v3 ACTIVE s:2025-09-10T10:03:12.700+0000 i:false
$ camunder get pi --bpmn-process-id=C87SimpleUserTask_Process --one-line --process-version=1
found: 3
2251799813685518 dev01 C87SimpleUserTask_Process v1/v1.0.0 ACTIVE s:2025-09-09T12:14:36.618+0000 i:false
2251799813685525 dev01 C87SimpleUserTask_Process v1/v1.0.0 ACTIVE s:2025-09-09T12:14:40.675+0000 i:false
2251799813685532 dev01 C87SimpleUserTask_Process v1/v1.0.0 ACTIVE s:2025-09-09T12:14:44.338+0000 i:false
```

Copyright © 2025 Adam Bogdan Boczek | [boczek.info](https://boczek.info)

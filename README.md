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

## Table of Contents
- [Camunder – a CLI for Camunda 8](#camunder--a-cli-for-camunda-8)
- [Quick Start with Camunda 8 Run](#quick-start-with-camunda-8-run)
- [Highlights](#highlights)
- [Supported Camunda 8 APIs](#supported-camunda-8-apis)
- [Configuration](#configuration)
    - [Choose authentication method](#choose-authentication-method)
        - [Authentication with OAuth2 (OIDC)](#authentication-with-oauth2-oidc)
        - [Authentication with API Cookie (development with Camunda 8 Run only)](#authentication-with-api-cookie-development-with-camunda-8-run-only)
    - [Connecting to Camunda 8 APIs](#connecting-to-camunda-8-apis)
        - [If you use Camunda 8 Run with API Cookie Authentication (Development only)](#if-you-use-camunda-8-run-with-api-cookie-authentication-development-only)
        - [If you use Camunda 8 with OAuth2 (oidc) authentication](#if-you-use-camunda-8-with-oauth2-oidc-authentication)
    - [Ways to provide settings](#ways-to-provide-settings)
    - [Precedence](#precedence)
    - [Default configuration file locations](#default-configuration-file-locations)
    - [File format](#file-format)
    - [Environment variables](#environment-variables)
    - [Security note](#security-note)
    - [Example: Show effective configuration](#example-show-effective-configuration)
- [Usage Help](#usage-help)
- [Camunder in Action](#camunder-in-action)
    - [Deleting an active process instance by cancelling it first](#deleting-an-active-process-instance-by-cancelling-it-first)
    - [Finding process instances with orphan parent process instances](#finding-process-instances-with-orphan-parent-process-instances)
    - [Listing process instances for a specific process definition (model) and its first version](#listing-process-instances-for-a-specific-process-definition-model-and-its-first-version)

 
## Quick Start with Camunda 8 Run

1. **Install Camunda 8.7 Run**  
   Download Camunda 8 Run from the [Camunda Releases](https://github.com/camunda/camunda/releases?q=%22c8run-8.7%22&expanded=true) page, unpack and start it with `./start.sh`.
2. **Install Camunder**  
   Download the latest release from the [Camunder Releases](https://github.com/grafvonb/camunder/releases) page and unpack it.
3. **Configure Camunder**  
   Create a configuration file (YAML) in the folder where you unpacked Camunder with the name `config.yaml` or in `$HOME/.camunder/config.yaml` with the minimal connection and authentication details:
    ```yaml
    auth:
      mode: "cookie"
      cookie:
        base_url: "http://localhost:8080"
    
    apis:
      version: "87"
      camunda_api:
        base_url: "http://localhost:8080/v2"
    ```
4. **Run Camunder**  
   Test the connection and list cluster topology:
   ```bash
   ./camunder get cluster-topology
   ```
   or use explicit path to config file:
   ```bash
   ./camunder get cluster-topology --config ./config-minimal.yaml
   ```
   You should see output like this:
    ```json
    {
      "Brokers": [
        {
          "Host": "192.168.178.88",
          "NodeId": 0,
          "Partitions": [
            {
              "Health": "healthy",
              "PartitionId": 1,
              "Role": "leader"
            }
          ],
          "Port": 26501,
          "Version": "8.7.12"
        }
      ],
      "ClusterSize": 1,
      "GatewayVersion": "8.7.12",
      "PartitionsCount": 1,
      "ReplicationFactor": 1,
      "LastCompletedChangeId": ""
    }  
    ```

## Highlights

Camunder simplifies various tasks related to Camunda 8, including these special use cases:

- **Delete active process instances by cancelling them first**
  ```bash
  ./camunder delete pi --key <process-instance-key> --cancel
  ```

- **List process instances that are children (sub-processes) of other process instances**
  ```bash
  ./camunder get pi --bpmn-process-id=<bpmn-process-id> --children-only
  ```

- **List process instances that are parents of other process instances**
  ```bash
  ./camunder get pi --bpmn-process-id=<bpmn-process-id> --parents-only
  ```

- **List process instances that are children of orphan parent process instances**  
  (i.e., their parent process instance no longer exists)
  ```bash
  ./camunder get pi --bpmn-process-id=<bpmn-process-id> --orphan-parents-only
  ```

- **List process instances for a specific process definition (model) and its first version**
  ```bash
  ./camunder get pi --bpmn-process-id=<bpmn-process-id> --process-version=1
  ```

- **List process instances with incidents**
  ```bash
  ./camunder get pi --incidents-only
  ```

- **List process instances without incidents**
  ```bash
  ./camunder get pi --no-incidents-only
  ```

- **Recursive search (walk) process instances with parent–child relationships**
    - List all child process instances of a given process instance
      ```bash
      ./camunder walk pi --mode children --start-key <process-instance-key>
      ```
    - List path from a given process instance to its root ancestor (top-level parent)
      ```bash
      ./camunder walk pi --mode parent --start-key <process-instance-key>
      ```
    - List the entire family (parent, grandparent, …) of a given process instance (traverse up and down the tree)
      ```bash
      ./camunder walk pi --mode family --start-key <process-instance-key>
      ```

- **List process instances in one line per instance (suitable for scripting)**  
  Works with all `get` commands.
  ```bash
  ./camunder get pi --one-line
  ```

- **List process instances just by their keys (suitable for scripting)**  
  Works with all `get` commands.
  ```bash
  ./camunder get pi --keys-only
  ```

- …and more to come:
- bulk operations (e.g., delete multiple process instances by filter)
- multiple Camunda 8 API versions support (currently 8.7, 8.8 to come)
- or submit a proposal or contribute code on [GitHub](https://github.com/grafvonb/camunder)

## Supported Camunda 8 APIs

- 8.7.x
- 8.8.x (to come after release in October 2025)

## Configuration

### Choose authentication method

Camunder supports two authentication methods for connecting to Camunda 8 APIs:
- OAuth2 (OIDC)
- API Cookie (development with Camunda 8 Run only)

#### Authentication with OAuth2 (OIDC)

Camunder supports OAuth2 (OIDC) authentication with client credentials flow.
You need to provide the following settings:
* Token URL
* Client ID
* Client Secret
* Scopes (optional, depending on your identity provider and API setup)

Here is an example configuration snippet for OAuth2 authentication for Camunda 8 running locally with Keycloak:
```yaml
auth:
  mode: "oauth2" # options: "oauth2", "cookie"
  oauth2:
    token_url: "http://localhost:18080/auth/realms/camunda-platform/protocol/openid-connect"
    client_id: "camunder"
    client_secret: "*******" # use environment variable CAMUNDER_AUTH_CLIENT_SECRET if possible
    scopes:
      camunda_api: "profile"
      operate_api: "profile"
      tasklist_api: "profile"
```

#### Authentication with API Cookie (development with Camunda 8 Run only)

This method is only suitable for local development with Camunda 8 Run, as it uses the API cookie set by the web interface.
You need to provide the following settings:
* API base URL
* (optional) Username and Password (if not set, defaults to `demo`/`demo`)

Here is an example configuration snippet for API Cookie authentication for Camunda 8 Run running locally:
```yaml
auth:
  mode: "cookie" # options: "oauth2", "cookie"
  cookie:
    base_url: "http://localhost:8090"
    username: "demo"
    password: "demo"
```

### Connecting to Camunda 8 APIs

To run Camunder, you need to configure the connection to your Camunda 8 APIs (Camunda, Operate, Tasklist) and authentication details.
With the introduction of Camunda 8.8 and unified APIs, some of these settings may become optional or redundant in the future.

Camunder expects the following API configurations:
* Camunda 8 API (required, formally known as Zeebe API)
* Operate API (optional, required for some commands, if not set defaults to Camunda 8 API)
* Tasklist API (optional, required for some commands, if not set defaults to Camunda 8 API)

#### If you use Camunda 8 Run with API Cookie Authentication (Development only)

After starting Camunda 8 Run, you can find the API endpoint in the terminal output (default is `http://localhost:8080/v2`). 
(current, October 2025, gotcha: if you use `--port` flag, the terminal output still shows port 8080, but it actually runs on the port you specified).

Provide the following settings in the config file, environment variables, or flags:
```yaml
apis:
  version: "87"
  camunda_api:
    base_url: "http://localhost:8080/v2"
  operate_api:
    base_url: "http://localhost:8080"
  tasklist_api:
    base_url: "http://localhost:8080"
```

#### If you use Camunda 8 with OAuth2 (oidc) authentication

### Ways to provide settings

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
$ ./camunder --show-config
config loaded: /Users/adam.boczek/Development/Workspace/Boczek/Projects/camunder/camunder/config.yaml
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
      "Timeout": 30000000000
    }
  },
  "Auth": {
    "Mode": "",
    "OAuth2": {
      "TokenURL": "http://localhost:18080/auth/realms/camunda-platform/protocol/openid-connect",
      "ClientID": "******",
      "ClientSecret": "******",
      "Scopes": {
        "camunda_api": "profile",
        "operate_api": "profile",
        "tasklist_api": "profile"
      }
    },
    "Cookie": {
      "BaseURL": "http://localhost:8090",
      "Username": "******",
      "Password": "******"
    }
  },
  "APIs": {
    "Version": "87",
    "Camunda": {
      "Key": "camunda_api",
      "BaseURL": "http://localhost:8086/v2"
    },
    "Operate": {
      "Key": "operate_api",
      "BaseURL": "http://localhost:8081"
    },
    "Tasklist": {
      "Key": "tasklist_api",
      "BaseURL": "http://localhost:8082"
    }
  },
  "HTTP": {
    "Timeout": "23s"
  }
}
```

## Usage Help
```bash
$ ./camunder help
Camunder is a CLI tool to interact with Camunda 8.

Usage:
  camunder [flags]
  camunder [command]

Available Commands:
  cancel      Cancel a resource of a given type by its key. Supported resource types are: process-instance (pi)
  completion  Generate the autocompletion script for the specified shell
  delete      Delete a resource of a given type by its key. Supported resource types are: process-instance (pi)
  expect      Expect a resource of a given type to change (e.g. its state) by its key. Supported resource types are: process-instance (pi)
  get         List resources of a resource type. Supported resource types are: cluster-topology (ct), process-definition (pd), process-instance (pi)
  help        Help about any command
  version     Print version information
  walk        Traverse (walk) the parent/child graph of resource type. Supported resource types are: process-instance (pi)

Flags:
      --auth-client-id string         auth client ID
      --auth-client-secret string     auth client secret
      --auth-scopes stringToString    auth scopes as key=value (repeatable or comma-separated) (default [])
      --auth-token-url string         auth token URL
  -a, --camunda-apis-version string   Camunda API version (supported: [8.7 8.8]) (default "8.7")
      --camunda-base-url string       Camunda API base URL
      --config string                 path to config file
  -h, --help                          help for camunder
      --http-timeout string           HTTP timeout (Go duration, e.g. 30s)
      --log-format string             log format (json, plain, text) (default "plain")
      --log-level string              log level (debug, info, warn, error) (default "info")
      --log-with-source               include source file and line number in logs
      --operate-base-url string       Operate API base URL
      --show-config                   print effective config (secrets redacted)
      --tasklist-base-url string      Tasklist API base URL
      --tenant string                 default tenant ID

Use "camunder [command] --help" for more information about a command.
```

### Camunder in Action
Look here for practical examples of using Camunder for common tasks and special use cases.

#### Deleting an active process instance by cancelling it first, if the instance is active
List all process instances for a specific process definition:
```
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
```
List only active process instances for the specific process definition in active state:
```
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
```
Try to delete an active process instance (will fail):
```
$ camunder delete pi --key 2251799813685511
trying to delete process instance with key 2251799813685511...
Error deleting process instance with key 2251799813685511: unexpected status 400: {"status":400,"message":"Process instances needs to be in one of the states [COMPLETED, CANCELED]","instance":"dae2c2ce-58dd-4396-a948-4d57463168ed","type":"Invalid request"}
```
Try to delete an active process instance by forcing cancellation first (will succeed):
```
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

#### Finding process instances with orphan (non-existing) parent process instances
List all process instances for a specific process definition:
```
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
```
List all process instances for a specific process definition that are children of other process instances (look at p:...):
```
$ camunder get pi --bpmn-process-id=C87SimpleUserTask_Process --one-line --children-only
filter: children-only=true
found: 4
2251799813685571 dev01 C87SimpleUserTask_Process v2 CANCELED s:2025-09-09T19:10:28.190+0000 e:2025-09-10T20:36:44.990+0000 p:2251799813685566 i:false
2251799813685582 dev01 C87SimpleUserTask_Process v2 CANCELED s:2025-09-09T19:10:33.364+0000 e:2025-09-09T20:27:55.530+0000 p:2251799813685577 i:false
2251799813685595 dev01 C87SimpleUserTask_Process v2 ACTIVE s:2025-09-09T22:01:27.621+0000 p:2251799813685590 i:false
2251799813685606 dev01 C87SimpleUserTask_Process v2 ACTIVE s:2025-09-09T22:01:33.533+0000 p:2251799813685601 i:false
```
List all process instances for a specific process definition that are children of orphan parent process instances (their parent process instance no longer exists):
```
$ camunder get pi --bpmn-process-id=C87SimpleUserTask_Process --one-line --orphan-parents-only
filter: orphan-parents-only=true
found: 2
2251799813685571 dev01 C87SimpleUserTask_Process v2 CANCELED s:2025-09-09T19:10:28.190+0000 e:2025-09-10T20:36:44.990+0000 p:2251799813685566 i:false
2251799813685582 dev01 C87SimpleUserTask_Process v2 CANCELED s:2025-09-09T19:10:33.364+0000 e:2025-09-09T20:27:55.530+0000 p:2251799813685577 i:false
```

#### Listing process instances for a process definition in a specific version
List all process instances for a specific process definition:
```
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
```
List only process instances for version 1 of the process definition:
```
$ camunder get pi --bpmn-process-id=C87SimpleUserTask_Process --one-line --process-version=1
found: 3
2251799813685518 dev01 C87SimpleUserTask_Process v1/v1.0.0 ACTIVE s:2025-09-09T12:14:36.618+0000 i:false
2251799813685525 dev01 C87SimpleUserTask_Process v1/v1.0.0 ACTIVE s:2025-09-09T12:14:40.675+0000 i:false
2251799813685532 dev01 C87SimpleUserTask_Process v1/v1.0.0 ACTIVE s:2025-09-09T12:14:44.338+0000 i:false
```

Copyright © 2025 Adam Bogdan Boczek | [boczek.info](https://boczek.info)

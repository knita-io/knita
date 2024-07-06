# Knita CLI Config

The Knita CLI is configured through a YAML-based config file. The default file location is `~/.knita.yaml`, where `~`
represents your user's home directory, and will vary by platform. The file location can optionally be overridden via the `--config` argument on the CLI.

```yaml
---
executors:
  # Local configures the built-in Executor.
  local:
    # Set to true to disable local builds.
    disabled: false
    # Labels will be advertised to the Broker and are used to match compatible builds to this Executor.
    # The Executor will always advertise its OS and Architecture via built-in labels. These will vary 
    # based on the host system the Executor is deployed to. OS will be one of 'linux', 'darwin' or 'windows',
    # and Architecture will be on of 'amd64', 'arm' or 'arm64'.
    labels:
      - hello
      - world
  # Remote configures an optional set of well-known remote Executors to distribute builds to.
  # Remote Executors are intended to be used by an individual or small team as a convenience
  # as they do not require a separate Queue server to be deployed. 
  remote:
      # Address to connect to the remote Executor on. IP and DNS names are valid.
    - address: 192.168.1.10:9091
      # Set to true to disable this remote Executor.
      disabled: false
```
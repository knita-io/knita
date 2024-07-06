# Knita Executor Config

The Knita Executor is configured through a YAML-based config file. The default file location varies per platform.
The file location can optionally be overridden via the `--config` argument on the Executor.

**Default config file location:**
- **Linux**: `/etc/knita/executor.yaml`
- **MacOS**: `/Library/Application Support/knita/executor.yaml`
- **Windows**: `%ProgramData%\knita\executor.yaml`

```yaml
---
# Name identifies the Executor in build logs. It should be unique across all Executors to avoid confusion.
# Defaults to the system host name if not set.
name: knita-exec-arm64-linux

# Bind Address configures the interface and port the Executor will bind to. The Executor must be routable
# by the Broker and the Knita CLI.
# Defaults to 127.0.0.1:9091 if not set.
bind_address: 0.0.0.0:9091

# Labels will be advertised to the Broker and are used to match compatible builds to this Executor.
# The Executor will always advertise its OS and Architecture via built-in labels. These will vary 
# based on the host system the Executor is deployed to. OS will be one of 'linux', 'darwin' or 'windows',
# and Architecture will be on of 'amd64', 'arm' or 'arm64'.
labels:
  - nvidia-h100
```



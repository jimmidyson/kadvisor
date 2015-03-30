# kAdvisor: Configurable metrics collection & publishing for your Kubernetes cluster

**_This project is a proof of concept to help me learn more about [Golang](https://golang.org/)._**

**_For Kubernetes cluster monitoring that is usable today, see [Heapster](https://github.com/GoogleCloudPlatform/heapster)._**

**_All feedback on this project greatly received - if you feel it's useful then please let me know
& help me to improve it & make it production quality._**

kAdvisor collects metrics from various configured sources & sends them to configured sinks.
kAdvisor supports multiple sources & sinks.

## Usage

```
A configurable metrics collector & publisher for Kubernetes

Usage:
  kadvisor [command]

Available Commands:
  start       Start kAdvisor
  help        Help about any command

Flags:
  -c, --config="": config file
      --default-poll=10s: poll duration
  -h, --help=false: help for kadvisor
      --sink=[]: sinks
      --source=[]: sources
  -v, --verbose=false: verbose logging


Use "kadvisor help [command]" for more information about a command.
```

### Source & sink configuration

Sources & sinks are specified in the same way, using a URL of the format:

    <source_or_sink_prefix>://<source_or_sink_url>

For example, to use Kubernetes as a source you would specify:

    --source=kubernetes://https://192.168.0.1

You can specify multiple sources/sinks by just repeating the `--source` & `--sink` flags.
For example, if you want to collect from multiple Kubernetes masters:

    --source=kubernetes://https://192.168.0.1 --source=kubernetes://https://192.168.0.2

Every source & sink has its own specific configuration & these are passed as query parameters
to the source/sink URL. For example, to configure the API version for the Kubernetes source
you would specify the `apiVersion` query parameter:

    --source=kubernetes://https://192.168.0.1?apiVersion=v1beta2

See the sources & sinks documentation for their specific flags. Sources & sinks should
validate their input at creation time (fail-fast approach).

## Sources
### Kubernetes

Collecting metrics from Kubernetes is done via the Kubelet which exposes container
statistics collected from [cAdvisor](https://github.com/google/cadvisor). kAdvisor
retrieves nodes & pods from the Kubernetes master, collecting & collating metrics before
passing them to the configured sinks to process.

You can configure the Kubernetes source as detailed above. The prefix for a Kubernetes
source is `kubernetes://` & the remainder of the URL is passed to the source initializer.
The hostname in the source URL is the hostname/IP of the Kubernetes master. The Kubernetes source
also supports the following flags:

* `apiVersion` - the specified API version (default: `v1beta3`)
* `insecure` - whether to ignore certificate validation failures (default: `false`)
* `auth` - the Kubernetes client auth file as detailed at https://github.com/GoogleCloudPlatform/kubernetes/blob/master/pkg/clientauth/clientauth.go (default: None)

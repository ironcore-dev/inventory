# Benchmark-scheduler

Benchmark-scheduler is a tool responsible for benchmark scheduling and updating results in the Kubernetes cluster.
It requires `machines.benchmark.onmetal.de` CRD in cluster (more information may be found [here](./development.md)).

## Usage
```
bench-scheduler -h                                                                                                 ok  01:01:46 PM 
NAME:
   bench-scheduler - Start benchmarks in a scheduler way

USAGE:
   bench-scheduler [global options] command [command options] [arguments...]

COMMANDS:
   run, start, do
   help, h         Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h  show help (default: false)

```

```
NAME:
   bench-scheduler run - run benchmark jobs

USAGE:
   bench-scheduler run [command options] [arguments...]

DESCRIPTION:
   run benchmark jobs

OPTIONS:
   --provider value, -p value    Specify provider for cluster interaction. Example [bench-scheduler run -p kubernetes]
   --config value, -c value      Specify config file with benchmarks. Example [bench-scheduler run -c examples/config.yaml]
   --help, -h                    show help (default: false)
```
### Flags

Following configuration parameters are available for both binaries
- `-c, --config`

  Path to config file.

  Provides information about benchmarks which should be run on machine.

  Accepts `string`.

- `-p, --provider`

  Upload Provider.

  Benchmarks result should be updated in specific fields on Inventory Custom Resource in kubernetes cluster.
  Possible variants are - kubernetes / http.

  Accepts `string`.

  Default value is http.

### Environment variables

---------------------
| Key         | Definition |
| ----------- | ----------- |
| GATEWAY      | This variable is http-api host name or ip. Example: 'http://127.0.0.1:8080'|
| KUBECONFIG   | Kubernetes config file: Example: `~/.kube/config`        |
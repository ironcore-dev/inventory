# Development
### Requirements
Following tools are required to work on that project.

- [make](https://www.gnu.org/software/make/) - to execute build goals.
- [golang](https://golang.org/) - to compile source code.
- [cgroups](https://www.kernel.org/doc/Documentation/cgroup-v2.txt) - is a Linux kernel feature that limits, accounts for, and isolates the resource usage.
- [curl](https://curl.se/) - to download resources.
- [docker](https://www.docker.com/) - to build container with the tool.
- [mlc](https://software.intel.com/content/www/us/en/develop/articles/intelr-memory-latency-checker.html) - memory benchmark utility

### Prerequisites
[Сgroups](https://www.kernel.org/doc/Documentation/cgroup-v2.txt) is required in order to work with the benchmark-scheduler application.

### CRD
Before usage, please install a set of CRD's in the Kubernetes cluster:
- [inventories](https://github.com/onmetal/metal-api/blob/main/config/crd/bases/machine.onmetal.de_inventories.yaml)
- [benchmarks](https://github.com/onmetal/metal-api/blob/main/config/crd/bases/benchmark.onmetal.de_machines.yaml)

CRD explanation may be found in the metal-api [api-reference](https://github.com/onmetal/metal-api/tree/main/docs/api-reference).
# Build

To build all binaries execute:
```shell
make 
```
This command will produce a `./dist` directory with all required files. 

### Setting up Dev

Here's a brief intro about what a developer must do to start developing
the project further:

1. Check cgroups status.

```shell
mount -l | grep cgroup
cgroup2 on /sys/fs/cgroup type cgroup2 (rw,nosuid,nodev,noexec,relatime,nsdelegate,memory_recursiveprot)
```

2. Clone repo

```shell
git clone https://github.com/onmetal/inventory.git
cd inventory/
```

### Make goals

- `compile` (`all`) - build a project distribution.
- `fmt` - apply `go fmt` to the project.
- `vet` - apply `go vet` to the project.
- `dl-pciids` - downloads/updates PCI IDs database.
- `docker-build` - builds docker image.
- `docker-push` - pushes docker image to hte registry.
- `clean` - deletes built artifacts.

## Versioning

We use [SemVer](http://semver.org/) for versioning. For the versions available, see the [link to tags on this repository](https://github.com/onmetal/benchmark-scheduler/tags).


# Documentation Setup

The documentation of the [inventory](https://github.com/onmetal/inventory) project is written primarily using Markdown.
All documentation related content can be found in the `/docs` folder. New content also should be added there.
[MkDocs](https://www.mkdocs.org/) and [MkDocs Material](https://squidfunk.github.io/mkdocs-material/) are then used to render the contents of the `/docs` folder to have a more user-friendly experience when browsing the projects' documentation.

Currently, following tools are provided:
- `inventory` - collects data about system hardware;
- `nic-updater` - collects only NIC data (LLDP and NDP), in order to keep it up to date.
- `benchmark` - collects info from Intel® Memory Latency Checker.
- `benchmark-scheduler` - benchmark tasks scheduler.

### Requirements
Following tools are required to work on that package.

- [make](https://www.gnu.org/software/make/) - to execute build goals.
- [golang](https://golang.org/) - to compile source code.
- [cgroups](https://www.kernel.org/doc/Documentation/cgroup-v2.txt) - to create CPU and other limits.
- [curl](https://curl.se/) - to download resources.
- [docker](https://www.docker.com/) - to build container with the tool.
- [mlc](https://software.intel.com/content/www/us/en/develop/articles/intelr-memory-latency-checker.html) - memory benchmark utility

### Prerequisites
To work with benchmark-scheduler application [cgroups](https://www.kernel.org/doc/Documentation/cgroup-v2.txt) are required.

# Build

To build all binaries just execute:
```shell
make 
```
This will produce a `./dist` directory with all required files. 

### Setting up Dev

Here's a brief intro about what a developer must do in order to start developing
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


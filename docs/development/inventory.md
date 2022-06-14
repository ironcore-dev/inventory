# Inventory

Inventory tool is responsible for hardware data collection and configuration. As soon as data collected, the application will update it in the Kubernetes cluster CRD. 
CRD `inventories.machine.onmetal.de` is required (more information may be found [here](./development.md)).

### Run

To execute, simply run
```shell
    sudo ./dist/inventory
```

Because of:
- [DMI](https://www.kernel.org/doc/html/v4.15/driver-api/firmware/other_interfaces.html) - Generates a standard framework   for managing and tracking components on a server, by abstracting these components from the software that manages them. 
- [ioctl](https://man7.org/linux/man-pages/man2/ioctl.2.html) - System call manipulates the underlying device parameters of special files.  In particular, many operating characteristics of character special files (e.g., terminals) maybe controlled with ioctl() requests.

Program usage requires additional privileges.

### Clean

To remove built distribution directory run:
```shell
    make clean
```
## Usage

### Flags

Following configuration parameters are available for all binaries:

- `-k, --kubeconfig`
  
    Path to kubeconfig file.
    
    Used to establish connection with k8s cluster API. If multiple contexts are available, 
  currently selected context will be used. 
  
    Accepts `string`.
  
    Default value is `/home/${username}/.kube/config` if home directory is present, empty string otherwise.

- `-g, --gateway`

    Gateway host. 

    Resource will be pushed to cluster through provided gateway host address.

    If provided, has a priority over kubeconfig.

    Accepts `string`.

    Default value is empty string.

- `-t, --timeout`

    Request timeout.

    Used if resource is pushed through gateway.

    Accepts `string`.

    Default value is empty string.

- `-n, --namespace`
  
    k8s namespace.
    
    Resource will be pushed to selected namespace.
    
    Accepts `string`.
    
    Default value is `default`.
  
- `-r, --root string`
  
    Path to root file system.
    
    Used to build paths and chroot to target filesystem, if current is not the one.
    
    Accepts `string`.
    
    Default value is `/`.

- `-p, --patch bool`
  
    CRD method update.
    
    Used Patch method instead of Post for CRD update.
    
    Accepts `bool`.
    
    Default value is `false`.
  
- `-v, --verbose`
  
    Verbose output. 
  
    Enables verbose output. May be used to troubleshoot the process if data is not collected for some reason.

### Data

There are currently 2 out of 4 planned tools:
- inventory
- planned: network analysis tool using LLDP
- benchmark
- planned: stresstest possibly using benchmark for transient error analysis

inventory is collecting data about:
- block device from sysfs `/sys/block`.
- block device partition tables and partitions from `/dev`.
- CPU from `/proc/cpuinfo`.
- memory from `/proc/meminfo`.
- NUMA from sysfs `/sys/devices/system/node`.
- system from SMBIOS via DMI.
- IPMI from `/dev`.
- NIC from sysfs `/sys/class/net`.
- LLDP from `/run/systemd/netif/lldp`.
- NDP from kernel routing tables via ioctl.
- PCI devices from sysfs `/sys/devices` using PCI IDs database.
- Virtualization devices from `/sys`, `/proc` and DMI.

benchmark is collecting data about:
- local and remote memory latency
- local and remote memory bandwidth
Planned:
- CPU test: integer, floating point, vector instructions
- Disk: read write performance, sequential and random
- Network latency and throughput

### Known issues

- Inventory tools will not work inside the WSL2 machine due to bug [#6874](https://github.com/microsoft/WSL/issues/6874)
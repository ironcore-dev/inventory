# inventory

inventory is a CLi tool that collects data on host about its hardware and configuration and, 
then pushes it to the k8s cluster in a form of custom resource. 

## Getting started

### Required tools

- [make](https://www.gnu.org/software/make/) - to execute build goals.
- [golang](https://golang.org/) - to compile source code.
- [curl](https://curl.se/) - to download resources.
- [docker](https://www.docker.com/) - to build container with the tool.

### Compile

To compile, execute

    make 

This will produce a `./dist` directory with all required files. 

### Run

To execute, simply run

    sudo ./dist/inventory

| This will not work inside the WSL2 machine due to bug [#6874](https://github.com/microsoft/WSL/issues/6874) |
| --- |

Program uses machine's DMI interface and ioctl, so it requires additional privileges.

### Clean

To remove built distribution directory run

    make clean

## Usage

### Flags

Following configuration parameters are available 

- `-k, --kubeconfig`
  
    Path to kubeconfig file.
    
    Used to establish connection with k8s cluster API. If multiple contexts are available, 
  currently selected context will be used. 
  
    Accepts `string`.
  
    Default value is `/home/${username}/.kube/config` if home directory is  present, empty string otherwise.
  
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
  
- `-v, --verbose`
  
    Verbose output. 
  
    Enables verbose output. May be used to troubleshoot the process if data is not collected for some reason.

### Data

Tool is collecting data about.
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

## Development

### Make goals

- `compile` (`all`) - build a project distribution.
- `fmt` - apply `go fmt` to the project.
- `vet` - apply `go vet` to the project.
- `dl-pciids` - downloads/updates PCI IDs database.
- `docker-build` - builds docker image.
- `docker-push` - pushes docker image to hte registry.
- `clean` - deletes built artifacts.

### Libraries

Following libraries are used to simplify the process of data collection and processing

- [`github.com/digitalocean/go-smbios`](https://github.com/digitalocean/go-smbios) to obtain raw SMBIOS data.
- [`github.com/diskfs/go-diskfs`](https://github.com/diskfs/go-diskfs) to get data about partition table.
- [`github.com/lunixbochs/struc`](https://github.com/lunixbochs/struc) to deserialize aligned binary structures.
- [`github.com/mdlayher/lldp`](https://github.com/mdlayher/lldp) to deserialize binary LLDP frames.
- [`github.com/u-root/u-root`](https://github.com/u-root/u-root) to get IPMI data.
- [`github.com/vishvananda/netlink`](https://github.com/vishvananda/netlink) to get kernel routing tables data.
- [`github.com/jeek120/cpuid`](https://github.com/jeek120/cpuid) to get CPUID.

### Resources

List of useful resources worth to check out if it is required to extend or alter the functionality

- [Accessing SMBIOS information with Go](https://mdlayher.com/blog/accessing-smbios-information-with-go/) - 
  in-depth explanation on SMBIOS data collection from DMI.
- [System Management BIOS](https://www.dmtf.org/standards/smbios) - page with SMBIOS specs.   
- [LLDP](https://wiki.wireshark.org/LinkLayerDiscoveryProtocol) - Wireshark wiki page that contains some info on LLDP packages.
  Also includes some sample files.
- [The PCI ID Repository](https://pci-ids.ucw.cz/) - PCI ID database. Submit here a record if hardware is not properly detected.
- [GPT](https://github.com/rekby/gpt) and [MBR](https://github.com/rekby/mbr) libraries - low-profile libs with no dependencies to read GPT and MBR partitions. 
  May be used if there will be a need to shrink binary size.

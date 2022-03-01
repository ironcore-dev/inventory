# Resources

List of useful resources worth checking out if it is required to extend or alter the functionality:

- [Accessing SMBIOS information with Go](https://mdlayher.com/blog/accessing-smbios-information-with-go/) - 
  an in-depth explanation of SMBIOS data collection from DMI.
- [System Management BIOS](https://www.dmtf.org/standards/smbios) - page with SMBIOS specs.   
- [LLDP](https://wiki.wireshark.org/LinkLayerDiscoveryProtocol) - Wireshark wiki page that contains some info on LLDP packages.
  Also includes some sample files.
- [The PCI ID Repository](https://pci-ids.ucw.cz/) - PCI ID database. Submit here a record if the hardware is not properly detected.
- [GPT](https://github.com/rekby/gpt) and [MBR](https://github.com/rekby/mbr) libraries - low-profile libs with no dependencies to read GPT and MBR partitions. 
  May be used if there will be a need to shrink the binary size.
- [`Cgroups`](github.com/containerd/cgroups) - to work with linux kernel features.
- [`Struc`](https://github.com/lunixbochs/struc) - to deserialize aligned binary structures.

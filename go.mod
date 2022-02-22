module github.com/onmetal/inventory

go 1.15

require (
	github.com/Jeffail/gabs/v2 v2.6.1
	github.com/containerd/cgroups v1.0.3
	github.com/digitalocean/go-smbios v0.0.0-20180907143718-390a4f403a8e
	github.com/diskfs/go-diskfs v1.1.1
	github.com/go-redis/redis/v8 v8.8.2
	github.com/google/uuid v1.3.0
	github.com/jeek120/cpuid v0.0.0-20200914054105-8fa8c861dea6
	github.com/lunixbochs/struc v0.0.0-20200707160740-784aaebc1d40
	github.com/mdlayher/ethernet v0.0.0-20190606142754-0394541c37b7 // indirect
	github.com/mdlayher/lldp v0.0.0-20150915211757-afd9f83164c5
	github.com/onmetal/k8s-inventory v0.0.2-0.20211117172137-e7a07f43bd63
	github.com/onmetal/metal-api v0.2.6
	github.com/onmetal/metal-api-gateway v0.2.2
	github.com/opencontainers/runtime-spec v1.0.2
	github.com/pkg/errors v0.9.1
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.7.0
	github.com/u-root/u-root v7.0.0+incompatible
	github.com/urfave/cli/v2 v2.3.0
	github.com/vishvananda/netlink v1.1.0
	github.com/vtolstov/go-ioctl v0.0.0-20151206205506-6be9cced4810 // indirect
	go.uber.org/zap v1.21.0
	gopkg.in/yaml.v2 v2.4.0
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
	k8s.io/apimachinery v0.23.3
	k8s.io/client-go v0.23.3
)

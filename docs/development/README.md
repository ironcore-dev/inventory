# Installation, using and developing

The documentation of the [inventory](https://github.com/onmetal/inventory) project is written primarily using Markdown.
All documentation related content can be found in the `/docs` folder. New content also should be added there.
[MkDocs](https://www.mkdocs.org/) and [MkDocs Material](https://squidfunk.github.io/mkdocs-material/) are then used to render the contents of the `/docs` folder to have a more user-friendly experience when browsing the projects' documentation.

Currently, following tools are provided:
- `inventory` - collects data about system hardware;
- `nic-updater` - collects only NIC data (LLDP and NDP), in order to keep it up to date.
- `benchmark` - collects info from Intel® Memory Latency Checker.
- `benchmark-scheduler` - benchmark tasks scheduler.

To find out more, please refer to: 
- [development](./development.md)
- [contribution](./contribution.md)
- [libraries](./libraries.md)
- [resources](./resources.md)
- [testing](./testing.md)
- [inventory usage](./inventory.md)
- [benchmark-scheduler usage](./benchmark_scheduler.md)
- [benchmark-scheduler config usage](./benchmark_scheduler_config.md)


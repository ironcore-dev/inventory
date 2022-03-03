## Tests

In order to run tests, choose one of the following variant:
```shell
go test -race ./...
```
or `make` command:
```shell
make test
```
`make` command will provide a cover profile as a result.
For the current moment, tests exist only for benchmark-scheduler.

### Permissions

Tests require root access (`sudo`). 
For example:
```shell
worker_test.go:26: can't get user info or user is not a root:
```
The Root is required to create CGroup (more information about CGroup, may be found [here](./development.md)).

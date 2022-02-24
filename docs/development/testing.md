## Tests

```shell
go test -race ./...
```

Some test requires sudo. For example
```shell
worker_test.go:26: can't get user info or user is not a root:
```
Root is required to create CGroup.

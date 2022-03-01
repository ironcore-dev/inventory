A full configuration example can be found [here](../../config/samples/config.yaml).

```
benchmarks:
  - name: disk-4k
    app: fio
    jsonpathInputSelector: "spec.blocks.*.name"
    output: "jobs.*.read.io_bytes"
    resources:
      cpu: "8"
    args:
      - '--rw=randread'
      - '--filename={{ inputSelector }}'
```

1. benchmarks - List of benchmarks that should be executed on the machine.
2. name - Test name.
3. app - An application that should be used for the benchmark.
4. jsonpathInputSelector - JSON path to the machine element which one should be tested with this benchmark. Special filed `{{ inputSelector }}` would be replaced with the selector.
5. output - Specify which data from the output need to be put into the cluster. Possible variants are:
  5.1 jsonpath - if the output result is JSON, it's possible to find out the required element with jsonpath. (default)
  5.2 text: works as simple grep tool. For instance, if output like this `someKey 123, anotherKey, 232`, it's will separate string and find out your key and value.
6. resources:
```
  cpu_set - CPUSet defines specific cores for application.
  cpu -  Hard cap limit (in usecs). Allowed CPU time in a given period.
  shares - Defines CPU core share which can be used by the application.
  cores - Defines the number of CPU cores that can be used by the application.
  period - Indicates that the group may consume available CPU in each period duration.
```
7. args - List of arguments for application.
Full config example can be found [here](../../config/samples/config.yaml).

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

1. benchmarks - list of benchmarks which should be executed on machine.
2. name - test name.
3. app - application which should be used for benchmark.
4. jsonpathInputSelector - jsonpath to the machine element which one should be tested with this benchmark. Special filed `{{ inputSelector }}` would be replaced with the selector.
5. output - specify which data from the output need to be putted into the cluster. Possible variants are:
  5.1 jsonpath - if output result is json, it's possible to find out required element with jsonpath. (default)
  5.2 text: works as simple grep tool. For instance if output like this `someKey 123, anotherKey, 232`, it's will separate string and find out your key and value.
6. resources:
```
	cpu_set - CPUSet defines specific cores for application.
	cpu -  hardcap limit (in usecs). Allowed cpu time in a given period.
	shares - defines CPU core share which can be used by application.
	cores - defines number of CPU cores which can be used by application.
	period - indicates that the group may consume available cpu in each period duration.
```
7. args - list of arguments for application.
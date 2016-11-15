# Consul Demonstration Code

This is some extremely sloppy example code for a talk given on 2016-11-15 with DevopsBoise.


# Seriously, don't use this for anything beyond basic concurrency examples.

There are six subcomponents here.

* consul - a trivial `start.sh` script to use the consul docker image and bind the API port.
* demo-frame - a simple web server/page to watch the results of the other modules.
* feature-flags - an example implementation of feature flagging with Consul.  Try "flags/magic = true".
* kv-service-registry - The simplest example, which the others are based on.  Sets keys to use the KV store for service registry.
* leader-election - an example of leader election using locks in Consul.  Kill some containers.
* locking - an example of multiple processes sharing a single worklock.  Includes a bash implementation as well.

## Build

If you really want to build this stuff, your basic instructions are:

```
GOOS=linux go build
docker build -t <component> .
```

For the most part, that should give you runnable containers.  Enjoy!



## Credits

The basic webserver/html tiny app pattern (and implementation) was borrowed from Kelsey Hightower's spectacular demos, in particular [app-healthz](https://github.com/kelseyhightower/app-healthz).

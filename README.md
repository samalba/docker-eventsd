docker-eventsd
==============

Events manager for a docker cluster

What is it?
-----------

docker-eventsd listens to any events on a cluster of docker hosts. When it receives an event
(container start, stop, die, etc...), it will trigger an action.

It allows to change the state of your cluster dynamically at runtime.

Possible use cases
------------------

- Mysql master-slave (promote the slave as master when the master goes down)
- Dynamic load-balancing (register new backends to the Load-Balancer when new containers start)
- Logging / Reporting the state of a docker cluster dynamically

How to use it?
--------------

Create a file `events.yml` or look at [this one](https://github.com/samalba/docker-eventsd/blob/master/events.yml).

The purpose of this file is to:

- describe your cluster of docker host (give each one a name and a daemon URL)
- declare different event handler with the following properties:
  - what kind of event
  - filter on an engine name (where the event comes from)
  - filter on a container name
  - filter on an image name
  - set the command handler (the command to trigger when all the above matches)

The command handler runs in a real shell with a custom environment to give info on where the event comes from:
- `FROM_ENGINE=` (contains the engine name)
- `FROM_CONTAINER=` (contains the container name)
- `ENGINE_FOO=` (contains the daemon URL of the engine foo)

The last thing is to launch the `docker-eventsd` process with the events.yml file as an argument:

```
$ go build && ./docker-eventsd events.yml
2014/10/12 11:24:04 Listening to events...
2014/10/12 11:24:04 Get http://192.168.59.105:2375/v1.10/events: dial tcp 192.168.59.105:2375: host is down
2014/10/12 11:24:06 Uncaught event: type create from foo@tcp://192.168.59.103:2375 on /grave_euclid
2014/10/12 11:24:06 Uncaught event: type start from foo@tcp://192.168.59.103:2375 on /grave_euclid
```

TODO
----

- Discover the cluster info: etcd, redis...
- Write custom events type (tcp timeout, http error, etc...)
- Custom events type has to come with the ability to add customer monitor (plugins that makes a regular check and trigger the custom event)

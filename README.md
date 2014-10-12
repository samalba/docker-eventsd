docker-eventsd
==============

Events manager for a docker cluster

What is it?
-----------

docker-eventsd listens to any events on a cluster of docker hosts. When it receives an event
(container start, stop, die, etc...), it will trigger an action.

It allows to change the state of your cluster dynamically at runtime.

Possible use case
-----------------

- Mysql master-slave (promote the slave as master when the master goes down)
- Dynamic load-balancing (register new backends to the Load-Balancer when new containers start)
- Logging / Reporting the state of a docker cluster dynamically

How to use it?
--------------

Create a file `events.yml` or look at [this one](https://github.com/samalba/docker-eventsd/blob/master/events.yml)


TODO
----

Discover the cluster info: etcd, redis...

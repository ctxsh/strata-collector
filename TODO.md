# TODO


## Current state
* Most of the collector/discovery sync issues have been taken care of.  The discovery service checks for the collector in the cache and grabs the channel directly from the collector pool.  We could still end up having a closed channel, so we need to make sure and check.  But with this solution, if the discovery service is sending data to the collector pool, it's locked and will need to wait on any shutdown so races should be uncommon.
* With this relationship, the collectors are mostly autonomous and do not need to know about the state of any discovery services that use them.  We can rethink the central registry and just pass a list of collectors to the discovery service. i.e. collectors no longer need any references to discovery. This could be a significant refactor which will recombine all of the services into a single package - could simplify things.
* We are ready to implement the filters and any configuration around that.
* We are ready to implement the first of the sinks.
* We still need to add the webhook, validation, and defaults for the collector.

## Encoders
* [] JSON metric encoder
* [] Fluentbit
* [] Statsd

## Filters
* [] Clip
* [] Clamp

## Sinks

In general we'll try to stay away from specific vendor implementations, though I'm not completely sure that we'll be able to do that.  Initially the pattern will be to utilize the collectors to encode the metrics and push them into a message queue.  Mostly this could be overkill for shops with limited numbers of metrics that are being ingested, but it gives some additional reliability on upstream failures and scale.  I've been toying with the idea about supporting a simple kafka or nsq install as part of the controller (most likely nsq).  This would allow a very low entry barrier for people that don't want to manage queues.  I think the queues should be the default and recommended workflow - with the vendor agents reading from the queues, with other higher latency push models such as HTTP/s and statsd sinks and potential vendor specific sinks.

### Short term
* [] NSQ sink package and potentially manage default NSQ service.
* [] Kafka sink package.
* [] Document how to use agents to connect to the message queues.

### Long term
* [] HTTP/s sink
* [] Statsd sink
* [] Fluent sink

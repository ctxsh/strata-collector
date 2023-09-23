# TODO

## Current state
* Most of the collector/discovery sync issues have been taken care of.  The discovery service checks for the collector in the cache and grabs the channel directly from the collector pool.  We could still end up having a closed channel, so we need to make sure and check.  But with this solution, if the discovery service is sending data to the collector pool, it's locked and will need to wait on any shutdown so races should be uncommon.
  * On second thought, even with the locks, could we still end up in a situation where there are metrics left in the channel after shutdown?  I don't think so since the discovery services will still be trying to lock the channel, and if our stop mechanism will wait until the channel is empty, we should be good.  However, there may be a race between removing the collector from the registry and a new collector coming online, so delete the registry from the map before stopping.
  * This method of managing the channels could definitely cause metric loss. In the situation where the discovery service is waiting for a new collector to come up and be registered, new metrics coming in would be discarded.  There's a couple ways to handle this: 1) allow the discovery service to buffer up to N endpoints - this is not ideal as we expect that the collectors should be collecting the metrics in a short amount of time from the time the resource is added to the channel so we maintain (even though we don't guaruntee) a collection interval.  2) Set up a retry for adding to the channel instead of dropping (with the retry timing out before the start of the next collection interval).  Works better since the stop, start, and registration of a modified collector should be pretty fast.  We can measure the delta between queuing and collection to ensure things are within a certain window (retry will need to be around the send itself).  3) Registry takes care of channels between the discovery services and collection pools.  This seems like it would be the best solution, however this may be a bit more complex.  Might start with 2 and then the goal will be to implement 3 in time.
  * Need to collect the time the resource is added to the collection channel and drop if outside a certain window.
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
* [] NATS sink package and potentially manage default NATS service.  Collectors will contain a option to spin up an intermediate queue and by default send all events into the queue.
* [] Kafka sink package.
* [] Document how to use agents to connect to the message queues.

### Long term
* [] HTTP/s sink
* [] Statsd sink
* [] Fluent sink


Tommorrow:
* [] Flags.
* [] Implement filters.
* [] Implement encoders.
* [] Buffer length for collection channels.
* [] Validation either inline or in the validation hooks.
* [] Start looking at testing using kubebuilder and controller-runtime as an example.
* [] Status updates outside of the normal interval runs.
* [] License headers.
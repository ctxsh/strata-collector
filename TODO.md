# TODO

## Current state
* We are ready to implement the filters and any configuration around that.
* We are ready to implement the first of the sinks.
* We still need to add the validation for discovery and collectors.

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
* [x] Need to collect the time the resource is added to the collection channel and drop if outside a certain window.
* [] Flags.
* [] Implement filters.
* [] Buffer length for collection channels.
* [] Validation either inline or in the validation hooks.
* [] Start looking at testing using kubebuilder and controller-runtime as an example.
* [] Status updates outside of the normal interval runs.
* [] License headers.

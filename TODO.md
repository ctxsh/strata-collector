# TODO

## Current state
* We are ready to implement the filters and any configuration around that.
* We are ready to implement the first of the sinks.
* We still need to add the validation for discovery and collectors.

## Encoders
* [x] JSON metric encoder
* [] Fluentbit
* [] Statsd

## Filters
* [x] Clip
* [x] Clamp

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

## Fixes
* [] Because of the way service and endpoints are intertwined, if endpoint resource discovery has been configured but service has not, endpoints will never be discovered.  Need to split them up.

## Brainstorm
* [] Datadog won't actually use the histogram values.  It expects you to pass in the value and the agent generates the statistical distribution suffixes in the agent as part of a "stateful" metric.  I need to look at the agent a bit more, but this seems like it could introduce some inaccuracies on metrics that may come in with an incomplete window (startup or shutdown). It also means that we may want to include our own datadog output. It's not necessarily a horrible thing as long as it's known and expected
* [] Can we track cardinality and add it to the metric as a way to alert users about tags that get to high?  Depending on the window, there could be some resource pressures, but if we are looking at a small enough window it should be pretty quick.  This could also lead into statistical information about the metrics being collected as well. This could be used as a way to filter/throttle metrics that we are scraping - either automatically or by explicit configuration.
* [] In the future, we can add more robust upstream services for analytics that the collector could utilize.  The question is, do we want to design things around this now - i.e. all the cardinality/stats (including things like avg etc.) coming from a scaled down version?  Might be able to have it both ways as long as the interfaces are modular enough to run both inline and as an api.  These systems would be used as a side channel which could statefully store additional information about a series.
* [] I like the idea having some sort of outlier detection available.  This could lead to smart sampling or filtering.  The user would be able to specify the upper and lower limits or we could try to determing automatically by watching the distribution.  The stats would need to be somewhat stateful since we'd need to reload on startup.
* [] Represent a batches of like metrics as a series for custom filtering/sampling/cardinality/stats/etc.  This probably means that we need to create a fanout with stats collectors in the middle which then push them back into the outputs (outputs become seperate workers).
* [] I think we can start treating logs and k8s events as first class and collecting them as well.  Structuring if needed.  This isn't going to be a huge ask, but it may require us to create a small agent to scrape the logs on each node.  The collector would then register itself to a log stream.

Tommorrow:
* [x] Need to collect the time the resource is added to the collection channel and drop if outside a certain window.
* [] Flags.
* [x] Implement filters.
* [x] Buffer length for collection channels.
* [] Validation either inline or in the validation hooks.
* [] Start looking at testing using kubebuilder and controller-runtime as an example.
* [] Status updates outside of the normal interval runs.
* [] License headers.

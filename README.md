# Strata Metrics Collector

The Strata Collector is a kubernetes native metrics collection tool.  It is meant to
work as part of a suite of tools meant to simplify the generation, collection, and transfer
of metrics.  Current client tooling includes:

* [https://github.com/ctxsh/strata-go](strata-go): Prometheus client wrapper for golang.

## Development

### Goals

In general the initial version of will be quite limited in what it can to.

* Manage configuration of the scrapers through a custom resource definition.
* Discover pods and services that need to be scraped through the same annotations
that are present with the prometheus k8s integration.
* Provide basic filtering (cut/cap/min/max) before the metrics get to the upstream.
* Initially we'll support the following sinks: kafka, http, statsd.
* Less resource utilization than some of the alternative options out there.
* Leader election to support redundancy

### Non-Goals

* Won't try to support any of the vendors that are out there right now.  We'll make
the assumption that there will be upstream consumers that will be able to pull/receive
from one of the sinks and send it on.  This could change in the future as the project
grows or if needed.

## Alternative options

* TODO

## Notes

Each one should have it's own discovery mechanism.  This will be the only thing that needs to
be updated.  Push workers will be their own config setup.  Should the pipeline have it's own as
well or should the channels and the buffering be part of the worker pool.  Multiple discovery
services can go into the same collection pool so we can size for upstream.

* Allows tuning and protection of upstream queues/apis/etc
* Allows us to more effectively use resources locally

What do we have now?

```
---
apiVersion: strata.ctx.sh/v1beta1
kind: Discovery
metadata:
  name: example
  namespace: default
  labels:
    service: example
spec:
  selector:
    matchLabels:
      sink: redpanda
  enabled: false
  intervalSeconds: 10
---
apiVersion: strata.ctx.sh/v1beta1
kind: Collector
metadata:
  name: example
  namespace: default
  labels:
    sink: redpanda
spec:
  enabled: false
  workers: 5
  buffer:
    size: 100000
    # different intermediates will be available.  starting out
    # with just channels, but we should add nsq and kafka (or variants)
    # that will allow for a more robust solution that can withstand
    # the upstreams going out.
    type: channel

```

If the collector has been updated, the discovery services will need to be stopped, 

Collector
Discover

Where does the buffer size end up?  Main config?

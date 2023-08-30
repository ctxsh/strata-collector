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


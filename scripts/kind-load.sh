#!/usr/bin/env bash

docker pull golang
kind load docker-image golang --name strata
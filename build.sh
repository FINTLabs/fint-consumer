#!/usr/bin/env bash

docker build -t fint-consumer --build-arg VERSION=0.$(date +%y%m%d.%H%M) .

#!/usr/bin/env sh


rm -rf build

gox -output="./build/{{.Dir}}-{{.OS}}" -rebuild -osarch="darwin/amd64 windows/amd64"
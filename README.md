# fint-consumer

[![Build Status](https://jenkins.fintlabs.no/buildStatus/icon?job=FINTLabs/fint-consumer/master)](https://jenkins.fintlabs.no/job/FINTLabs/fint-consumer/master)
[![MIT license](http://img.shields.io/badge/license-MIT-brightgreen.svg)](http://opensource.org/licenses/MIT)

## Description
Generates `Java` consumer code from EA XMI export. This utility is mainly for internal FINT use, but if you find it useful, please use it!

## Usage

```
$ fint-consumer
NAME:
   fint-consumer - Generates consumer code from EA XMI export. This utility is mainly for internal FINT use, but if you find it usefull, please use it!

USAGE:
   fint-consumer [global options] command [command options] [arguments...]

VERSION:
   2.0.0

AUTHOR:
   FINTLabs

COMMANDS:
     generate      generates consumer code
     listPackages  list Java packages
     listTags      list tags
     listBranches  list branches
     setup         setup a consumer project
     help, h       Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --owner value          Git repository containing model (default: "FINTLabs") [$GITHUB_OWNER]
   --repo value           Git repository containing model (default: "fint-informasjonsmodell") [$GITHUB_PROJECT]
   --filename value       File name containing information model (default: "FINT-informasjonsmodell.xml") [$MODEL_FILENAME]
   --tag value, -t value  the tag (version) of the model to generate (default: "latest")
   --force, -f            force downloading XMI for GitHub.
   --help, -h             show help
   --version, -v          print the version
```

The downloaded XMI file is put in the `$HOME/.fint-consumer/.cache`. If you don't use the 
`force` flag and the file exists in the cache directory it uses this one. 

## Install

### Binaries

### Binaries

Precompiled binaries are available as [Docker images](https://cloud.docker.com/u/fint/repository/docker/fint/fint-consumer)

Mount the directory where you want the generated source code to be written as `/src`.

Linux / MacOS:
```bash
docker run -v $(pwd):/src fint/fint-consumer:2.1.0 <ARGS>
```

Windows PowerShell:
```ps1
docker run -v ${pwd}:/src fint/fint-consumer:2.1.0 <ARGS>
```

### Go

To install, use `go get`:

```bash
go get -d github.com/FINTLabs/fint-consumer
go install github.com/FINTLabs/fint-consumer
```

#### Release 

```bash
docker build --build-arg VERSION=2.5.2 -t fint/fint-consumer:v2.5.2 .
```

Verfify releases:
```bash
docker images | grep fint-consumer
```


## Author

[FINTLabs](https://fintlabs.github.io)

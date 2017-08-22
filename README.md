# fint-consumer



## Description
Generates `Java` consumer code from EA XMI export. This utility is mainly for internal FINT use, but if you find it usefull, please use it!

## Usage

```
$ fint-consumer
NAME:
   fint-consumer - Generates consumer code from EA XMI export. This utility is mainly for internal FINT use, but if you find it usefull, please use it!

USAGE:
   fint-consumer [global options] command [command options] [arguments...]

VERSION:
   1.0.0

AUTHOR:
   FINTProsjektet

COMMANDS:
     generate      generates consumer code
     listPackages  list Java packages
     listTags      list tags
     listBranches  list branches
     help, h       Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --tag value, -t value  the tag (version) of the model to generate (default: "latest")
   --force, -f            force downloading XMI for GitHub.
   --help, -h             show help
   --version, -v          print the version
```

The downloaded XMI file is put in the `$HOME/.fint-consumer/.cache`. If you don't use the 
`force` flag and the file exists in the cache directory it uses this one. 

## Install

### Binaries

Precompiled binaries can be downloaded [here](https://github.com/FINTprosjektet/fint-consumer/releases/latest)

* Download for your os
* Rename to fint-consumer
* Copy to where ever you want it.

*Example macOS*
```shell
$ mv fint-consumer-darwin fint-consumer
$ chmod +x fint-consumer
$ sudo mv fint-consumer /usr/local/bin
```

### Go

To install, use `go get`:

```bash
go get -d github.com/FINTprosjektet/fint-consumer
go install github.com/FINTprosjektet/fint-consumer
```

## Author

[FINTProsjektet](https://fintprosjektet.github.io)

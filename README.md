# fint-consumer



## Description
Generates `Java` and `C#` models from EA XMI export. This utility is mainly for internal FINT use, but if you 
find it usefull, please use it!

## Usage

```
$ fint-consumer

```

The downloaded XMI file is put in the `$HOME/.fint-consumer/.cache`. If you don't use the 
`force` flag and the file exists in the cache directory it uses this one. 

## Install

### Binaries

Precompiled binaries can be downloaded [here](https://github.com/FINTprosjektet/fint-consumer/releases/latest)

### Go

To install, use `go get`:

```bash
go get -d github.com/FINTprosjektet/fint-consumer
go install github.com/FINTprosjektet/fint-consumer
```

## Author

[FINTProsjektet](https://fintprosjektet.github.io)

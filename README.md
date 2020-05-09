# go-trace-json
A command line tool to dump go traces into a json file


Why? We wanted to get details from a trace file without having to startup a webserver.


Current supported features:
- Outputting a list of goroutine groups


## Setup

1. Build with `go build .`

## Usage

Run with 

```
go-trace-json --input-file <inputfile>
```

Optional flags:
`--output-file <outputfile>` (Default : trace.json)

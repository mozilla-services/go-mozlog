# mozlog [![GoDoc](https://godoc.org/go.mozilla.org/mozlog?status.svg)](https://godoc.org/go.mozilla.org/mozlog) [![Build Status](https://travis-ci.org/mozilla-services/go-mozlog.svg?branch=master)](https://travis-ci.org/mozilla-services/go-mozlog)
A logging library which conforms to Mozilla's logging standard.

## Example Usage
```
import "go.mozilla.org/mozlog"

func init() {
    mozlog.Logger.LoggerName = "ApplicationName"
}
```

## Testing logs

The utility at `go.mozilla.org/mozlog/cmd/mozlog-tester` can take an input
mozlog string and evaluate it.
```
$ go get go.mozilla.org/mozlog/cmd/mozlog-tester

$ echo
'{"Time":"2017-02-17T17:18:36Z","Type":"app.log","Hostname":"gator3","EnvVersion":"2.0","Pid":24318,"Severity":4,"Fields":{"animal":"walrus","msg":"A group of walrus emerges from the ocean","size":10}}' | mozlog-tester

error: nanosecond Timestamp is missing
error: Logger is missing
Timestamp=0 Time="2017-02-17T17:18:36Z" Type="app.log" Logger="" Hostname="gator3" EnvVersion="2.0" Pid=24318 Severity=4 Fields[animal="walrus", msg="A group of walrus emerges from the ocean", size=10]
```

/*
Package ln is the Natural Logger for Go

`ln` provides a simple interface to logging, and metrics, and
obviates the need to utilize purpose built metrics packages, like
`go-metrics` for simple use cases.

The design of `ln` centers around the idea of key-value pairs, which
can be interpreted on the fly, but "Filters" to do things such as
aggregated metrics, and report said metrics to, say Librato, or
statsd.

"Filters" are like WSGI, or Rack Middleware. They are run "top down"
and can abort an emitted log's output at any time, or continue to let
it through the chain. However, the interface is slightly different
than that. Rather than encapsulating the chain with partial function
application, we utilize a simpler method, namely, each plugin defines
an `Apply` function, which takes as an argument the log event, and
performs the work of the plugin, only if the Plugin "Applies" to this
log event.

If `Apply` returns `false`, the iteration through the rest of the
filters is aborted, and the log is dropped from further processing.
*/
package ln

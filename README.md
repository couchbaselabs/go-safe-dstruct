[![Build Status](https://drone.io/github.com/tleyden/go-safe-dstruct/status.png)](https://drone.io/github.com/tleyden/go-safe-dstruct/latest)

A collection of data structures that are safe to use from multiple goroutines.

## Mapserver

[![GoDoc](https://godoc.org/github.com/tleyden/go-safe-dstruct/mapserver?status.png)](http://godoc.org/github.com/tleyden/go-safe-dstruct/mapserver)

This is a map that is designed in the same manner as one would write an Erlang *gen_server*.  There is a goroutine which wraps the non-goroutine-safe map and serializes all access to it.  It is not particularly performant, so don't use it if you need performance.  

## Queue 

[![GoDoc](https://godoc.org/github.com/tleyden/go-safe-dstruct/queue?status.png)](http://godoc.org/github.com/tleyden/go-safe-dstruct/queue)

A thread safe queue that uses mutex locks internally.

## License

Apache 2


A collection of data structures that are safe to use from multiple goroutines.

## Mapserver

This is a map that is designed in the same manner as one would write an Erlang *gen_server*.  There is a goroutine which wraps the non-goroutine-safe map and serializes all access to it.

It is used in [isync](https://github.com/tleyden/isync), which tries to avoid the use of any explicit locks.


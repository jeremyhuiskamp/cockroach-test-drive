This is the go client that simulates database traffic.

It is compiled while building the docker container, but can also
be compiled and run locally.

The code is in `src/client` in order to be able to satisfy `dep`'s
requirements for `$GOPATH`.

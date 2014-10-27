mptcp [![Build Status](https://travis-ci.org/mdlayher/mptcp.svg?branch=master)](https://travis-ci.org/mdlayher/mptcp) [![GoDoc](http://godoc.org/github.com/mdlayher/mptcp?status.svg)](http://godoc.org/github.com/mdlayher/mptcp)
=====

Package mptcp provides detection functionality for active, multipath TCP
connections from a remote client to the current host.  MIT Licensed.

This package is inspired by the original, PHP-based multipath TCP detection
functions, courtesy of Christoph Paasch and [multipath-tcp.org](http://multipath-tcp.org/).

An example binary called `mptcphttp` is provided, which demonstrates basic
multipath TCP detection functionality.  Please see
[cmd/mptcphttp/README.md](https://github.com/mdlayher/mptcp/blob/master/cmd/mptcphttp/README.md)
for details.

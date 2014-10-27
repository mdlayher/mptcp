Usage
=====

To install and use `mptcphttp`, simply run:

```
$ go install github.com/mdlayher/mptcp/...
```

The `mptcphttp` binary is now installed in your `$GOPATH`.  It can be run
as follows:

```
$ mptcphttp -host :8080
mptcphttp: 2014/10/27 18:00:00 binding to: :8080
```

You can now test your multipath TCP capability by simply using `curl` or a
similar tool against `mptcphttp`.

```
$ curl http://localhost:8080/
YES
```

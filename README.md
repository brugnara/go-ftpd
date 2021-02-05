Go FTP
======
> A simple FTP Server, in Go

As a default behaviour, `./public` folder with examples is served.

```bash
go run *.go

#
go run *.go --path /home/$USER
```

Then connect to it:

```bash
nc localhost 8021
```

While this is not a real FTP server, it handles basic commands like:

- cd
- ls

The program accepts a `-path ./public` as the root path to point the server to.

It's guaranteed the client can not escape the root path, as it is intended as
a jail for it.

This is a WiP or even better, an experiment.

Feedbacks are very welcome.

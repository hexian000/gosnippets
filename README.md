# gosnippets

Minimized Go code snippets.

## Snippets

- /
  - `formats`: human readable value formatter, like `2.3k`, `3.45 MiB`, `1d23:59:59` etc.
  - `routines`: joinable goroutine group.
  - `slog`: general purposed logger backed by stdout, syslog or Android logd.
  - `systemd`: systemd daemon notifier, like `sd_notify(3)`.

- /net
  - `net/hlistener`: TCP listener which is hardened for authenticated services.
  - `net/flowmeter.go`: data usage meter for `net.Conn`.

# sing-snell

This fork uses `github.com/metacubex/sing` buffers, addresses and network interfaces. Import it as `github.com/Ember-Moth/sing-snell`. Request completion callbacks and optional vectorised/batch writer discovery are owned by this package; they do not require SagerNet/sing.

[Surge snell](https://kb.nssurge.com/surge-knowledge-base/release-notes/snell) implemented in Go.

## Features

* All complete features except the v5 QUIC proxy.
* Behavior as consistent with the official implementation as possible.
* Performance at least on par with the official implementation.

## Why

Surge believes that being closed-source and not proliferated can keep this protocol
covert, but this is already impossible in 2026. Considering that snell still has
advantages that other random-traffic protocols do not possess, such as multiplexing
support with complete TCP semantics and traffic-characteristic diversity, we made this
instead of reinventing the wheel.

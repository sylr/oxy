Oxy
===

This is a fork of [Oxy](https://github.com/vulcand/oxy/).

Oxy is a Go library with HTTP handlers that enhance HTTP standard library:

* [Buffer](https://pkg.go.dev/abstraction.fr/oxy/v2/buffer) retries and buffers requests and responses
* [Stream](https://pkg.go.dev/abstraction.fr/oxy/v2/stream) passes-through requests, supports chunked encoding with configurable flush interval
* [Forward](https://pkg.go.dev/abstraction.fr/oxy/v2/forward) forwards requests to remote location and rewrites headers
* [Roundrobin](https://pkg.go.dev/abstraction.fr/oxy/v2/roundrobin) is a round-robin load balancer
* [Circuit Breaker](https://pkg.go.dev/abstraction.fr/oxy/v2/cbreaker) Hystrix-style circuit breaker
* [Connlimit](https://pkg.go.dev/abstraction.fr/oxy/v2/connlimit) Simultaneous connections limiter
* [Ratelimit](https://pkg.go.dev/abstraction.fr/oxy/v2/ratelimit) Rate limiter (based on tokenbucket algo)
* [Trace](https://pkg.go.dev/abstraction.fr/oxy/v2/trace) Structured request and response logger

It is designed to be fully compatible with http standard library, easy to customize and reuse.

Status
------

* Initial design is completed
* Covered by tests

Quickstart
----------

Every handler is ``http.Handler``, so writing and plugging in a middleware is easy. Let us write a simple reverse proxy as an example:

Simple reverse proxy
====================

```go
import (
	"net/http"
	"abstraction.fr/oxy/v2/forward"
	"abstraction.fr/oxy/v2/testutils"
)

// Forwards incoming requests to whatever location URL points to, adds proper forwarding headers
fwd, _ := forward.New()

redirect := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
	// let us forward this request to another server
	req.URL = testutils.ParseURI("http://localhost:63450")
	fwd.ServeHTTP(w, req)
})

// that's it! our reverse proxy is ready!
s := &http.Server{
	Addr:    ":8080",
	Handler: redirect,
}

s.ListenAndServe()
```

As a next step, let us add a round robin load-balancer:


```go
import (
	"net/http"
	"abstraction.fr/oxy/v2/forward"
	"abstraction.fr/oxy/v2/roundrobin"
)

// Forwards incoming requests to whatever location URL points to, adds proper forwarding headers
fwd, _ := forward.New()
lb, _ := roundrobin.New(fwd)

lb.UpsertServer(url1)
lb.UpsertServer(url2)

s := &http.Server{
	Addr:    ":8080",
	Handler: lb,
}

s.ListenAndServe()
```

What if we want to handle retries and replay the request in case of errors? `buffer` handler will help:


```go
import (
	"net/http"
	"abstraction.fr/oxy/v2/forward"
	"abstraction.fr/oxy/v2/buffer"
	"abstraction.fr/oxy/v2/roundrobin"
)

// Forwards incoming requests to whatever location URL points to, adds proper forwarding headers

fwd, _ := forward.New()
lb, _ := roundrobin.New(fwd)

// buffer will read the request body and will replay the request again in case if forward returned status
// corresponding to nework error (e.g. Gateway Timeout)
buffer, _ := buffer.New(lb, buffer.Retry(`IsNetworkError() && Attempts() < 2`))

lb.UpsertServer(url1)
lb.UpsertServer(url2)

// that's it! our reverse proxy is ready!
s := &http.Server{
	Addr:    ":8080",
	Handler: buffer,
}
s.ListenAndServe()
```

Logging
=======

As of v2, oxy let's you provide your own logger instead of forcing the use of logrus.
To do so you have to provide a struct that comply with the minimal interface `utils.Logger`.

github.com/sirupsen/logrus
--------------------------

```go
import (
	"abstraction.fr/oxy/v2/cbreaker"
	"github.com/sirupsen/logrus"
)

stdLogger := logrus.StandardLogger()
stdLogger.SetLevel(logrus.DebugLevel)

logrusLogger := stdLogger.WithField("lib", "vulcand/oxy")

cbLogger := cbreaker.Logger(logrusLogger)

cb, err := cbreaker.New(next, "NetworkErrorRatio() > 0.3", cbLogger)
```

go.uber.org/zap
---------------

```go
import (
	"abstraction.fr/oxy/v2/cbreaker"
	"go.uber.org/zap/zap"
	"go.uber.org/zap/zapcore"
)

zapAtomLevel := zap.NewAtomicLevel()
zapAtomLevel.SetLevel(zapcore.DebugLevel)

zapEncoderCfg := zap.NewProductionEncoderConfig()
zapEncoderCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder

zapCore := zapcore.NewCore(zapcore.NewConsoleEncoder(zapEncoderCfg), zapcore.Lock(os.Stdout), zapAtomLevel)
zapLogger := zap.New(zapCore).With(zap.String("lib", "vulcand/oxy"))

zapSugaredLogger := zapLogger.Sugar()
zapDebug := func() bool {
	return zapAtomLevel.Enabled(zapcore.DebugLevel)
}

cbLogger := cbreaker.Logger(zapSugaredLogger)

cb, err := cbreaker.New(next, "NetworkErrorRatio() > 0.3", cbLogger)
```

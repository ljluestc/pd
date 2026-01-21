# Fix Wrong Usage of gRPC Interface in ReportBuckets

## Problem Summary
In `pd/server/grpc_service.go`, the `recv()` methods for `bucketHeartbeatServer` and `tsoServer` were wrapping `io.EOF` errors with a stack trace using `errors.WithStack(err)`. 

This wrapping caused the error comparison `err == io.EOF` to fail at the call sites (e.g., in `ReportBuckets`), treating a graceful stream termination as an unexpected error.

## Root Cause Analysis
The `github.com/pingcap/errors` package's `WithStack` function creates a new error wrapping the original. Standard equality checks (`==`) fail when comparing the wrapped error against the sentinel `io.EOF`.

**Problematic Code (Before):**
```go
req, err := b.stream.Recv()
if err != nil { // io.EOF is caught here
    atomic.StoreInt32(&b.closed, 1)
    return nil, errors.WithStack(err) // Wraps io.EOF, breaking equality checks
}
```

**Caller Expectation:**
```go
// grpc_service.go:1082 (approx)
if err == io.EOF { // Fails because err is struct{stack, cause=io.EOF}
    return nil
}
```

## Solution
I modified `bucketHeartbeatServer.recv()` and `tsoServer.recv()` to explicitly check for `io.EOF` and return it unwrapped.

**Fix Applied:**
```go
req, err := b.stream.Recv()
// ...
if err != nil {
    atomic.StoreInt32(&b.closed, 1)
    if err == io.EOF {
        return nil, io.EOF // Return unwrapped sentinel
    }
    return nil, errors.WithStack(err)
}
```

## Verification

### 1. Reproduction Test
A new test `TestBucketHeartbeatServerRecvEOF` was added to `server/grpc_service_test.go`. It mocks the gRPC stream to return `io.EOF` and asserts that `recv()` returns the exact sentinel error.

**Run the reproduction test:**
```bash
go test -v ./server -run TestBucketHeartbeatServerRecvEOF
```

### 2. Regression Testing
Run the full suite of tests for the server package to ensure no side effects:

```bash
go test -v ./tests/server/...
```

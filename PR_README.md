# Fix wrong usage of gRPC interface in ReportBuckets

## What is changed?

This PR fixes an issue where the `io.EOF` error returned by the gRPC stream in `bucketHeartbeatServer.recv()` was being incorrectly wrapped with a stack trace using `errors.WithStack(err)`. 

## Why is this change needed?

The caller of `recv()`, `ReportBuckets`, explicitly checks for `err == io.EOF` to determine if the stream has ended gracefully. 
```go
if err == io.EOF {
    return nil
}
```
Because the `io.EOF` was wrapped, this equality check failed, treating the EOF as a regular error.

## Verification

Added a new regression test `TestBucketHeartbeatServerRecvEOF` in `server/grpc_service_test.go` which simulates the `Recv` method returning `io.EOF` and asserts that the `recv` wrapper returns it unwrapped.

Run the test with:
```bash
go test -v ./server -run TestBucketHeartbeatServerRecvEOF
```

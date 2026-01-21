package server

import (
	"io"
	"testing"
	"context"

	"github.com/pingcap/kvproto/pkg/pdpb"
	"google.golang.org/grpc"
)

type mockReportBucketsServer struct {
	grpc.ServerStream
	recvFunc func() (*pdpb.ReportBucketsRequest, error)
}

func (m *mockReportBucketsServer) SendAndClose(*pdpb.ReportBucketsResponse) error { return nil }
func (m *mockReportBucketsServer) Recv() (*pdpb.ReportBucketsRequest, error) {
	if m.recvFunc != nil {
		return m.recvFunc()
	}
	return nil, io.EOF
}
func (m *mockReportBucketsServer) Context() context.Context { return context.Background() }

func TestBucketHeartbeatServerRecvEOF(t *testing.T) {
	stream := &mockReportBucketsServer{
		recvFunc: func() (*pdpb.ReportBucketsRequest, error) {
			return nil, io.EOF
		},
	}
	s := &bucketHeartbeatServer{stream: stream}
	_, err := s.recv()

	t.Logf("Error type: %T, value: %v", err, err)
	if err == io.EOF {
		t.Log("Error is strictly equal to io.EOF")
	} else {
		t.Log("Error is NOT strictly equal to io.EOF")
	}
	
	// We expect io.EOF exactly, but due to the bug it is wrapped.
    // This test should FAIL if the bug is present.
	if err != io.EOF {
		t.Fatalf("recv() returned error %v (type %T), want exactly io.EOF", err, err)
	}
}

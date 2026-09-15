package source

import (
	sourcepb "nfxnews/protos/gen/source"

	"google.golang.org/grpc"
)

func New(conn *grpc.ClientConn) sourcepb.SourceServiceClient {
	return sourcepb.NewSourceServiceClient(conn)
}

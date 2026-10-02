package news

import (
	newspb "nfxnews/protos/gen/news"

	"google.golang.org/grpc"
)

func New(conn *grpc.ClientConn) newspb.NewsServiceClient {
	return newspb.NewNewsServiceClient(conn)
}

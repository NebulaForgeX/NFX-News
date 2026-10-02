package crawl

import (
	crawlpb "nfxnews/protos/gen/crawl"

	"google.golang.org/grpc"
)

func New(conn *grpc.ClientConn) crawlpb.CrawlServiceClient {
	return crawlpb.NewCrawlServiceClient(conn)
}

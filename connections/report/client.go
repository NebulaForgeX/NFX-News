package report

import (
	reportpb "nfxnews/protos/gen/report"

	"google.golang.org/grpc"
)

func New(conn *grpc.ClientConn) reportpb.ReportServiceClient {
	return reportpb.NewReportServiceClient(conn)
}

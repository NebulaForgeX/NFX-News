package grpc

import (
	"nfxnews/modules/report/application/resource"
	reportapp "nfxnews/modules/report/application/report"
	grpcHandler "nfxnews/modules/report/interfaces/grpc/handler"
	"nfxnews/pkgs/grpcx/interceptor"
	"nfxnews/pkgs/security/token"
	"nfxnews/pkgs/security/token/servertoken"
	healthpb "nfxnews/protos/gen/common/health"
	reportpb "nfxnews/protos/gen/report"

	"google.golang.org/grpc"
)

type Deps interface {
	AppSvc() *reportapp.Service
	ResourceSvc() *resource.Service
	ServerTokenVerifier() token.Verifier
}

func NewServer(d Deps) *grpc.Server {
	s := grpc.NewServer(grpc.ChainUnaryInterceptor(interceptor.UnaryErrorHandler(), servertoken.UnaryAuthInterceptor(d.ServerTokenVerifier())))
	reportpb.RegisterReportServiceServer(s, grpcHandler.NewReportHandler(d.AppSvc()))
	healthpb.RegisterHealthServiceServer(s, grpcHandler.NewHealthHandler(d.ResourceSvc(), "report"))
	return s
}

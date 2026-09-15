package grpc

import (
	"nfxnews/modules/crawl/application/resource"
	crawlapp "nfxnews/modules/crawl/application/crawl"
	grpcHandler "nfxnews/modules/crawl/interfaces/grpc/handler"
	"nfxnews/pkgs/grpcx/interceptor"
	"nfxnews/pkgs/security/token"
	"nfxnews/pkgs/security/token/servertoken"
	healthpb "nfxnews/protos/gen/common/health"
	crawlpb "nfxnews/protos/gen/crawl"

	"google.golang.org/grpc"
)

type Deps interface {
	AppSvc() *crawlapp.Service
	ResourceSvc() *resource.Service
	ServerTokenVerifier() token.Verifier
}

func NewServer(d Deps) *grpc.Server {
	s := grpc.NewServer(grpc.ChainUnaryInterceptor(interceptor.UnaryErrorHandler(), servertoken.UnaryAuthInterceptor(d.ServerTokenVerifier())))
	crawlpb.RegisterCrawlServiceServer(s, grpcHandler.NewCrawlHandler(d.AppSvc()))
	healthpb.RegisterHealthServiceServer(s, grpcHandler.NewHealthHandler(d.ResourceSvc(), "crawl"))
	return s
}

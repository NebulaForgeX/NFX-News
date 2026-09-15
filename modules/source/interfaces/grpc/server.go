package grpc

import (
	"nfxnews/modules/source/application/resource"
	sourceapp "nfxnews/modules/source/application/source"
	grpcHandler "nfxnews/modules/source/interfaces/grpc/handler"
	"nfxnews/pkgs/grpcx/interceptor"
	"nfxnews/pkgs/security/token"
	"nfxnews/pkgs/security/token/servertoken"
	healthpb "nfxnews/protos/gen/common/health"
	sourcepb "nfxnews/protos/gen/source"

	"google.golang.org/grpc"
)

type Deps interface {
	AppSvc() *sourceapp.Service
	ResourceSvc() *resource.Service
	ServerTokenVerifier() token.Verifier
}

func NewServer(d Deps) *grpc.Server {
	s := grpc.NewServer(grpc.ChainUnaryInterceptor(interceptor.UnaryErrorHandler(), servertoken.UnaryAuthInterceptor(d.ServerTokenVerifier())))
	sourcepb.RegisterSourceServiceServer(s, grpcHandler.NewSourceHandler(d.AppSvc()))
	healthpb.RegisterHealthServiceServer(s, grpcHandler.NewHealthHandler(d.ResourceSvc(), "source"))
	return s
}

package grpc

import (
	"nfxnews/modules/system/application/resource"
	systemapp "nfxnews/modules/system/application/system"
	grpcHandler "nfxnews/modules/system/interfaces/grpc/handler"
	"nfxnews/pkgs/grpcx/interceptor"
	"nfxnews/pkgs/security/token"
	"nfxnews/pkgs/security/token/servertoken"
	healthpb "nfxnews/protos/gen/common/health"
	systemstatepb "nfxnews/protos/gen/system/system_state"

	"google.golang.org/grpc"
)

type Deps interface {
	AppSvc() *systemapp.Service
	ResourceSvc() *resource.Service
	ServerTokenVerifier() token.Verifier
}

func NewServer(d Deps) *grpc.Server {
	s := grpc.NewServer(grpc.ChainUnaryInterceptor(interceptor.UnaryErrorHandler(), servertoken.UnaryAuthInterceptor(d.ServerTokenVerifier())))
	systemstatepb.RegisterSystemStateServiceServer(s, grpcHandler.NewSystemHandler(d.AppSvc()))
	healthpb.RegisterHealthServiceServer(s, grpcHandler.NewHealthHandler(d.ResourceSvc(), "system"))
	return s
}

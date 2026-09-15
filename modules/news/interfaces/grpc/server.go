package grpc

import (
	newsapp "nfxnews/modules/news/application/news"
	"nfxnews/modules/news/application/resource"
	grpcHandler "nfxnews/modules/news/interfaces/grpc/handler"
	"nfxnews/pkgs/grpcx/interceptor"
	"nfxnews/pkgs/security/token"
	"nfxnews/pkgs/security/token/servertoken"
	healthpb "nfxnews/protos/gen/common/health"
	newspb "nfxnews/protos/gen/news"

	"google.golang.org/grpc"
)

type Deps interface {
	AppSvc() *newsapp.Service
	ResourceSvc() *resource.Service
	ServerTokenVerifier() token.Verifier
}

func NewServer(d Deps) *grpc.Server {
	s := grpc.NewServer(grpc.ChainUnaryInterceptor(interceptor.UnaryErrorHandler(), servertoken.UnaryAuthInterceptor(d.ServerTokenVerifier())))
	newspb.RegisterNewsServiceServer(s, grpcHandler.NewNewsHandler(d.AppSvc()))
	healthpb.RegisterHealthServiceServer(s, grpcHandler.NewHealthHandler(d.ResourceSvc(), "news"))
	return s
}

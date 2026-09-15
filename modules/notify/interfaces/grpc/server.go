package grpc

import (
	"nfxnews/modules/notify/application/resource"
	notifyapp "nfxnews/modules/notify/application/notify"
	grpcHandler "nfxnews/modules/notify/interfaces/grpc/handler"
	"nfxnews/pkgs/grpcx/interceptor"
	"nfxnews/pkgs/security/token"
	"nfxnews/pkgs/security/token/servertoken"
	healthpb "nfxnews/protos/gen/common/health"
	notifypb "nfxnews/protos/gen/notify"

	"google.golang.org/grpc"
)

type Deps interface {
	AppSvc() *notifyapp.Service
	ResourceSvc() *resource.Service
	ServerTokenVerifier() token.Verifier
}

func NewServer(d Deps) *grpc.Server {
	s := grpc.NewServer(grpc.ChainUnaryInterceptor(interceptor.UnaryErrorHandler(), servertoken.UnaryAuthInterceptor(d.ServerTokenVerifier())))
	notifypb.RegisterNotifyServiceServer(s, grpcHandler.NewNotifyHandler(d.AppSvc()))
	healthpb.RegisterHealthServiceServer(s, grpcHandler.NewHealthHandler(d.ResourceSvc(), "notify"))
	return s
}

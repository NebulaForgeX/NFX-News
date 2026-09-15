package grpc

import (
	"nfxnews/modules/mcp/application/resource"
	mcpapp "nfxnews/modules/mcp/application/mcp"
	grpcHandler "nfxnews/modules/mcp/interfaces/grpc/handler"
	"nfxnews/pkgs/grpcx/interceptor"
	"nfxnews/pkgs/security/token"
	"nfxnews/pkgs/security/token/servertoken"
	healthpb "nfxnews/protos/gen/common/health"
	mcppb "nfxnews/protos/gen/mcp"

	"google.golang.org/grpc"
)

type Deps interface {
	AppSvc() *mcpapp.Service
	ResourceSvc() *resource.Service
	ServerTokenVerifier() token.Verifier
}

func NewServer(d Deps) *grpc.Server {
	s := grpc.NewServer(grpc.ChainUnaryInterceptor(interceptor.UnaryErrorHandler(), servertoken.UnaryAuthInterceptor(d.ServerTokenVerifier())))
	mcppb.RegisterMCPServiceServer(s, grpcHandler.NewMCPHandler(d.AppSvc()))
	healthpb.RegisterHealthServiceServer(s, grpcHandler.NewHealthHandler(d.ResourceSvc(), "mcp"))
	return s
}

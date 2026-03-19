package server

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	pb "primer/proto"
)

// NewGatewayMux creates a grpc-gateway mux wired to the StagePrimer gRPC endpoint.
func NewGatewayMux(ctx context.Context, grpcEndpoint string) (*runtime.ServeMux, error) {
	mux := runtime.NewServeMux(runtime.WithErrorHandler(gatewayErrorHandler))
	if err := pb.RegisterStagePrimerHandlerFromEndpoint(
		ctx,
		mux,
		grpcEndpoint,
		[]grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())},
	); err != nil {
		return nil, err
	}
	return mux, nil
}

func gatewayErrorHandler(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, r *http.Request, err error) {
	st := status.Convert(err)
	payload := map[string]string{
		"error": st.Message(),
		"code":  st.Code().String(),
	}

	w.Header().Set("Content-Type", marshaler.ContentType(payload))
	w.WriteHeader(runtime.HTTPStatusFromCode(st.Code()))
	_ = marshaler.NewEncoder(w).Encode(payload)
}

// MountGateway attaches an HTTP handler under the provided prefix.
func (s *Server) MountGateway(prefix string, handler http.Handler) {
	if s == nil || s.router == nil || handler == nil {
		return
	}
	if prefix == "" {
		prefix = "/"
	}
	s.router.Any(prefix, gin.WrapH(handler))
	if !strings.HasSuffix(prefix, "/") {
		s.router.Any(prefix+"/*any", gin.WrapH(handler))
	}
}

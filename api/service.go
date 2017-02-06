package api

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	gw "github.com/alanchchen/ethermis/api/ethereum"
	"github.com/tylerb/graceful"
)

type Service interface {
	Start() error
	Stop() error
}

func New(controller Controller) Service {
	opts := []grpc.ServerOption{
		grpc.Creds(credentials.NewClientTLSFromCert(demoCertPool, fmt.Sprintf("%s:%d", host, port)))}

	grpcServer := grpc.NewServer(opts...)
	gw.RegisterEthereumServer(grpcServer, controller)

	mux := runtime.NewServeMux()
	dcreds := credentials.NewTLS(&tls.Config{
		ServerName: fmt.Sprintf("%s:%d", host, port),
		RootCAs:    demoCertPool,
	})
	dopts := []grpc.DialOption{grpc.WithTransportCredentials(dcreds)}
	err := gw.RegisterEthereumHandlerFromEndpoint(context.Background(), mux, fmt.Sprintf("%s:%d", host, port), dopts)
	if err != nil {
		return nil
	}

	s := &service{
		server: &graceful.Server{
			Timeout: 10 * time.Second,
			Server: &http.Server{
				Addr:    fmt.Sprintf("%s:%d", host, port),
				Handler: grpcHandlerFunc(grpcServer, mux),
				TLSConfig: &tls.Config{
					Certificates: []tls.Certificate{*demoKeyPair},
					NextProtos:   []string{"h2"},
				},
			},
		},
	}

	return s
}

// grpcHandlerFunc returns an http.Handler that delegates to grpcServer on incoming gRPC
// connections or otherHandler otherwise. Copied from cockroachdb.
func grpcHandlerFunc(grpcServer *grpc.Server, otherHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO(tamird): point to merged gRPC code rather than a PR.
		// This is a partial recreation of gRPC's internal checks https://github.com/grpc/grpc-go/pull/514/files#diff-95e9a25b738459a2d3030e1e6fa2a718R61
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
		} else {
			otherHandler.ServeHTTP(w, r)
		}
	})
}

// ----------------------------------------------------------------------------

type service struct {
	server *graceful.Server
}

func (s *service) Start() error {
	return s.server.ListenAndServeTLSConfig(s.server.TLSConfig)
}

func (s *service) Stop() error {
	s.server.Stop(10 * time.Second)
	return nil
}

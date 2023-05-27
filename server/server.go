package server

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/intelops/go-common/logging"
	"github.com/intelops/vault-cred/config"
	"github.com/intelops/vault-cred/internal/api"
	"github.com/intelops/vault-cred/proto/pb/vaultcredpb"
	"google.golang.org/grpc"

	"google.golang.org/grpc/reflection"
)

func Start() {
	log := logging.NewLogger()

	log.Info("staring vaultcred server")
	vaultCredServer, err := api.NewVaultCredServ()
	if err != nil {
		log.Fatal("failed to start vaultserv", err)
	}

	cfg, err := config.FetchConfiguration()
	if err != nil {
		log.Fatal("Fetching application configuration failed", err)
	}

	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal("Failed to listen", err)
	}

	grpcServer := grpc.NewServer()
	vaultcredpb.RegisterVaultCredServer(grpcServer, vaultCredServer)
	log.Infof("Server listening at %s", addr)

	// Register reflection service on gRPC server.
	reflection.Register(grpcServer)
	go func() {
		if err := grpcServer.Serve(listener); err != nil {
			log.Error("failed to start vaultserv", err)
		}
	}()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	<-signals

	grpcServer.Stop()
	log.Debug("exiting vaultcred server")
}

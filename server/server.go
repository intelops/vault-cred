package server

import (
	"fmt"
	"github.com/intelops/vault-cred/internal/job"
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

	log.Info("staring vault-cred server")
	vaultCredServer, err := api.NewVaultCredServ(log)
	if err != nil {
		log.Fatal("failed to start vault-cred", err)
	}

	cfg, err := config.FetchConfiguration()
	if err != nil {
		log.Fatal("Fetching application configuration failed", err)
	}

	j, err := job.NewVaultSealWatcher(log, cfg.VaultSealWatchInterval)
	if err != nil {
		log.Fatal("failed to init job", err)
	}

	s := job.NewScheduler(log)
	err = s.AddJob("vault-seal-watcher", j)
	if err != nil {
		log.Fatal("failed to add job", err)
	}
	s.Start()

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
			log.Fatalf("failed to start vault-cred", err)
		}
	}()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	<-signals

	grpcServer.Stop()
	log.Debug("exiting vault-cred server")
}

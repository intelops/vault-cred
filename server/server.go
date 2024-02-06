package server

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/intelops/vault-cred/internal/job"

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

	s := initScheduler(log, cfg)
	s.Start()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	<-signals

	s.Stop()
	grpcServer.Stop()
	log.Debug("exiting vault-cred server")
}

func initScheduler(log logging.Logger, cfg config.Configuration) (s *job.Scheduler) {
	s = job.NewScheduler(log)
	if cfg.VaultSealWatchInterval != "" {
		sj, err := job.NewVaultSealWatcher(log, cfg.VaultSealWatchInterval)
		if err != nil {
			log.Fatal("failed to init seal watcher job", err)
		}
		err = s.AddJob("vault-seal-watcher", sj)
		if err != nil {
			log.Fatal("failed to add seal watcher job", err)
		}
	}

	if cfg.VaultPolicyWatchInterval != "" {
		pj, err := job.NewVaultPolicyWatcher(log, cfg.VaultPolicyWatchInterval)
		if err != nil {
			log.Fatal("failed to init policy watcher job", err)
		}

		err = s.AddJob("vault-policy-watcher", pj)
		if err != nil {
			log.Fatal("failed to add policy watcher job", err)
		}
	}

	if cfg.VaultCredSyncInterval != "" {
		pj, err := job.NewVaultCredSync(log, cfg.VaultCredSyncInterval)
		if err != nil {
			log.Fatal("failed to init cred sync job", err)
		}

		err = s.AddJob("vault-cred-sync", pj)
		if err != nil {
			log.Fatal("failed to add cred sync job", err)
		}
	}
	return
}

package api

import (
	"context"

	"fmt"

	"github.com/intelops/go-common/logging"
	"github.com/intelops/vault-cred/config"
	"github.com/intelops/vault-cred/internal/client"
	"github.com/intelops/vault-cred/proto/pb/vaultcredpb"
	"github.com/kelseyhightower/envconfig"

	"github.com/pkg/errors"
	"google.golang.org/grpc"

	"google.golang.org/grpc/credentials/insecure"
	v1 "k8s.io/api/core/v1"
)

var (
	kadAppRolePrefix = "vault-approle-"
)

type VaultCredServ struct {
	vaultcredpb.UnimplementedVaultCredServer
	conf config.VaultEnv
	log  logging.Logger
}
type Config struct {
	Address    string
	CaCert     string
	Cert       string
	Key        string
	ServicName string

	TlsEnabled bool `envconfig:"TLS_ENABLED" default:"true"`
}

func FetchConfig() (Config, error) {
	cfg := Config{}
	err := envconfig.Process("", &cfg)
	return cfg, err
}

func NewVaultCredServ(log logging.Logger) (*VaultCredServ, error) {

	conf, err := config.GetVaultEnv()
	if err != nil {
		return nil, err
	}

	return &VaultCredServ{
		conf: conf,
		log:  log,
	}, nil
}

func CredentialMountPath() string {
	return "secret"
}

func PrepareCredentialSecretPath(credentialType, credEntityName, credIdentifier string) string {
	return fmt.Sprintf("%s/%s/%s", credentialType, credEntityName, credIdentifier)
}

func (v *VaultCredServ) GetCredential(ctx context.Context, request *vaultcredpb.GetCredentialRequest) (*vaultcredpb.GetCredentialResponse, error) {
	vc, err := client.NewVaultClientForServiceAccount(ctx, v.log, v.conf)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to initiize vault client")
	}

	secretPath := PrepareCredentialSecretPath(request.CredentialType, request.CredEntityName, request.CredIdentifier)
	credentail, err := vc.GetCredential(ctx, CredentialMountPath(), secretPath)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to get credential")
	}

	v.log.Infof("get credential request processed for %s", secretPath)
	return &vaultcredpb.GetCredentialResponse{Credential: credentail}, nil
}

func (v *VaultCredServ) PutCredential(ctx context.Context, request *vaultcredpb.PutCredentialRequest) (*vaultcredpb.PutCredentialResponse, error) {
	vc, err := client.NewVaultClientForServiceAccount(ctx, v.log, v.conf)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to initiize vault client")
	}

	secretPath := PrepareCredentialSecretPath(request.CredentialType, request.CredEntityName, request.CredIdentifier)
	err = vc.PutCredential(ctx, CredentialMountPath(), secretPath, request.Credential)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to write credential")
	}

	v.log.Infof("write credential request processed for %s", secretPath)
	return &vaultcredpb.PutCredentialResponse{}, nil
}

func (v *VaultCredServ) DeleteCredential(ctx context.Context, request *vaultcredpb.DeleteCredentialRequest) (*vaultcredpb.DeleteCredentialResponse, error) {
	vc, err := client.NewVaultClientForServiceAccount(ctx, v.log, v.conf)
	if err != nil {
		return nil, err
	}

	secretPath := PrepareCredentialSecretPath(request.CredentialType, request.CredEntityName, request.CredIdentifier)
	err = vc.DeleteCredential(ctx, CredentialMountPath(), secretPath)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to delete credential")
	}

	v.log.Infof("delete credential request processed for %s", secretPath)
	return &vaultcredpb.DeleteCredentialResponse{}, nil
}

func (v *VaultCredServ) ConfigureVaultSecret(ctx context.Context, request *vaultcredpb.ConfigureVaultSecretRequest) (*vaultcredpb.ConfigureVaultSecretResponse, error) {
	v.log.Infof("Configure Vault Secret Request recieved for secret ", request.SecretName)

	secretPathsData := map[string]string{}
	secretPaths := []string{}
	for _, secretPathData := range request.SecretPathData {
		secretPathsData[secretPathData.SecretKey] = secretPathData.SecretPath
		secretPaths = append(secretPaths, secretPathData.SecretPath)
	}

	appRoleName := kadAppRolePrefix + request.SecretName
	token, err := v.GetAppRoleToken(appRoleName, secretPaths)
	if err != nil {
		v.log.Errorf("failed to create app role token for %s, %v", appRoleName, err)
		return &vaultcredpb.ConfigureVaultSecretResponse{Status: vaultcredpb.StatusCode_INTERNRAL_ERROR}, err
	}

	k8sclient, err := client.NewK8SClient(v.log)
	if err != nil {
		v.log.Errorf("failed to initalize k8s client, %v", err)
		return &vaultcredpb.ConfigureVaultSecretResponse{Status: vaultcredpb.StatusCode_INTERNRAL_ERROR}, err
	}

	cred := map[string][]byte{"token": []byte(token)}
	vaultTokenSecretName := "vault-token-" + request.SecretName
	err = k8sclient.CreateOrUpdateSecret(ctx, request.Namespace, vaultTokenSecretName, v1.SecretTypeOpaque, cred, nil)
	if err != nil {
		v.log.Errorf("failed to create cluter vault token secret, %v", err)
		return &vaultcredpb.ConfigureVaultSecretResponse{Status: vaultcredpb.StatusCode_INTERNRAL_ERROR}, err
	}

	vaultAddressStr := fmt.Sprintf(v.conf.Address, request.DomainName)
	secretStoreName := "ext-store-" + request.SecretName
	err = k8sclient.CreateOrUpdateSecretStore(ctx, secretStoreName, request.Namespace, vaultAddressStr, vaultTokenSecretName, "token")
	if err != nil {
		v.log.Errorf("failed to create cluter vault token secret, %v", err)
		return &vaultcredpb.ConfigureVaultSecretResponse{Status: vaultcredpb.StatusCode_INTERNRAL_ERROR}, err
	}
	v.log.Infof("created secret store %s/%s", request.Namespace, secretStoreName)

	externalSecretName := "ext-secret-" + request.SecretName
	err = k8sclient.CreateOrUpdateExternalSecret(ctx, externalSecretName, request.Namespace, secretStoreName,
		request.SecretName, "", secretPathsData)
	if err != nil {
		v.log.Errorf("failed to create vault external secret, %v", err)
		return &vaultcredpb.ConfigureVaultSecretResponse{Status: vaultcredpb.StatusCode_INTERNRAL_ERROR}, err
	}
	v.log.Infof("created external secret %s/%s", request.Namespace, externalSecretName)
	return &vaultcredpb.ConfigureVaultSecretResponse{Status: vaultcredpb.StatusCode_OK}, nil
}

func (v *VaultCredServ) GetAppRoleToken(appRoleName string, credentialPaths []string) (string, error) {
	addr := "vault:8200"
	vc, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return "", fmt.Errorf("failed to connect vauld-cred server, %v", err)
	}
	vcClient := vaultcredpb.NewVaultCredClient(vc)

	tokenData, err := vcClient.CreateAppRoleToken(context.Background(), &vaultcredpb.CreateAppRoleTokenRequest{
		AppRoleName: appRoleName,
		SecretPaths: credentialPaths,
	})
	if err != nil {
		return "", fmt.Errorf("failed to generate app role token for %s, %v", appRoleName, err)
	}
	return tokenData.Token, nil
}

package job

import (
	//	"github.com/argoproj/argo-cd/v2/util/notification/k8s"
	"fmt"

	"strings"

	"github.com/intelops/go-common/logging"
	"github.com/intelops/vault-cred/config"
	"github.com/intelops/vault-cred/internal/client"
	// "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
)

type VaultPolicyWatcher struct {
	log       logging.Logger
	conf      config.VaultEnv
	frequency string
}

func NewVaultPolicyWatcher(log logging.Logger, frequency string) (*VaultPolicyWatcher, error) {
	conf, err := config.GetVaultEnv()
	if err != nil {
		return nil, err
	}
	return &VaultPolicyWatcher{
		log:       log,
		frequency: frequency,
		conf:      conf,
	}, nil
}

func (v *VaultPolicyWatcher) CronSpec() string {
	return v.frequency
}

func (v *VaultPolicyWatcher) Run() {
	v.log.Info("started vault policy watcher")
	vault, err := client.NewVaultClientForVaultToken(v.log, v.conf)
	if err != nil {
		v.log.Errorf("%s", err)
	}

	newclient, _ := client.NewK8SClient(v.log)

	//Get all the configmaps
	configmap, err := client.FilterVaultConfigMaps(newclient.Client)

	if err != nil {
		v.log.Errorf("error while filtering: %v", err)
	}

	var saDataList []client.SAData
	//for each configmap retrieve the SA and Credential Access
	//create rolename and policy name based on each SA

	for _, cm := range configmap {
		serviceAccount := cm.Data["serviceAccount"]
		credentialAccess := cm.Data["credentialAccess"]
		roleName := serviceAccount + "-role"
		policyName := serviceAccount + "-policy"
		v.log.Info(roleName)
		v.log.Info(policyName)
		accessList := strings.Split(credentialAccess, "\n")
		var credentialAccessList []string
		for _, access := range accessList {
			parts := strings.SplitN(access, ":", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])
				credentialAccessList = append(credentialAccessList, fmt.Sprintf("%s: %s", key, value))
			}
		}

		saData := client.SAData{
			ServiceAccount:   serviceAccount,
			CredentialAccess: credentialAccessList,
		}
		// fmt.Println("SAData", saData)
		// fmt.Println("CredAccess", saData.CredentialAccess)
		saDataList = append(saDataList, saData)
		v.log.Infof("Service Account", saDataList)
		//based on the configmap ,prepare PolicyData
		policyData := client.PreparePolicyData(saDataList)
		v.log.Infof("PolicyData", policyData)
		res, err := vault.CheckVaultPolicyExistsAndEqual(policyName, policyData)
		if err != nil {
			v.log.Errorf("failed to checking the policy: %v", err)
		}
		if !res {
			vault.DeleteVaultPolicy(policyName)
			err = vault.CreatePolicy(policyName, policyData)
			if err != nil {
				v.log.Error("Error while creating policy", err)
			}
			vault.DeleteVaultRole(roleName)
			err = vault.CreateVaultRole(roleName, policyName, serviceAccount)
			if err != nil {
				v.log.Error("Error while creating vault role", err)
			}
		} 

	}

	// Svc iam creates vault-policy-iam
	// Svc capten creates vault-policy-capen
	// get config maps with prefix "vault-policy-<>"
	// for each config map
	// read SA name, credentail-acess : [] {"credential-type", "access type"}
	//"serviceAccount" : "iam-sa"
	//"credentailAccess" : "system-cred:read, client-cert:admin, xyz: read"

	// prepare role "iam-sa-role" and policy "iam-sa-policy"  with path
	// paths are based on credential types
	// secret/data/system-cred/* with read,  secret/data/client-cert/* with read, write, list, delete
	// secret/data/xyz/* read
	// create policy

	// is policy exists, is any differnet from vault-policy data
	// update delta
}

![VAULT-CRED](.readme_assets/Vault-Cred.png)   


The open source platform for read,write and delete credentials on vault server.

[![Docker Image CI](https://github.com/intelops/vault-cred/actions/workflows/docker-image.yaml/badge.svg)](https://github.com/intelops/vault-cred/actions/workflows/docker-image.yaml)
[![CodeQL](https://github.com/intelops/vault-cred/actions/workflows/github-code-scanning/codeql/badge.svg)](https://github.com/intelops/vault-cred/actions/workflows/github-code-scanning/codeql)
[![Helm Chart Publish Image CI](https://github.com/intelops/vault-cred/actions/workflows/helm_release.yml/badge.svg)](https://github.com/intelops/vault-cred/actions/workflows/helm_release.yml)

<hr>

## Vault-Cred

To store any sensitive data like service-based credentials (username,Password),certificate based credential or any generic type of credential like token based credential and so on.

## Table of Contents
- [How to install and run Vault-Cred](#how-to-install-and-run-Vault-Cred)
- [How VaultCred works](#how-VaultCred-works)
- [Use Cases](#use-cases)
- [Contributing](#contributing)
- [Code of Conduct](#code-of-conduct)
- [Support](#support)
- [License](#license)
- [Join our Slack channel](#join-our-slack-channel)

## How to install and run vault-cred

```bash
git clone git@github.com:intelops/vault-cred.git && cd vault-cred/charts/vault-cred
```

```bash
helm install vault-cred -n default .
```
#### Prerequisites
* A Kubernetes cluster 
* Helm binary
* Vault    (https://github.com/hashicorp/vault-helm) 

## How VaultCred works

After installation of vault-cred ,it initializes the vault if vault is not initialized and also generate unseal keys and root token.If vault is already deployed and initialized,then it stores the unsealkeys and root token in a secret named vault-server with the key prefix unsealkey and roottoken.

```bash     
kubectl get po -n default

NAME                                   READY   STATUS    RESTARTS   AGE
vault-cred-5777789576-hpg9r            1/1     Running   0          20h
vault-0                                1/1     Running   0          5h46m
```
You can check the logs by giving below command

```bash     
kubectl logs -f vault-cred-5777789576-hpg9r -n default
```


Vault-Cred monitors the vault service frequently,and monitors whether the vault is unsealed or not.If vault is sealed,vault-cred automatically unseal it by taking the keys from the secret vault-server.


Vault-Cred can also automate the creation of vault policy and role.Vault-Cred continuously monitors for configmap with the prefix vault-policy and vault-role.If it found any configmap,name with the prefix vault-policy,then creates vault-policy with the data and  similarly if it found any configmap ,name with the prefix vault-role,then it creates vault-role with the data.

Sample Vault-policy-Configmap given below:

```yaml
apiVersion: v1
data:
  policyData: |
    path "secret/data/cluster-cred/*" {
      capabilities = ["read"]
    }
    path "auth/kubernetes/login" {
      capabilities = ["create","read","update"]
    }
  policyName: vault-policy-cluster-read
kind: ConfigMap
metadata:
  name: vault-policy-cluster-read
  namespace: default
```

Sample Vault-role-Configmap given below:

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: vault-role-sample
data:
  roleName: vault-role-sample
  policyNames: vault-policy-cluster-read
  servieAccounts: default
  servieAccountNameSpaces: default
```

If any sensitive data needed to be stored in vault,you can store the sensitive data in a secret named vault-cred-sync-data .For storing service based credential,you can use below format in the secret

```bash
SERVICE-CRED-<uniquevalue>: `echo {"entityName":"db", "userName":"xxx","password":"xxx"} | base64 -w 0`
```

For certificate based credential ,use the below format in storing the credential in the secret
```bash
 CERTS-<uniquevalue>: `echo '{"entityName":"xxx", "certIndetifier":"xxx","caCert":"xxx", "cert": "xxx", "key":"xxx"}' | base64 -w 0`
```

for storing generic credential,use the below format in storing the credential in the secret
```bash
GENERIC-1: `echo '{"credentialType":"cluster-cred","entityName":"xxx", "credIndetifier":"xxx", "credential":{"token":"xxx","id":"1"}}' | base64 -w 0`
```
With the above mentioned echo command,encode and create a secret with the key prefix generic,service-cred,certs .

From this secret,vault-cred stores the credential,taking the credentialtype,entityname and credIdentifier as a secret path .


## Use Cases

* Automate Vault Unsealing
* Continuous monitoring of ConfigMap to create Vault Policy and Vault Role.
* Storing service based credential,certificate and any generic credential.









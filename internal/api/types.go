package api

const (
	vaultPolicyReadPath  = `path "secret/data/%s" {capabilities = ["read"]}`
	vaultPolicyWritePath = `path "secret/data/%s" {capabilities = ["create","read","update","delete","list"]}`
)

package vault

import (
	"github.com/hashicorp/vault/api"
)

func (v *VaultClient) CheckAndMountKVMount(mountPath string) error {
	sysMounts, err := v.c.Sys().ListMounts()
	if err != nil {
		return err
	}

	mount, found := sysMounts["secret/"]
	if found && mount.Options["version"] == "2" {
		v.log.Debug("kv secret mount with version 2 mounted")
		return nil
	}

	mountInput := &api.MountInput{
		Type: "kv-v2",
	}
	err = v.c.Sys().Mount("secret/", mountInput)
	if err != nil {
		return err
	}
	return nil
}
